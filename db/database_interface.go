package db

import (
	"database/sql"

	"github.com/askarkasimov/yg-colonel/models"
	_ "github.com/lib/pq"
)

type iDatabase interface {
	GetOneAvailableExpression() (models.Expression, error)
	AddExpression(e models.ExpressionAdding) (int64, error)
	AllExpressions() ([]models.Expression, error)
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

// adding expression and giving its id back
func (d *database) AddExpression(e models.ExpressionAdding) (int64, error) {
	var id int64 = 0
	err := d.db.QueryRow("INSERT INTO expressions (vanilla) VALUES ($1) RETURNING id", e.Expression).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// taking the oldest untaken expression
func (d *database) GetOneAvailableExpression() (models.Expression, error) {
	var expression models.Expression

	err := d.db.QueryRow(
		"SELECT id, extract(epoch from incomingDate)::INT, vanilla, answer, progress FROM expressions WHERE progress='waiting' ORDER BY incomingDate LIMIT 1").Scan(
		&expression.Id, &expression.IncomingDate, &expression.Vanilla, &expression.Answer, &expression.Progress)
	if err != nil {
		return models.Expression{}, err
	}

	return expression, nil
}

// just taking all expressions (untaken first)
func (d *database) AllExpressions() ([]models.Expression, error) {
	var expressions []models.Expression

	rows, err := d.db.Query("SELECT id, extract(epoch from incomingDate)::INT, vanilla, answer, progress FROM expressions ORDER BY progress DESC")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var expression models.Expression

		err = rows.Scan(&expression.Id, &expression.IncomingDate, &expression.Vanilla, &expression.Answer, &expression.Progress)
		if err != nil {
			return nil, err
		}

		expressions = append(expressions, expression)
	}

	return expressions, nil
}
