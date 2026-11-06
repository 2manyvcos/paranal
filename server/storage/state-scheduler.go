package storage

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/2manyvcos/paranal/server/db"
	"github.com/2manyvcos/paranal/server/events"
	"github.com/2manyvcos/paranal/utils"
)

type stateInstruction struct {
	Collector string
	States    []State
	Done      chan<- error
}

type stateScheduler struct {
	db          db.DB
	eventBroker *events.EventBroker
	scheduler   *utils.JobScheduler[stateInstruction]
}

func (s *stateScheduler) PutCollectorStates(collector string, states []State) error {
	done := make(chan error)
	s.scheduler.Schedule(stateInstruction{
		Collector: collector,
		States:    states,
		Done:      done,
	})

	return <-done
}

func (s *stateScheduler) Shutdown() {
	s.scheduler.ShutdownSync()
}

func stateDedupeKey(s State) (key struct{ Registry, ID string }) {
	key.Registry = s.Registry
	key.ID = s.ID
	return
}

type stateRev struct {
	Collector string `json:"collector"`
	Data      any    `json:"data"`
}

func (s *stateScheduler) handleInstruction(collector string, states []State) error {
	dedupedStates := utils.DedupeEnd(states, stateDedupeKey)
	records := make([]db.State, len(dedupedStates))
	for i, state := range dedupedStates {
		// TODO: If rev is also used outside (for E-Tags), this can be kept here. Otherwise, this logic should be moved into the db package.
		rev, err := utils.HashJSONValue(stateRev{
			Collector: collector,
			Data:      state.Data,
		})
		if err != nil {
			return fmt.Errorf("deriving state revision: %w", err)
		}

		data, err := json.Marshal(state.Data)
		if err != nil {
			return fmt.Errorf("serializing state data: %w", err)
		}

		records[i] = db.State{
			Registry: state.Registry,
			ID:       state.ID,
			Rev:      rev,
			Data:     data,
		}
	}

	modifiedStates, err := s.db.PutCollectorStates(collector, records)
	if err != nil {
		return err
	}

	eventBundle := make([]events.Event, len(modifiedStates))
	for i, modifiedState := range modifiedStates {
		eventBundle[i] = events.NewChangeEvent(fmt.Sprintf("v1/registries/%s/%s", url.PathEscape(modifiedState.Registry), url.PathEscape(modifiedState.ID)))
	}
	s.eventBroker.Publish(eventBundle...)

	return nil
}

func stateInstructionDedupeKey(i stateInstruction) string {
	return i.Collector
}

func (s *stateScheduler) handle(instructions []stateInstruction) {
	for _, instruction := range utils.DedupeEnd(instructions, stateInstructionDedupeKey) {
		err := s.handleInstruction(instruction.Collector, instruction.States)
		if err != nil {
			instruction.Done <- err
		}

		close(instruction.Done)
	}
}

func newStateScheduler(db db.DB, eventBroker *events.EventBroker) *stateScheduler {
	var s stateScheduler
	s.db = db
	s.eventBroker = eventBroker
	s.scheduler = utils.NewJobScheduler(s.handle)

	s.scheduler.Start()

	return &s
}
