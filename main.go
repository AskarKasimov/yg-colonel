package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/askarkasimov/yg-colonel/db"
	docs "github.com/askarkasimov/yg-colonel/docs"
	models "github.com/askarkasimov/yg-colonel/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @BasePath /api/v1

// @Summary One available expression for worker
// @Tags worker
// @Success 200 {object} models.Expression
// @Failure 404 {object} models.Error "no rows now OR no such worker id"
// @Failure 500 {object} models.Error "unprocessed error"
// @Router /worker/want_to_calculate [get]
func ProvideCalculation(g *gin.Context) {
	if len(g.Request.Header["Authorization"]) == 0 {
		g.JSON(http.StatusBadRequest, models.Error{ErrorMessage: "No Authorization HTTP header"})
		return
	}

	workerId, err := uuid.Parse(g.Request.Header["Authorization"][0])
	if err != nil {
		g.JSON(http.StatusBadRequest, models.Error{ErrorMessage: "Not valid UUID"})
		return
	}

	err = db.DB().WakeUp(workerId)
	if err != nil {
		g.JSON(http.StatusNotFound, models.Error{ErrorMessage: err.Error()})
		return
	}

	if err == sql.ErrNoRows {
		g.JSON(http.StatusNotFound, models.Error{ErrorMessage: "No worker with such ID. Create it"})
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

	id, err := db.DB().GetWorkerIdByName(worker.Name)

	if err == nil {
		g.JSON(http.StatusOK, id)
		return
	}

	createdId, err := db.DB().NewWorker(worker.Name, worker.NumberOfGoroutines)
	if err != nil {
		g.JSON(http.StatusInternalServerError, models.Error{ErrorMessage: err.Error()})
		return
	}
	g.JSON(http.StatusOK, createdId)
}

// @Summary Add a solve on expression
// @Tags expression
// @Accept json
// @Param solve body models.ExpressionSolving true "solve of expression"
// @Success 200 {string} string "id of just created expression"
// @Failure 400 {object} models.Error "incorrect body"
// @Failure 500 {object} models.Error "unprocessed error"
// @Router /expression/solve [post]
func SolveExpression(g *gin.Context) {
	if len(g.Request.Header["Authorization"]) == 0 {
		g.JSON(http.StatusBadRequest, models.Error{ErrorMessage: "No Authorization HTTP header"})
		return
	}

	workerId, err := uuid.Parse(g.Request.Header["Authorization"][0])
	if err != nil {
		g.JSON(http.StatusBadRequest, models.Error{ErrorMessage: "Not valid UUID"})
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
// @Param expression body models.ExpressionAdding true "expression to calculate"
// @Success 200 {string} string "id of just created expression"
// @Failure 400 {object} models.Error "incorrect body"
// @Failure 500 {object} models.Error "unprocessed error"
// @Router /expression/add [post]
func AddExpression(g *gin.Context) {
	var req models.ExpressionAdding
	err := g.ShouldBind(&req)
	if err != nil {
		err2 := g.ShouldBindJSON(&req)
		if err2 != nil {
			g.JSON(http.StatusBadRequest, models.Error{ErrorMessage: err.Error()})
			return
		}
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
		g.JSON(http.StatusOK, []models.ExpressionGeneral{})
		return
	}

	g.JSON(http.StatusOK, expressions)
}

// @Summary Get info about 1 expression
// @Tags expression
// @Param id path int true "Expression ID"
// @Success 200 {object} models.Expression
// @Failure 500 {object} models.Error "unprocessed error"
// @Router /expression/{id} [get]
func GetExpressionInfo(g *gin.Context) {
	expressionId, err := uuid.Parse(g.Param("expressionId"))
	if err != nil {
		g.JSON(http.StatusBadRequest, models.Error{ErrorMessage: "Not valid UUID"})
		return
	}

	expression, err := db.DB().GetExpressionById(expressionId)
	if err != nil {
		g.JSON(http.StatusInternalServerError, models.Error{ErrorMessage: err.Error()})
		return
	}

	g.JSON(http.StatusOK, expression)
}

func AllWorkers(g *gin.Context) {
	workers, err := db.DB().AllWorkers()
	if err != nil {
		g.JSON(http.StatusInternalServerError, models.Error{ErrorMessage: err.Error()})
		return
	}

	g.JSON(http.StatusOK, workers)
}

func MainPage(g *gin.Context) {
	expressions, err := db.DB().AllExpressions()
	if err != nil {
		g.JSON(http.StatusInternalServerError, models.Error{ErrorMessage: err.Error()})
		return
	}

	workers, err := db.DB().AllWorkers()
	if err != nil {
		g.JSON(http.StatusInternalServerError, models.Error{ErrorMessage: err.Error()})
		return
	}

	g.HTML(http.StatusOK, "index.html", gin.H{
		"Expressions": expressions,
		"Workers":     workers,
	})
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.LoadHTMLGlob("templates/*")
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		expression := v1.Group("/expression")
		{
			expression.POST("/add", AddExpression)
			expression.GET("/all", AllExpressions)
			expression.POST("/solve", SolveExpression)
			expression.GET("/:expressionId", GetExpressionInfo)
		}

		worker := v1.Group("/worker")
		{
			worker.GET("/want_to_calculate", ProvideCalculation)
			worker.POST("/register", WorkerRegistration)
			worker.GET("/all", AllWorkers)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	ui := r.Group("/")
	{
		ui.GET("/", MainPage)
	}

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
			if time.Unix(worker.LastHeartbeat, 0).Before(time.Now().Add(-3 * time.Minute)) {
				log.Printf("%s IS OFFLINE NOW", worker.Name)
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
		time.Sleep(2 * time.Minute)
	}
}
