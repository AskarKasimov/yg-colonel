package main

import (
	"database/sql"
	"net/http"

	"github.com/askarkasimov/yg-colonel/db"
	docs "github.com/askarkasimov/yg-colonel/docs"
	models "github.com/askarkasimov/yg-colonel/models"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @BasePath /api/v1

// @Summary One available expression for worker
// @Tags worker
// @Success 200 {object} models.Expression
// @Failure 404 {object} models.Error "no rows now"
// @Failure 500 {object} models.Error "unprocessed error"
// @Router /worker/want_to_calculate [get]
func ProvideCalculation(g *gin.Context) {
	expression, err := db.DB().GetOneAvailableExpression()
	if err == sql.ErrNoRows {
		g.JSON(http.StatusNotFound, models.Error{ErrorMessage: err.Error()})
		return
	}
	if err != nil {
		g.JSON(http.StatusInternalServerError, models.Error{ErrorMessage: err.Error()})
		return
	}

	g.JSON(http.StatusOK, expression)
}

// @Summary Add an expression
// @Tags expression
// @Accept json
// @Param expression body models.ExpressionAdding true "expression to calculate"
// @Success 200 {string} string "id of just created expression"
// @Failure 400 {object} models.Error "incorrect body"
// @Failure 500 {object} models.Error "unprocessed error"
// @Router /expression/add [post]
func AddExpression(g *gin.Context) {
	var req models.ExpressionAdding

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, models.Error{ErrorMessage: err.Error()})
		return
	}

	id, err := db.DB().AddExpression(req)
	if err != nil {
		g.JSON(http.StatusInternalServerError, models.Error{ErrorMessage: err.Error()})
		return
	}

	g.JSON(http.StatusOK, id)
}

// @Summary Get all expressions
// @Tags expression
// @Success 200 {object} []models.Expression
// @Failure 500 {object} models.Error "unprocessed error"
// @Router /expression/all [get]
func AllExpressions(g *gin.Context) {
	expressions, err := db.DB().AllExpressions()

	if err != nil {
		g.JSON(http.StatusInternalServerError, models.Error{ErrorMessage: err.Error()})
		return
	}

	if expressions == nil {
		g.JSON(http.StatusOK, []models.Expression{})
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
