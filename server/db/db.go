package db

type DB interface {
	PutCollectorStates(collector string, states []State) (modified []StateID, err error)

	Shutdown()
}
