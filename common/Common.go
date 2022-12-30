package common

import (
	"fmt"
	"strings"
	"time"
	"math/rand"
)

func GetDataDirectory() []string {
	directory := []string{"Volumes", "ramdisk"}
	return directory
}

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

	if string_value == "nil" {
		return true
	}

	if string_value == "map[]" {
		return true
	}

	rep := fmt.Sprintf("%T", object)

	if string_value == "%!s("+rep+"=<nil>)" {
		return true
	}

	if string_value == "&map[value:%!s(" + rep + "=<nil>)]" {
		return true
	}

	return false
}

func IsTime(object interface{}, decimal_places int) bool {
	if IsNil(object) {
		return false
	}

	time, time_errors := GetTime(object, decimal_places)
	if time_errors != nil {
		return false
	} else if IsNil(time) {
		return false
	}

	return true
}


func IsMap(object interface{}) bool {
	if IsNil(object) {
		return false
	}

	type_of := GetType(object)
	if type_of == "json.Map" || type_of == "*json.Map" {
		return true
	}

	return false
}

func IsString(object interface{}) bool {
	if IsNil(object) {
		return false
	}

	type_of := GetType(object)
	if type_of == "string" || type_of == "*string" {
		return true
	}

	return false
}


func IsBool(object interface{}) bool {
	if IsNil(object) {
		return false
	}

	type_of := GetType(object)
	if type_of == "bool" || type_of == "*bool" {
		return true
	}

	return false
}

func IsBoolTrue(object interface{}) bool {
	if !IsBool(object) {
		return false
	}

	type_of := GetType(object)
	switch type_of {
		case "bool":
			return object.(bool) == true
		case "*bool":
			return *(object.(*bool)) == true
		default:
			return false
	}
}

func IsBoolFalse(object interface{}) bool {
	if !IsBool(object) {
		return false
	}

	type_of := GetType(object)
	switch type_of {
		case "bool":
			return object.(bool) == false
		case "*bool":
			return *(object.(*bool)) == false
		default:
			return false
	}
}

func IsArray(object interface{}) bool {
	if IsNil(object) {
		return false
	}

	type_of := GetType(object)
	if type_of == "json.Array" || type_of == "*json.Array" {
		return true
	}

	return false
}

func IsInteger(object interface{}) bool {
	if IsNil(object) {
		return false
	}

	type_of := GetType(object)
	switch type_of {
	case "*int", 
		  "int",
		  "*uint", 
		  "uint",
		  "*int64",
		  "int64",
		  "*uint64",
		  "uint64",
		  "*int32",
		  "int32",
		  "*uint32",
		  "uint32",
		  "*int16",
		  "int16",
		  "*uint16",
		  "uint16",
		  "*int8",
		  "int8",
		  "*uint8",
		  "uint8":
		return true
	default: 
		return false
	}
}


func GetTime(object interface{}, decimal_places int) (*time.Time, []error) {
	var errors []error
	var result *time.Time

	if object == nil {
		return nil, nil
	}

	// time package does not allow zeroes and (so it should!, database allows zero)
	zero_mapping := map[string]int{"zero":0, 
		"0000-00-00 00:00:00.000000000":0, 
		"0000-00-00 00:00:00.00000000":0, 
		"0000-00-00 00:00:00.0000000":0, 
		"0000-00-00 00:00:00.000000":0, 
		"0000-00-00 00:00:00.00000":0, 
		"0000-00-00 00:00:00.0000":0, 
		"0000-00-00 00:00:00.000":0, 
		"0000-00-00 00:00:00.00":0, 
		"0000-00-00 00:00:00.0":0, 
		"0000-00-00 00:00:00":0, 
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
		value := *(object.(*string))
		if value == "now" {
			result = GetTimeNow()
		} else {
			_, zero_map_value_found := zero_mapping[value]
			if zero_map_value_found {
				result = GetTimeZero()
			}
		}

		if !IsNil(result) {
			return result, nil
		}
		
		if decimal_places < 0 || decimal_places > 9 {
			errors = append(errors, fmt.Errorf("error: common.GetTime decimal places not supported [0,9] actual %d", decimal_places))
			return nil, errors
		}

		switch decimal_places {
		case 9:
			temp_time, temp_time_error := time.Parse("2006-01-02 15:04:05.000000000", *(object.(*string)))
			if temp_time_error != nil {
				errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
			} else {
				result = &temp_time
			}
		case 8:
			temp_time, temp_time_error := time.Parse("2006-01-02 15:04:05.00000000", *(object.(*string)))
			if temp_time_error != nil {
				errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
			} else {
				result = &temp_time
			}
		case 7:
			temp_time, temp_time_error := time.Parse("2006-01-02 15:04:05.0000000", *(object.(*string)))
			if temp_time_error != nil {
				errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
			} else {
				result = &temp_time
			}
		case 6:
			temp_time, temp_time_error := time.Parse("2006-01-02 15:04:05.000000", *(object.(*string)))
			if temp_time_error != nil {
				errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
			} else {
				result = &temp_time
			}
		case 5:
			temp_time, temp_time_error := time.Parse("2006-01-02 15:04:05.00000", *(object.(*string)))
			if temp_time_error != nil {
				errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
			} else {
				result = &temp_time
			}
		case 4:
			temp_time, temp_time_error := time.Parse("2006-01-02 15:04:05.0000", *(object.(*string)))
			if temp_time_error != nil {
				errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
			} else {
				result = &temp_time
			}
		case 3:
			temp_time, temp_time_error := time.Parse("2006-01-02 15:04:05.000", *(object.(*string)))
			if temp_time_error != nil {
				errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
			} else {
				result = &temp_time
			}
		case 2:
			temp_time, temp_time_error := time.Parse("2006-01-02 15:04:05.00", *(object.(*string)))
			if temp_time_error != nil {
				errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
			} else {
				result = &temp_time
			}
		case 1:
			temp_time, temp_time_error := time.Parse("2006-01-02 15:04:05.0", *(object.(*string)))
			if temp_time_error != nil {
				errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
			} else {
				result = &temp_time
			}
		case 0:
			temp_time, temp_time_error := time.Parse("2006-01-02 15:04:05", *(object.(*string)))
			if temp_time_error != nil {
				errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
			} else {
				result = &temp_time
			}
		}

		if !IsNil(result) {
			return result, nil
		}

		errors = append(errors,  fmt.Errorf("error: common.GetTimeNow value not supported %s", *(object.(*string))))
	case "string":
		value := (object.(string))
		if value == "now" {
			result = GetTimeNow()
		} else {
			_, zero_map_value_found := zero_mapping[value]
			if zero_map_value_found {
				result = GetTimeZero()
			}
		}
		if !IsNil(result) {
			return result, nil
		}
		
		if decimal_places < 0 || decimal_places > 9 {
			errors = append(errors, fmt.Errorf("error: common.GetTime decimal places not supported [0,9] actual %d", decimal_places))
			return nil, errors
		}

		switch decimal_places {
		case 9:
			temp_time, temp_time_error := time.Parse("2006-01-02 15:04:05.000000000", (object.(string)))
			if temp_time_error != nil {
				errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
			} else {
				result = &temp_time
			}
		case 8:
			temp_time, temp_time_error := time.Parse("2006-01-02 15:04:05.00000000", (object.(string)))
			if temp_time_error != nil {
				errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
			} else {
				result = &temp_time
			}
		case 7:
			temp_time, temp_time_error := time.Parse("2006-01-02 15:04:05.0000000", (object.(string)))
			if temp_time_error != nil {
				errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
			} else {
				result = &temp_time
			}
		case 6:
			temp_time, temp_time_error := time.Parse("2006-01-02 15:04:05.000000", (object.(string)))
			if temp_time_error != nil {
				errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
			} else {
				result = &temp_time
			}
		case 5:
			temp_time, temp_time_error := time.Parse("2006-01-02 15:04:05.00000", (object.(string)))
			if temp_time_error != nil {
				errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
			} else {
				result = &temp_time
			}
		case 4:
			temp_time, temp_time_error := time.Parse("2006-01-02 15:04:05.0000", (object.(string)))
			if temp_time_error != nil {
				errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
			} else {
				result = &temp_time
			}
		case 3:
			temp_time, temp_time_error := time.Parse("2006-01-02 15:04:05.000", (object.(string)))
			if temp_time_error != nil {
				errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
			} else {
				result = &temp_time
			}
		case 2:
			temp_time, temp_time_error := time.Parse("2006-01-02 15:04:05.00", (object.(string)))
			if temp_time_error != nil {
				errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
			} else {
				result = &temp_time
			}
		case 1:
			temp_time, temp_time_error := time.Parse("2006-01-02 15:04:05.0", (object.(string)))
			if temp_time_error != nil {
				errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
			} else {
				result = &temp_time
			}
		case 0:
			temp_time, temp_time_error := time.Parse("2006-01-02 15:04:05", (object.(string)))
			if temp_time_error != nil {
				errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
			} else {
				result = &temp_time
			}
		}

		if !IsNil(result) {
			return result, nil
		}
		
		errors = append(errors,  fmt.Errorf("error: common.GetTimeNow value not supported %s", *(object.(*string))))
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
	now := time.Now().UTC()
	return &now
}

func GetTimeZero() *time.Time {
	zero := time.Time{}
	return &zero
}

func FormatTime(value time.Time, decimal_places int) (*string, []error) {
	var errors []error
	if decimal_places < 0 || decimal_places > 9 {
		errors = append(errors, fmt.Errorf("error: common.GetTimeNowString decimal places not support [0,9]: %d", decimal_places))
		return nil, errors
	}

	var result string
	if decimal_places == 0 {
		result = value.Format("2006-01-02 15:04:05")
	} else {
		post_fix := "."
		for i := 0; i < decimal_places; i++ {
			post_fix += "0"
		}
		result = value.Format("2006-01-02 15:04:05" + post_fix)
	}
	return &result, nil
}

func GetTimeNowString(decimal_places int) (*string, []error) {
	var errors []error
	if decimal_places < 0 || decimal_places > 9 {
		errors = append(errors, fmt.Errorf("error: common.GetTimeNowString decimal places not support [0,9]: %d", decimal_places))
		return nil, errors
	}

	var result string
	time_now := *GetTimeNow()
	if decimal_places == 0 {
		result = time_now.Format("2006-01-02 15:04:05")
	} else {
		post_fix := "."
		for i := 0; i < decimal_places; i++ {
			post_fix += "0"
		}
		result = time_now.Format("2006-01-02 15:04:05" + post_fix)
	}
	return &result, nil
}

func GetTimeZeroStringSQL(decimal_places int) (*string, []error) {
	var errors []error
	if decimal_places < 0 || decimal_places > 9 {
		errors = append(errors, fmt.Errorf("error: common.GetTimeZeroString decimal places not support [0,9]: %d", decimal_places))
		return nil, errors
	}

	result := "0000-00-00 00:00:00"
	if decimal_places > 0 {
		result += "."
		for i := 0; i < decimal_places; i++ {
			result += "0"
		}
	}

	return &result, nil
}
