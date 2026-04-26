package request

import (
	"encoding/json"
	"fmt"
)

type AboutUpdate struct {
	Values map[string]string `json:"values"`
}

func (a *AboutUpdate) UnmarshalJSON(data []byte) error {
	var wrapped struct {
		Values map[string]any `json:"values"`
	}
	if err := json.Unmarshal(data, &wrapped); err != nil {
		return err
	}
	if wrapped.Values != nil {
		a.Values = toAboutValues(wrapped.Values)
		return nil
	}

	var direct map[string]any
	if err := json.Unmarshal(data, &direct); err != nil {
		return err
	}
	a.Values = toAboutValues(direct)
	return nil
}

func toAboutValues(values map[string]any) map[string]string {
	result := make(map[string]string, len(values))
	for key, value := range values {
		if value == nil {
			result[key] = ""
			continue
		}
		result[key] = fmt.Sprint(value)
	}
	return result
}
