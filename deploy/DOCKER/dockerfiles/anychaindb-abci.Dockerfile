#
# Copyright (C) 2018 eeonevision
#
# Permission is hereby granted, free of charge, to any person obtaining a copy of
# this software and associated documentation files (the "Software"), to deal in
# the Software without restriction, including without limitation the rights to
# use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
# the Software, and to permit persons to whom the Software is furnished to do so,
# subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
# FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
# COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
# IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
# CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
#

# Description:
#   Builds an image with AnychainDB business logic installed.
#
# Run:
#   $ docker run anychaindb-abci

# Stage Zero. Build sources
FROM golang:latest

RUN mkdir -p $GOPATH/src/github.com/anychaindb/anychaindb && \
	go get github.com/tools/godep && \
	go get github.com/tinylib/msgp && \
    cd $GOPATH/src/github.com/anychaindb/anychaindb && \
    git clone https://github.com/anychaindb/anychaindb . && \
    git checkout master && \
	cd $GOPATH/src/github.com/anychaindb/anychaindb/state && \
	go generate && \
	cd $GOPATH/src/github.com/anychaindb/anychaindb/transaction && \
	go generate && \
	cd $GOPATH/src/github.com/anychaindb/anychaindb && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 godep go install ./... && \
    cd - && \
    rm -rf $GOPATH/src/github.com/anychaindb/anychaindb

# Stage One. AnychainDB ABCI
FROM alpine:latest

RUN apk add --no-cache ca-certificates bash curl jq

WORKDIR /usr/bin/

COPY --from=0 /go/bin/tmlc-abci .

ENTRYPOINT [ "tmlc-abci", "--help" ]
