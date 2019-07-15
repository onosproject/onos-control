ARG ONOS_CONTROL_BASE_VERSION=latest

FROM onosproject/onos-control-base:$ONOS_CONTROL_BASE_VERSION as base

FROM golang:1.12.6-alpine3.9 as debugBuilder

RUN apk upgrade --update --no-cache && apk add git && \
    go get -u github.com/go-delve/delve/cmd/dlv && \
    go get -u github.com/atomix/atomix-cli/cmd/atomix

FROM alpine:3.9

RUN apk upgrade --update --no-cache && apk add bash bash-completion libc6-compat

COPY --from=base /go/src/github.com/onosproject/onos-control/build/_output/onos /usr/local/bin/onos
COPY --from=base /go/src/github.com/onosproject/onos-control/build/_output/onos-control-debug /usr/local/bin/onos-control
COPY --from=debugBuilder /go/bin/dlv /usr/local/bin/dlv
COPY --from=debugBuilder /go/bin/atomix /usr/local/bin/atomix

RUN echo "#!/bin/sh" >> /usr/local/bin/onos-control-debug && \
    echo "atomix controller set \$ATOMIX_CONTROLLER" >> /usr/local/bin/onos-control-debug && \
    echo "atomix config set namespace \$ATOMIX_NAMESPACE" >> /usr/local/bin/onos-control-debug && \
    echo "atomix config set app \$ATOMIX_APP" >> /usr/local/bin/onos-control-debug && \
    echo "dlv --listen=:40000 --headless=true --accept-multiclient=true --api-version=2 exec /usr/local/bin/onos-control -- \"\$@\"" >> /usr/local/bin/onos-control-debug && \
    chmod +x /usr/local/bin/onos-control-debug

RUN addgroup -S onos-control && adduser -S -G onos-control onos-control
USER onos-control
WORKDIR /home/onos-control

ENV ATOMIX_CONTROLLER=atomix-controller.kube-system.svc.cluster.local:5679
ENV ATOMIX_NAMESPACE=default

RUN onos init && \
    cp /etc/profile /home/onos-control/.bashrc && \
    onos completion bash > /home/onos-control/.onos/bash_completion.sh && \
    echo "source /home/onos-control/.onos/bash_completion.sh" >> /home/onos-control/.bashrc && \
    onos config set address 127.0.0.1:5150 && \
    atomix init && \
    atomix completion bash > /home/onos-control/.atomix/bash_completion.sh && \
    echo "source /home/onos-control/.atomix/bash_completion.sh" >> /home/onos-control/.bashrc

ENTRYPOINT ["onos-control-debug"]
