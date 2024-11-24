# GoGateway API

# Index
- [Prerequisites](#prerequisites)
- [Setup](#setup)
  - [Clone the Repository](#clone-the-repository)
  - [Create Configuration File](#create-configuration-file)
  - [Build and Run the Docker Containers](#build-and-run-the-docker-containers)
    - [Test the router](#test-the-router)
- [Client](#client)
- [Server](#server)
  - [Generate a JWT Token](#generate-a-jwt-token)
    - [JWT Auth in Depth](#jwt-auth-in-depth)
      - [You can also modify the claims](https://datatracker.ietf.org/doc/html/rfc7519#section-4)
  - [Use the JWT Token to Authenticate](#use-the-jwt-token-to-authenticate)
- [Client](#client-1)
- [Server](#server-1)
  - [Test the Rate Limiting](#test-the-rate-limiting)
- [Project Structure](#project-structure)
- [License](#license)

## Prerequisites

- [Docker](https://docs.docker.com/get-started/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Go](https://go.dev/doc/install)

## Setup

### Clone the Repository

```sh
git clone https://github.com/zvdy/gogateway-api.git
cd gogateway-api
```

### Create Configuration File

Create a `config.yaml` file in the root directory with the following content:

```yaml
ServerAddress: ":8080"
RedisAddress: "redis:6379"
RedisPassword: ""
RedisDB: 0
```

### Build and Run the Docker Containers

Use Docker Compose to build and run the containers:

```sh
docker-compose up --build
```

#### Test the router

```sh
# client
curl http://localhost:8080/health
{"status":"healthy"}
```

```sh
#server
api-gateway        | [GIN] 2024/11/24 - 19:18:49 | 200 |     151.849µs |      172.27.0.1 | GET      "/health"
```

### Generate a JWT Token

Create a `gen_token.go` file in the root directory with the following content:

> [!WARNING]  
> If you modify the claims, the script will have to be modified too.

```go
package main

import (
	"fmt"
	"log"

	"github.com/zvdy/gogateway-api/internal/auth"
)

func main() {
	token, err := auth.GenerateJWT("test-user-id")
	if err != nil {
		log.Fatalf("Failed to generate JWT: %v", err)
	}
	fmt.Println("Generated JWT:", token)
}
```

Run the script to generate a JWT token:

```sh
go run gen_token.go
Generated JWT: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzI1NjE5NTksInVzZXJfaWQiOiJ0ZXN0LXVzZXItaWQifQ.79bjRnzgN7ub1757ecWOj0-cx4uo0dnjGU1ZLCPGfaQ
```

#### JWT Auth in Depth

`internal/auth/auth.go` contains the logic to generate the tokens. at the moment.

The key is hardcoded, and as a personal recommendation, I suggest using `openssl rand -base64 32`.
```go
var secretKey = []byte("vyPcARcpHaov7o7aU1kDcLHjR0ZR9+UWx/TqtCvhl+g=") // ❯ openssl rand -base64 32
```

You can also modify the claims [rfc7519](https://datatracker.ietf.org/doc/html/rfc7519#section-4).

```go
func GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
```

I suggest using [this web](https://dinochiesa.github.io/jwt/) To verify jwt signatures and general work with JWT.


### Use the JWT Token to Authenticate

Replace `<JWT_TOKEN>` with the token generated in the previous step and use it to authenticate and access the API Gateway:

```sh
# client
export TOKEN='<JWT_TOKEN>'
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/test
Hello from backend
```

```sh
# server
api-gateway        | [GIN] 2024/11/24 - 19:17:20 | 200 |     684.511µs |      172.27.0.1 | GET      "/api/test"
```

### Test the Rate Limiting

Create a loadtest.sh file in the root directory with the following content:
```sh
#!/bin/bash

TOKEN='<JWT_TOKEN>'
URL='http://localhost:8080/api/test'

for i in {1..110}; do
  echo "Request $i"
  curl -H "Authorization: Bearer $TOKEN" $URL
  echo
done
```

Run the script to test the rate limiting:

```sh
chmod +x loadtest.sh
./loadtest.sh
...
...
...
Request 101
{"error":"rate limit exceeded"}
```

## Project Structure

```
.
├── cmd
│   └── main.go
├── config.yaml
├── Dockerfile
├── Dockerfile.backend
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal
│   ├── auth
│   │   └── auth.go
│   ├── cache
│   │   ├── cache.go
│   │   └── middleware.go
│   ├── logging
│   │   └── logger.go
│   ├── proxy
│   │   └── proxy.go
│   ├── ratelimit
│   │   ├── middleware.go
│   │   └── ratelimiter.go
│   └── routes
│       ├── handlers.go
│       └── router.go
├── pkg
│   ├── middleware
│   │   ├── cors.go
│   │   └── recovery.go
│   └── utils
│       ├── http.go
│       └── json.go
├── backend
│   └── backend.go
├── gen_token.go
├── loadtest.sh
├── README.md
└── test
    └── integration
        ├── auth_test.go
```

## License

This project is licensed under the [MIT License](LICENSE).
