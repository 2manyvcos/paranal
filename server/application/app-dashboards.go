package application

import (
	"maps"
	"slices"
	"strings"

	dashboard_spec "github.com/2manyvcos/paranal/spec/dashboard"
)

func (a *App) ListDashboards() []dashboard_spec.Dashboard {
	a.lock.Lock()
	defer a.lock.Unlock()

	dashboards := slices.Collect(maps.Values(a.dashboards))
	slices.SortFunc(dashboards, func(a, b dashboard_spec.Dashboard) int {
		return strings.Compare(a.ID, b.ID)
	})

	return dashboards
}

func (a *App) GetDashboard(id string) *dashboard_spec.Dashboard {
	a.lock.Lock()
	defer a.lock.Unlock()

	dashboard, ok := a.dashboards[id]
	if !ok {
		return nil
	}

	return &dashboard
}
