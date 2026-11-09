package db

type DB interface {
	CleanupCollectorStates(collectors []string) (modified []StateID, err error)
	PutCollectorStates(collector string, states []State) (modified []StateID, err error)

	Shutdown()
}
