# Use the official Golang image as a base image
FROM golang:1.21.3 as build

# Set the working directory inside the container
WORKDIR /app

# Copy your Go files and download dependencies
COPY . /app

# Build your Go application
RUN go build -o app main.go

# Use a minimal Alpine Linux image as the final base image
FROM alpine:latest

# Set an environment variable for the port (you can change the value as needed)
ENV PORT=3000

# Copy the compiled Go application from the previous stage
COPY --from=build /app /app

# Expose the port that your Go application will listen on
EXPOSE $PORT

# Start your Go application
CMD ["/app"]
