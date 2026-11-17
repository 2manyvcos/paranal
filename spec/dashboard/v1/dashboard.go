package dashboard_spec_v1

import (
	"fmt"
	"net/url"

	dashboard_spec "github.com/2manyvcos/paranal/spec/dashboard"
	"github.com/go-viper/mapstructure/v2"
)

type dashboard struct {
	ID            string            `mapstructure:"id"`
	Routes        []route           `mapstructure:"routes"`
	Theme         map[string]any    `mapstructure:"theme"`
	Vars          map[string]_var   `mapstructure:"vars"`
	DefaultLayout map[string]any    `mapstructure:"defaultLayout"`
	Layout        map[string]any    `mapstructure:"layout"`
	Content       []map[string]any  `mapstructure:"content"`
	Scripts       map[string]string `mapstructure:"scripts"`
}

type route struct {
	Host string `mapstructure:"host"`
	Path string `mapstructure:"path"`
}

type _var struct {
	Value   any  `mapstructure:"value"`
	Exposed bool `mapstructure:"exposed"`
}

func (dashboardParser) ParseDashboard(provider string, input any) (dashboard_spec.Dashboard, error) {
	var d dashboard
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: &d})
	if err != nil {
		return dashboard_spec.Dashboard{}, err
	}
	err = decoder.Decode(input)
	if err != nil {
		return dashboard_spec.Dashboard{}, err
	}

	var result dashboard_spec.Dashboard
	result.Version = Version

	if d.ID == "" {
		return dashboard_spec.Dashboard{}, fmt.Errorf("missing \"id\"")
	}
	result.ID = fmt.Sprintf("%v@%v", url.QueryEscape(d.ID), url.QueryEscape(provider))

	if len(d.Routes) == 0 {
		return dashboard_spec.Dashboard{}, fmt.Errorf("missing \"routes\"")
	}
	result.Routes = make([]dashboard_spec.Route, len(d.Routes))
	for i, route := range d.Routes {
		result.Routes[i], err = parseRoute(route)
		if err != nil {
			return dashboard_spec.Dashboard{}, fmt.Errorf("parsing \"routes\": %w", err)
		}
	}

	result.Theme = d.Theme

	result.Vars = make(map[string]dashboard_spec.Var, len(d.Vars))
	for i, _var := range d.Vars {
		result.Vars[i] = dashboard_spec.Var{
			Value:   _var.Value,
			Exposed: _var.Exposed,
		}
	}

	result.DefaultLayout = d.DefaultLayout
	result.Layout = d.Layout
	result.Content = d.Content
	result.Scripts = d.Scripts

	return result, nil
}

func parseRoute(r route) (dashboard_spec.Route, error) {
	var result dashboard_spec.Route

	result.Host = r.Host

	if r.Path == "" {
		return dashboard_spec.Route{}, fmt.Errorf("missing \"path\"")
	}
	result.Path = r.Path

	return result, nil
}
