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
		} else if current_value == "`" {
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

func GenerateRandomLetters(length uint64, uppercase bool, lowercase bool) (string) {
	rand.Seed(time.Now().UnixNano())
	
	letters_to_use := ""
	uppercase_letters :=  "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercase_letters := "abcdefghijklmnopqrstuvwxyz"

	if uppercase {
		letters_to_use += uppercase_letters
	} 

	if lowercase {
		letters_to_use += lowercase_letters
	}

	var sb strings.Builder

	l := len(letters_to_use)

	for i := uint64(0); i < length; i++ {
		c := letters_to_use[rand.Intn(l)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func IsFunc(object interface{}) bool {
	if IsNil(object) {
		return false
	}

	type_of := fmt.Sprintf("%T", object)
	if type_of == "*func(json.Map) []error" || 
	   type_of == "func(json.Map) []error" {
		return true
	}

	return false
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

	if string_value == "&map[value:<nil>]" {
		return true
	}

	if string_value == "%!s(*json.Value=<nil>)" {
		return true
	}

	rep := fmt.Sprintf("%T", object)

	if string_value == "%!s("+rep+"=<nil>)" {
		return true
	}

	if string_value == "&map[value:%!s(" + rep + "=<nil>)]" {
		return true
	}

	string_value = strings.ReplaceAll(string_value, "<nil>", "")
	string_value = strings.ReplaceAll(string_value, " ", "")
	
	if strings.HasPrefix(string_value, "&map[value:%!s(")  && strings.HasSuffix(string_value, "=&{})]") {
		return true
	}	

	return false
}

func IsTime(object interface{}, decimal_places int) bool {
	if IsNil(object) {
		return false
	}

	time, time_errors := GetTimeWithDecimalPlaces(object, decimal_places)
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

func IsValue(object interface{}) bool {
	if IsNil(object) {
		return false
	}

	type_of := GetType(object)
	if type_of == "json.Value" || type_of == "*json.Value" {
		return true
	}

	return false
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


func IsFloat(object interface{}) bool {
	return IsFloat32(object) || IsFloat64(object)
}


func IsFloat32(object interface{}) bool {
	if IsNil(object) {
		return false
	}

	type_of := GetType(object)
	switch type_of {
		case "*float32", 
			"float32":
		return true
	default: 
		return false
	}
}

func IsFloat64(object interface{}) bool {
	if IsNil(object) {
		return false
	}

	type_of := GetType(object)
	switch type_of {
		case "*float64", 
			"float64":
		return true
	default: 
		return false
	}
}

func IsInteger(object interface{}) bool {
	return IsInt(object) || IsUInt(object)
}
			
func IsInt(object interface{}) bool {
	if IsNil(object) {
		return false
	}

	if IsInt8(object) || IsInt16(object) || IsInt32(object) || IsInt64(object) {
		return true
	}

	type_of := GetType(object)
	switch type_of {
	case "*int", 
		  "int":
		return true
	default: 
		return false
	}
}

func IsUInt(object interface{}) bool {
	if IsNil(object) {
		return false
	}

	if IsUInt8(object) || IsUInt16(object) || IsUInt32(object) || IsUInt64(object) {
		return true
	}

	type_of := GetType(object)
	switch type_of {
	case "*uint", 
		  "uint":
		return true
	default: 
		return false
	}
}

func IsUInt8(object interface{}) bool {
	if IsNil(object) {
		return false
	}

	type_of := GetType(object)
	switch type_of {
	case "*uint8",
		  "uint8":
		return true
	default: 
		return false
	}
}

func IsUInt16(object interface{}) bool {
	if IsNil(object) {
		return false
	}

	type_of := GetType(object)
	switch type_of {
	case "*uint16",
		  "uint16":
		return true
	default: 
		return false
	}
}

func IsUInt32(object interface{}) bool {
	if IsNil(object) {
		return false
	}

	type_of := GetType(object)
	switch type_of {
	case "*uint32",
		  "uint32":
		return true
	default: 
		return false
	}
}

func IsUInt64(object interface{}) bool {
	if IsNil(object) {
		return false
	}

	type_of := GetType(object)
	switch type_of {
	case "*uint64",
		  "uint64":
		return true
	default: 
		return false
	}
}

func IsInt8(object interface{}) bool {
	if IsNil(object) {
		return false
	}

	type_of := GetType(object)
	switch type_of {
	case "*int8",
		  "int8":
		return true
	default: 
		return false
	}
}

func IsInt16(object interface{}) bool {
	if IsNil(object) {
		return false
	}

	type_of := GetType(object)
	switch type_of {
	case "*int16",
		  "int16":
		return true
	default: 
		return false
	}
}

func IsInt32(object interface{}) bool {
	if IsNil(object) {
		return false
	}

	type_of := GetType(object)
	switch type_of {
	case "*int32",
		  "int32":
		return true
	default: 
		return false
	}
}

func IsInt64(object interface{}) bool {
	if IsNil(object) {
		return false
	}

	type_of := GetType(object)
	switch type_of {
	case "*int64",
		  "int64":
		return true
	default: 
		return false
	}
}

func IsNumber(object interface{}) bool {
	return IsInt(object) || IsFloat(object) || IsUInt(object)
}

func GetTime(object interface{}) (*time.Time, []error) {
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
	
		temp_time, temp_time_error := time.Parse("2006-01-02 15:04:05.000000000", *(object.(*string)))
		if temp_time_error != nil {
			errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
		} else {
			result = &temp_time
		}

		if !IsNil(result) {
			return result, nil
		}

		temp_time, temp_time_error = time.Parse("2006-01-02 15:04:05.00000000", *(object.(*string)))
		if temp_time_error != nil {
			errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
		} else {
			result = &temp_time
		}

		if !IsNil(result) {
			return result, nil
		}

		temp_time, temp_time_error = time.Parse("2006-01-02 15:04:05.0000000", *(object.(*string)))
		if temp_time_error != nil {
			errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
		} else {
			result = &temp_time
		}


		if !IsNil(result) {
			return result, nil
		}

		temp_time, temp_time_error = time.Parse("2006-01-02 15:04:05.000000", *(object.(*string)))
		if temp_time_error != nil {
			errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
		} else {
			result = &temp_time
		}

		if !IsNil(result) {
			return result, nil
		}

		temp_time, temp_time_error = time.Parse("2006-01-02 15:04:05.00000", *(object.(*string)))
		if temp_time_error != nil {
			errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
		} else {
			result = &temp_time
		}

		if !IsNil(result) {
			return result, nil
		}

		temp_time, temp_time_error = time.Parse("2006-01-02 15:04:05.0000", *(object.(*string)))
		if temp_time_error != nil {
			errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
		} else {
			result = &temp_time
		}

		if !IsNil(result) {
			return result, nil
		}

		temp_time, temp_time_error = time.Parse("2006-01-02 15:04:05.000", *(object.(*string)))
		if temp_time_error != nil {
			errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
		} else {
			result = &temp_time
		}

		if !IsNil(result) {
			return result, nil
		}

		temp_time, temp_time_error = time.Parse("2006-01-02 15:04:05.00", *(object.(*string)))
		if temp_time_error != nil {
			errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
		} else {
			result = &temp_time
		}

		if !IsNil(result) {
			return result, nil
		}

		temp_time, temp_time_error = time.Parse("2006-01-02 15:04:05.0", *(object.(*string)))
		if temp_time_error != nil {
			errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
		} else {
			result = &temp_time
		}

		if !IsNil(result) {
			return result, nil
		}

		temp_time, temp_time_error = time.Parse("2006-01-02 15:04:05", *(object.(*string)))
		if temp_time_error != nil {
			errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
		} else {
			result = &temp_time
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

		temp_time, temp_time_error := time.Parse("2006-01-02 15:04:05.000000000", (object.(string)))
		if temp_time_error != nil {
			errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
		} else {
			result = &temp_time
		}

		if !IsNil(result) {
			return result, nil
		}

		temp_time, temp_time_error = time.Parse("2006-01-02 15:04:05.00000000", (object.(string)))
		if temp_time_error != nil {
			errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
		} else {
			result = &temp_time
		}

		if !IsNil(result) {
			return result, nil
		}

		temp_time, temp_time_error = time.Parse("2006-01-02 15:04:05.0000000", (object.(string)))
		if temp_time_error != nil {
			errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
		} else {
			result = &temp_time
		}

		if !IsNil(result) {
			return result, nil
		}

		temp_time, temp_time_error = time.Parse("2006-01-02 15:04:05.000000", (object.(string)))
		if temp_time_error != nil {
			errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
		} else {
			result = &temp_time
		}

		if !IsNil(result) {
			return result, nil
		}

		temp_time, temp_time_error = time.Parse("2006-01-02 15:04:05.00000", (object.(string)))
		if temp_time_error != nil {
			errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
		} else {
			result = &temp_time
		}

		if !IsNil(result) {
			return result, nil
		}

		temp_time, temp_time_error = time.Parse("2006-01-02 15:04:05.0000", (object.(string)))
		if temp_time_error != nil {
			errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
		} else {
			result = &temp_time
		}

		if !IsNil(result) {
			return result, nil
		}

		temp_time, temp_time_error = time.Parse("2006-01-02 15:04:05.000", (object.(string)))
		if temp_time_error != nil {
			errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
		} else {
			result = &temp_time
		}

		if !IsNil(result) {
			return result, nil
		}

		temp_time, temp_time_error = time.Parse("2006-01-02 15:04:05.00", (object.(string)))
		if temp_time_error != nil {
			errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
		} else {
			result = &temp_time
		}

		if !IsNil(result) {
			return result, nil
		}

		temp_time, temp_time_error = time.Parse("2006-01-02 15:04:05.0", (object.(string)))
		if temp_time_error != nil {
			errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
		} else {
			result = &temp_time
		}

		if !IsNil(result) {
			return result, nil
		}

		temp_time, temp_time_error = time.Parse("2006-01-02 15:04:05", (object.(string)))
		if temp_time_error != nil {
			errors = append(errors, fmt.Errorf("error: common.GetTime parsing error %s", fmt.Sprintf("%s",temp_time_error)))
		} else {
			result = &temp_time
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


func GetTimeWithDecimalPlaces(object interface{}, decimal_places int) (*time.Time, []error) {
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
	type_of := fmt.Sprintf("%T", object)
	return type_of
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

func MapPointerToStringArrayValueToInterface(a *[]string) *[]interface{} {
	if IsNil(a) {
		return nil
	}

	interface_array := make([]interface{}, len(*a))
	for index, value := range *a {
		interface_array[index] = value
	}
	
	return &interface_array
}
