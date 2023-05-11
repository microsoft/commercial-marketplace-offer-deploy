# Setup

- Fork the https://github.com/microsoft/commercial-marketplace-offer-deploy repo and clone to your local machine
- Copy `configs/.env.development.local` to `./bin` and rename it to `.env.local`

# Starting the API Server

- Run `make apiserver-local` from the root directory to build and start the API server
    - Note: This will create a `./bin` folder with the apiserver executable and .env file
- In a new terminal window, test the API server with `curl http://localhost:8080`. You will see logs streamed from the API Server terminal window.

# Environment Testing

- In the terminal, navigate to the `./test` folder
- Test your development environment by running `go test -v`