package main

import (
	"log"
	"net/http"

	"github.com/askarkasimov/yg-colonel/db"
	docs "github.com/askarkasimov/yg-colonel/docs"
	models "github.com/askarkasimov/yg-colonel/models"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}

}

// @BasePath /api/v1

// @Summary One available expression for worker
// @Tags worker
// @Success 200 {object} models.Expression
// @Failure 404 {object} models.Error
// @Router /worker/want_to_calculate [get]
func ProvideCalculation(g *gin.Context) {
	expression, err := db.DB().GetAvailableExpression()
	if err != nil {
		g.JSON(http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}

	g.JSON(http.StatusOK, expression)
}

// @Summary Add an expression
// @Tags expression
// @Accept json
// @Param expression body models.ExpressionAdding true "expression to calculate"
// @Success 200 {string} string "id of just created expression"
// @Failure 404 {object} models.Error
// @Router /expression/add [post]
func AddExpression(g *gin.Context) {
	var req models.ExpressionAdding
	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, models.Error{Message: err.Error()})
		return
	}

	id, err := db.DB().AddExpression(req)
	if err != nil {
		g.JSON(http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}

	g.JSON(http.StatusOK, id)
}

// @Summary Get all expressions
// @Tags expression
// @Success 200 {object} []models.Expression
// @Failure 500 {object} models.Error
// @Router /expression/all [get]
func AllExpressions(g *gin.Context) {
	expressions, err := db.DB().AllExpressions()
	if err != nil {
		g.JSON(http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}

	g.JSON(http.StatusOK, expressions)
}

func main() {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		expression := v1.Group("/expression")
		{
			expression.POST("/add", AddExpression)
			expression.GET("/all", AllExpressions)
		}

		worker := v1.Group("/worker")
		{
			worker.GET("/want_to_calculate", ProvideCalculation)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")
}
