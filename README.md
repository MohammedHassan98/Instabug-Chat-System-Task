# Chat System

This project is a chat system built using Go, Redis, MySQL, and Elasticsearch. It provides a RESTful API for managing applications, chats, and messages.

## Prerequisites

Before you begin, ensure you have the following installed:

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Go](https://golang.org/doc/install)

## Getting Started

### 1. Clone the Repository

git clone https://github.com/yourusername/chat-system.git
cd chat-system

### 2. Set Up Environment Variables

Create a `.env` file in the root directory of the project and add the following environment variables:

```
DB_USER=root
DB_PASSWORD=root
DB_HOST=db
DB_PORT=3306
DB_NAME=chat_system
REDIS_HOST=redis
ELASTICSEARCH_URL=http://elasticsearch:9200
```

### 3. Build and Run the Application

You can use Docker Compose to build and run the application along with its dependencies (MySQL, Redis, and Elasticsearch).

```
docker-compose up --build
```

### 4. Access the Application

Once the application is running, you can access the API at `http://localhost:8080`.

### 5. API Endpoints

Here are some of the available API endpoints:

- **Create Application**: `POST /applications`
- **Get All Applications**: `GET /applications`
- **Create Chat**: `POST /applications/{token}/chats`
- **Create Message**: `POST /chats/{chatNumber}/messages`
- **Get Messages**: `GET /applications/{token}/chats/{chatNumber}/messages`

### 6. Stopping the Application

To stop the application, press `CTRL + C` in the terminal where Docker Compose is running.

### 7. Running Migrations

Migrations are automatically run when the application starts. If you need to run them manually, you can do so by calling the migration function in the code.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

# Chat System

This project is a chat system built using Go, Redis, MySQL, and Elasticsearch. It provides a RESTful API for managing applications, chats, and messages.

## Prerequisites

Before you begin, ensure you have the following installed:

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Go](https://golang.org/doc/install)

## Getting Started

### 1. Clone the Repository

git clone https://github.com/yourusername/chat-system.git
cd chat-system

### 2. Set Up Environment Variables

Create a `.env` file in the root directory of the project and add the following environment variables:

```
DB_USER=root
DB_PASSWORD=root
DB_HOST=db
DB_PORT=3306
DB_NAME=chat_system
REDIS_HOST=redis
ELASTICSEARCH_URL=http://elasticsearch:9200
```

### 3. Build and Run the Application

You can use Docker Compose to build and run the application along with its dependencies (MySQL, Redis, and Elasticsearch).

```
docker-compose up --build
```

### 4. Access the Application

Once the application is running, you can access the API at `http://localhost:8080`.

### 5. API Endpoints

Here are some of the available API endpoints:

- **Create Application**: `POST /applications`
- **Get All Applications**: `GET /applications`
- **Create Chat**: `POST /applications/{token}/chats`
- **Create Message**: `POST /chats/{chatNumber}/messages`
- **Get Messages**: `GET /applications/{token}/chats/{chatNumber}/messages`

### 6. Stopping the Application

To stop the application, press `CTRL + C` in the terminal where Docker Compose is running.

### 7. Running Migrations

Migrations are automatically run when the application starts. If you need to run them manually, you can do so by calling the migration function in the code.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
