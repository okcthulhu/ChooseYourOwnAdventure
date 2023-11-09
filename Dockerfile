# Use the official Go image as a parent image
FROM golang:1.21.3

# Set the working directory in the image to /app
WORKDIR /app

# Copy the go mod and sum files, download dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the local package files to the container’s workspace
COPY . .

# Build the Go app
RUN go build -o main .

# Expose the port the app runs on.
EXPOSE 8080

# Run the Go app when the container launches
CMD ["./main"]
