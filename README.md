# Product Distributor

## RUN
`go run cmd/distributor/main.go`

## TODOs:
- [ ] Add Configuration file for running the application
- [ ] Add Dockerfile
- [ ] Add More Tests
- [ ] Add a simple UI for the application
- [ ] Add benchmarks
- [ ] Limit requests to GET/POST strictly

## API Endpoints:
- Get All: https://api.4li.org/packages
- Remove: https://api.4li.org/packages/remove -b {"id": "package_id"}
- Add: https://api.4li.org/packages/add -b {"id": "package_id", "quantity": 1}

- Order: https://api.4li.org/orders/add -b {"quantity": 1}