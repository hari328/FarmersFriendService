package api

import (
	"net/http"
	"github.com/FarmersFriendService/model"
	"encoding/json"
)

func (api *Api) ListProblems(w http.ResponseWriter, r *http.Request) {
	rows, err := api.Db.Query("SELECT * FROM problems")
	if err != nil {
		panic(err)
	}

	problems := make([]model.Problem, 0)

	var problem model.Problem
	for rows.Next() {
		err = rows.Scan(&problem.Id, &problem.FarmerId, &problem.ProblemDesc, &problem.PostedDate, &problem.IsSolved)
		if err != nil {
			panic(err)
		}

		problems = append(problems, problem)
	}

	rows.Close()

	farmerDetails,err := json.Marshal(problems)
	if err != nil {
		panic(err)
	}

	w.Write([]byte(farmerDetails))
}
