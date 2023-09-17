# Use the official Golang image as a base image
FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Copy everything from the current directory to the working directory inside the container
COPY . .

# Run go get to download dependencies
RUN go get -d -v ./...

# This command will run when the container starts
CMD ["go", "run", "main.go", "serve", "--http", "0.0.0.0:8090"]
