# Loyalty Points API
This is a loyalty points API built with GO and the Gin web framework. The API allows customers to earn points for purchases, and redeem those points for rewards.

## Requirements
- GO (version 1.16 or higher)
- Gin web framework (version 1.7 or higher)
- PostgreSQL database 

## Setup Dev Environment
### Traditional
- Clone the repository: `git clone https://github.com/your-username/loyalty-points-api.git`
- Install dependencies: `go mod download`
- Copy the `.env.sampple` file to `.env` and update the values as needed.
- Start the API: `go run main.go`

### Using Docker

![images](https://user-images.githubusercontent.com/60213982/224079894-df3edad3-cea7-45c4-9c3b-5017926a54b2.png)

## Usage

![images](https://user-images.githubusercontent.com/60213982/224079894-df3edad3-cea7-45c4-9c3b-5017926a54b2.png)

### Endpoints
The following endpoints will be available:

![images](https://user-images.githubusercontent.com/60213982/224079894-df3edad3-cea7-45c4-9c3b-5017926a54b2.png)

All endpoints except for account creation and authentication require the Authorization header with a valid JWT token.

## Documentation

The documentation for this project is generated using [aglio](https://github.com/danielgtaylor/aglio), an API Blueprint renderer that supports multiple themes and outputs static HTML that can be served by any web host.

### Generating Docs

To generate the documentation for this project, follow these steps:

1. Install Aglio by running `npm install -g aglio` in your terminal.

2. Once installed, navigate to the root of the project and run the command `aglio -i docs/index.apib --theme-template triple -o docs/index.html`.

3. This will generate the API documentation in the `docs` folder.

4. Open the `docs/index.html` file in your web browser to view the documentation.

### Updating Docs

If you make changes to the API documentation, you'll need to regenerate the docs by following the steps above.

## License
This project is licensed under the MIT License. See the [LICENSE](https://github.com/Njoguu/loyalty-points-system-api/blob/main/LICENSE) file for details.
