package state

import (
	"log"
	"maps"
	"slices"

	"github.com/2manyvcos/paranal/server/schema"
)

func (s *State) ListMaintenanceTaskStates() []schema.MaintenanceTaskState {
	result := make([]schema.MaintenanceTaskState, 0, len(s.maintenanceTasks))
	taskNames := slices.Collect(maps.Keys(s.maintenanceTasks))
	slices.Sort(taskNames)
	for _, taskName := range taskNames {
		task := s.maintenanceTasks[taskName]
		state := schema.MaintenanceTaskState{
			Name: taskName,
		}
		state.LastRun, state.Running, state.NextRun = task.State()
		result = append(result, state)
	}
	return result
}

func (s *State) GetMaintenanceTaskState(taskName string) *schema.MaintenanceTaskState {
	task, ok := s.maintenanceTasks[taskName]
	if !ok {
		return nil
	}
	state := schema.MaintenanceTaskState{
		Name: taskName,
	}
	state.LastRun, state.Running, state.NextRun = task.State()
	return &state
}

func (s *State) RunMaintenanceTask(taskName string) error {
	task, ok := s.maintenanceTasks[taskName]
	if !ok {
		return schema.ErrNotFound
	}
	if _, running, _ := task.State(); running {
		return schema.ErrAlreadyRunning
	}
	err := task.Run()
	if err != nil {
		log.Printf("Error running maintenance task \"%s\" - %s\n", taskName, err)
	}
	return err
}
