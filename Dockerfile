FROM golang:1.10-alpine3.7

## Prepare system tools dependencies
RUN set -xe \
  && DEP_VERSION=0.4.1 \
  && apk add --no-cache git mercurial bzr curl make gcc musl-dev \
  && curl -fL https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-linux-amd64 -o /usr/local/bin/dep \
  && chmod a+x /usr/local/bin/dep

ENV SRC_DIR=/go/src/github.com/virtengine/vertice
WORKDIR ${SRC_DIR}

## Install library dependencies
COPY Gopkg.* ${SRC_DIR}/
RUN dep ensure -vendor-only -v
## Copy all the source code to directory
COPY . ${SRC_DIR}/

ARG GO_EXTRAFLAGS='-v'

## Testing
RUN go test ${GO_EXTRAFLAGS} ./...

## Building
RUN set -xe \
  && export BUILD_DATE=$(date +%Y-%m-%d_%H:%M:%S%Z) \
  && export COMMIT_HASH=$(git rev-parse HEAD) \
  && export LIBGO_COMMIT_HASH=$(cd vendor/github.com/virtengine/libgo && git rev-parse HEAD) \
  && go build ${GO_EXTRAFLAGS} \
  -ldflags="-X main.date=${BUILD_DATE} -X main.commit=${COMMIT_HASH}_lib_${LIBGO_COMMIT_HASH}" \
  -o vertice ./cmd/vertice

## Command to start server
CMD [ "/go/src/github.com/virtengine/vertice/vertice", \
  "-v", \
  "start", \
  "--config", "/go/src/github.com/virtengine/vertice/conf/vertice.conf" \
  ]
