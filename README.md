# Loyalty Points API
[![LICENSE](https://img.shields.io/github/license/Njoguu/loyalty-program-api?color=blue)](https://github.com/Njoguu/loyalty-points-system-api/blob/main/LICENSE) 
[![Go](https://img.shields.io/github/go-mod/go-version/Njoguu/loyalty-program-api)](https://github.com/Njoguu/loyalty-program-api)
[![Github Issues](https://img.shields.io/github/issues-raw/Njoguu/loyalty-program-api)](https://github.com/Njoguu/loyalty-program-api/issues) 
[![Github pull requests](https://img.shields.io/github/issues-pr-raw/Njoguu/loyalty-program-api?color=yellow)](https://github.com/Njoguu/loyalty-program-api/pulls)
<br>
[![Build](https://github.com/Njoguu/loyalty-program-api/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/Njoguu/loyalty-program-api/actions/workflows/build.yml)
[![CodeQL Analysis](https://github.com/Njoguu/loyalty-program-api/actions/workflows/codeql_analysis.yml/badge.svg?branch=main)](https://github.com/Njoguu/loyalty-program-api/actions/workflows/codeql_analysis.yml) <br>
[![pages-build-deployment](https://github.com/Njoguu/loyalty-program-api/actions/workflows/pages/pages-build-deployment/badge.svg?branch=main)](https://github.com/Njoguu/loyalty-program-api/actions/workflows/pages/pages-build-deployment)
<br>

This is a loyalty points API built with GO and the Gin web framework. The API allows customers to earn points for purchases, and redeem those points for rewards.

## Requirements
- GO (version 1.18 or higher)

## Supported databases
- PostgreSQL
> **Note** loyalty-program-api uses [GORM](https://github.com/go-gorm/gorm) as its ORM

## Features
The API has put in place the following features:
- [x] Built on top of [Gin](https://github.com/gin-gonic/gin)
- [x] Uses the supported database without writing any extra configuration files
- [x] Reading Environment variables using [godotenv](https://github.com/joho/godotenv)
- [x] Caching responses
- [ ] Logging
- [ ] Follows CORS policy
- [x] DB Migration Support
- [ ] Comprehensive Error Handling in API Services
- [x] Basic auth
- [x] Password hashing with `bcrypt`
- [ ] Simple firewall (whitelist/blacklist IP)
- [ ] Webhook Integration for Real-time Deployment Status Notifications
- [x] Request data validation
- [x] Email verification (sending verification email)
- [x] Forgot password recovery
- [x] Render `HTML` templates
- [ ] Forward error logs and crash reports.
- [x] User logout (Send an expired cookie to user’s browser or client which invalidates the user’s ‘session’)

## Project Structure
The project follows the following directory structure:
```
.
├── controllers
│   ├── admin
│   │   ├── admin-auth.go
│   │   └── admin-endpoints.go
│   ├── get-products.go
│   ├── googleauth.go
│   ├── pass-reset-controller.go
│   ├── points.go
│   ├── transactions-history.go
│   └── user-auth.go
├── db
│   ├── migrations
│   │   ├── 000001_init_schema.down.sql
│   │   └── 000001_init_schema.up.sql
├── docs
│   ├── admins
│   │   └── admins.apib
│   ├── auth.apib
│   ├── index.apib
│   ├── index.html
│   ├── points.apib
│   ├── products.apib
│   ├── reset-password.apib
│   └── transaction-history.apib
├── middlewares
│   └── caching.go
│   ├── logger.go
│   └── middlewares.go
├── migrate
│   └── migrate.go
├── models
│   ├── admins.go
│   ├── products.go
│   ├── setup.go
│   ├── transactions.go
│   └── user.go
├── templates
│   ├── base.html
│   ├── passwordReset.html
│   ├── styles.html
│   ├── templates.html
│   └── verificationCode.html
├── utils
│   ├── mail
│   │   ├── email.go
│   │   └── encode.go
│   ├── token
│   │   └── token.go
│   ├── googleOAuth.go
│   └── utils.go
├── .gitignore
├── Dockerfile
├── LICENSE
├── Makefile
├── README.md
├── docker-compose.yml
├── env.sample
├── go.mod
├── go.sum
└── main.go

```

## Environment Installation
### Development Setup Instructions
Follow the steps below to set up the development environment for the Loyalty Program API:

1. Clone the Loyalty Program API repository from the version control system of your choice (e.g., GitHub).

```shell
git clone https://github.com/Njoguu/loyalty-program-api.git
cd loyalty-program-api
```

2. Copy the example environment file and update the configuration variables as needed.

```shell
cp .env.example .env
```

3. Open the .env file in a text editor and provide the necessary values for the environment variables. Make sure to set the database connection details, JWT secret, and any other required variables.

4. Start the API and its dependencies using Docker Compose.

```shell
docker-compose up
```
Docker Compose will build and start the Loyalty Program API, along with the required database and any other defined services. The API will be accessible at http://localhost:8000.

5. Make changes to the API code as needed. The API server will automatically reload whenever code changes are detected, allowing for a seamless development experience.

#### Commands

- ##### Application Lifecycle

  - Install dependencies

  ```sh
  $ go get . || go mod || make install
  ```

  - Build application

  ```sh
  $ go build -o loyalty-program-api || make build
  ```

  - Start application server in development

  ```sh
  $ go run main.go | make start
  ```

* ##### Docker Lifecycle
  - Start postgres container:
  
  ```sh
  $ make postgres
  ```
  
  - Create simple_bank database:
    
  ```sh
  $ make createdb
  ```
  - Create a new db migration:

  ```sh
  $ make create-migrations
  ```

  - Run db migration up all versions:
  
  ```sh
  $ make migrateup
  ```
  
  - Run db migration down all versions:

  ```sh
  $ make migratedown
  ```
    
  - Build container
  
  ```sh
  $ docker-compose build | make dcb
  ```
  
  - Run container
  
  ```sh
  $ docker-compose up -d --build | make dcu
  ```
  
  - Stop container
  
  ```sh
  $ docker-compose down | make dcd
  ```

## API Routes

### User Authentication

Users have the following routes available for their authentication:

| Path          | Method | Required JSON | Header                                | Description                                                             |
| ------------- | ------ | ------------- | ------------------------------------- | ----------------------------------------------------------------------- |
| /api/auth/register   | POST   | username, firstname, lastname, gender, email, password, city, phone_number         |                                       | Create User account and send activation email                      |
| /api/auth/verify-email/{secret_code}   | GET   |          |                                       | Verify email of correspondent account and grant login access |
| /api/auth/login | POST   | email, password               |  | Validate logging in of user                                                            |
| /api/auth/logout  | GET   |               | Authorizaiton: "Bearer token" | Invalidate the user’s ‘session’.                                                  |



## Documentation

The [documentation](https://njoguu.github.io/loyalty-program-api/) for this project is generated using [aglio](https://github.com/danielgtaylor/aglio), an API Blueprint renderer that supports multiple themes and outputs static HTML that can be served by any web host.

### Generating Docs

To generate the documentation for this project, follow these steps:

1. Install Aglio by running `npm install -g aglio` in your terminal.

2. Once installed, navigate to the root of the project and run the command `aglio -i docs/index.apib -o docs/index.html`.

3. This will generate the API documentation in the `docs` folder.

4. Open the `docs/index.html` file in your web browser to view the documentation.

### Updating Docs

If you make changes to the API documentation, you'll need to regenerate the docs by following the steps above.

## License
This project is licensed under the MIT License. See the [LICENSE](https://github.com/Njoguu/loyalty-points-system-api/blob/main/LICENSE) file for details.
