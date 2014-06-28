/*
Package flagtag provides support for creating command line flags by tagging appropriate struct fields with the 'flag' tag.
*/
package flagtag

import (
	"errors"
	"flag"
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
// - Invalid default values, error of type ErrInvalidDefault.
// - nil pointer provided.
// - nil interface provided.
// - interface to nil value provided.
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
			tag := parseTag(t, field.Tag.Get("flagopt"))
			if tag.Name == "" {
				// tag is invalid, since there is no name
				return errors.New("field '" + field.Name + "': invalid flag name: empty string")
			}
			switch fieldType.Kind() {
			case reflect.Ptr:
				// unwrap pointer
				if fieldValue.IsNil() {
					return errors.New("field '" + field.Name + "' (tag '" + tag.Name + "'): cannot use nil pointer")
				}
				fieldType = fieldType.Elem()
				fieldValue = fieldValue.Elem()
			case reflect.Interface:
				// check if interface is valid
				if fieldValue.IsNil() {
					return errors.New("field '" + field.Name + "' (tag '" + tag.Name + "'): cannot use nil interface")
				}
				var value = reflect.ValueOf(fieldValue.Interface())
				switch value.Type().Kind() {
				case reflect.Ptr, reflect.Interface:
					if value.IsNil() {
						return errors.New("field '" + field.Name + "' (tag '" + tag.Name + "'): cannot use nil interface value")
					}
				}
			}
			if !fieldValue.CanSet() {
				return errors.New("field '" + field.Name + "' (tag '" + tag.Name + "') is unexported or unaddressable: cannot use this field")
			}
			if !tag.Options.SkipFlagValue && registerFlagByValueInterface(fieldValue, &tag) {
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
	var value flag.Value
	switch fieldValue.Type().Kind() {
	case reflect.Interface:
		var ok bool
		value, ok = fieldValue.Interface().(flag.Value)
		if !ok {
			return false
		}
	default:
		var ok bool
		value, ok = fieldValue.Addr().Interface().(flag.Value)
		if !ok {
			return false
		}
	}
	flag.Var(value, tag.Name, tag.Description)
	if tag.DefaultValue != "" {
		// a default value is provided, first call value.Set() with the provided default value
		value.Set(tag.DefaultValue)
	}
	return true
}

// registerFlagByPrimitive registers a single field as one of the primitive flag types. Types are matched by
// kind, so types derived from one of the basic types are still eligible for a flag.
//
// If it is not possible to register a flag because of an unknown data type, an error will be returned.
// If the specified default value is invalid, an error of type ErrInvalidDefault will be returned.
func registerFlagByPrimitive(fieldName string, fieldValue reflect.Value, tag *flagTag) error {
	var fieldType = fieldValue.Type()
	// Check time.Duration first, since it will also match one of the basic kinds.
	if durationVar, ok := fieldValue.Addr().Interface().(*time.Duration); ok {
		// field is a time.Duration
		defaultVal, err := time.ParseDuration(tag.DefaultValue)
		if err != nil {
			return &ErrInvalidDefault{fieldName, tag.Name, err}
		}
		flag.DurationVar(durationVar, tag.Name, defaultVal, tag.Description)
		return nil
	}
	// Check basic kinds.
	// TODO convert to detected kind without using unsafe
	var fieldPtr = unsafe.Pointer(fieldValue.UnsafeAddr())
	switch fieldType.Kind() {
	case reflect.String:
		flag.StringVar((*string)(fieldPtr), tag.Name, tag.DefaultValue, tag.Description)
	case reflect.Bool:
		defaultVal, err := strconv.ParseBool(tag.DefaultValue)
		if err != nil {
			return &ErrInvalidDefault{fieldName, tag.Name, err}
		}
		flag.BoolVar((*bool)(fieldPtr), tag.Name, defaultVal, tag.Description)
	case reflect.Float64:
		defaultVal, err := strconv.ParseFloat(tag.DefaultValue, 64)
		if err != nil {
			return &ErrInvalidDefault{fieldName, tag.Name, err}
		}
		flag.Float64Var((*float64)(fieldPtr), tag.Name, defaultVal, tag.Description)
	case reflect.Int:
		defaultVal, err := strconv.ParseInt(tag.DefaultValue, 0, fieldType.Bits())
		if err != nil {
			return &ErrInvalidDefault{fieldName, tag.Name, err}
		}
		flag.IntVar((*int)(fieldPtr), tag.Name, int(defaultVal), tag.Description)
	case reflect.Int64:
		defaultVal, err := strconv.ParseInt(tag.DefaultValue, 0, 64)
		if err != nil {
			return &ErrInvalidDefault{fieldName, tag.Name, err}
		}
		flag.Int64Var((*int64)(fieldPtr), tag.Name, defaultVal, tag.Description)
	case reflect.Uint:
		defaultVal, err := strconv.ParseUint(tag.DefaultValue, 0, fieldType.Bits())
		if err != nil {
			return &ErrInvalidDefault{fieldName, tag.Name, err}
		}
		flag.UintVar((*uint)(fieldPtr), tag.Name, uint(defaultVal), tag.Description)
	case reflect.Uint64:
		defaultVal, err := strconv.ParseUint(tag.DefaultValue, 0, 64)
		if err != nil {
			return &ErrInvalidDefault{fieldName, tag.Name, err}
		}
		flag.Uint64Var((*uint64)(fieldPtr), tag.Name, defaultVal, tag.Description)
	default:
		return errors.New("unsupported data type (kind '" + strconv.FormatUint(uint64(fieldType.Kind()), 10) + "') for field '" + fieldName + "' (tag '" + tag.Name + "')")
	}
	return nil
}

// getStructValue checks that the provided config instance is actually a struct not a nil value.
func getStructValue(config interface{}) (reflect.Value, error) {
	var zero reflect.Value
	if config == nil {
		return zero, errors.New("config cannot be nil")
	}
	ptr := reflect.ValueOf(config)
	if ptr.IsNil() {
		return zero, errors.New("config cannot point to nil")
	}
	val := reflect.Indirect(ptr)
	if val.Kind() != reflect.Struct {
		return zero, errors.New("config instance is not a struct")
	}
	return val, nil
}

// parseTag parses a string of text and separates the various sections of the 'flag'-tag.
func parseTag(value string, optvalue string) flagTag {
	parts := strings.SplitN(value, ",", 3)
	for len(parts) < 3 {
		parts = append(parts, "")
	}
	var flag = flagTag{Name: parts[0], DefaultValue: parts[1], Description: parts[2]}
	if optvalue != "" {
		if strings.Contains(optvalue, "skipFlagValue") {
			flag.Options.SkipFlagValue = true
		}
	}
	return flag
}

// flagTag contains the parsed tag values.
type flagTag struct {
	Name         string
	DefaultValue string
	Description  string
	Options      struct {
		SkipFlagValue bool
	}
}

// ErrInvalidDefault is an error type for the case of invalid defaults.
type ErrInvalidDefault struct {
	field string
	tag   string
	err   error
}

// Error returns the error explaining the bad default value.
func (e *ErrInvalidDefault) Error() string {
	return "invalid default value for field '" + e.field + "' (tag '" + e.tag + "'): " + e.err.Error()
}
