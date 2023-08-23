# Product Distributor

## RUN
`go run cmd/distributor/main.go`

This will run an HTTP server on port 8080, you need to change the urls in index.html to match localhost or any other host you're running on

After running the Http server the application will be accesible at given address

Default: http://localhost:8080

## TODOs:
- [ ] Add Configuration file for running the application
- [ ] Add Dockerfile
- [ ] Add More Tests
- [ ] Add a simple UI for the application
- [ ] Add benchmarks
- [ ] Limit requests to GET/POST strictly
- [ ] Add Dockerfile
- [ ] Add authorization for http methods

## API Endpoints:
- Get All: https://api.4li.org/packages
- Remove: https://api.4li.org/packages/remove -b {"id": "package_id"}
- Add: https://api.4li.org/packages/add -b {"id": "package_id", "quantity": 1}

- Order: https://api.4li.org/orders/add -b {"quantity": 1}
