package shared

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GetFloatFromMap(data map[string]interface{}, fieldName string) (float64, bool) {
	if data[fieldName] == nil {
		return 0, false
	}

	var field float64
	var err error

	if item, ok := data[fieldName].(string); ok {
		field, err = strconv.ParseFloat(item, 64)
	} else if item, ok := data[fieldName].(float64); ok {
		field = item
	} else {
		field, err = strconv.ParseFloat(fmt.Sprintf("%v", data[fieldName]), 64)
	}

	return field, err == nil
}

func GetUintFromMap(data map[string]interface{}, fieldName string) (uint64, bool) {
	if data[fieldName] == nil {
		return 0, false
	}

	var field uint64
	var err error

	if item, ok := data[fieldName].(string); ok {
		field, err = strconv.ParseUint(item, 10, 64)
	} else if item, ok := data[fieldName].(float64); ok {
		field = uint64(item)
	} else {
		field, err = strconv.ParseUint(fmt.Sprintf("%v", data[fieldName]), 10, 64)
	}

	return field, err == nil
}

func GetStringFromMap(data map[string]interface{}, fieldName string) (string, bool) {
	if data[fieldName] == nil {
		return "", false
	}

	field := fmt.Sprintf("%v", data[fieldName])
	return field, true
}

func GetBooleanFromMap(data map[string]interface{}, fieldName string) (bool, bool) {
	if data[fieldName] == nil {
		return false, false
	}

	field, ok := data[fieldName].(bool)
	if !ok {
		field = data[fieldName] == "true"
	}

	return field, ok
}

func GetTimeRFC3339NanoFromMap(data map[string]interface{}, fieldName string) (time.Time, error) {
	if data[fieldName] == nil {
		return time.Time{}, fmt.Errorf("field %v is not exist", fieldName)
	}

	field := fmt.Sprintf("%v", data[fieldName])
	if strings.Contains(field, "Z") {
		return time.Parse(time.RFC3339, field)
	}

	if len(field) >= 25 && strings.Contains(field, "+0") {
		return time.Parse(time.RFC3339Nano, field)
	}

	if len(field) == 10 && !strings.Contains("T", field) {
		field = fmt.Sprintf("%sT00:00:00", field)
	} else if len(field) > 10 && len(field) <= 16 {
		field = fmt.Sprintf("%s:00", field)
	}

	if !strings.Contains(field, "+0") {
		field = fmt.Sprintf("%s+07:00", field)
	}

	tm, err := time.Parse(time.RFC3339Nano, field)

	if err != nil {
		return time.Time{}, err
	}

	return tm, nil
}

func GetMapFromMap(data map[string]interface{}, fieldName string) (map[string]interface{}, bool) {
	if data[fieldName] == nil {
		return nil, false
	}

	field, ok := data[fieldName].(map[string]interface{})
	return field, ok
}

func GetArrayMapFromMap(data map[string]interface{}, fieldName string) (field []map[string]interface{}, ok bool) {
	if data[fieldName] == nil {
		return nil, false
	}

	field, ok = data[fieldName].([]map[string]interface{})

	if !ok {
		field = nil
		nField, ok := data[fieldName].([]interface{})

		if ok {
			for _, item := range nField {
				nItem, ok := item.(map[string]interface{})

				if !ok {
					return field, ok
				}

				field = append(field, nItem)
			}

			return field, ok
		}
	}

	return field, ok
}
