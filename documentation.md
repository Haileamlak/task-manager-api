### **Task Management API: Refactored Architecture Documentation**

---

#### **Overview**

This documentation details the refactored architecture of the Task Management API, adhering to Clean Architecture principles. The API's main purpose is to manage tasks while providing authentication and authorization via JWT. The design ensures separation of concerns, scalability, maintainability, and testability. This guide also includes design decisions made during refactoring and guidelines for future development.

---

### **1. Architecture Overview**

The Task Management API is structured into four primary layers:

1. **Domain Layer**: The core of the application, containing business logic and entities.
2. **Use Cases Layer**: Defines application-specific business rules. This layer orchestrates the flow of data between the entities and the outer layers.
3. **Infrastructure Layer**: Manages all external dependencies, such as databases, frameworks, and third-party services.
4. **Delivery Layer**: Responsible for the communication with the outside world, handling HTTP requests, and delivering responses.

---

### **2. Layers and Responsibilities**

#### **2.1 Domain Layer**

- **Purpose**: This is the core of the application, containing business entities and domain logic.
- **Components**:
  - **Entities**: Represent the fundamental business objects, such as `Task` and `User`.
  - **Value Objects**: Objects that represent domain concepts like status, which do not have an identity (e.g., task status).
  
- **Design Decisions**:
  - Business rules and logic reside within this layer, making it independent of external factors such as databases or UI.

- **Example**:
  
  ```go
  type Task struct {
      ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
      Title       string             `json:"title" bson:"title"`
      Description string             `json:"description" bson:"description"`
      DueDate     time.Time          `json:"due_date" bson:"due_date"`
      Status      string             `json:"status" bson:"status"`
  }
  ```

#### **2.2 Use Cases Layer**

- **Purpose**: Contains the application-specific business logic. It orchestrates the execution of various actions within the system.
- **Components**:
  - **Interactors**: Coordinate data flow between the Domain and Infrastructure layers.
  
- **Design Decisions**:
  - This layer only contains logic that pertains to the specific application use cases. It does not depend on the Infrastructure or Delivery layers.
  
- **Example**:
  
  ```go
  type UserUsecase interface {
      Register(username, password string) (*Domain.User, error)
      Login(username, password string) (string, error)
      PromoteUser(userID string) error
  }
  ```

#### **2.3 Infrastructure Layer**

- **Purpose**: Deals with external systems and frameworks, including data storage, third-party services, and libraries.
- **Components**:
  - **Repositories**: Handle data persistence and retrieval.
  - **Services**: Implement infrastructure-related concerns such as JWT token generation, password hashing, etc.
  
- **Design Decisions**:
  - The implementation of repositories and services are isolated from the core business logic to allow easier swapping of external dependencies.
  
- **Example**:
  
  ```go
  type JWTService interface {
      GenerateToken(userID string, username string, role string) (string, error)
      ValidateToken(token string) (*jwt.Token, error)
  }
  ```

#### **2.4 Delivery Layer**

- **Purpose**: Manages interaction with the outside world, such as handling HTTP requests and sending responses.
- **Components**:
  - **Controllers**: Handle requests, validate input, and call the appropriate use cases.
  - **Middleware**: Manages cross-cutting concerns such as authentication and logging.
  
- **Design Decisions**:
  - Controllers do not contain business logic. They act as intermediaries between user inputs and the use cases.
  
- **Example**:
  
  ```go
  func (ctrl *UserController) Login(c *gin.Context) {
      var request struct {
          Username string `json:"username"`
          Password string `json:"password"`
      }
  
      if err := c.ShouldBindJSON(&request); err != nil {
          c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
          return
      }
  
      token, err := ctrl.userUsecase.Login(request.Username, request.Password)
      if err != nil {
          c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
          return
      }
  
      c.JSON(http.StatusOK, gin.H{"token": token})
  }
  ```

---

### **3. Design Decisions**

#### **3.1 Adopting Clean Architecture**

- **Why**: Clean Architecture was chosen to ensure the separation of concerns, allowing each layer to focus on its specific responsibilities. This makes the system easier to maintain, extend, and test.
  
#### **3.2 JWT for Authentication**

- **Why**: JWT was selected for stateless, scalable, and secure user authentication. This approach decouples user sessions from the server, allowing for horizontal scaling.
  
#### **3.3 Password Hashing with bcrypt**

- **Why**: bcrypt provides a strong level of security for password storage, making it resistant to brute-force attacks.
  
#### **3.4 MongoDB for Data Persistence**

- **Why**: MongoDB was chosen for its flexibility and scalability, particularly suited for handling the document-based structure of tasks and users.
  
#### **3.5 Dependency Injection**

- **Why**: Dependency Injection allows for easy swapping of components (e.g., changing the database) without affecting other parts of the system. It also simplifies unit testing by allowing mock implementations.

---

### **4. Guidelines for Future Development**

#### **4.1 Adding New Features**

1. **Define the Feature**: Start by clearly defining the new feature or enhancement in terms of business logic.
2. **Domain Layer**:
   - Create or update domain entities and value objects as needed.
   - Ensure the domain layer remains independent of external dependencies.
3. **Use Cases Layer**:
   - Implement or update interactors to handle the new feature's logic.
   - Ensure the use cases layer coordinates the flow between the domain and infrastructure layers without directly depending on either.
4. **Infrastructure Layer**:
   - Update repositories or services to support the new feature.
   - Ensure changes are isolated from business logic.
5. **Delivery Layer**:
   - Create or update controllers to expose the new feature via the API.
   - Update routes, middleware, and input validation as necessary.

#### **4.2 Testing**

1. **Unit Tests**:
   - Focus on testing individual components (e.g., use cases, services).
   - Use mock implementations of dependencies (e.g., mock the repository in use case tests).
2. **Integration Tests**:
   - Test interactions between layers (e.g., ensure that the Delivery layer correctly calls the Use Cases layer).
   - Use in-memory databases or isolated environments to test data-related logic.
3. **End-to-End Tests**:
   - Test the full flow of the application, from HTTP request to database interactions.
   - Validate that the API behaves as expected for different user roles.

#### **4.3 Security Considerations**

1. **Environment Variables**:
   - Store sensitive information such as JWT secret keys in environment variables.
   - Ensure these variables are not exposed in version control.
2. **HTTPS**:
   - Enforce HTTPS in production environments to protect data in transit.
3. **Data Validation**:
   - Perform input validation at the Delivery layer to prevent SQL injection, XSS, and other attacks.
4. **Role-Based Access Control**:
   - Use the JWT token to enforce role-based access control across the API.

#### **4.4 Documentation**

- Maintain up-to-date API documentation for all endpoints, including request/response examples, error codes, and authentication details.
- Document significant architectural decisions and trade-offs to provide context for future developers.
- Provide onboarding documentation for new developers to understand the system architecture, development guidelines, and testing practices.

---

### **5. Conclusion**

The refactored Task Management API follows Clean Architecture principles to ensure that the system is modular, maintainable, and scalable. By isolating business logic from external dependencies, the architecture allows for easy enhancements and robust testing. This documentation serves as a comprehensive guide for current and future development, ensuring consistency and adherence to best practices.