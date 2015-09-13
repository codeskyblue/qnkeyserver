package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/qiniu/api.v6/conf"
	"github.com/qiniu/api.v6/rs"
)

func genUptoken(bucket string, key string) string {
	policy := rs.PutPolicy{
		Scope: bucket + ":" + key,
	}
	policy.Expires = uint32(time.Now().Unix()) + 1800
	policy.FsizeLimit = 20 << 20 // 20M
	return policy.Token(nil)
}

func init() {
	conf.ACCESS_KEY = os.Getenv("QINIU_AK")
	conf.SECRET_KEY = os.Getenv("QINIU_SK")
}

func main() {
	http.HandleFunc("/uptoken", func(w http.ResponseWriter, r *http.Request) {
		privateToken := r.FormValue("private_token")
		log.Println(privateToken)
		bucket := r.FormValue("bucket")
		key := r.FormValue("key")
		if privateToken == os.Getenv("APP_TOKEN") {
			io.WriteString(w, genUptoken(bucket, key))
			return
		}
		http.Error(w, "auth denied", 500)
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
