package models

type Error struct {
	ErrorMessage string `json:"errorMessage"`
}

type ExpressionAdding struct {
	Expression string `json:"expression" binding:"required"`
}

type Expression struct {
	Id           int64  `json:"id"`
	IncomingDate int64  `json:"incomingDate"`
	Vanilla      string `json:"vanilla"`
	Answer       string `json:"answer"`
	Progress     string `json:"progress"`
}
