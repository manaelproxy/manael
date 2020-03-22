FROM alpine:3.11

RUN apk --update-cache add ca-certificates

COPY manael /usr/local/bin/

CMD ["manael"]
