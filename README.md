# Building basic RESTful (CRUD) with Golang

>This project focuses on developing a basic RESTful CRUD API’s using Golang & MySql as backend database.




## Introduction
We would be developing an application that exposes a basic REST-API server for CRUD operations for managing files and directorys

### Run docker 
    docker pull roy990427/hp_test:latest
    docker run -d -p 8080:8080 hp_test
### Run server  
	go mod download
	go build -o app
	go run main.go
### Run test 
    docker run -d -p 8080:8080 hp_test
    docker exec -it {docker-container-ID} /bin/sh
	go test -v 
