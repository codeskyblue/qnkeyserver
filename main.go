package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/qiniu/api.v6/conf"
	"github.com/qiniu/api.v6/rs"
	redis "gopkg.in/redis.v3"
)

func genUptoken(bucket string, key string) string {
	policy := rs.PutPolicy{
		Scope: bucket + ":" + key,
	}
	policy.Expires = uint32(time.Now().Unix()) + 1800
	policy.FsizeLimit = 20 << 20 // 20M
	return policy.Token(nil)
}

var rdx *redis.Client

func init() {
	conf.ACCESS_KEY = os.Getenv("QINIU_AK")
	conf.SECRET_KEY = os.Getenv("QINIU_SK")
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	rdx = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	if _, err := rdx.Ping().Result(); err != nil {
		log.Fatal(err)
	}
}

const (
	DEFAULT_BUCKET = "gobuild5"
	DEFAULT_DOMAIN = "qn-gobuild5.qbox.me"
)

func main() {
	http.HandleFunc("/uptoken", func(w http.ResponseWriter, r *http.Request) {
		privateToken := r.FormValue("private_token")
		key := r.FormValue("key")
		//log.Println(privateToken)
		bucket := r.FormValue("bucket")
		if bucket == "" {
			bucket = DEFAULT_BUCKET
		}
		if privateToken == os.Getenv("APP_TOKEN") {
			io.WriteString(w, genUptoken(bucket, key))
			return
		}
		username := rdx.Get("token:" + privateToken + ":user").Val()
		if username == "" {
			http.Error(w, "token not exists", 500)
			return
		}
		if !strings.HasPrefix(key, "/gorelease/") {
			http.Error(w, "key prefix must be /gorelease/", 500)
			return
		}
		parts := strings.Split(key, "/")
		//parts := filepath.SplitList(key)
		if len(parts) < 4 {
			http.Error(w, "key too short", 500)
			return
		}
		hkey := "orgs:" + username + ":repos"
		repoPath := parts[2] + "/" + parts[3]
		log.Println(repoPath)
		if !rdx.HExists(hkey, repoPath).Val() {
			http.Error(w, "repo is not ownered by this token, or not updated", 500)
			return
		}
		// todo: update gorelease repo
		rdx.HSetNX(hkey, repoPath, DEFAULT_BUCKET)
		io.WriteString(w, genUptoken(DEFAULT_BUCKET, key))
		return
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Qiniu Key Server, request /uptoken to get uptoken")
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	log.Println("Listen port", port)
	http.ListenAndServe(":"+port, nil)
}
