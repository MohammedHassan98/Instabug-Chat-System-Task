# Chat System

This project is a chat system built using Go, Redis, MySQL, and Elasticsearch. It provides a RESTful API for managing applications, chats, and messages.

## Prerequisites

Before you begin, ensure you have the following installed:

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Go](https://golang.org/doc/install)

## Getting Started

### 1. Clone the Repository

git clone https://github.com/mohammedhassan98/chat-system.git
cd chat-system

### 2. Build and Run the Application

You can use Docker Compose to build and run the application along with its dependencies (MySQL, Redis, and Elasticsearch).

```
docker-compose up --build
```

Or you can use docker compose plugin using this command

```
docker compose up --build
```

### 3. Access the Application

Once the application is running, you can access the API at `http://localhost:8080`.

### 4. API Endpoints

Here are some of the available API endpoints:

| Method | Endpoint                                            | Description          |
| ------ | --------------------------------------------------- | -------------------- |
| POST   | `/applications`                                     | Create Application   |
| GET    | `/applications`                                     | Get All Applications |
| POST   | `/applications/{token}/chats`                       | Create Chat          |
| POST   | `/chats/{chatNumber}/messages`                      | Create Message       |
| GET    | `/applications/{token}/chats/{chatNumber}/messages` | Get Messages         |
| GET    | `/chats/{chatNumber}/messages/search`               | Search Messages      |

### 5. Stopping the Application

To stop the application, press `CTRL + C` in the terminal where Docker Compose is running.

### 6. Running Migrations

Migrations are automatically run when the application starts. If you need to run them manually, you can do so by calling the migration function in the code.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
