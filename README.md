# Miti Arts Backend

Miti Arts is an online store that sells high-quality timber-based furniture. This repository contains the backend logic for managing furniture products and user accounts.

## Technologies Used:
- **Gin Framework**: A web framework for Go (Golang) to build high-performance REST APIs.
- **GORM**: An ORM library for Golang to interact with the PostgreSQL database efficiently.
- **PostgreSQL**: Relational database for storing furniture and user data.
- **JWT (JSON Web Token)**: For secure user authentication.
- **Go Modules**: Dependency management for Go.

## Features:
- **CRUD for Furniture**: Add, update, delete, and fetch details about furniture items.
- **User Authentication**: Secure user registration, login, and JWT-based authentication.
- **Search and Filtering**: Filter and search furniture based on price, material, and name.
- **Error Handling**: Clean and user-friendly error responses for better UX.
- **Logging**: Custom logging for tracking requests and errors.

## API Endpoints

### Admin
#### Login
- **Endpoint**: `POST /user/login`
- **Request Body**:
  ```json
  {
    "email": "your-admin-email",
    "password": "your-admin-password"
  }
  ```
- **Response**:
  - JWT Token containing email and role.

#### Invite Vendor
- **Endpoint**: `POST /admin/invite`
- **Request Body**:
  ```json
  {
    "VendorEmail": "vendor@example.com",
    "VendorFirstName": "FirstName",
    "VendorOtherName": "OtherName"
  }
  ```
- **Response**:
  - Email is sent containing a token as a query parameter.

### Client
#### Register
- **Endpoint**: `POST /user/register`
- **Request Body**:
  ```json
  {
    "ClientFirstName": "John",
    "ClientOtherName": "Doe",
    "ClientEmail": "client@example.com",
    "ClientPassword": "securepassword#123"
  }
  ```

#### Login
- **Endpoint**: `POST /user/login`
- **Request Body**:
  ```json
  {
    "email": "client@example.com",
    "password": "your-password"
  }
  ```
- **Response**:
  - JWT Token containing email and role.

### Vendor
#### Register
- **Endpoint**: `POST /vendor/register`
- **Request Headers**:
  - `Authorization: Bearer <token>` (Token received in the invite email)
- **Request Body**:
  ```json
  {
    "VendorPassword": "securePassword#123",
    "VendorTin": 999991130,
    "ShopName": "Viva Business Group"
  }
  ```

#### Login
- **Endpoint**: `POST /user/login`
- **Request Body**:
  ```json
  {
    "email": "vendor@example.com",
    "password": "your-password"
  }
  ```
- **Response**:
  - JWT Token containing email and role.

---

## Database & ORM
This backend uses **GORM** as the ORM (Object-Relational Mapping) tool to manage database operations efficiently with PostgreSQL. The models for `User`, `Vendor`, `Product`, `Order`, and `Wishlist` are structured using GORM to handle relationships, constraints, and migrations.

## Installation & Setup
1. Clone the repository:
   ```sh
   git clone git@github.com:Mr-Ndi/Miti_Art.git
   cd Miti_Art
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```

### 3. Set Up Environment Variables

You’ll need to configure the environment variables for your application. Create a `.env` file in the root of your project and define the following variables:

#### Database Connection
```env
DATABASE_URL=your_database_url
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_NAME=your_db_name
```

#### Secret Key for JWT Authentication
```env
SECRET_KEY=your_jwt_secret_key
```

#### Application Port
```env
APP_PORT=8080
```

#### Admin Credentials (for email sending)
```env
ADMIN_EMAIL=admin@example.com
ADMIN_EMAIL_PASS=your_admin_password
```

#### Link for Invitation Emails or Verification
```env
LINK=http://localhost:8080/verify
```

4. Run migrations to set up the database:
   ```sh
   go run main.go migrate
   ```

5. Start the server:
   ```sh
   go run main.go
   ```

## Contributing

### Forking the Repository
To contribute to this project:
1. Fork the repository on GitHub.
2. Clone your forked repository:
   ```sh
   git clone git@github.com:your-username/Miti_Art.git
   cd Miti_Art
   ```
3. Create a new branch for your feature or bug fix:
   ```sh
   git checkout -b feature-branch-name
   ```
4. Make your changes and commit them:
   ```sh
   git add .
   git commit -m "Description of changes"
   ```
5. Push your changes to your forked repository:
   ```sh
   git push origin feature-branch-name
   ```
6. Open a pull request on the original repository for review.
