package models

import "database/sql"

type Error struct {
	ErrorMessage string `json:"errorMessage"`
}

type ExpressionAdding struct {
	Expression string `json:"expression" form:"expression" binding:"required"`
}

type ExpressionSolving struct {
	Id     int64  `json:"id" binding:"required"`
	Answer string `json:"answer" binding:"required"`
}

type Expression struct {
	Id           int64  `json:"id"`
	IncomingDate int64  `json:"incomingDate"`
	Vanilla      string `json:"vanilla"`
	Answer       string `json:"answer"`
	Progress     string `json:"progress"`
}

type ExpressionGeneral struct {
	Id           int64          `json:"id"`
	IncomingDate int64          `json:"incomingDate"`
	Vanilla      string         `json:"vanilla"`
	Answer       string         `json:"answer"`
	Progress     string         `json:"progress"`
	WorkerName   sql.NullString `json:"workerName"`
}

type WorkerAdding struct {
	Name string `json:"name" binding:"required"`
}

type Worker struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	IsAlive       bool   `json:"isAlive"`
	LastHeartbeat int64  `json:"lastHeartbeat"`
}
