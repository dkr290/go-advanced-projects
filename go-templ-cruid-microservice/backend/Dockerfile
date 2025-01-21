FROM golang:1.23.0-alpine3.20 AS builder
ENV CGO_ENABLED=0
COPY . /build
WORKDIR /build
RUN go build -o backend -ldflags "-X main.build=main" ./


#Multistage build 
FROM alpine:3.20.0

#Add maintainer info

LABEL maintainer="Danail Surudzhiyski"
#add user and group in the alpine image

RUN addgroup -g 1000 -S backend && \
  adduser -u 1000 -h /cmd -G backend -S backend


COPY --from=builder --chown=backend:backend /build/backend .

WORKDIR /
USER backend
CMD "./backend"
