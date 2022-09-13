FROM golang:1.18-alpine
# Define build env
ENV GOOS linux
ENV CGO_ENABLED 0
# Add a work directory
WORKDIR /app
# Cache and install dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN apk update && apk add bash

# Copy app files
COPY . .
# Build app
RUN go build -o app

# Expose port
EXPOSE 8888
ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["./app"]