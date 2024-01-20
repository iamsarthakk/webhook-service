# Use the official Go image as the base image
FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Copy the necessary files into the container
COPY . .

# Build the Go application
RUN go build -o webhook-receiver

# Expose the port that the application will run on
EXPOSE 8080

# Set environment variables (adjust as needed)
ENV BATCH_SIZE=10
ENV BATCH_INTERVAL=60
ENV POST_ENDPOINT=http://localhost:8080/log

# Command to run the application
CMD ["./webhook-receiver"]
