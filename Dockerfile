FROM golang:1.23-alpine as builder

# Add Maintainer Info
LABEL maintainer="Sam Zhou <sam@mixmedia.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go version \
 && export GO111MODULE=on \
 && export GOPROXY=https://goproxy.io,direct \
 && go mod vendor \
 && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o motors-backup


######## Start a new stage from scratch #######
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata dumb-init gettext-envsubst subversion \
    && update-ca-certificates \
    && envsubst --version \
    && rm -rf /var/cache/apk/*

WORKDIR /app


# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/motors-backup .

ENV DB_PORT=3306 \
 DB_HOST=localhost \
 DB_USER= \
 DB_PASSWORD= \
 DB_NAME=

ENTRYPOINT ["dumb-init", "--"]

CMD /app/motors-backup

