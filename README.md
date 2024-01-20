# Webhook Service

Webhook Service is a simple HTTP server written in Go that collects and sends batches of JSON payloads. It is designed for receiving and processing webhook data.

## Features

- Receive JSON payloads through HTTP POST requests.
- Collect payloads in memory and send them in batches at regular intervals.
- Retry logic for failed batch deliveries.

## Getting Started

### Prerequisites

- Go (Golang) installed on your machine
- Docker (optional)

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/your-username/webhook-service.git
    cd webhook-service
    ```

2. Build the application:

    ```bash
    go build -o webhook-service
    ```

### Configuration

Configure the application using the `config.env` file. Ensure the necessary environment variables are set.

Example `config.env`:

```env
BATCH_SIZE=10
BATCH_INTERVAL=60
POST_ENDPOINT=http://localhost:8080/log_data
# Add other configuration variables as needed
```

## Usage

Run the application:

```bash
./webhook-service
```
By default, the application runs on port 8080.


## Docker

You can also use Docker to build and run the application:

```bash
docker build -t webhook-service .
docker run -p 8080:8080 webhook-service
```

## API Endpoints

### Health Check

- **Endpoint:** `/healthz`
- **Method:** `GET`
- **Description:** Check the health status of the service.

### Log Payload

- **Endpoint:** `/log`
- **Method:** `POST`
- **Description:** Receive JSON payloads. Payloads are collected and sent in batches.

### Send Data (Example)

- **Endpoint:** `/send-data`
- **Method:** `POST`
- **Description:** Simulate sending data to the configured endpoint, this will receive the 
 batch of data on condition bing met. 

## Contributing

Contributions are welcome! Please follow the [Contribution Guidelines](CONTRIBUTING.md).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
