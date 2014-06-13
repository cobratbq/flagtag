/*
Package flagtag provides support for creating command line flags by tagging appropriate struct fields with the 'flag' tag.
*/
package flagtag

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// MustConfigureAndParse is like ConfigureAndParse, the only difference is that
// it will panic in case of an error.
func MustConfigureAndParse(config interface{}) {
	if err := ConfigureAndParse(config); err != nil {
		panic(err)
	}
}

// MustConfigure is like Configure, the only difference is that it will panic
// in case of an error.
func MustConfigure(config interface{}) {
	if err := Configure(config); err != nil {
		panic(err)
	}
}

// ConfigureAndParse will first attempt to configure the flags according to the
// provided config type. If any error occurs, this error will be returned and
// the command line arguments will not be parsed. If no error occurs, the
// command line arguments will be parsed and the config type will contain the
// result.
// Using this function may remove the need to even import the flag package at
// all.
func ConfigureAndParse(config interface{}) error {
	if err := Configure(config); err != nil {
		return err
	}
	flag.Parse()
	return nil
}

// Configure will configure the flag parameters according to the tags of the
// provided data type. It is allowed to call this method multiple times with
// different data types. (As long as flag's Parse() method has not been called
// yet.)
// Fields without a 'flag' tag or with an empty 'flag' tag will be ignored.
//
// The 'flag' tag consists of 3 parts, similar to the *Var-functions of the
// flag package. Parts are separated by a comma. The parts are:
//  - 1st: flag name
//  - 2nd: default value
//  - 3rd: usage description
//
// Example:
//  `flag:"verbose,false,Enable verbose output."`.
//
// This will create a flag 'verbose', which defaults to 'false' and shows usage
// information "Enable verbose output.".
//
// If an error occurs, this error will be returned and the configuration of
// other struct fields will be aborted.
func Configure(config interface{}) error {
	val, err := getStructValue(config)
	if err != nil {
		return err
	}
	return configure(val)
}

// configure (recursively) configures flags as they are discovered in the provided type and value.
// In case of an error, the error is returned. Possible errors are:
// - Invalid default values.
// - nil pointer provided.
// - Tagged variable uses unsupported data type.
func configure(structValue reflect.Value) error {
	var structType = structValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldType := field.Type
		fieldValue := structValue.Field(i)
		t := field.Tag.Get("flag")
		if t == "" {
			// if field is not tagged then we do not need to flag the type itself
			if fieldType.Kind() == reflect.Struct {
				// kind is a struct => recurse into inner struct
				if err := configure(fieldValue); err != nil {
					return err
				}
			}
		} else {
			// field is tagged, continue investigating what kind of flag to create
			tag := parseTag(t)
			if tag.Name == "" {
				// tag is invalid, since there is no name
				return fmt.Errorf("field '%s': invalid flag name: empty string", field.Name)
			}
			if fieldType.Kind() == reflect.Ptr {
				// unwrap pointer
				if fieldValue.IsNil() {
					return fmt.Errorf("field '%s' (tag '%s'): cannot use nil pointer", field.Name, tag.Name)
				}
				fieldType = fieldType.Elem()
				fieldValue = fieldValue.Elem()
			}
			if !fieldValue.CanSet() {
				return fmt.Errorf("field '%s' (tag '%s') is unexported or unaddressable: cannot use this field", field.Name, tag.Name)
			}
			// TODO create a tag hint for ignoring the ValueInterface check
			if registerFlagByValueInterface(fieldValue, &tag) {
				// no error during registration => Var-flag registered => continue with next field
				continue
			}
			if err := registerFlagByPrimitive(field.Name, fieldValue, &tag); err != nil {
				return err
			}
		}
	}
	return nil
}

// registerFlagByValueInterface checks if the provided type can be treated as flag.Value.
// If so, a flag.Value flag is set and true is returned. If no flag is set, false is returned.
func registerFlagByValueInterface(fieldValue reflect.Value, tag *flagTag) bool {
	if value, ok := fieldValue.Addr().Interface().(flag.Value); ok {
		// field type implements flag.Value interface, register as such
		// TODO (how to) set default value? (i.e. ignore default value?)
		flag.Var(value, tag.Name, tag.Description)
		return true
	}
	return false
}

// registerFlagByPrimitive registers a single field as one of the primitive flag types. Types are matched by
// kind, so types derived from one of the basic types are still eligible for a flag.
//
// If it is not possible to register a flag because of an unknown type, an error will be returned.
// If the default value is invalid, an error will be returned.
func registerFlagByPrimitive(fieldName string, fieldValue reflect.Value, tag *flagTag) error {
	var fieldType = fieldValue.Type()
	// Check time.Duration first, since it will also match one of the basic kinds.
	if durationVar, ok := fieldValue.Addr().Interface().(*time.Duration); ok {
		// field is a time.Duration
		defaultVal, err := time.ParseDuration(tag.DefaultValue)
		if err != nil {
			return fmt.Errorf("invalid default value for field '%s' (tag '%s'): %s", fieldName, tag.Name, err.Error())
		}
		flag.DurationVar(durationVar, tag.Name, defaultVal, tag.Description)
		return nil
	}
	// Check basic kinds.
	var fieldPtr = unsafe.Pointer(fieldValue.UnsafeAddr())
	switch fieldType.Kind() {
	case reflect.String:
		flag.StringVar((*string)(fieldPtr), tag.Name, tag.DefaultValue, tag.Description)
	case reflect.Bool:
		defaultVal, err := strconv.ParseBool(tag.DefaultValue)
		if err != nil {
			return fmt.Errorf("invalid default value for field '%s' (tag '%s'): %s", fieldName, tag.Name, err.Error())
		}
		flag.BoolVar((*bool)(fieldPtr), tag.Name, defaultVal, tag.Description)
	case reflect.Float64:
		defaultVal, err := strconv.ParseFloat(tag.DefaultValue, 64)
		if err != nil {
			return fmt.Errorf("invalid default value for field '%s' (tag '%s'): %s", fieldName, tag.Name, err.Error())
		}
		flag.Float64Var((*float64)(fieldPtr), tag.Name, defaultVal, tag.Description)
	case reflect.Int:
		// TODO parse exact number of available bits, or always 64?
		defaultVal, err := strconv.ParseInt(tag.DefaultValue, 0, fieldType.Bits())
		if err != nil {
			return fmt.Errorf("invalid default value for field '%s' (tag '%s'): %s", fieldName, tag.Name, err.Error())
		}
		flag.IntVar((*int)(fieldPtr), tag.Name, int(defaultVal), tag.Description)
	case reflect.Int64:
		defaultVal, err := strconv.ParseInt(tag.DefaultValue, 0, 64)
		if err != nil {
			return fmt.Errorf("invalid default value for field '%s' (tag '%s'): %s", fieldName, tag.Name, err.Error())
		}
		flag.Int64Var((*int64)(fieldPtr), tag.Name, defaultVal, tag.Description)
	case reflect.Uint:
		// TODO parse exact number of available bits, or always 64?
		defaultVal, err := strconv.ParseUint(tag.DefaultValue, 0, fieldType.Bits())
		if err != nil {
			return fmt.Errorf("invalid default value for field '%s' (tag '%s'): %s", fieldName, tag.Name, err.Error())
		}
		flag.UintVar((*uint)(fieldPtr), tag.Name, uint(defaultVal), tag.Description)
	case reflect.Uint64:
		defaultVal, err := strconv.ParseUint(tag.DefaultValue, 0, 64)
		if err != nil {
			return fmt.Errorf("invalid default value for field '%s' (tag '%s'): %s", fieldName, tag.Name, err.Error())
		}
		flag.Uint64Var((*uint64)(fieldPtr), tag.Name, defaultVal, tag.Description)
	default:
		return fmt.Errorf("unsupported data type for field '%s' (tag '%s')", fieldName, tag.Name)
	}
	return nil
}

// getStructValue checks that the provided config instance is actually a struct not a nil value.
func getStructValue(config interface{}) (reflect.Value, error) {
	var zero reflect.Value
	if config == nil {
		return zero, fmt.Errorf("config cannot be nil")
	}
	ptr := reflect.ValueOf(config)
	if ptr.IsNil() {
		return zero, fmt.Errorf("config cannot point to nil")
	}
	val := reflect.Indirect(ptr)
	if val.Kind() != reflect.Struct {
		return zero, fmt.Errorf("config instance is not a struct")
	}
	return val, nil
}

// parseTag parses a string of text and separates the various sections of the 'flag'-tag.
func parseTag(value string) flagTag {
	parts := strings.SplitN(value, ",", 3)
	for len(parts) < 3 {
		parts = append(parts, "")
	}
	return flagTag{Name: parts[0], DefaultValue: parts[1], Description: parts[2]}
}

// flagTag contains the parsed tag values
type flagTag struct {
	Name         string
	DefaultValue string
	Description  string
}
