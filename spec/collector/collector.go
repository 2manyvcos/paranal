package collector_spec

type Collector struct {
	Version  string
	ID       string
	Triggers []Trigger
	Source   string
}

type Trigger struct {
	Trigger     string
	Crontab     string
	PathPattern string
}

type Fact interface {
	Kind() string
}

var _ Fact = (*State)(nil)

type State struct {
	Registry string
	ID       string
	Data     map[string]any
}

func (State) Kind() string {
	return "state"
}

type CollectorSpecParser interface {
	ParseCollector(provider string, input any) (Collector, error)
	ParseFact(input any) (Fact, error)
}
