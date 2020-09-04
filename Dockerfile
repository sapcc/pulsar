FROM golang:1.13.4-alpine3.10 as builder
RUN apk add --no-cache make git

WORKDIR /go/src/github.com/sapcc/pulsar

# Copy vendor.
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# Copy go source.
COPY main.go main.go
COPY cmd/ cmd/
COPY pkg/ pkg/

# Copy misc.
COPY Makefile Makefile
COPY VERSION VERSION
COPY .git/ .git/
RUN make all

FROM alpine:3.10
MAINTAINER Arno Uhlig <arno.uhlig@sap.com>
LABEL source_repository="https://github.com/sapcc/pulsar"

RUN apk add --no-cache ca-certificates curl tini bash
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.15.0/bin/linux/amd64/kubectl && chmod +x ./kubectl && mv ./kubectl /usr/local/bin/kubectl && kubectl version --client
COPY --from=builder /go/src/github.com/sapcc/pulsar/bin/linux/pulsar /usr/local/bin/
ENTRYPOINT ["tini", "--"]
CMD ["pulsar"]
