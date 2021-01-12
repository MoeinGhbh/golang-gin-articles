
go get -u github.com/gin-gonic/gin

## make venor and add refrence to the project
go mod vendor
go build -mod=vendor



## make execute file to run web applicatoin 
go build -o app
./app


## Debug project
go to the edit configuration item in Run
select pakage from dropdown menu and select folder which is location of project

F8 F7


# Go gin app
This is the code from the article [Building Go Web Applications and Microservices Using Gin](https://semaphoreci.com/community/tutorials/building-go-web-applications-and-microservices-using-gin).
