# Miti Arts Backend

Miti Arts is an online store that sells high-quality timber-based furniture. This repository contains the backend logic for managing furniture products and user accounts.

## Technologies Used:
- **Gin Framework**: A web framework for Go (Golang) to build high-performance REST APIs.
- **GORM**: An ORM library for Golang to interact with the PostgreSQL database efficiently.
- **PostgreSQL**: Relational database for storing furniture and user data.
- **JWT (JSON Web Token)**: For secure user authentication.
- **Go Modules**: Dependency management for Go.

## API Documentation

### Authentication
All protected routes require a JWT token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

### Client Endpoints

#### Public Routes
1. **Register Client**
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

2. **Get All Furniture**
   - **Endpoint**: `GET /user/furniture`
   - **Response**: List of all furniture items

3. **Get Furniture Details**
   - **Endpoint**: `GET /user/furniture/:id`
   - **Response**: Detailed information about a specific furniture item

#### Protected Routes (Requires Authentication)
1. **Create Order**
   - **Endpoint**: `POST /user/order/:id`
   - **Description**: Creates a new order for a specific furniture item
   - **URL Parameters**: `id` (furniture item ID)

2. **Add to Wishlist**
   - **Endpoint**: `POST /user/wished-item/:id`
   - **Description**: Adds a furniture item to user's wishlist
   - **URL Parameters**: `id` (furniture item ID)

3. **List User Orders**
   - **Endpoint**: `POST /user/my-orders`
   - **Description**: Retrieves all orders made by the authenticated user

### Vendor Endpoints

#### Public Routes
1. **Register Vendor**
   - **Endpoint**: `POST /vendor/register`
   - **Headers**: 
     - `Authorization: Bearer <invitation_token>`
   - **Request Body**:
     ```json
     {
       "VendorPassword": "securePassword#123",
       "VendorTin": 999991130,
       "ShopName": "Viva Business Group"
     }
     ```

#### Protected Routes (Requires Authentication)
1. **Upload Product**
   - **Endpoint**: `POST /vendor/upload`
   - **Description**: Upload a new furniture product

2. **Get My Products**
   - **Endpoint**: `GET /vendor/my-products`
   - **Description**: List all products uploaded by the vendor

3. **Get Single Product**
   - **Endpoint**: `POST /vendor/my-product/:id`
   - **URL Parameters**: `id` (product ID)

4. **Get Required Products**
   - **Endpoint**: `GET /vendor/required-product`
   - **Description**: List all orders for vendor's products

5. **Edit Product**
   - **Endpoint**: `POST /vendor/edit-product/:id`
   - **URL Parameters**: `id` (product ID)

6. **Delete Product**
   - **Endpoint**: `POST /vendor/remove-product/:id`
   - **URL Parameters**: `id` (product ID)

### Admin Endpoints

#### Public Routes
1. **Login**
   - **Endpoint**: `POST /user/login`
   - **Request Body**:
     ```json
     {
       "email": "admin@example.com",
       "password": "your-password"
     }
     ```
   - **Response**: JWT token with admin role

#### Protected Routes (Requires Admin Authentication)
1. **Invite Vendor**
   - **Endpoint**: `POST /admin/invite`
   - **Request Body**:
     ```json
     {
       "VendorEmail": "vendor@example.com",
       "VendorFirstName": "FirstName",
       "VendorOtherName": "OtherName"
     }
     ```

2. **View Clients**
   - **Endpoint**: `GET /admin/view-clients`
   - **Description**: List all registered clients

3. **View Vendors**
   - **Endpoint**: `GET /admin/view-vendors`
   - **Description**: List all registered vendors

4. **View Orders**
   - **Endpoint**: `GET /admin/view-orders`
   - **Description**: List all orders in the system

5. **View Products**
   - **Endpoint**: `GET /admin/view-products`
   - **Description**: List all products in the system

6. **Edit Vendor**
   - **Endpoint**: `POST /admin/edit-vendor`
   - **Description**: Update vendor information

7. **Edit Client**
   - **Endpoint**: `POST /admin/edit-client`
   - **Description**: Update client information

8. **Delete Vendor**
   - **Endpoint**: `POST /admin/eliminate-vendor`
   - **Description**: Remove a vendor from the system

9. **Delete Client**
   - **Endpoint**: `POST /admin/eliminate-client`
   - **Description**: Remove a client from the system

## Error Responses
All endpoints return errors in the following format:
```json
{
  "error": "Error message description"
}
```

Common HTTP Status Codes:
- 200: Success
- 201: Created
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 500: Internal Server Error

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

Create a `.env` file in the root of your project with the following variables:

```env
DATABASE_URL=your_database_url
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_NAME=your_db_name
SECRET_KEY=your_jwt_secret_key
APP_PORT=8080
ADMIN_EMAIL=admin@example.com
ADMIN_EMAIL_PASS=your_admin_password
LINK=http://localhost:8080/verify
```

4. Run migrations:
   ```sh
   go run main.go migrate
   ```

5. Start the server:
   ```sh
   go run main.go
   ```

## Contributing
Please refer to the original contributing guidelines in the repository.
