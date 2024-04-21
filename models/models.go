package models

import (
	"database/sql"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Error struct {
	ErrorMessage string `json:"errorMessage"`
}

type ExpressionAdding struct {
	Expression string `json:"expression" form:"expression" binding:"required"`
}

type ExpressionSolving struct {
	Id     uuid.UUID `json:"id" binding:"required"`
	Answer string    `json:"answer" binding:"required"`
}

type Expression struct {
	Id           uuid.UUID `json:"id"`
	IncomingDate int64     `json:"incomingDate"`
	Vanilla      string    `json:"vanilla"`
	Answer       string    `json:"answer"`
	Progress     string    `json:"progress"`
}

type ExpressionGeneral struct {
	Id           uuid.UUID      `json:"id"`
	IncomingDate int64          `json:"incomingDate"`
	Vanilla      string         `json:"vanilla"`
	Answer       string         `json:"answer"`
	Progress     string         `json:"progress"`
	WorkerName   sql.NullString `json:"workerName"`
}

type WorkerAdding struct {
	Name               string `json:"name" binding:"required"`
	NumberOfGoroutines int    `json:"number_of_goroutines"`
}

type Worker struct {
	Id                 uuid.UUID `json:"id"`
	Name               string    `json:"name"`
	NumberOfGoroutines int       `json:"number_of_goroutines"`
	IsAlive            bool      `json:"isAlive"`
	LastHeartbeat      int64     `json:"lastHeartbeat"`
}

type JWT struct {
	UserId uuid.UUID `json:"id"`
	jwt.RegisteredClaims
}

type User struct {
	Id       uuid.UUID `json:"id"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
}
