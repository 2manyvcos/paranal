package storage

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/2manyvcos/paranal/server/db"
	"github.com/2manyvcos/paranal/server/events"
	collector_spec "github.com/2manyvcos/paranal/server/spec/collector"
	"github.com/2manyvcos/paranal/utils"
)

func (s *Storage) CleanupCollectorStates(collectors []string) error {
	s.stateLock.Lock()
	defer s.stateLock.Unlock()

	_, err := s.db.CleanupCollectorStates(collectors)
	return err
}

func stateDedupeKey(s collector_spec.State) (key struct{ Registry, ID string }) {
	key.Registry = s.Registry
	key.ID = s.ID
	return
}

type stateRevPayload struct {
	Registry  string `json:"registry"`
	ID        string `json:"id"`
	Collector string `json:"collector"`
	Data      any    `json:"data"`
}

func (s *Storage) PutCollectorStates(collector string, states []collector_spec.State) error {
	s.stateLock.Lock()
	defer s.stateLock.Unlock()

	dedupedStates := utils.DedupeEnd(states, stateDedupeKey)
	records := make([]db.State, len(dedupedStates))
	for i, state := range dedupedStates {
		rev, err := utils.HashJSONValue(stateRevPayload{
			Registry:  state.Registry,
			ID:        state.ID,
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

func (s *Storage) ListStates(registry string) ([]collector_spec.State, error) {
	records, err := s.db.ListStates(registry)
	if err != nil {
		return nil, err
	}

	states := make([]collector_spec.State, len(records))
	for i, record := range records {
		var state collector_spec.State
		state.Registry = record.Registry
		state.ID = record.ID

		err := json.Unmarshal(record.Data, &state.Data)
		if err != nil {
			return nil, fmt.Errorf("parsing state data: %w", err)
		}

		states[i] = state
	}

	return states, nil
}

func (s *Storage) GetState(registry string, id string) (collector_spec.State, error) {
	record, err := s.db.GetState(registry, id)
	if err != nil {
		return collector_spec.State{}, err
	}

	var state collector_spec.State
	state.Registry = record.Registry
	state.ID = record.ID

	err = json.Unmarshal(record.Data, &state.Data)
	if err != nil {
		return collector_spec.State{}, fmt.Errorf("parsing state data: %w", err)
	}

	return state, nil
}
