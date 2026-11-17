FROM node:lts AS client

WORKDIR /app

COPY package.json package-lock.json .
RUN npm install

COPY . .
RUN npm run client:build

FROM golang AS build

WORKDIR /app

COPY . .
COPY --from=client /app/client/dist client/dist

ENV GOFLAGS="-buildvcs=false"
ENV CGO_ENABLED=0
RUN go build -o /usr/local/bin/paranal ./cmd/paranal

FROM alpine

COPY --chmod=755 --from=build /usr/local/bin/paranal /usr/local/bin/
RUN mkdir -m 0700 -p \
  /etc/paranal \
  /etc/paranal/config \
  /etc/paranal/scripts \
  /var/lib/paranal

WORKDIR /etc/paranal

ENV PARANAL_CONFIG_ROOT=/etc/paranal
ENV PARANAL_DATA_ROOT=/var/lib/paranal

VOLUME /var/lib/paranal
EXPOSE 8080
CMD ["paranal"]
