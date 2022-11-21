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

ENV NEO4J_HOST=neo4j+s://03181944.databases.neo4j.io
ENV NEO4J_USER=neo4j
ENV NEO4J_PASSWORD=5V9YyRWxq3OqfFmEUGLL1BqPHARTRrSYVA19i3d_j6w
ENV API_URL=0.0.0.0:8080

# Copy app files
COPY . .
# Build app
RUN go build -o app

# Expose port
EXPOSE 8080
CMD ["./app"]