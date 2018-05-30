FROM alpine:3.6

RUN apk --no-cache add ca-certificates

COPY build/bin/backup-etcd /app/backup-etcd

ENTRYPOINT [ "/app/backup-etcd" ]