package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"net/http"
	"xm/gateway"
	"xm/internal/grpc"
	"xm/internal/web"
	"xm/proto/gen"
)

const (
	clientIDKey  = "client-id"
	clientApiKey = "api-key"
)

var (
	//TODO: store api keys in database
	supportedApiKeys = map[string]string{
		"default-client": "33imUw5xxF",
	}
)

var (
	httpAddr = flag.String("http", ":8080", "Http address")
	cmpAddr  = flag.String("cmp_addr", ":8001", "Company service address")
)

func main() {
	flag.Parse()

	router := web.NewRouter()

	api := router.Group("api")
	api.GET("healthz", healthz())

	cc := grpc.CreateClientConnection(*cmpAddr)
	client := gen.NewCompaniesClient(cc)

	companiesApi := api.Group("v1/companies")
	companiesApi.GET("healthz", healthz())
	companiesApi.GET("/:id", gateway.FindCompanyByIDHandler(client))

	companiesApi.Use(BasicAuthMiddleware())
	{
		companiesApi.POST("", gateway.AddCompanyHandler(client))
		companiesApi.PUT("/:id", gateway.ModifyCompanyHandler(client))
		companiesApi.DELETE("/:id", gateway.DeleteCompanyHandler(client))
	}

	web.ServeHttp(*httpAddr, "gateway", router)
}

func healthz() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, struct {
			Status string `json:"status"`
		}{
			Status: "OK",
		})
	}
}

// BasicAuthMiddleware will compare
func BasicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := c.Request.Header.Get(clientIDKey)
		clientKey := c.Request.Header.Get(clientApiKey)

		apiKey, ok := supportedApiKeys[clientID]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "invalid api key")
			return
		}

		if apiKey != clientKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "invalid credentials")
			return
		}

		c.Next()
	}
}
