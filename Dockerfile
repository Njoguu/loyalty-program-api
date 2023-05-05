# Use an official Go runtime as a parent image
FROM golang:1.20.4-alpine3.17

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /go/src/app
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port 8080 for the API
EXPOSE 8080

# Set the command to run when the container starts
CMD ["./main"]
