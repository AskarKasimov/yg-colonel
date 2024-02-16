package main

import (
	"log"
	"net/http"

	docs "github.com/askarkasimov/yg-colonel/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}

}

type Expression struct {
	Vanilla string
	Answer  string
}

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func ProvideCalculation(g *gin.Context) {
	g.JSON(http.StatusOK, "helloworld")
}

func AddExpression(g *gin.Context) {

}

func main() {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		expression := v1.Group("/expression")
		{
			expression.POST("/add", AddExpression)
		}

		worker := v1.Group("/worker")
		{
			worker.GET("/want_to_calculate", ProvideCalculation)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")
}
