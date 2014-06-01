package flag

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

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
		// TODO support Duration
		// TODO support Var (any variable via flag.Value interface)
		// TODO support for smaller int, uint, float types?
		// TODO how to handle unsupported types (return error, leave message, etc.)
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
