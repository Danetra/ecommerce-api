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