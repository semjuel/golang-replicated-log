FROM golang:1.15-alpine

ENV SERVICE_PATH="replicated_log"

RUN apk update && apk upgrade \
	&& apk add --no-cache ca-certificates \
	bash \
	git

# Install go package manager - dep.
RUN go get -d -u github.com/golang/dep && \
    cd $(go env GOPATH)/src/github.com/golang/dep && \
    DEP_LATEST=$(git describe --abbrev=0 --tags) && \
    git checkout $DEP_LATEST && \
    go install -ldflags="-X main.version=$DEP_LATEST" ./cmd/dep

RUN mkdir -p $GOPATH/src/${SERVICE_PATH}

# COPY source files
COPY . $GOPATH/src/${SERVICE_PATH}
COPY ./.env.example $GOPATH/src/${SERVICE_PATH}/.env

WORKDIR $GOPATH/src/${SERVICE_PATH}

# Install all dependencies.
RUN dep ensure

CMD ["go", "run", "secondary/main.go"]
