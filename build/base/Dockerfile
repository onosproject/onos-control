ARG ONOS_BUILD_VERSION=stable

FROM onosproject/golang-build:$ONOS_BUILD_VERSION
ENV GO111MODULE=on
COPY . /go/src/github.com/onosproject/onos-control
RUN cd /go/src/github.com/onosproject/onos-control && GOFLAGS=-mod=vendor make build
