# Loyalty Points API
This is a loyalty points API built with GO and the Gin web framework. The API allows customers to earn points for purchases, and redeem those points for rewards.

## Requirement
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

## License
This project is licensed under the MIT License. See the [LICENSE](https://github.com/Njoguu/loyalty-points-system-api/blob/main/LICENSE) file for details.
