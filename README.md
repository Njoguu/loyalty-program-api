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
- GO (version 1.17 or higher)
- Gin web framework (version 1.8 or higher)
- PostgreSQL database 

## Environment
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

  - Start application in development

  ```sh
  $ go run main.go | make start
  ```

* ##### Docker Lifecycle

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

## Usage

![images](https://user-images.githubusercontent.com/60213982/224079894-df3edad3-cea7-45c4-9c3b-5017926a54b2.png)

### Endpoints
The following endpoints will be available:

![images](https://user-images.githubusercontent.com/60213982/224079894-df3edad3-cea7-45c4-9c3b-5017926a54b2.png)

All endpoints except for account creation and authentication require the Authorization header with a valid JWT token.

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
