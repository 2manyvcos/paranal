package schema

import "fmt"

type ServiceUptimeStatusQuery struct {
	Name      *string
	ServiceID *string
}

type ServiceUptimeStatus struct {
	Name      string
	ServiceID string
	Status    int
}

const (
	ServiceUptimeStatusStatusDown = iota + 1
	ServiceUptimeStatusStatusUp
)

var (
	ServiceUptimeStatusStatusNames = map[int]string{
		ServiceUptimeStatusStatusDown: "down",
		ServiceUptimeStatusStatusUp:   "up",
	}
	ServiceUptimeStatusStatusCodes map[string]int
)

func init() {
	ServiceUptimeStatusStatusCodes = make(map[string]int, len(ServiceUptimeStatusStatusNames))
	for code, name := range ServiceUptimeStatusStatusNames {
		ServiceUptimeStatusStatusCodes[name] = code
	}
}

func (r ServiceUptimeStatus) Valid() error {
	if r.ServiceID == "" {
		return fmt.Errorf("invalid service ID")
	}
	if _, ok := ServiceUptimeStatusStatusNames[r.Status]; !ok {
		return fmt.Errorf("invalid status")
	}
	return nil
}
