ARG ONOS_CONTROL_BASE_VERSION=latest

FROM onosproject/onos-control-base:$ONOS_CONTROL_BASE_VERSION as base

FROM alpine:3.9
RUN apk add libc6-compat

USER nobody

COPY --from=base /go/src/github.com/onosproject/onos-control/build/_output/onos-control /usr/local/bin/onos-control

ENTRYPOINT ["onos-config"]
