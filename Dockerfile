# ========================
#       BUILD STAGE
# ========================
# Use an official Go runtime as a parent image
FROM golang:1.20.4-alpine3.17 AS builder

# Set the working directory 
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Build the Go api
RUN go build -o main main.go

# Download migrate binary
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz

# ========================
#       RUN STAGE
# ========================
FROM alpine:3.17

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY .env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration

# Expose port 8080 for the API
EXPOSE 8000

# Set the command to run when the container starts
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]
