FROM golang:1.23.0-alpine3.20 as builder
ENV CGO_ENABLED=0

# WORKDIR /build
# COPY go.mod ./
# RUN go mod download && go mod tidy
# COPY . .
# RUN go build main.go
COPY . /build
# Build the service binary.
WORKDIR /build
RUN go build -o gopypi -ldflags "-X main.build=main" ./cmd/web
#Multistage build 

FROM alpine:3.20.0
# Add maintainer info
LABEL maintainer="Danail Surudzhiyski"
#add user and group in the alpine image
RUN addgroup -g 1000 -S gopypi && \
  adduser -u 1000 -h /cmd -G gopypi -S gopypi
#Copy from stage 0 builder only the binary files
COPY --from=builder --chown=gopypi:gopypi /build/gopypi .
COPY --from=builder /build/templates ./templates
# Set environment variables (optional)
WORKDIR /
USER gopypi
CMD "./gopypi" "--username=$USERNAME" "--password=$PASSWORD"
