## Simple Restful Api

### Requirements
   - Golang version ~1.10
   - Mysql version ~5.6
   - Docker
   
### Folder struct   

```
.
├── cmd                 # Folder contains main service packages
│   └── apis            # Api services
│       └── user-api
├── conf                # Contains config files
├── db                  # Database migration
├── pkg                 # Common packages
│   ├── core
│   │   ├── servehttp
│   │   └── utils
│   ├── dtos            # Data transfer object (DTO) For API response
│   ├── models          # Database model
│   └── transformers    # Helpers to transform between models and DTOs
├── tests               # Contains tests
└── vendor
````
   
### How to run
   - Build and run:
        - You need create a database in you local.
        - Import data in `db/simple.sql` to you database.
        - Update config to connect your database in `conf/app.yaml`
        - Build app to executable file with command `go build -o app ./cmd/apis/user-api`
        - Run your app `./app -cf=./conf/app.yaml`
   - Use Docker:
        - Just run `docker-compose up` in project root folder.
   - Check url `http://127.0.0.1:8080/`