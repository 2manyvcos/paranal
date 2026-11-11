package dashboard_spec

type Dashboard struct {
	Version       string
	ID            string
	Routes        []Route
	Theme         map[string]any
	Vars          map[string]Var
	DefaultLayout map[string]any
	Layout        map[string]any
	Content       []map[string]any
	Scripts       map[string]string
}

type Route struct {
	Host string
	Path string
}

type Var struct {
	Value   any
	Exposed bool
}

type DashboardSpecParser interface {
	ParseDashboard(provider string, input any) (Dashboard, error)
}
