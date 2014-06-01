package flagtag

import (
	"flag"
	"testing"
)

func TestConfigureNil(t *testing.T) {
	if Configure(nil) == nil {
		t.Fatal("Expected an error, since nil cannot be parsed.")
	}
}

func TestConfigureNilPointer(t *testing.T) {
	var c *struct{}
	var p interface{}
	p = c
	if Configure(p) == nil {
		t.Fatal("Expected an error, since pointer is nil.")
	}
}

func TestConfigureNonStructPointer(t *testing.T) {
	var i = 42
	if Configure(&i) == nil {
		t.Fatal("Expected an error, since config value is not a struct.")
	}
}

func TestEmptyStruct(t *testing.T) {
	var s = struct{}{}
	if Configure(&s) != nil {
		t.Fatal("Expected correct processing of empty struct.")
	}
}

func TestNonTaggedStruct(t *testing.T) {
	var s = struct {
		v string
	}{}
	if Configure(&s) != nil {
		t.Fatal("Expected correct processing of struct without any tags.")
	}
}

func TestNonRelevantlyTaggedStruct(t *testing.T) {
	var s = struct {
		v string `json:"some,value"`
	}{}
	if Configure(&s) != nil {
		t.Fatal("Expected correct processing of struct without any relevant tags.")
	}
}

func TestEmptyTagName(t *testing.T) {
	var s = struct {
		v string `flag:",,"`
	}{}
	if Configure(&s) == nil {
		t.Fatal("Expected error because no flag name was specified.")
	}
}

func TestIncompleteTag(t *testing.T) {
	var s = struct {
		v string `flag:"a"`
	}{}
	if Configure(&s) != nil {
		t.Fatal("Expected correct processing even though tag does not contain all parts.")
	}
	tag := flag.Lookup("a")
	if tag == nil {
		t.Fatal("Cannot find configured flag.")
	}
	if tag.Name != "a" {
		t.Error("Expected another tag name.")
	}
	if tag.DefValue != "" {
		t.Error("Expected empty string as default value, since we didn't specify any.")
	}
	if tag.Usage != "" {
		t.Error("Expected empty string as usage information, since we didn't specify any.")
	}
}

func TestTagString(t *testing.T) {
	var s = struct {
		v string `flag:"s,hello world,This sets the string flag."`
	}{}
	if Configure(&s) != nil {
		t.Fatal("Expected correct configuration without any errors.")
	}
	f := flag.Lookup("s")
	if f == nil {
		t.Fatal("Could not find configured flag.")
	}
	if f.Name != "s" {
		t.Error("Configured flag has incorrect name.")
	}
	if f.DefValue != "hello world" {
		t.Error("Configured flag has incorrect default value.")
	}
	if f.Usage != "This sets the string flag." {
		t.Error("Configured flag has incorrect usage description.")
	}
}

func TestTagBool(t *testing.T) {
	var s = struct {
		v bool `flag:"b,true,This sets the bool flag."`
	}{}
	if Configure(&s) != nil {
		t.Fatal("Expected correct configuration without any errors.")
	}
	f := flag.Lookup("b")
	if f == nil {
		t.Fatal("Could not find configured flag.")
	}
	if f.Name != "b" {
		t.Error("Configured flag has incorrect name.")
	}
	if f.DefValue != "true" {
		t.Error("Configured flag has incorrect default value.")
	}
	if f.Usage != "This sets the bool flag." {
		t.Error("Configured flag has incorrect usage description.")
	}
}

func TestTagBoolInvalidDefault(t *testing.T) {
	var s = struct {
		v bool `flag:"b2,foo,This sets the bool flag."`
	}{}
	if Configure(&s) == nil {
		t.Fatal("Expected error due to incorrect default value.")
	}
}

func TestTagFloat64(t *testing.T) {
	var s = struct {
		v float64 `flag:"f,0.2345,This sets the float flag."`
	}{}
	if Configure(&s) != nil {
		t.Fatal("Expected correct configuration without any errors.")
	}
	f := flag.Lookup("f")
	if f == nil {
		t.Fatal("Could not find configured flag.")
	}
	if f.Name != "f" {
		t.Error("Configured flag has incorrect name.")
	}
	if f.DefValue != "0.2345" {
		t.Error("Configured flag has incorrect default value.")
	}
	if f.Usage != "This sets the float flag." {
		t.Error("Configured flag has incorrect usage description.")
	}
}

func TestTagFloat64InvalidDefault(t *testing.T) {
	var s = struct {
		v float64 `flag:"f2,abcde,This sets the float64 flag."`
	}{}
	if Configure(&s) == nil {
		t.Fatal("Expected error due to incorrect default value.")
	}
}

func TestTagInt(t *testing.T) {
	var s = struct {
		v int `flag:"i,64,This sets the int flag."`
	}{}
	if Configure(&s) != nil {
		t.Fatal("Expected correct configuration without any errors.")
	}
	f := flag.Lookup("i")
	if f == nil {
		t.Fatal("Could not find configured flag.")
	}
	if f.Name != "i" {
		t.Error("Configured flag has incorrect name.")
	}
	if f.DefValue != "64" {
		t.Error("Configured flag has incorrect default value.")
	}
	if f.Usage != "This sets the int flag." {
		t.Error("Configured flag has incorrect usage description.")
	}
}

func TestTagIntInvalidDefault(t *testing.T) {
	var s = struct {
		v int `flag:"i2,0.33333,This sets the int flag."`
	}{}
	if Configure(&s) == nil {
		t.Fatal("Expected error due to incorrect default value.")
	}
}

func TestTagInt64(t *testing.T) {
	var s = struct {
		v int64 `flag:"i64,-6400000000,This sets the int64 flag."`
	}{}
	if Configure(&s) != nil {
		t.Fatal("Expected correct configuration without any errors.")
	}
	f := flag.Lookup("i64")
	if f == nil {
		t.Fatal("Could not find configured flag.")
	}
	if f.Name != "i64" {
		t.Error("Configured flag has incorrect name.")
	}
	if f.DefValue != "-6400000000" {
		t.Error("Configured flag has incorrect default value.")
	}
	if f.Usage != "This sets the int64 flag." {
		t.Error("Configured flag has incorrect usage description.")
	}
}

func TestTagInt64InvalidDefault(t *testing.T) {
	var s = struct {
		v int64 `flag:"i64-2,abcdefgh,This sets the int64 flag."`
	}{}
	if Configure(&s) == nil {
		t.Fatal("Expected error due to incorrect default value.")
	}
}

func TestTagUint(t *testing.T) {
	var s = struct {
		v uint `flag:"u,3200,This sets the uint flag."`
	}{}
	if Configure(&s) != nil {
		t.Fatal("Expected correct configuration without any errors.")
	}
	f := flag.Lookup("u")
	if f == nil {
		t.Fatal("Could not find configured flag.")
	}
	if f.Name != "u" {
		t.Error("Configured flag has incorrect name.")
	}
	if f.DefValue != "3200" {
		t.Error("Configured flag has incorrect default value.")
	}
	if f.Usage != "This sets the uint flag." {
		t.Error("Configured flag has incorrect usage description.")
	}
}

func TestTagUintInvalidDefault(t *testing.T) {
	var s = struct {
		v uint `flag:"u2,-200,This sets the uint flag."`
	}{}
	if Configure(&s) == nil {
		t.Fatal("Expected error due to incorrect default value.")
	}
}

func TestTagUint64(t *testing.T) {
	var s = struct {
		v uint64 `flag:"u64,6400000000,This sets the uint64 flag."`
	}{}
	if Configure(&s) != nil {
		t.Fatal("Expected correct configuration without any errors.")
	}
	f := flag.Lookup("u64")
	if f == nil {
		t.Fatal("Could not find configured flag.")
	}
	if f.Name != "u64" {
		t.Error("Configured flag has incorrect name.")
	}
	if f.DefValue != "6400000000" {
		t.Error("Configured flag has incorrect default value.")
	}
	if f.Usage != "This sets the uint64 flag." {
		t.Error("Configured flag has incorrect usage description.")
	}
}

func TestTagUint64InvalidDefault(t *testing.T) {
	var s = struct {
		v uint64 `flag:"u64-2,abcdefgh,This sets the uint64 flag."`
	}{}
	if Configure(&s) == nil {
		t.Fatal("Expected error due to incorrect default value.")
	}
}
