# Go JWT Authentication Mini-Project

A beginner-friendly Go web application that demonstrates building a secure REST API with user registration, login, and protected routes using JWT (JSON Web Tokens). It connects to a PostgreSQL database and features live-reloading capabilities using **Air**.

## 🏗️ Architecture Stack

- **Go (Golang)**: The core programming language.
- **Echo (`labstack/echo/v4`)**: A high-performance, minimalist Go web framework used for routing and handling HTTP requests.
- **GORM (`gorm.io/gorm`)**: A popular Object-Relational Mapping library for Go, used to interact with the PostgreSQL database.
- **PostgreSQL**: The relational database used to store user information securely.
- **JWT (`golang-jwt/jwt/v5`)**: Used for generating and parsing stateless authentication tokens.
- **Bcrypt (`golang.org/x/crypto/bcrypt`)**: Used for securely hashing user passwords before storing them.

### Application Flow
1. **Registration (`POST /register`)**: Users provide a name, email, and password. The application hashes the password using `bcrypt` and creates a new user record in the database.
2. **Login (`POST /login`)**: Users provide their email and password. The application verifies the credentials and returns a signed **JWT Token**.
3. **Protected Route (`GET /user/profile`)**: Requires the client to send the generated JWT token in the `Authorization` header (`Bearer <token>`). A custom Auth Middleware validates the token before allowing access to the profile data.

## 🚀 Getting Started

### Prerequisites
- Go installed on your system.
- PostgreSQL database running locally with the following credentials (as configured in `main.go`):
  - User: `postgres`
  - Password: `1234`
  - Database Name: `testdb`
  - Port: `5432`

### Commands to Run

**1. Clone the repository**
```bash
git clone https://github.com/abhishek15032000/gorm-middleware-token-using.git
cd gorm-middleware-token-using
```

**2. Download dependencies**
```bash
go mod tidy
```

**3. Run the application (Standard)**
```bash
go run main.go
```
The server will start on `http://localhost:8091`.

## 🌬️ Live Re-loading with Air

This project uses **Air** to automatically rebuild and restart the Go application whenever a code change is saved. 

### What is Air?
[Air](https://github.com/air-verse/air) is a command-line utility. Instead of manually stopping the server and restarting it with `go run main.go` after every change, `air` watches your directories and does it automatically, providing a smoother developer experience.

### Configuration (`.air.toml`)
The behavior of Air is controlled by the `.air.toml` file generated at the project root. Key configurations include:
- **`include_ext`**: Which file extensions trigger a reload (e.g., `.go`, `.html`).
- **`exclude_dir`**: Directories Air shouldn't monitor (e.g., `tmp` which holds the compiled binary).
- **`cmd`**: The command used to compile your application.
- **`bin`**: The built executable that Air will run.

### Running with Air
If you haven't installed Air yet:
```bash
go install github.com/air-verse/air@latest
```

To run your project with live-reloading enabled:
```bash
go run github.com/air-verse/air@latest
```
*(If your `GOPATH/bin` is accessible in your system path, you can simply type `air` in your terminal)*
