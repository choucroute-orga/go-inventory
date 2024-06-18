FROM golang:1.22-alpine AS builder

WORKDIR /usr/src/app/

RUN apk update

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=./go-inventory/go.sum,target=go.sum \
    --mount=type=bind,source=./go-inventory/go.mod,target=go.mod \
    go mod download

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,rw,source=./go-inventory,target=. \
    go build -ldflags "-s -w" -o /go/bin/inventory/ ./

FROM alpine

WORKDIR /usr/src/app/

COPY --from=builder /go/bin/inventory/ ./

EXPOSE ${INVENTORY_PORT}
ENTRYPOINT [ "./inventory" ]