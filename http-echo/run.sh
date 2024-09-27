#! /bin/sh

set -e

# 1. build the http server
go build  -o http-echo -ldflags "-w -extldflags '-static'" -tags netgo http-echo.go

# 2. build the fenc program
go build -o fenc -ldflags "-w -extldflags '-static'"  fenc.go

# 3. generate key
openssl rand 32 > key.bin

# 4. encrypt the http server
./fenc -file http-echo -key key.bin -operation encryption

# 5. rm the biaries
rm -f http-echo

docker build -t quay.io/eesposit/http-echo-ema .

docker push quay.io/eesposit/http-echo-ema