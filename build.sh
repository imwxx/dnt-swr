#/bin/bash

mkdir -p www/dnt-swr/{bin,logs}
CGO_ENABLED=0 GOOS=linux go build -i -o www/dnt-swr/bin/dnt-api main.go
strip www/dnt-swr/bin/dnt-api
cp dntact.sh www/dnt-swr/bin/
cp  -r script www/dnt-swr/ && chmod +x www/dnt-swr/bin/* && cp processor.py www/dnt-swr/
tar -cvf www.tar www && rm -rf www
