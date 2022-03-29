# SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0

ARG ONOS_CONTROL_BASE_VERSION=latest

FROM onosproject/onos-control-base:$ONOS_CONTROL_BASE_VERSION as base

FROM alpine:3.9
RUN apk add libc6-compat

USER nobody

COPY --from=base /go/src/github.com/onosproject/onos-control/build/_output/onos-control /usr/local/bin/onos-control

ENTRYPOINT ["onos-control"]
