## About
```
this program was use for sync docker img from harbor to HuaWei Cloud SWR
build it by golang, depend on docker, redis and mysql
dnt-api:
    accept notification from harbor when the imgages pushed, and set the notification into redis queue.
processor:
    get job from redis and sync job by processor, the processor dep more machines where it can connect to redis, make rsync job quickly
```
### 1. build it
mkdir /data/gopath/src
export GOPATH="/data/gopath/"
cd /data/gopath/src/ && git clone https://github.com/imwxx/dnt-swr && cd dnt-swr && bash build.sh
then you can use www.tar to deploy it.

### 2. deployment
deploy it by yourself or use www.tar.
#### 2.1 install redis & mysql
source db.sql;
#### 2.2 run api
change the conf/appconf.yml by yourself
/www/dnt/bin/dnt-api -f /www/dnt/conf/appconf.yml
#### 2.3 run processor, and the machine should install docker ,and must be docker login harbor and swr, check it in /root/.docker/config.json
python /www/dnt-swr/processor.py
#### 2.4 show you my env
```
[root@localhost dnt-swr]# pwd
/www/dnt-swr
[root@localhost dnt-swr]# tree
.
├── bin
│   ├── dntact.sh
│   ├── dnt-api
├── conf
│   └── appconf.yaml
├── logs
│   ├── dnt-swr.info.log
│   ├── dnt-swr.info.log.1
│   ├── dnt-swr.info.log.10
│   ├── dnt-swr.info.log.2
│   ├── dnt-swr.info.log.3
│   ├── dnt-swr.info.log.4
│   ├── dnt-swr.info.log.5
│   ├── dnt-swr.info.log.6
│   ├── dnt-swr.info.log.7
│   ├── dnt-swr.info.log.8
│   ├── dnt-swr.info.log.9
│   ├── dnt-swr.info.stderr
│   └── dnt-swr-processor.info.log
├── processor.py
└── script
    └── img.sh

6 directories, 29 files
```
### 3. about api
get more infomation from dntact.sh or source code

### 4. add notification into harbor like this
```
[root@localhost config]# grep -A 10 swr-harbor_sync_tool registry/config.yml 
  - name: swr-harbor_sync_tool
    disabled: false
    url: http://192.168.1.123:9999/service/notifications
    timeout: 3000ms
    threshold: 5
    backoff: 1s

[root@localhost config]# pwd
/www/harbor/common/config
```
