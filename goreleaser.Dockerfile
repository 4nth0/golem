FROM alpine:3.12.0
RUN apk --update add --no-cache ca-certificates tzdata
WORKDIR /root/

COPY golem .

CMD [ "./golem", "run" ]
