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
#   Builds an image with AnychainDB node (Tendermint) installed.
#

# Stage Zero. Build sources
FROM golang:latest

RUN mkdir -p /go/src/github.com/tendermint/tendermint && \
    cd /go/src/github.com/tendermint/tendermint && \
    git clone https://github.com/tendermint/tendermint . && \
    git checkout master && \
    make get_tools && \
    make get_vendor_deps && \
    make install && \
    cd - && \
    rm -rf /go/src/github.com/tendermint/tendermint

# Stage One. Tendermint image for AnychainDB platform
FROM alpine:latest

RUN apk add --no-cache ca-certificates bash curl jq

ENV DATA_ROOT /tendermint
ENV TMHOME $DATA_ROOT

RUN addgroup tmuser && \
    adduser -S -G tmuser tmuser

RUN mkdir -p "$DATA_ROOT/config"

ARG config

# validator genesis and config files if exists
COPY $config ${DATA_ROOT}/config/

VOLUME $DATA_ROOT

WORKDIR /usr/bin/

COPY --from=0 /go/bin/tendermint .

ENTRYPOINT [ "tendermint", "version" ]
