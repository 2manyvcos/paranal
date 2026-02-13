package state

import (
	"fmt"
	"time"

	"github.com/go-viper/mapstructure/v2"
)

type ByUptimeStatusOrder []ServiceUptimeStatus

func (a ByUptimeStatusOrder) Len() int      { return len(a) }
func (a ByUptimeStatusOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByUptimeStatusOrder) Less(i, j int) bool {
	return a[i].Order < a[j].Order
}

type ByVersionOrder []ServiceVersion

func (a ByVersionOrder) Len() int           { return len(a) }
func (a ByVersionOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByVersionOrder) Less(i, j int) bool { return a[i].Order < a[j].Order }

type ByActionGroupOrder []ServiceActionGroup

func (a ByActionGroupOrder) Len() int           { return len(a) }
func (a ByActionGroupOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByActionGroupOrder) Less(i, j int) bool { return a[i].Order < a[j].Order }

type ByActionOrder []ServiceAction

func (a ByActionOrder) Len() int           { return len(a) }
func (a ByActionOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByActionOrder) Less(i, j int) bool { return a[i].Order < a[j].Order }

func decodeServiceInstruction(input, output any) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: output, Squash: true, IgnoreUntaggedFields: true})
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}

func decodeTime(input any) (time.Time, error) {
	if input == nil {
		return time.Now(), nil
	}
	switch v := input.(type) {
	case string:
		if t, err := time.Parse(time.RFC3339Nano, v); err == nil {
			return t, nil
		} else if t, err = time.Parse(time.RFC3339, v); err == nil {
			return t, nil
		} else if t, err = time.Parse(time.RFC1123Z, v); err == nil {
			return t, nil
		} else if t, err = time.Parse(time.RFC1123, v); err == nil {
			return t, nil
		} else if t, err = time.Parse(time.RFC850, v); err == nil {
			return t, nil
		} else if t, err = time.Parse(time.RFC822Z, v); err == nil {
			return t, nil
		} else if t, err = time.Parse(time.RFC822, v); err == nil {
			return t, nil
		} else if t, err = time.Parse(time.DateTime, v); err == nil {
			return t, nil
		} else if t, err = time.Parse(time.DateOnly, v); err == nil {
			return t, nil
		} else {
			return time.Time{}, fmt.Errorf("invalid time")
		}

	case float64:
		return time.Unix(int64(v), 0), nil

	default:
		return time.Time{}, fmt.Errorf("invalid time")
	}
}
