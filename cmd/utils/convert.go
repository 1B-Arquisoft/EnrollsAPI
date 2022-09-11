package u

import "encoding/json"

func StructToMap(obj interface{}) (newMap map[string]interface{}) {
	data, err := json.Marshal(obj) // Convert to a json string

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &newMap) // Convert to a map
	return
}
