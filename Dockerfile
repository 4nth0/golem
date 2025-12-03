FROM golang:1.16.15-alpine3.15 AS builder

ARG GIT_TAG_NAME
ARG LD_FLAGS="-s -w -X github.com/4nth0/golem/pkg/version.Version=$GIT_TAG_NAME"

WORKDIR /project
COPY . .
RUN GOOS=linux go build -ldflags "$LD_FLAGS" \
    -o golem ./

FROM alpine:3.12.0

RUN apk --update add --no-cache ca-certificates tzdata

RUN addgroup -g 1000 golem && \
    adduser -D -u 1000 -G golem golem

WORKDIR /home/golem
COPY --from=builder /project/golem golem
RUN chown golem:golem golem && chmod +x golem

USER golem

CMD ["./golem", "run"]