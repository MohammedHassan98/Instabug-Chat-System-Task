# Chat System API

The **Chat System API** is a scalable, RESTful service designed to manage applications, chats, and messages efficiently. Built with Ruby on Rails (with optional Go services for specific endpoints), the API ensures robust concurrency handling and optimized performance for high-volume requests.

## Features

1. **Applications Management**:

   - Create, update, and view applications identified by tokens.

2. **Chats Management**:

   - Handle numbered chats within applications, ensuring unique numbering per application.

3. **Messages Management**:

   - Manage sequentially numbered messages for each chat.

4. **Search Functionality**:

   - Search messages using partial matching with Elasticsearch.

5. **Concurrency Handling**:

   - Use queuing systems and Redis to prevent race conditions and support high-load scenarios.

6. **Optimized Data**:

   - Chats and messages counts are maintained in the database with minimal lag (<1 hour).

7. **Containerized Deployment**:
   - The entire stack runs seamlessly with `docker-compose up`.

## Technologies Used

- **Backend**: Go
- **Database**: MySQL with proper indexing for performance
- **Search**: Elasticsearch
- **Caching & Queueing**: Redis

## Additional Information

The project includes:

- Swagger Documatation you can Reach via http://localhost:8080/swagger/index.html After Building the app composer
- A detailed README file with setup and usage instructions.

## Prerequisites

Before you begin, ensure you have the following installed:

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Go](https://golang.org/doc/install)

## Getting Started

### 1. Clone the Repository

git clone https://github.com/MohammedHassan98/Instabug-Chat-System-Task.git
cd Instabug-Chat-System-Task/

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

| Method | Endpoint                                            | Description           |
| ------ | --------------------------------------------------- | --------------------- |
| POST   | `/applications`                                     | Create Application    |
| GET    | `/applications`                                     | Get All Applications  |
| PUT    | `/applications/{token}`                             | Update Application    |
| GET    | `/applications/{token}/chats`                       | Get Application Chats |
| POST   | `/chats/{token}`                                    | Create Chat           |
| POST   | `/chats/{chatNumber}/messages`                      | Create Message        |
| GET    | `/applications/{token}/chats/{chatNumber}/messages` | Get Messages          |
| GET    | `/chats/{chatNumber}/messages/search`               | Search Messages       |

### 5. Stopping the Application

To stop the application, press `CTRL + C` in the terminal where Docker Compose is running.

### 6. Running Migrations

Migrations are automatically run when the application starts. If you need to run them manually, you can do so by calling the migration function in the code.

### Thanks For Your Time
