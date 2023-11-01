# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory to /app
WORKDIR /app

# Copy all the files from the current directory to the container's /app directory
COPY . /app

# Set the environment variable PORT to 3000
ENV PORT=3000

RUN ls
RUN pwd
RUN mkdir ./files
RUN mkdir ./files/output

# Build the Go application inside the container
RUN go build -o myapp

# Run the application when the container starts
CMD ["./myapp"]