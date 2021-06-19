FROM golang:1.15.2-alpine3.12 as builder

ARG GIT_TAG_NAME
ARG LD_FLAGS="-s -w -X github.com/4nth0/golem/pkg/version.Version=$GIT_TAG_NAME"

WORKDIR /project
COPY . .
RUN GOOS=linux go build -ldflags "$LD_FLAGS" \
    -o golem ./

FROM alpine:3.12.0
RUN apk --update add --no-cache ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /project/golem golem

CMD ./golem run
