package common

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// ConvertValueToFloat64 converts interface value to float64
func ConvertValueToFloat64(value interface{}) (float64, error) {
	var resultValue float64

	switch val := value.(type) {
	case nil:
		return float64(0), errors.New("Unable to convert metric value: value is nil type")
	case uint:
		resultValue = float64(int64(val))
	case uint8:
		resultValue = float64(int64(val))
	case uint16:
		resultValue = float64(int64(val))
	case uint32:
		resultValue = float64(int64(val))
	case int:
		resultValue = float64(int64(val))
	case int8:
		resultValue = float64(int64(val))
	case int16:
		resultValue = float64(int64(val))
	case int32:
		resultValue = float64(int64(val))
	case int64:
		resultValue = float64(int64(val))
	case uint64:
		// Prometheus does not support writing uint64
		if val < uint64(9223372036854775808) {
			resultValue = float64(int64(val))
		} else {
			resultValue = float64(9223372036854775807)
		}
	case float32:
		resultValue = float64(val)
	case float64:
		// NaNs are invalid values in Prometheus, skip measurement
		if math.IsNaN(val) || math.IsInf(val, 0) {
			return float64(0), errors.New("Unable to convert metric value: value is a Nan or Inf")
		}

		resultValue = val
	case bool:
		resultValue = float64(0)
		if val == true {
			resultValue = float64(1)
		}
	case string:
		resultValue = float64(0)
		if strings.ToLower(val) == "up" {
			resultValue = float64(1)
		}
	default:
		return float64(0), fmt.Errorf("Unable to convert metric value: invalid type, type: '%v'", val)
	}

	return resultValue, nil
}
