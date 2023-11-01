# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory to /app
WORKDIR /app

# Copy all the files from the current directory to the container's /app directory
COPY . /app

# Set the environment variable PORT to 3000
ENV PORT=3000

RUN mkdir ./files

RUN mkdir ./files/output

# Make sure the ffmpeg binary is executable
RUN chmod +x ./ffmpeg


# Build the Go application inside the container
RUN go build -o myapp

RUN ls

# Set the entry point for the container
CMD ["./myapp"]
