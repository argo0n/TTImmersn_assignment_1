# TikTok Tech Immersion Assignment - Instant Messaging System

This project consists of an instant messageing system implented using a gRPC and HTTP server, with a MySQL database for persisting messages. The services are orchestrated using Docker Compose.

## Getting Started

### Prerequisites

Ensure you have the following installed:

- Docker
- Go

### Build & Run
To start all services (RPC server, HTTP server, etcd and MySQL database), simply run this in the base directory:

```shell
docker-compose build && docker-compose up -d
```

## API Documentation

### 1. GET /api/pull

Pull messages from a specified chat.

**Request Parameters**

| Parameter | Type    | Description                                                                                      | Required |
|-----------|---------|--------------------------------------------------------------------------------------------------|----------|
| chat      | string  | The chat identifier in the format "sender:receiver"                                              | Yes      |
| cursor    | integer | The ID of the message to start pulling from                                                      | No       |
| limit     | integer | The maximum number of messages to return                                                         | No       |
| reverse   | boolean | The order of the messages, in reverse or not. If it's not specified, it will default to `false`. | No       |

**Example Request**

```shell
curl -X GET -H "Content-Type: application/json" -d '{"chat":"a:b", "cursor":15, "limit":10, "reverse": false}' 'http://localhost:8080/api/pull'
```

**Example Response**

```json
{
  "messages": [
    {
      "chat": "a:b",
      "text": "Message16",
      "sender": "1",
      "send_time": 1686584832
    },
    {
      "chat": "a:b",
      "text": "Message17",
      "sender": "1",
      "send_time": 1686584832
    },
    {
      "chat": "a:b",
      "text": "Message18",
      "sender": "1",
      "send_time": 1686584832
    }
  ],
  "has_more": true,
  "next_cursor": 18
}
```

### 2. POST /api/send

Send a message in a specified chat.

**Request Parameters**

| Parameter | Type   | Description                                         | Required |
|-----------|--------|-----------------------------------------------------|----------|
| chat      | string | The chat identifier in the format "sender:receiver" | Yes      |
| text      | string | The message text                                    | Yes      |
| sender    | string | The ID of the sender                                | Yes      |

**Example Request**

```shell
curl -X POST -H "Content-Type: application/json" -d '{"chat":"a:b", "text":"Hello World", "sender":"1"}' 'http://localhost:8080/api/send'
```

**Example Response**

```http
HTTP/1.1 200 OK
```

No body or content is returned for this endpoint.