package class

import (
	"fmt"
	"strings"
	"time"
	"math/rand"
)


func EscapeString(value string, string_quote_value string) (string, error) {
	if !(string_quote_value == "'" || string_quote_value == "\"") {
		return "", fmt.Errorf(fmt.Sprintf("string_quote_value not supported: %s available values are ' or \"", string_quote_value))
	}

	var result strings.Builder
	runes := []rune(value)
	length := len(runes)
	i := 0
	for ;i < length ; i++ {
		current_value := string(runes[i])
		if current_value == "\\" {
			result.WriteString("\\\\")
		} else if current_value == string_quote_value {
			result.WriteString("\\")
			result.WriteString(current_value)
		} else {
			result.WriteString(current_value)
		}
	}

	return result.String(), nil
}

func CloneString(value *string) *string {
	if value == nil {
		return nil
	}

	temp := strings.Clone(*value)
	return &temp
}

func GenerateRandomLetters(length uint64, upper_case *bool) (*string) {
	rand.Seed(time.Now().UnixNano())
	
	var letters_to_use string
	uppercase_letters :=  "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercase_letters := "abcdefghijklmnopqrstuvwxyz"

	if upper_case == nil {
		letters_to_use = uppercase_letters + lowercase_letters
	} else if *upper_case {
		letters_to_use = uppercase_letters
	} else {
		letters_to_use = lowercase_letters
	}

	var sb strings.Builder

	l := len(letters_to_use)

	for i := uint64(0); i < length; i++ {
		c := letters_to_use[rand.Intn(l)]
		sb.WriteByte(c)
	}

	value := sb.String()
	return &value
}

func IsNil(object interface{}) bool {
	if object == nil {
		return true
	}
	
	string_value := fmt.Sprintf("%s", object) 

	if string_value == "<nil>" {
		return true
	}

	if string_value == "map[]" {
		return true
	}

	rep := fmt.Sprintf("%T", object)

	if string_value == "%!s("+rep+"=<nil>)" {
		return true
	}

	return false
}

func IsTime(object interface{}) bool {
	if IsNil(object) {
		return false
	}

	time, time_errors := GetTime(object)
	if time_errors != nil {
		return false
	} else if IsNil(time) {
		return false
	}

	return true
}

func GetTime(object interface{}) (*time.Time, []error) {
	var errors []error
	var result *time.Time

	if object == nil {
		return nil, nil
	}

	rep := fmt.Sprintf("%T", object)
	switch rep {
	case "*time.Time":
		value := *(object.(*time.Time))
		result = &value
	case "time.Time":
		value := object.(time.Time)
		result = &value
	case "*string":
		//todo: parse for null
		value1, value_errors1 := time.Parse("2006-01-02 15:04:05.000000", *(object.(*string)))
		value2, value_errors2 := time.Parse("2006-01-02 15:04:05", *(object.(*string)))
		var value3 *time.Time
		if *(object.(*string)) == "now" {
			value3 = GetTimeNow()
		} else {
			value3 = nil
		}

		if value_errors1 != nil && value_errors2 != nil && value3 == nil {
			errors = append(errors, value_errors1)
		}

		if value_errors1 == nil {
			result = &value1
		}

		if value_errors2 == nil {
			result = &value2
		}

		if value3 != nil {
			result = value3
		}

	case "string":
		//todo: parse for null
		value1, value_errors1 := time.Parse("2006-01-02 15:04:05.000000", (object.(string)))
		value2, value_errors2 := time.Parse("2006-01-02 15:04:05", (object.(string)))
		var value3 *time.Time
		if (object.(string)) == "now" {
			value3 = GetTimeNow()
		} else {
			value3 = nil
		}

		if value_errors1 != nil && value_errors2 != nil && value3 == nil {
			errors = append(errors, value_errors1)
		}

		if value_errors1 == nil {
			result = &value1
		}

		if value_errors2 == nil {
			result = &value2
		}

		if value3 != nil {
			result = value3
		}

	default:
		errors = append(errors, fmt.Errorf("error: json.Map.GetTime: type %s is not supported please implement", rep))
	}

	if len(errors) > 0 {
		return nil, errors
	}

	return result, nil
}

func GetType(object interface{}) string {
	if IsNil(object) {
		return "nil"
	}
	return fmt.Sprintf("%T", object)
}



func Contains(sl []string, name string) bool {
	for _, value := range sl {
		if value == name {
			return true
		}
	}
	return false
}

func GetTimeNow() *time.Time {
	now := time.Now()
	return &now
}

func FormatTime(value time.Time) string {
	return value.Format("2006-01-02 15:04:05.000000")
}

func GetTimeNowString() string {
	return (*GetTimeNow()).Format("2006-01-02 15:04:05.000000")
}
