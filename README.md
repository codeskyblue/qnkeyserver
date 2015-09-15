# qnkeyserver
qiniu key server

Need environ

	QINIU_AK
	QINIU_SK
	APP_TOKEN

	REDIS_ADDR=localhost:6379
	REDIS_PASSWORD=abcdefg

This repo is used by [gorelease](https://github.com/gorelease/gorelease)

## Get uptoken

	http GET qntoken.herokuapp.com/uptoken private_token==abcdefg bucket==gorelease key==/gorelease/codeskyblue/gosuv/master.zip
