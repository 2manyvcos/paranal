package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/2manyvcos/paranal/server/helper"
	"github.com/2manyvcos/paranal/server/schema"
	"github.com/2manyvcos/paranal/utils"
)

type MaintenanceTask struct {
	Name    string     `json:"name"`
	LastRun *time.Time `json:"lastRun"`
	Running bool       `json:"running"`
	NextRun *time.Time `json:"nextRun"`
}

func GetMaintenanceTasks(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	var runningQuery *bool
	q := req.URL.Query()
	if v, ok := utils.LoadQueryBool(q, "running"); ok {
		runningQuery = &v
	}

	taskStates := app.ListMaintenanceTaskStates()

	responsePayload := make([]MaintenanceTask, 0, len(taskStates))
	for _, taskState := range taskStates {
		if runningQuery != nil && taskState.Running != *runningQuery {
			continue
		}
		responsePayload = append(responsePayload, MaintenanceTask{
			Name:    taskState.Name,
			LastRun: taskState.LastRun,
			Running: taskState.Running,
			NextRun: taskState.NextRun,
		})
	}
	res.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(res).Encode(responsePayload)
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}

func GetMaintenanceTasksByName(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	taskName := req.PathValue("taskName")
	if taskName == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	taskState := app.GetMaintenanceTaskState(taskName)
	if taskState == nil {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(res).Encode(MaintenanceTask{
		Name:    taskState.Name,
		LastRun: taskState.LastRun,
		Running: taskState.Running,
		NextRun: taskState.NextRun,
	})
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}

func PostMaintenanceTasksByNameRun(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	taskName := req.PathValue("taskName")
	if taskName == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := app.RunMaintenanceTask(taskName)
	if errors.Is(err, schema.ErrNotFound) {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if errors.Is(err, schema.ErrAlreadyRunning) {
		http.Error(res, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusAccepted)
}
