FROM golang AS builder

COPY . /app
WORKDIR /app

ENV GOFLAGS="-buildvcs=false"
ENV CGO_ENABLED=0
RUN go build -o /usr/local/bin/paranal ./cmd/paranal

FROM alpine

COPY --chmod=755 --from=builder /usr/local/bin/paranal /usr/local/bin/

RUN apk add tzdata

ENV PARANAL_SCRIPTS_PATH=/etc/paranal/scripts
ENV PARANAL_SERVER_CERT_FILE=/etc/paranal/ssl/cert.pem
ENV PARANAL_SERVER_KEY_FILE=/etc/paranal/ssl/key.pem
ENV PARANAL_DB_PATH=/etc/paranal/paranal.db

VOLUME /etc/paranal
EXPOSE 8080
CMD ["paranal"]
