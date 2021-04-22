package gocloud

import (
	"errors"
	"fmt"
	"strconv"
)

type Mp map[string]interface{}

func (c Mp) GetString(key string) string {
	v, ok := c[key]
	if !ok {
		return ""
	}
	return fmt.Sprint(v)
}
func (c Mp) GetInt(key string) (int64, error) {
	v, ok := c[key]
	if !ok {
		return 0, errors.New("not found")
	}
	switch v.(type) {
	case int:
		return v.(int64), nil
	case string:
		return strconv.ParseInt(v.(string), 10, 64)
	case int64:
		return v.(int64), nil
	case float32:
		return int64(v.(float32)), nil
	case float64:
		return int64(v.(float64)), nil
	}
	return 0, errors.New("not found")
}
func (c Mp) GetFloat(key string) (float64, error) {
	v, ok := c[key]
	if !ok {
		return 0, errors.New("not found")
	}
	switch v.(type) {
	case int:
		return float64(v.(int)), nil
	case string:
		return strconv.ParseFloat(v.(string), 64)
	case int64:
		return float64(v.(int64)), nil
	case float32:
		return float64(v.(float32)), nil
	case float64:
		return v.(float64), nil
	}
	return 0, errors.New("not found")
}
func (c Mp) GetBool(key string) bool {
	v, ok := c[key]
	if !ok {
		return false
	}
	switch v.(type) {
	case bool:
		return v.(bool)
	}
	return false
}
