package db

import (
	"database/sql"

	"github.com/askarkasimov/yg-colonel/models"
	_ "github.com/lib/pq"
)

type iDatabase interface {
	GetOneAvailableExpression(workerId int64) (models.Expression, error)
	AddExpression(e models.ExpressionAdding) (int64, error)
	AllExpressions() ([]models.ExpressionGeneral, error)
	IsWorkerAlive(workerId int64) (bool, error)
	WakeUp(workerId int64) error
	GetWorkerIdByName(name string) (int64, error)
	NewWorker(name string) (int64, error)
	AllAliveWorkers() ([]models.Worker, error)
	FallAsleep(workerId int64) error
	GetActiveExpressionsFromWorker(workerId int64) ([]models.Expression, error)
	MakeExpressionAvailableAgain(expressionId int64) error
	SolveExpression(workerId, expressionId int64, solution string) error
	GetExpressionById(expressionId int64) (models.Expression, error)
	AllWorkers() ([]models.Worker, error)
}

type database struct {
	db *sql.DB
}

var db iDatabase

func DB() iDatabase { return db }

func init() {
	connStr := "user=admin password=admin dbname=yg sslmode=disable host=postgres"
	newConn, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	db = &database{db: newConn}
}

// adding expression into DB and giving its id back
func (d *database) AddExpression(e models.ExpressionAdding) (int64, error) {
	var id int64
	err := d.db.QueryRow("INSERT INTO expressions (vanilla) VALUES ($1) RETURNING id", e.Expression).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// taking the oldest untaken expression
func (d *database) GetOneAvailableExpression(workerId int64) (models.Expression, error) {
	var expression models.Expression

	err := d.db.QueryRow(
		"SELECT id, extract(epoch from incomingDate)::INT, vanilla, answer, progress FROM expressions WHERE progress='waiting' ORDER BY incomingDate LIMIT 1").Scan(
		&expression.Id, &expression.IncomingDate, &expression.Vanilla, &expression.Answer, &expression.Progress)
	if err != nil {
		return models.Expression{}, err
	}

	ex, err := d.db.Exec("UPDATE expressions SET progress='processing' WHERE id=$1", expression.Id)
	if err != nil {
		return models.Expression{}, err
	}

	h, err := ex.RowsAffected()
	if h == 0 {
		return models.Expression{}, sql.ErrNoRows
	}
	if err != nil {
		return models.Expression{}, err
	}

	ex, err = d.db.Exec("INSERT INTO workers_and_expressions (workerId, expressionId) VALUES ($1, $2)", workerId, expression.Id)
	if err != nil {
		return models.Expression{}, err
	}

	h, err = ex.RowsAffected()
	if h == 0 {
		return models.Expression{}, sql.ErrNoRows
	}
	if err != nil {
		return models.Expression{}, err
	}

	return expression, nil
}

func (d *database) AllAliveWorkers() ([]models.Worker, error) {
	var workers []models.Worker

	rows, err := d.db.Query("SELECT id, name, isAlive, extract(epoch from lastHeartbeat)::INT FROM workers WHERE isAlive=true")
	if err != nil {
		return []models.Worker{}, err
	}

	for rows.Next() {
		var worker models.Worker
		err = rows.Scan(&worker.Id, &worker.Name, &worker.IsAlive, &worker.LastHeartbeat)
		if err != nil {
			return []models.Worker{}, err
		}
		workers = append(workers, worker)
	}
	return workers, nil
}

// finding out if the worker with given id is alive
func (d *database) IsWorkerAlive(workerId int64) (bool, error) {
	var isAlive bool

	err := d.db.QueryRow("SELECT isAlive FROM workers WHERE id=$1", workerId).Scan(&isAlive)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return isAlive, nil
}

// setting isAlive field to true
func (d *database) WakeUp(workerId int64) error {
	ex, err := d.db.Exec("UPDATE workers SET isAlive=true, lastHeartbeat=now() WHERE id=$1", workerId)
	if err != nil {
		return err
	}

	h, err := ex.RowsAffected()
	if h == 0 {
		return sql.ErrNoRows
	}
	if err != nil {
		return err
	}

	return nil
}

func (d *database) GetWorkerIdByName(name string) (int64, error) {
	var id int64

	err := d.db.QueryRow("SELECT id FROM workers WHERE name=$1", name).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (d *database) NewWorker(name string) (int64, error) {
	var id int64
	err := d.db.QueryRow("INSERT INTO workers (name, isAlive) VALUES ($1, true) RETURNING id", name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (d *database) FallAsleep(workerId int64) error {
	ex, err := d.db.Exec("UPDATE workers SET isAlive=false WHERE id=$1", workerId)
	if err != nil {
		return err
	}

	h, err := ex.RowsAffected()
	if h == 0 {
		return sql.ErrNoRows
	}
	if err != nil {
		return err
	}

	return nil
}

func (d *database) GetActiveExpressionsFromWorker(workerId int64) ([]models.Expression, error) {
	var expressions []models.Expression

	rows, err := d.db.Query("SELECT expressions.id, expressions.vanilla, expressions.answer, expressions.progress, extract(epoch from expressions.incomingDate)::INT FROM expressions JOIN workers_and_expressions ON expressions.id=workers_and_expressions.expressionId WHERE workers_and_expressions.workerId=$1 AND expressions.progress='processing'", workerId)
	if err != nil {
		return []models.Expression{}, err
	}

	for rows.Next() {
		var expression models.Expression
		err = rows.Scan(&expression.Id, &expression.Vanilla, &expression.Answer, &expression.Progress, &expression.IncomingDate)
		if err != nil {
			return []models.Expression{}, err
		}
		expressions = append(expressions, expression)
	}
	return expressions, nil
}

func (d *database) MakeExpressionAvailableAgain(expressionId int64) error {
	ex, err := d.db.Exec("UPDATE expressions SET progress='waiting' WHERE id=$1", expressionId)
	if err != nil {
		return err
	}

	h, err := ex.RowsAffected()
	if h == 0 {
		return sql.ErrNoRows
	}
	if err != nil {
		return err
	}

	return nil
}

func (d *database) SolveExpression(workerId, expressionId int64, solution string) error {
	ex, err := d.db.Exec("UPDATE expressions SET progress='done', answer=$1 FROM workers_and_expressions WHERE expressions.id=$2 AND expressions.progress='processing' AND workers_and_expressions.workerId=$3", solution, expressionId, workerId)
	if err != nil {
		return err
	}

	h, err := ex.RowsAffected()
	if h == 0 {
		return sql.ErrNoRows
	}
	if err != nil {
		return err
	}

	return nil
}

// just taking all expressions (untaken first)
func (d *database) AllExpressions() ([]models.ExpressionGeneral, error) {
	var expressions []models.ExpressionGeneral

	rows, err := d.db.Query("SELECT expressions.id, extract(epoch from expressions.incomingDate)::INT, expressions.vanilla, expressions.answer, expressions.progress, workers.name FROM expressions LEFT JOIN workers_and_expressions ON workers_and_expressions.expressionId=expressions.id LEFT JOIN workers ON workers_and_expressions.workerId=workers.id ORDER BY expressions.progress DESC")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var expression models.ExpressionGeneral

		err = rows.Scan(&expression.Id, &expression.IncomingDate, &expression.Vanilla, &expression.Answer, &expression.Progress, &expression.WorkerName)
		if err != nil {
			return nil, err
		}

		expressions = append(expressions, expression)
	}

	return expressions, nil
}

func (d *database) GetExpressionById(expressionId int64) (models.Expression, error) {
	var expression models.Expression
	err := d.db.QueryRow("SELECT id, extract(epoch from incomingDate)::INT, vanilla, answer, progress FROM expressions WHERE id=$1", expressionId).Scan(&expression.Id, &expression.IncomingDate, &expression.Vanilla, &expression.Answer, &expression.Progress)
	if err != nil {
		return models.Expression{}, err
	}
	return expression, nil
}

func (d *database) AllWorkers() ([]models.Worker, error) {
	var workers []models.Worker

	rows, err := d.db.Query("SELECT id, name, isAlive, extract(epoch from lastHeartbeat)::INT FROM workers")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var worker models.Worker

		err = rows.Scan(&worker.Id, &worker.Name, &worker.IsAlive, &worker.LastHeartbeat)
		if err != nil {
			return nil, err
		}

		workers = append(workers, worker)
	}

	return workers, nil
}
