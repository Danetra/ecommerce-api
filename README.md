# 🛒 Ecommerce API

Ecommerce API adalah RESTful API backend yang dibangun menggunakan **Golang (Gin Framework)** untuk mendukung sistem e-commerce modern.  
API ini menyediakan fitur **autentikasi**, **manajemen pengguna & role**, **produk & kategori**, serta **transaksi**, dengan fokus pada **keamanan, performa, dan skalabilitas**.

Project ini dirancang mengikuti **Go project best practice**, mendukung **Swagger API documentation**, **JWT authentication**, dan **database migration**, serta siap dideploy ke **Railway**.

---

## 🚀 Features

- 🔐 **Authentication & Authorization**
    - Register & login user
    - JWT-based authentication
    - Role-based access control (Admin / User)
    - Middleware untuk protected routes

- 👤 **User & Role Management**
    - CRUD user
    - Manajemen role
    - Validasi akses berbasis role

- 🛒 **Product & Category Management**
    - CRUD produk
    - CRUD kategori produk
    - Relasi produk dengan kategori

- 💳 **Transaction Management**
    - Create transaksi pembelian
    - Support payment method
    - Reference number untuk transaksi
    - Relasi user ↔ produk ↔ transaksi

- 🗄️ **Database Migration**
    - PostgreSQL
    - SQL-based migration menggunakan `sql-migrate`
    - Auto-run migration saat aplikasi dijalankan

- 📄 **API Documentation**
    - Swagger UI
    - Auto-generated menggunakan `swaggo`
    - Endpoint `/swagger/index.html`

---

## 🧱 Project Structure

ecommerce-api/
├─ config/
│ └─ database.go                                            # Config Database
├─ controllers/
│ └─ authController.go                                      # Auth Controller (Login & Register)
│ └─ productCategoryController.go                           # CRUD Product Category
│ └─ productController.go                                   # CRUD Product
│ └─ roleController.go                                      # CRUD Product
│ └─ transactionController.go                               # Create Transaction, Get All Transaction, Get History Transaction, Transaction Payment
│ └─ userController.go                                      # Get User By Id & Update User
├─ helpers/
│ └─ baseUrl.go
│ └─ fileUrl.go                                             # Configuration upload image
│ └─ jwt.go                                                 # helpers for jwt
├─ middleware/
│ └─ jwt.go                                                 # middleware for auth
├─ migrations/                                              # migration table
│ └─ 20260114215500_create_roles_table.sql
│ └─ 20260114215900_create_users_table.sql
│ └─ 20260114234200_create_product_categories_table.sql
│ └─ 20260114235000_create_products_table.sql
│ └─ 20260114235500_create_transactions_table.sql
│ └─ 20260114235700_add_payment_method_and_reference_number_at_transactions_table.sql
├─ models/                                                  # models
│ └─ auth.go
│ └─ product.go
│ └─ product_category.go
│ └─ role.go
│ └─ transaction.go
│ └─ user.go
├─ requests/                                                # request 
│ └─ UserRequest.go
├─ responses/                                               # response handler
│ └─ ProductResponse.go
│ └─ RoleResponse.go
│ └─ TransactionResponse.go
│ └─ UserDetailResponse.go
├─ routes/                                                  # routes
│ └─ api.go
│ └─ productCategoryRoute.go
│ └─ productRoute.go
│ └─ roleRoute.go
│ └─ transactionRoute.go
│ └─ userRoute.go
├─ .env                                                     # Environment project
└─ env.example                                              # Environment example
└─ go.mod                                                   # go dependencies
└─ main.go                                                  # main application entry point
└─ README.md



---

## 🛠️ Tech Stack

| Layer | Technology |
|------|------------|
| Language | Go (Golang) |
| Framework | Gin |
| Database | PostgreSQL |
| Authentication | JWT |
| Migration | sql-migrate |
| API Docs | Swagger (swaggo) |
| Deployment | Railway |

---

## ⚙️ Environment Variables

Buat file `.env` di root project


---

## ▶️ Run Project Locally

### Install dependencies
go mod tidy

### Generate Swagger documentation
swag init

### Generate Air
air init

### Run application
air / go run main.go

Akses aplikasi:
- API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html

---

## 🗄️ Database Migration

Migration akan dijalankan **secara otomatis** saat aplikasi start.

Jika ingin menjalankan manual:

sql-migrate up

---

## 🏗️ Build Application
go build


---

## 🚀 Deployment (Railway)

**Build Command**

go build -ldflags="-w -s" -o out


**Start Command**

./out


---

## 🔐 Authentication Usage

Gunakan JWT pada header request:
Authorization: Bearer <your_token>


---

## 📄 API Documentation

Swagger UI tersedia di endpoint:

/swagger/index.html


---

## 📌 Future Improvements

- Refresh token
- Pagination & filtering
- Unit & integration testing
- Docker support
- CI/CD pipeline

---

## 👨‍💻 Author

**Daniel Eka Putra**  
Backend Engineer  
Golang • REST API • PostgreSQL • JWT



