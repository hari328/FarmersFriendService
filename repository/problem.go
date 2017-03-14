package repository

import (
	"time"
)

type Problem struct {
	Id int								`json:"problemId"`
	FarmerId int					`json:"farmerId"`
	ProblemDesc string		`json:"problemDesc"`
	PostedDate time.Time	`json:"postedDate"`
	IsSolved bool					`json:"isSolved"`
}