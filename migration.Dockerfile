FROM alpine:3.13

RUN apk update && \
    apk add bash curl && \
    rm -rf /var/cache/apk/*

# Скачиваем goose
ADD https://github.com/pressly/goose/releases/download/v3.14.0/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose

WORKDIR /app

COPY .env .
COPY migration.sh .
COPY migrations ./migrations/

RUN apk add dos2unix && dos2unix migration.sh

RUN chmod +x migration.sh

ENTRYPOINT ["bash", "migration.sh"] 
