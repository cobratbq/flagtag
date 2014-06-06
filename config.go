package flagtag

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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
// - 1st: flag name
// - 2nd: default value
// - 3rd: usage description
// Example tag: `flag:"verbose,false,Enable verbose output."`.
// This will create a flag 'verbose', which defaults to 'false' and shows usage
// information "Enables default output.".
//
// If an error occurs, this error will be returned and the configuration of
// other struct fields will be aborted.
func Configure(config interface{}) error {
	val, err := getStructValue(config)
	if err != nil {
		return err
	}
	return configure(val.Type(), val.UnsafeAddr())
}

func configure(structType reflect.Type, baseAddr uintptr) error {
	for i := 0; i < structType.NumField(); i++ {
		f := (reflect.StructField)(structType.Field(i))
		t := f.Tag.Get("flag")
		if t == "" {
			// if field is not tagged then we do not need to flag the type itself
			if f.Type.Kind() == reflect.Struct {
				// kind is a struct => recurse into inner struct
				if err := configure(f.Type, baseAddr+f.Offset); err != nil {
					return err
				}
			}
		} else {
			// field is tagged, continue investigating what kind of flag to create
			tag := parseTag(t)
			if tag.Name == "" {
				// tag is invalid, since there is no name
				return fmt.Errorf("invalid flag name: empty string")
			}
			var fieldptr = unsafe.Pointer(baseAddr + f.Offset)
			var fieldtype = f.Type
			if fieldtype.Kind() == reflect.Interface {
				// unwrap interface indirection
				var ifValue = reflect.NewAt(fieldtype, fieldptr).Elem()
				if ifValue.Interface() == nil {
					// nil interface, return error
					return fmt.Errorf("cannot use nil interface for flag target")
				}
				// non-nil interface, so continue investigation
				if reflect.TypeOf(ifValue.Interface()).Kind() == reflect.Ptr && reflect.ValueOf(ifValue.Interface()).IsNil() {
					// interface with nil pointer, return error
					return fmt.Errorf("cannot use interface that contains nil pointer")
				} else {
					// actual interface with legitimate content
					fieldtype = reflect.TypeOf(ifValue.Interface())
					fieldptr = unsafe.Pointer(ifValue.UnsafeAddr())
				}
			}
			if fieldtype.Kind() == reflect.Ptr {
				// unwrap pointer indirection
				var ptrTarget = reflect.NewAt(fieldtype, fieldptr).Elem()
				if !ptrTarget.IsNil() {
					fieldtype = fieldtype.Elem()
					fieldptr = unsafe.Pointer(ptrTarget.Pointer())
				}
			}
			if registerFlagByValueInterface(fieldtype, fieldptr, &tag) {
				// no error during registration => Var-flag registered => continue with next field
				continue
			}
			// TODO support Duration
			if err := registerFlagByPrimitive(f.Name, fieldtype, fieldptr, &tag); err != nil {
				return err
			}
		}
	}
	return nil
}

func registerFlagByValueInterface(fieldType reflect.Type, fieldPointer unsafe.Pointer, tag *flagTag) bool {
	// TODO does this implementation support all variants such as:
	//  -> interface
	var iface = reflect.NewAt(fieldType, fieldPointer).Interface()
	if value, ok := iface.(flag.Value); ok {
		// field type implements flag.Value interface, register as such
		// TODO (how to) set default value? (i.e. ignore default value?)
		flag.Var(value, tag.Name, tag.Description)
		return true
	}
	return false
}

func registerFlagByPrimitive(fieldName string, fieldType reflect.Type, fieldPtr unsafe.Pointer, tag *flagTag) error {
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

func parseTag(value string) flagTag {
	parts := strings.SplitN(value, ",", 3)
	for len(parts) < 3 {
		parts = append(parts, "")
	}
	return flagTag{Name: parts[0], DefaultValue: parts[1], Description: parts[2]}
}

type flagTag struct {
	Name         string
	DefaultValue string
	Description  string
}
