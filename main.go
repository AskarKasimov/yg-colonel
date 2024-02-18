package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

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
	if len(g.Request.Header["Authorization"]) == 0 {
		g.JSON(http.StatusBadRequest, models.Error{ErrorMessage: "No Authorization HTTP header"})
		return
	}

	workerId, err := strconv.ParseInt(g.Request.Header["Authorization"][0], 10, 64)
	if err != nil {
		g.JSON(http.StatusBadRequest, models.Error{ErrorMessage: "Error with parsing ID from Authorization HTTP header"})
		return
	}

	isAlive, err := db.DB().IsWorkerAlive(workerId)

	if err == sql.ErrNoRows {
		g.JSON(http.StatusNotFound, models.Error{ErrorMessage: "No worker with such ID. Create it"})
		return
	}
	if err != nil {
		g.JSON(http.StatusInternalServerError, models.Error{ErrorMessage: err.Error()})
		return
	}
	if !isAlive {
		g.JSON(http.StatusNotFound, models.Error{ErrorMessage: "Worker is not alive! Ensure your heartbeats work"})
		return
	}

	expression, err := db.DB().GetOneAvailableExpression(workerId)
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

// @Summary Handling heartbeats
// @Tags worker
// @Success 200 {object} models.Expression
// @Failure 400 {object} models.Error "parsing ID error"
// @Failure 404 {object} models.Error "no worker with such ID"
// @Failure 500 {object} models.Error "unprocessed error"
// @Router /worker/heartbeat [post]
func Heartbeat(g *gin.Context) {
	if len(g.Request.Header["Authorization"]) == 0 {
		g.JSON(http.StatusBadRequest, models.Error{ErrorMessage: "No Authorization HTTP header"})
		return
	}

	workerId, err := strconv.ParseInt(g.Request.Header["Authorization"][0], 10, 64)
	if err != nil {
		g.JSON(http.StatusBadRequest, models.Error{ErrorMessage: "Error with parsing ID from Authorization HTTP header"})
		return
	}

	err = db.DB().WakeUp(workerId)
	if err != nil {
		g.JSON(http.StatusInternalServerError, models.Error{ErrorMessage: err.Error()})
		return
	}

	g.JSON(http.StatusOK, "OK")
}

// @Summary Registrating worker in orchestrator
// @Tags worker
// @Accept json
// @Param expression body models.ExpressionAdding true "expression to calculate"
// @Success 200 {string} string "id of just created expression"
// @Failure 400 {object} models.Error "incorrect body"
// @Failure 500 {object} models.Error "unprocessed error"
// @Router /worker/register [post]
func WorkerRegistration(g *gin.Context) {
	var worker models.WorkerAdding

	if err := g.ShouldBindJSON(&worker); err != nil {
		g.JSON(http.StatusBadRequest, models.Error{ErrorMessage: err.Error()})
		return
	}

	id, _ := db.DB().GetWorkerIdByName(worker.Name)

	if id != 0 {
		g.JSON(http.StatusOK, id)
		return
	}

	createdId, err := db.DB().NewWorker(worker.Name)
	if err != nil {
		g.JSON(http.StatusInternalServerError, models.Error{ErrorMessage: err.Error()})
		return
	}
	g.JSON(http.StatusOK, createdId)
}

func SolveExpression(g *gin.Context) {
	if len(g.Request.Header["Authorization"]) == 0 {
		g.JSON(http.StatusBadRequest, models.Error{ErrorMessage: "No Authorization HTTP header"})
		return
	}

	workerId, err := strconv.ParseInt(g.Request.Header["Authorization"][0], 10, 64)
	if err != nil {
		g.JSON(http.StatusBadRequest, models.Error{ErrorMessage: "Error with parsing ID from Authorization HTTP header"})
		return
	}

	var ans models.ExpressionSolving

	if err := g.ShouldBindJSON(&ans); err != nil {
		g.JSON(http.StatusBadRequest, models.Error{ErrorMessage: err.Error()})
		return
	}

	err = db.DB().SolveExpression(workerId, ans.Id, ans.Answer)
	if err != nil {
		g.JSON(http.StatusInternalServerError, models.Error{ErrorMessage: err.Error()})
		return
	}

	g.JSON(http.StatusOK, "OK")
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
			expression.POST("/solve", SolveExpression)
		}

		worker := v1.Group("/worker")
		{
			worker.GET("/want_to_calculate", ProvideCalculation)
			worker.POST("/heartbeat", Heartbeat)
			worker.POST("/register", WorkerRegistration)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	go chekingWorkers()

	r.Run(":8080")
}

func handleError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func chekingWorkers() {
	for {
		workers, err := db.DB().AllAliveWorkers()
		handleError(err)
		for _, worker := range workers {
			if time.Unix(worker.LastHeartbeat, 0).Before(time.Now().Add(-1 * time.Minute)) {
				log.Println(worker.Name, " IS OFFLINE NOW")
				err = db.DB().FallAsleep(worker.Id)
				handleError(err)

				activeExpressions, err := db.DB().GetActiveExpressionsFromWorker(worker.Id)
				handleError(err)

				for _, expression := range activeExpressions {
					err = db.DB().MakeExpressionAvailableAgain(expression.Id)
					handleError(err)
				}
			}
		}
		time.Sleep(1 * time.Minute)
	}
}
