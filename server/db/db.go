package db

type DB interface {
	CleanupCollectorStates(collectors []string) (modified []StateID, err error)
	PutCollectorStates(collector string, states []State) (modified []StateID, err error)

	ListStates(registry string) ([]State, error)
	GetState(registry string, id string) (State, error)

	Shutdown()
}
