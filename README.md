## Simple Cinema Booking System

### Requirements
   - Golang version ~1.14
   - Redis
   - Docker
   
### Folder struct   

```
.
├── Dockerfile
├── LICENSE
├── README.md
├── cmd
│   ├── api
│   └── grpc
├── conf
│   └── cinema.toml
├── docker-compose.yml
├── exmsgs
│   └── seat
├── go.mod
├── go.sum
├── pkg
│   ├── core
│   ├── dtos
│   └── transformers
├── protos
│   └── seat
└── tests
````

### Data struct
Since this is a simple application, I use Redis to store the data.
For real applications, the data will be stored in other databases.
For example, seat reservation info will be stored in MySQL,
and seat hold information can be stored in Redis as the session has a timeout.

In this app, I used to Hash data type of Redis. Has 3 main keys:
- `booked`: store booked seats (reservation).
- `pending`: store pending seats and wait for a confirmation to change to booked status.
- `session`: store selected seats of sessions.

Even a session convert to booked status, the application will remove it from `Session` key, 
remove seats from `Pending` key and move them to `Booked` key.

   
### How to run
   - Build and run:
        - Update config to connect your database in `conf/cinema.yaml`.
        - Build app to executable file with command `go build -o app ./cmd/api/seat-api`.
        - Run your app `./app`.
   - Use Docker:
        - Just run `docker-compose up` in project root folder.
   - Check url `http://127.0.0.1:8080/seat/available`.
   - Run test files:
     - `tests/seat_helper_test.go`
     - `cmd/grpc/seat-svc/server_test.go`