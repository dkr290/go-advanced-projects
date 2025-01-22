FROM golang:1.23.0-alpine3.20 AS builder
ENV CGO_ENABLED=0
COPY . /build
WORKDIR /build
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go get -u github.com/a-h/templ
RUN /go/bin/templ generate view --path=/build/view  
RUN go build  -o bin/frontend main.go


#Multistage build 
FROM alpine:3.20.0

#Add maintainer info

LABEL maintainer="Danail Surudzhiyski"
#add user and group in the alpine image

RUN addgroup -g 1000 -S frontend && \
  adduser -u 1000 -h /cmd -G frontend -S frontend


COPY --from=builder --chown=frontend:frontend /build/bin/frontend .

WORKDIR /
USER frontend
CMD "./frontend"
