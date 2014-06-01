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
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}
	ptr := reflect.ValueOf(config)
	if ptr.IsNil() {
		return fmt.Errorf("config cannot point to nil")
	}
	val := reflect.Indirect(ptr)
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("config instance is not a struct")
	}
	structType := val.Type()
	for i := 0; i < structType.NumField(); i++ {
		f := (reflect.StructField)(structType.Field(i))
		t := f.Tag.Get("flag")
		if t == "" {
			// TODO no tag available
			fmt.Printf("Skipping '%s' ...\n", f.Name)
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
				// TODO invalid default, skipping
				fmt.Printf("Invalid value, skipping %s ...\n", f.Name)
				continue
			}
			flag.BoolVar((*bool)(fieldptr), tag.Name, defaultVal, tag.Description)
		case reflect.Float64:
			defaultVal, err := strconv.ParseFloat(tag.DefaultValue, f.Type.Bits())
			if err != nil {
				// TODO invalid default, skipping
				fmt.Printf("Invalid value, skipping %s ...\n", f.Name)
				continue
			}
			flag.Float64Var((*float64)(fieldptr), tag.Name, defaultVal, tag.Description)
		case reflect.Int:
			defaultVal, err := strconv.ParseInt(tag.DefaultValue, 10, f.Type.Bits())
			// TODO distinguish between oct (0o1234567), dec (1234), hex (0x0123456789abcdef)
			if err != nil {
				// TODO invalid default, skipping
				fmt.Printf("Invalid value, skipping %s ...\n", f.Name)
				continue
			}
			flag.IntVar((*int)(fieldptr), tag.Name, int(defaultVal), tag.Description)
		case reflect.Int64:
			defaultVal, err := strconv.ParseInt(tag.DefaultValue, 10, f.Type.Bits())
			// TODO distinguish between oct (0o1234567), dec (1234), hex (0x0123456789abcdef)
			if err != nil {
				// TODO invalid default, skipping
				fmt.Printf("Invalid value, skipping %s ...\n", f.Name)
				continue
			}
			flag.Int64Var((*int64)(fieldptr), tag.Name, defaultVal, tag.Description)
		case reflect.Uint:
			defaultVal, err := strconv.ParseUint(tag.DefaultValue, 10, f.Type.Bits())
			// TODO distinguish between oct (0o1234567), dec (1234), hex (0x0123456789abcdef)
			if err != nil {
				// TODO invalid default, skipping
				fmt.Printf("Invalid value, skipping %s ...\n", f.Name)
				continue
			}
			flag.UintVar((*uint)(fieldptr), tag.Name, uint(defaultVal), tag.Description)
		case reflect.Uint64:
			defaultVal, err := strconv.ParseUint(tag.DefaultValue, 10, 64)
			if err != nil {
				// TODO invalid default, skipping
				fmt.Printf("Invalid value, skipping %s ...\n", f.Name)
				continue
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
