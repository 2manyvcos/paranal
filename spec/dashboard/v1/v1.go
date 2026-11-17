package dashboard_spec_v1

import dashboard_spec "github.com/2manyvcos/paranal/spec/dashboard"

const Version = "v1.0"

var Parser dashboard_spec.DashboardSpecParser = dashboardParser{}

type dashboardParser struct{}
