### Start Mongodb

```shell
$ docker compose up
```

### Companies service

> * Responsible for companies management (CRUD).
> * Main internal logic is located in `service.go` and is exposed over gRPC (could be exposed with HTTP or something
    else too).
> * gRPC server is located in `grpc.go`.
> * Other protocols could be implemented as well, so the main logic in `service.go` is never touched.
> * Available on port: **8001**

```shell
go run cmd/companies/main.go
```

### Gateway

> * Main entry point for our backend system. Exposes HTTP endpoints for all available services (in our case Companies
    only).
> * Gateway is mainly used by clients (web frontend, mobile apps) over HTTP
> * Gateway internally chooses appropriate protocol for communication with backend system/services (Companies in our
    case).
> * Currently, it uses only gRPC client for communication with Companies Service.
> * It could also use HTTP to access Companies service (if it was exposed over HTTP directly).
    This is to demonstrate that Companies could be accessed with two different protocols (gRPC + HTTP).
> * Flow: Clients (web, mobile) -> HTTP -> Gateway -> gRPC/HTTP/others -> Backend services
> * Available on port: **8080**

```shell
go run cmd/gateway/main.go
```

### Run tests

```shell
make test-cover
```

### TODO:

- [ ] Improve test coverage (error handling, invalid data)
- [ ] Improve data validation
- [ ] Use relational database for data storage
- [ ] API documentation and specs (Swagger)
- [ ] **proto** package could, and very likely should, be moved to its own Git repository for better reusability
- [ ] Use appropriate Auth service. Current one is only for demo purpose.


### Endpoints

* Create company
> POST http://localhost:8080/api/v1/companies
```json
{
    "name": "Company 1",
    "description": "Company 1 description",
    "employees_num": 1,
    "registered": true,
    "company_type": 3
}
```

* Modify company
> PUT http://localhost:8080/api/v1/companies/[uuid]
```json
{
    "description": "Company 1 description changed",
    "employees_num": 2,
    "registered": true,
    "company_type": 1
}
```

* Delete company
> DELETE http://localhost:8080/api/v1/companies/[uuid]
```
// no request body
```
* Get company
> GET http://localhost:8080/api/v1/companies/[uuid]
```
// no request body
```
