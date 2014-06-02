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
	val, err := checkType(config)
	if err != nil {
		return err
	}
	structType := val.Type()
	for i := 0; i < structType.NumField(); i++ {
		f := (reflect.StructField)(structType.Field(i))
		t := f.Tag.Get("flag")
		if t == "" {
			continue
		}
		tag := parseTag(t)
		if tag.Name == "" {
			return fmt.Errorf("invalid flag name: empty string")
		}
		var fieldptr = unsafe.Pointer(val.UnsafeAddr() + f.Offset)
		// TODO support Duration
		// TODO support Var (any variable via flag.Value interface)
		// TODO support nested structs.
		// TODO support smaller int, uint, float types? (how to handle overflow?)
		switch f.Type.Kind() {
		case reflect.String:
			flag.StringVar((*string)(fieldptr), tag.Name, tag.DefaultValue, tag.Description)
		case reflect.Bool:
			defaultVal, err := strconv.ParseBool(tag.DefaultValue)
			if err != nil {
				return fmt.Errorf("invalid default value for field '%s': %s", f.Name, err.Error())
			}
			flag.BoolVar((*bool)(fieldptr), tag.Name, defaultVal, tag.Description)
		case reflect.Float64:
			defaultVal, err := strconv.ParseFloat(tag.DefaultValue, 64)
			if err != nil {
				return fmt.Errorf("invalid default value for field '%s': %s", f.Name, err.Error())
			}
			flag.Float64Var((*float64)(fieldptr), tag.Name, defaultVal, tag.Description)
		case reflect.Int:
			// TODO parse exact number of available bits, or always 64?
			defaultVal, err := strconv.ParseInt(tag.DefaultValue, 0, f.Type.Bits())
			if err != nil {
				return fmt.Errorf("invalid default value for field '%s': %s", f.Name, err.Error())
			}
			flag.IntVar((*int)(fieldptr), tag.Name, int(defaultVal), tag.Description)
		case reflect.Int64:
			defaultVal, err := strconv.ParseInt(tag.DefaultValue, 0, 64)
			if err != nil {
				return fmt.Errorf("invalid default value for field '%s': %s", f.Name, err.Error())
			}
			flag.Int64Var((*int64)(fieldptr), tag.Name, defaultVal, tag.Description)
		case reflect.Uint:
			// TODO parse exact number of available bits, or always 64?
			defaultVal, err := strconv.ParseUint(tag.DefaultValue, 0, f.Type.Bits())
			if err != nil {
				return fmt.Errorf("invalid default value for field '%s': %s", f.Name, err.Error())
			}
			flag.UintVar((*uint)(fieldptr), tag.Name, uint(defaultVal), tag.Description)
		case reflect.Uint64:
			defaultVal, err := strconv.ParseUint(tag.DefaultValue, 0, 64)
			if err != nil {
				return fmt.Errorf("invalid default value for field '%s': %s", f.Name, err.Error())
			}
			flag.Uint64Var((*uint64)(fieldptr), tag.Name, defaultVal, tag.Description)
		default:
			return fmt.Errorf("unsupported data type for field '%s'", f.Name)
		}
	}
	return nil
}

func checkType(config interface{}) (reflect.Value, error) {
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
