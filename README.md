# Task Manager API

Task Manager API is a robust and scalable RESTful API designed for managing tasks. It supports user authentication and authorization, task management, and role-based access control. The API is built using Go and follows Clean Architecture principles for maintainability, testability, and scalability.

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
- [Project Structure](#project-structure)
- [Environment Variables](#environment-variables)
- [Running the Application](#running-the-application)
- [Running Tests](#running-tests)
- [API Endpoints](#api-endpoints)
- [License](#license)

## Features

- **User Authentication**: JWT-based authentication with role management.
- **Task Management**: Create, update, delete, and manage tasks.
- **Role-Based Access Control**: Fine-grained control over user permissions based on roles.
- **Secure Password Handling**: Secure password storage and validation using bcrypt.
- **Database Integration**: MongoDB as the database for storing tasks and user information.
- **Unit Testing**: Comprehensive test suite for ensuring code quality.

## Tech Stack

- **Language**: Go (Golang)
- **Framework**: Gin for HTTP routing and middleware
- **Database**: MongoDB
- **Authentication**: JWT (JSON Web Tokens)
- **Environment Management**: godotenv
- **Testing**: Testify, httptest

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- [Go](https://golang.org/doc/install) (1.18 or higher)
- [MongoDB](https://www.mongodb.com/try/download/community) (running locally or via a cloud service)
- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

### Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/your-username/task-manager-api.git
   cd task-manager-api
   ```

2. **Install dependencies:**

   ```bash
   go mod tidy
   ```

3. **Set up the environment variables:**

   Create a `.env` file in the root directory with the following contents:

   ```env
   PORT=8080
   MONGO_URI=mongodb://localhost:27017
   JWT_SECRET=your_jwt_secret
   ```

## Project Structure

```bash
task-manager-api/
├── Delivery/           # HTTP handlers and request/response structures and main.go file the entry point
├── Domain/             # Core business logic and entities
├── Infrastructure/     # External services (database, JWT, password hashing)
├── Repositories/       # Data persistence and retrieval logic
├── Usecases/           # Application-specific business rules
├── documentation.md    # The testing strategy for the API
├── .env                # Environment variables file
├── go.mod              # Go module dependencies
├── README.md           # Project documentation
```

## Environment Variables

The following environment variables are used in the project:

- `PORT`: The port on which the server will run.
- `MONGO_URI`: The URI for connecting to your MongoDB instance.
- `JWT_SECRET`: Secret key used for signing JWT tokens.

## Running the Application

To start the server, use the following command:

```bash
go run main.go
```

The server will start on `http://localhost:8080` by default.

## Running Tests

The project includes a comprehensive suite of unit tests to ensure code quality. The tests cover various components of the application, including the delivery, infrastructure, domain, and use case layers.

To run the tests, use the following command:

```bash
go test ./...
```

### Testing Considerations

- **Test Coverage**: The test suite is designed to provide coverage for critical components, ensuring the robustness of the API.
- **Issues Encountered**: Ensure that your MongoDB service is running locally, or mock the database interactions if necessary. The `repository` tests are excluded from the CI workflow due to dependency on a local MongoDB instance.

## API Endpoints

The following are the main API endpoints:

- **User Authentication**
  - `POST /auth/login`: User login
  - `POST /auth/register`: User registration

- **Task Management**
    ***All Users***
      - `GET /tasks`: Retrieve all tasks
      - `GET /task/:id` Retrieve a task by ID

    ***Admins only***
      - `POST /tasks`: Create a new task
      - `PUT /tasks/:id`: Update an existing task
      - `DELETE /tasks/:id`: Delete a task

For detailed API documentation, refer to the [API Documentation](#).
