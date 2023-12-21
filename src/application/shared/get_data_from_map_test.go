package shared

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/**
 *	Best Case Scenario Get Float Field from Go Map
 *
 *	@author Mufid Jamaluddin
 **/
func TestGetFloatFromMap_BestCase(t *testing.T) {
	testCases := []map[string]interface{}{
		{
			"id": float64(9),
		},
		{
			"id": float64(8),
		},
		{
			"id": int64(14),
		},
		{
			"id": 17,
		},
		{
			"id": "30",
		},
		{
			"id": "3.56",
		},
	}

	expecteds := []float64{
		9,
		8,
		14,
		17,
		30,
		3.56,
	}

	for idx, testCase := range testCases {
		expected := expecteds[idx]

		num, ok := GetFloatFromMap(testCase, "id")

		assert.True(t, ok, "Must be true")
		assert.Equal(t, expected, num, "Must be same")
	}
}

/**
 *	Worst Case Scenario Get Float Field from Go Map
 *
 *	@author Mufid Jamaluddin
 **/
func TestGetFloatFromMap_WorstCase(t *testing.T) {
	testCases := []map[string]interface{}{
		{
			"id": "9k",
		},
		{
			"id": ",8",
		},
		{
			"name": "Mufid",
		},
		{},
		{
			"age": 24,
		},
	}

	expected := float64(0)

	for _, testCase := range testCases {
		num, ok := GetFloatFromMap(testCase, "id")

		assert.False(t, ok, "Must be false")
		assert.Equal(t, expected, num, "Must be zero")
	}
}

/**
 *	Best Case Scenario Get Uint Field from Go Map
 *
 *	@author Mufid Jamaluddin
 **/
func TestGetUintFromMap_BestCase(t *testing.T) {
	testCases := []map[string]interface{}{
		{
			"id": float64(9),
		},
		{
			"id": float64(8),
		},
		{
			"id": uint64(14),
		},
		{
			"id": 17,
		},
		{
			"id": "30",
		},
		{
			"id": "6",
		},
	}

	expecteds := []uint64{
		9,
		8,
		14,
		17,
		30,
		6,
	}

	for idx, testCase := range testCases {
		expected := expecteds[idx]

		num, ok := GetUintFromMap(testCase, "id")

		assert.True(t, ok, "Must be true")
		assert.Equal(t, expected, num, "Must be same")
	}
}

/**
 *	Worst Case Scenario Get Uint Field from Go Map
 *
 *	@author Mufid Jamaluddin
 **/
func TestGetUintFromMap_WorstCase(t *testing.T) {
	testCases := []map[string]interface{}{
		{
			"id": "9k",
		},
		{
			"id": ",8",
		},
		{
			"name": "Mufid",
		},
		{},
		{
			"age": 24,
		},
		{
			"id": "24.89",
		},
	}

	expected := uint64(0)

	for _, testCase := range testCases {
		num, ok := GetUintFromMap(testCase, "id")

		assert.False(t, ok, "Must be false")
		assert.Equal(t, expected, num, "Must be zero")
	}
}

/**
 *	Best Case Scenario Get String Field from Go Map
 *
 *	@author Mufid Jamaluddin
 **/
func TestGetStringFromMap_BestCase(t *testing.T) {
	testCases := []map[string]interface{}{
		{
			"id": float64(9),
		},
		{
			"id": float64(8),
		},
		{
			"id": uint64(14),
		},
		{
			"id": 17,
		},
		{
			"id": "30",
		},
		{
			"id": "6",
		},
		{
			"id": "IKA0009",
		},
		{
			"id": "5.35",
		},
	}

	expecteds := []string{
		"9",
		"8",
		"14",
		"17",
		"30",
		"6",
		"IKA0009",
		"5.35",
	}

	for idx, testCase := range testCases {
		expected := expecteds[idx]

		num, ok := GetStringFromMap(testCase, "id")

		assert.True(t, ok, "Must be true")
		assert.Equal(t, expected, num, "Must be same")
	}
}

/**
 *	Worst Case Scenario Get String Field from Go Map
 *
 *	@author Mufid Jamaluddin
 **/
func TestGetStringFromMap_WorstCase(t *testing.T) {
	testCases := []map[string]interface{}{
		{
			"name": "Mufid",
		},
		{},
		{
			"age": 24,
		},
	}

	expected := ""

	for _, testCase := range testCases {
		num, ok := GetStringFromMap(testCase, "id")

		assert.False(t, ok, "Must be false")
		assert.Equal(t, expected, num, "Must be zero")
	}
}

/**
 *	Best Case Scenario Get Map[String]Interface Field from Go Map
 *
 *	@author Mufid Jamaluddin
 **/
func TestGetMapFromMap_BestCase(t *testing.T) {
	testCases := []map[string]interface{}{
		{
			"data": map[string]interface{}{
				"name": "Mufid",
			},
		},
		{
			"data": map[string]interface{}{
				"age":  "24",
				"name": "Jamaluddin",
			},
		},
		{
			"data": map[string]interface{}{
				"waksa": 989,
			},
		},
	}

	expecteds := []map[string]interface{}{
		{
			"name": "Mufid",
		},
		{
			"age":  "24",
			"name": "Jamaluddin",
		},
		{
			"waksa": 989,
		},
	}

	for idx, testCase := range testCases {
		expected := expecteds[idx]

		res, ok := GetMapFromMap(testCase, "data")

		assert.True(t, ok, "Must be true")
		for key, value := range res {
			assert.Equal(t, expected[key], value, fmt.Sprintf("Must be equal for key %s", key))
		}
	}
}

/**
 *	Worst Case Scenario Get Map[String]Interface Field from Go Map
 *
 *	@author Mufid Jamaluddin
 **/
func TestGetMapFromMap_WorstCase(t *testing.T) {
	testCases := []map[string]interface{}{
		{
			"dataa": map[string]interface{}{
				"name": "Mufid",
			},
		},
		{
			"datta": map[string]interface{}{
				"age":  "24",
				"name": "Jamaluddin",
			},
		},
		{
			"daata": map[string]interface{}{
				"waksa": 989,
			},
		},
	}

	var expected map[string]interface{} = nil

	for _, testCase := range testCases {
		res, ok := GetMapFromMap(testCase, "data")

		assert.False(t, ok, "Must be false")
		assert.Equal(t, expected, res, "Must be nil")
	}
}

/**
 *	Best Case Scenario Get []Map[String]Interface Field from Go Map
 *
 *	@author Mufid Jamaluddin
 **/
func TestGetArrayMapFromMap_BestCase1(t *testing.T) {
	testCases := []map[string]interface{}{
		{
			"data": []map[string]interface{}{
				{
					"name": "Mufid",
				},
				{
					"age":  24,
					"name": "Jamaluddin",
					"city": "Sumedang",
				},
				{
					"address": "Jl Pion Sumedang",
				},
			},
		},
		{
			"data": []map[string]interface{}{},
		},
		{
			"data": []map[string]interface{}{
				{},
				{
					"age":  24,
					"name": "Jamaluddin",
					"city": "Sumedang",
				},
				{
					"address": "Jl Pion Sumedang",
					"name":    "Monica Cecilia",
				},
			},
		},
	}

	expecteds := [][]map[string]interface{}{
		{
			{
				"name": "Mufid",
			},
			{
				"age":  24,
				"name": "Jamaluddin",
				"city": "Sumedang",
			},
			{
				"address": "Jl Pion Sumedang",
			},
		},
		{},
		{
			{},
			{
				"age":  24,
				"name": "Jamaluddin",
				"city": "Sumedang",
			},
			{
				"address": "Jl Pion Sumedang",
				"name":    "Monica Cecilia",
			},
		},
	}

	for idx, testCase := range testCases {
		expected, err := json.Marshal(expecteds[idx])
		assert.Nil(t, err)

		res, ok := GetArrayMapFromMap(testCase, "data")
		assert.True(t, ok, "Must be true")

		resJson, err := json.Marshal(res)
		assert.Nil(t, err)
		assert.JSONEq(t, string(expected), string(resJson), "Must be same")
	}
}

/**
 *	Best Case Scenario Get []Interface Field from Go Map
 *
 *	@author Mufid Jamaluddin
 **/
func TestGetArrayMapFromMap_BestCase2(t *testing.T) {
	testCases := []map[string]interface{}{
		{
			"data": []map[string]interface{}{
				{
					"name": "Mufid",
				},
				{
					"age":  24,
					"name": "Jamaluddin",
					"city": "Sumedang",
				},
				{
					"address": "Jl Pion Sumedang",
				},
			},
		},
		{
			"data": []map[string]interface{}{},
		},
		{
			"data": []map[string]interface{}{
				{},
				{
					"age":  24,
					"name": "Jamaluddin",
					"city": "Sumedang",
				},
				{
					"address": "Jl Pion Sumedang",
					"name":    "Monica Cecilia",
				},
			},
		},
	}

	expecteds := [][]map[string]interface{}{
		{
			{
				"name": "Mufid",
			},
			{
				"age":  24,
				"name": "Jamaluddin",
				"city": "Sumedang",
			},
			{
				"address": "Jl Pion Sumedang",
			},
		},
		{},
		{
			{},
			{
				"age":  24,
				"name": "Jamaluddin",
				"city": "Sumedang",
			},
			{
				"address": "Jl Pion Sumedang",
				"name":    "Monica Cecilia",
			},
		},
	}

	for idx, testCase := range testCases {
		expected, err := json.Marshal(expecteds[idx])
		assert.Nil(t, err)

		res, ok := GetArrayMapFromMap(testCase, "data")
		assert.True(t, ok, "Must be true")

		resJson, err := json.Marshal(res)
		assert.Nil(t, err)
		assert.JSONEq(t, string(expected), string(resJson), "Must be same")

		// ------------------------------ TEST WITH ARRAY INTERFACE ------------------------------

		arrOfInterface := []interface{}{}

		for _, item := range testCase["data"].([]map[string]interface{}) {
			arrOfInterface = append(arrOfInterface, item)
		}

		res2, ok := GetArrayMapFromMap(map[string]interface{}{"data": arrOfInterface}, "data")
		assert.True(t, ok, "Must be true")

		resJson2, err := json.Marshal(res2)
		assert.Nil(t, err)

		if string(expected) != "[]" {
			assert.JSONEq(t, string(expected), string(resJson2), "Must be same")
		} else {
			assert.Equal(t, "null", string(resJson2), "Must be same")
		}
	}
}

/**
 *	Worst Case Scenario Get []Map[String]Interface Field from Go Map
 *
 *	@author Mufid Jamaluddin
 **/
func TestGetArrayMapFromMap_WorstCase(t *testing.T) {
	testCases := []map[string]interface{}{
		{
			"data": map[string]interface{}{
				"name": "Mufid",
			},
		},
		{
			"data": map[string]interface{}{
				"age":  "24",
				"name": "Jamaluddin",
			},
		},
		{
			"data": map[string]interface{}{
				"waksa": 989,
			},
		},
	}

	var expected []map[string]interface{} = nil

	for _, testCase := range testCases {
		res, ok := GetArrayMapFromMap(testCase, "data")

		assert.False(t, ok, "Must be false")
		assert.Equal(t, expected, res, "Must be nil")
	}
}

/**
 *	Best Case Scenario Get Time.Time Field from Go Map
 *
 *	@author Mufid Jamaluddin
 **/
func TestGetTimeRFC3339NanoFromMap_BestCase(t *testing.T) {
	testCases := [...]map[string]interface{}{
		{
			"date": "2023-03-31T10:13:39.724207+07:00",
		},
		{
			"date": "2023-03-31T10:13:39.724207",
		},
		{
			"date": "2023-03-30T10:13:39",
		},
		{
			"date": "2023-03-30T01:13",
		},
		{
			"date": "2023-03-30",
		},
	}

	time1, _ := time.Parse(time.RFC3339Nano, "2023-03-31T10:13:39.724207+07:00")
	time2, _ := time.Parse(time.RFC3339Nano, "2023-03-31T10:13:39.724207+07:00")
	time3, _ := time.Parse(time.RFC3339Nano, "2023-03-30T10:13:39+07:00")
	time4, _ := time.Parse(time.RFC3339Nano, "2023-03-30T01:13:00+07:00")
	time5, _ := time.Parse(time.RFC3339Nano, "2023-03-30T00:00:00+07:00")

	expecteds := [...]time.Time{
		time1,
		time2,
		time3,
		time4,
		time5,
	}

	for idx, testCase := range testCases {
		expected := expecteds[idx]

		num, err := GetTimeRFC3339NanoFromMap(testCase, "date")

		assert.Nil(t, err, "Must be nil, not error")
		assert.Equal(t, expected, num, "Must be same")
	}
}

/**
 *	Worst Case Scenario Get Time.Time Field from Go Map
 *
 *	@author Mufid Jamaluddin
 **/
func TestGetTimeRFC3339NanoFromMap_WorstCase(t *testing.T) {
	testCases := []map[string]interface{}{
		{
			"name": "Mufid",
		},
		{},
		{
			"age": 24,
		},
		{
			"date": 95,
		},
		{
			"date": "2023-03-30 01:13",
		},
	}

	expected := time.Time{}

	for _, testCase := range testCases {
		num, err := GetTimeRFC3339NanoFromMap(testCase, "date")

		assert.NotNil(t, err, "Must be error")
		assert.Equal(t, expected, num, "Must be empty time struct")
	}
}
