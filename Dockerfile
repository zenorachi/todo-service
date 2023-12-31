FROM alpine:latest
LABEL authors="zenorachi"

WORKDIR /root

# Copy binary
COPY ./.bin/app /root/.bin/app

# Copy configs
COPY ./.env /root/
COPY ./configs/main.yml /root/configs/main.yml

# Install psql-client
RUN apk add --no-cache postgresql-client

CMD ["sh", "-c", "sh ./scripts/db-connection/wait-db.sh && ./.bin/app"]