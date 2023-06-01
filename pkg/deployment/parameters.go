package deployment

func getParams(input map[string]interface{}) map[string]interface{} {
	var paramsValue map[string]interface{}
	if input != nil {
		if p, ok := input["parameters"]; ok {
			paramsValue = p.(map[string]interface{})
		} else {
			paramsValue = input
		}
	}

	return ensureValueOnParams(paramsValue)
}

func ensureValueOnParams(input map[string]interface{}) map[string]interface{} {
	outputMap := make(map[string]interface{})
	if input == nil {
		return outputMap
	}
	for k, v := range input {
		if v != nil {
			_, isMap := v.(map[string]interface{})
			if isMap {
				_, hasValueKey := v.(map[string]interface{})["value"]
				if !hasValueKey {
					outputMap[k] = map[string]interface{}{"value": v}
				} else {
					outputMap[k] = v
				}
			} else {
				outputMap[k] = map[string]interface{}{"value": v}
			}
		}
	}
	return outputMap
}
