package main

import (
	"fmt"
	"strconv"
)

func check(errs ...error) {
	if len(errs) > 0 && errs[0] != nil {

		for i := range errs {
			fmt.Println(errs[i])
		}

		panic(errs[0])
	}
}

func GetOptional[T any](mp *map[string]any, key string) *T {
	if mp == nil {
		return nil
	}

	val := (*mp)[key]

	if val == nil {
		return nil
	}

	_val, ok := val.(T)

	if !ok {
		panic(fmt.Sprintf("%s must be of type %T", key, _val))
	}

	return &_val
}

func toString(x any) string {
	stringValue, ok := x.(string)
	if ok {
		return stringValue
	}

	intValue, ok := x.(int)
	if ok {
		return strconv.Itoa(intValue)
	}

	boolValue, ok := x.(bool)
	if ok {
		if boolValue {
			return "true"
		} else {
			return "false"
		}
	}

	floatValue, ok := x.(float64)
	if ok {
		return fmt.Sprintf("%f", floatValue)
	}

	panic(fmt.Sprintf("Invalid type of variable (%T)", x))
}
