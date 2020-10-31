FROM golang:1.15.3-alpine AS builder

ADD control-plane/go.mod control-plane/go.sum /build/control-plane/
ADD protocol/go.mod protocol/go.sum /build/protocol/
ADD common/go.mod common/go.sum /build/common/

ENV CGO_ENABLED=0

WORKDIR /build/control-plane
RUN go mod download

COPY control-plane/ /build/control-plane/
COPY protocol/ /build/protocol
COPY common/ /build/common
RUN go build -o /build/kuly .

FROM scratch

COPY --from=builder /build/kuly /

CMD ["/kuly"]
