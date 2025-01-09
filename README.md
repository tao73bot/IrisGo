# This is Simple CRM using IRIS framework of GoLang

## I follow this folder structure

```paintext
project/
│
├── controllers/             # Logic for handling incoming requests and responses
├── db/                      # Database-related files and configurations
├── docs/                    # Documentation for the project
├── middlewares/             # Middleware logic (e.g., request validation, authentication)
├── models/                  # Data models and structs
├── routes/                  # Route definitions for the API
├── utils/                   # Utility functions and helpers
│
├── main.go                  # Go application entry point
├── go.mod                   # Go module file for dependencies
├── go.sum                   # Checksums for Go module dependencies
├── .gitignore               # Git ignore file to exclude unwanted files/folders
├── .env                     # Environment variables (should not be in version control)
└── README.md                # Project documentation
```

## Table of Contents

1. [Installation](#installation)  
2. [Running the Project](#running-the-project)  
3. [API Endpoints](#api-endpoints)  
   - [User Endpoints](#user-endpoints)  
   - [Task Endpoints](#task-endpoints)  
4. [Authentication](#authentication)  
5. [Environment Variables](#environment-variables)  
6. [Database](#database)  
7. [Testing](#testing)  
8. [Contributing](#contributing)  
9. [License](#license)

## Installation

### Prerequisites

Ensure the following are installed on your system:

- Go 1.12 or higher
- Postgresql for database
- Postman or similar Api testing tools

### Steps

1. **Clone the repository**  
   Clone the project repository to your local machine:
   ```bash
   git clone https://github.com/your-repo-name.git
   cd your-repo-name
   ```

2. **Install dependencies**  
    Install the required Go packages:
    ```bash
    go mod tidy
    ```

## Running the Project

Once the installation is complete, follow these steps to run the Iris application:

1. **Start the server**  
   Use `run` command to start the IRIS server:
   ```bash
   go run main.go
   ```

- To get live reload facilities install `air` on your device. To install go there `https://github.com/air-verse/air`
    ```bash
    air
    ```


2. **Access the application**  
    By default,GO project access point is:
    - **Base URL** `localhost:8080`

    But You can change your port based on your demand.

3. **Explore Api documentations**  
    Setup swagger for interactive API documentation at these endpoints:
    - **Swagger UI** : http://127.0.0.1:8080/swagger/index.html
    
    To setup Swagger visit `https://github.com/iris-contrib/swagger`


## API Endpoints
### User Endpoints

`POST /auth/signup`
- **Description**: Create a new user.
- **Parameters**: None
- **Request body**:
```json
{
    "email": "user@example.com",
    "password": "string",
    "name": "string"
}
```
- **Response**:
```json
{
    "user": user 
}
```

`POST /auth/signin`
- **Description**: Login as user.
- **Parameters**: None
- **Request body**:
```json
{
    "email": "user@example.com",
    "password": "string"
}
```
- **Response**:
```json
{
    "token": "string",
    "refreshToken": "string"
}
```

### Task Endpoints




## Authentication

- **Mechanism**: Bearer Token (JWT)
- **Endpoints**:
    - `POST /auth/login` to get token
    - Use the token in the `Authorization` header:
    ```
    Authorization: Bearer <your_token>
    ```

## Environment Variables
Create a `.env` file in the root directory with the following keys:
```.env
PORT = myport

DB_HOST = localhost
DB_USER = postgres
DB_PASSWORD = mypass
DB_NAME = mydb
DB_PORT = postgrsPORT

JWT_SECRET = myJWTSecrect

REDIS_URL = localhost://6379

EmailSender = your_mail@gmail.com
EmailPassword = your_app_password
EmailHost = smtp.gmail.com
EmailPort = 587

```

- **Config the settings using `.env` file**

## Database
- **Database Type**: `Postgres`
- **ORM**: I user`gorm` for ORM support. To setup gorm visit `https://gorm.io/docs/`
- **Setup**: Use database `URI` from `.env` file