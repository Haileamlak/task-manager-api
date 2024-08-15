# Task Manager API - Test Suite Documentation

## Overview

This document outlines the testing strategy for the Task Manager API, detailing how to run tests locally, check test coverage, and address issues encountered during testing. The test suite covers various layers of the application, including repositories, usecases, controllers, infrastructure, domain, and delivery.

## Test Suite Components

The test suite for the Task Manager API is divided into the following layers:

1. **Domain Tests**:

   - Tests the core domain entities and value objects, ensuring they are correctly implemented and validated.

   To run the domain tests, execute the following command:

   ```bash
   go test ./Domain -v
   ```

2. **Repository Tests**:
   - Tests the database repository layer, including CRUD operations and query functions.

     To run the repository tests, execute the following command:

     ```bash
     go test ./Repositories -v
     ```

     **Note**: Repository tests require a running MongoDB instance. Ensure that the MongoDB connection string is correctly set in the environment variables before running the tests.

3. **Usecase Tests**:

   Usecase tests validate the core business logic of the Task Manager API. These tests cover various functionalities such as task creation, updating, deletion, and retrieval. Additionally, they also test user registration and authentication processes.

   To run the usecase tests, execute the following command:

   ```bash
   go test ./Usecases -v
   ```

4. **Infrastructure Tests**:

   - Tests the external dependencies and services used by the application, such as JWT token generation and validation. These tests ensure that the infrastructure components are working as expected.

   To run the infrastructure tests, execute the following command:

   ```bash
   go test ./Infrastructure -v
   ```

5. **Delivery Tests**:
   - Tests the HTTP handlers and controllers responsible for processing requests and returning responses. These tests validate the API endpoints and their functionality.

    To run the delivery tests, execute the following command:
    
    ```bash
    go test ./Delivery -v
    ```

It is recommended to regularly run these tests to maintain the integrity and reliability of the Task Manager API.

## Prerequisites

Before running the tests, ensure you have the following installed:

- Go (version 1.19+)
- MongoDB (if running repository tests locally)
- Git

## Running Tests Locally

### 1. Clone the Repository

First, clone the repository and navigate to the project directory:

```bash
git clone github.com/haileamlak/task-manager-api.git
cd task-manager-api
```

### 2. Set Up Environment Variables

Ensure all necessary environment variables are set up before running the tests. For example, set the `GIN_MODE` environment variable:

```bash
export GIN_MODE=release
```

### 3. Run All Tests

You can run all tests across all layers using the following command:

```bash
go test ./... -v
```

The `-v` flag enables verbose mode, which provides detailed output for each test.

### 4. Run Specific Tests

To run tests for a specific package or test function, use the `-run` flag:

```bash
# Run all tests in the Infrastructure package
go test ./Infrastructure -v

# Run a specific test function
go test -run ^TestJWTService$ ./Infrastructure -v
```

### 5. Exclude Specific Tests (e.g., Repository Tests)

If you need to exclude tests that depend on external resources (like MongoDB), you can use build tags or manually skip them by commenting them out.

### 6. Running Tests with Coverage

To generate a test coverage report, use the following commands:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

Open `coverage.html` in your web browser to view the coverage report.

## Test Coverage Metrics

Test coverage is an important metric for assessing how much of your code is covered by tests. To check the overall test coverage:

1. **Run Tests with Coverage:**

   ```bash
   go test ./... -cover
   ```

2. **Generate Detailed Coverage Report:**

   ```bash
   go test ./... -coverprofile=coverage.out
   go tool cover -html=coverage.out
   ```

   This command will generate an HTML file (`coverage.html`) that you can open in your web browser for detailed coverage information.

## Issues Encountered During Testing

### 1. Interface Conversion Panic

- **Issue**: A panic occurred due to an invalid interface conversion when testing JWT token validation.
- **Solution**: Ensure the mock function returns a correctly typed `*jwt.Token` or `nil` as required.

### 2. Gin Logger and Recovery Middleware Warnings

- **Issue**: Warnings were displayed about creating an Engine instance with Logger and Recovery middleware during testing.
- **Solution**: The warnings are harmless, but to avoid them, use `gin.New()` instead of `gin.Default()` in test setups.

### 3. Local MongoDB Dependency

- **Issue**: Repository tests depend on a local MongoDB instance, causing issues in the CI pipeline.
- **Solution**: Exclude these tests from the CI workflow and provide instructions for running them locally.

## Continuous Integration (CI) Setup

The CI pipeline is configured to automatically run unit tests whenever new code is pushed to the repository. This ensures that all commits maintain the project's code quality standards.

### Modifying the CI Pipeline

To remove specific tests (e.g., repository tests) from the CI workflow:

1. Open the CI configuration file (e.g., `.github/workflows/ci.yml`).
2. Exclude the tests by adding a flag or condition that skips those tests during the CI run.

Example:

```yaml
- name: Run Tests
  run: go test ./... -v | grep -v "repository_test.go"
```

## Conclusion

The test suite for the Task Manager API is comprehensive and covers all major aspects of the application. By following the instructions provided in this document, developers can ensure that their code is well-tested and maintainable. Regularly running tests and monitoring test coverage will help maintain a high level of code quality throughout the development lifecycle.

---