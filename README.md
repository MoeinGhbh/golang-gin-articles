

## install gin

go get -u github.com/gin-gonic/gin


The route handler has a pointer to the context (gin.Context) as its parameter. 
This context contains all the information about the request that the handler might
need to process it. For example, it includes information about the headers,
cookies, etc.