package flagtag

import (
	"flag"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestConfigureNil(t *testing.T) {
	if Configure(nil) == nil {
		t.Fatal("Expected an error, since nil cannot be parsed.")
	}
}

func TestConfigureUntaggedStruct(t *testing.T) {
	var s = struct {
		A int
		B uint
		C string
	}{}
	if err := Configure(&s); err != nil {
		t.Fatal("Unexpected error: " + err.Error())
	}
}

func TestConfigureNilPointer(t *testing.T) {
	var c *struct{}
	if Configure(c) == nil {
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
		V string
	}{}
	if Configure(&s) != nil {
		t.Fatal("Expected correct processing of struct without any tags.")
	}
}

func TestNonRelevantlyTaggedStruct(t *testing.T) {
	var s = struct {
		V string `json:"some,value"`
	}{}
	if Configure(&s) != nil {
		t.Fatal("Expected correct processing of struct without any relevant tags.")
	}
}

func TestUnexportedTaggedField(t *testing.T) {
	var s = struct {
		v string `flag:"unexported,unexported,Tagging an unexported value."`
	}{}
	if Configure(&s) == nil {
		t.Fatal("Expected error regarding unexported field, but got nothing.")
	}
}

func TestEmptyTagName(t *testing.T) {
	var s = struct {
		V string `flag:",,"`
	}{}
	if Configure(&s) == nil {
		t.Fatal("Expected error because no flag name was specified.")
	}
}

func TestIncompleteTag(t *testing.T) {
	var s = struct {
		V string `flag:"a"`
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
		V string `flag:"s,hello world,This sets the string flag."`
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
		V bool `flag:"b,true,This sets the bool flag."`
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
		V bool `flag:"b2,foo,This sets the bool flag."`
	}{}
	if Configure(&s) == nil {
		t.Fatal("Expected error due to incorrect default value.")
	}
}

func TestTagFloat64(t *testing.T) {
	var s = struct {
		V float64 `flag:"f,0.2345,This sets the float flag."`
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
		V float64 `flag:"f2,abcde,This sets the float64 flag."`
	}{}
	if Configure(&s) == nil {
		t.Fatal("Expected error due to incorrect default value.")
	}
}

func TestTagInt(t *testing.T) {
	var s = struct {
		V int `flag:"i,64,This sets the int flag."`
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
		V int `flag:"i2,0.33333,This sets the int flag."`
	}{}
	if Configure(&s) == nil {
		t.Fatal("Expected error due to incorrect default value.")
	}
}

func TestTagInt64(t *testing.T) {
	var s = struct {
		V int64 `flag:"i64,-6400000000,This sets the int64 flag."`
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
		V int64 `flag:"i64-2,abcdefgh,This sets the int64 flag."`
	}{}
	if Configure(&s) == nil {
		t.Fatal("Expected error due to incorrect default value.")
	}
}

func TestTagUint(t *testing.T) {
	var s = struct {
		V uint `flag:"u,3200,This sets the uint flag."`
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
		V uint `flag:"u2,-200,This sets the uint flag."`
	}{}
	if Configure(&s) == nil {
		t.Fatal("Expected error due to incorrect default value.")
	}
}

func TestTagUint64(t *testing.T) {
	var s = struct {
		V uint64 `flag:"u64,6400000000,This sets the uint64 flag."`
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
		V uint64 `flag:"u64-2,abcdefgh,This sets the uint64 flag."`
	}{}
	if Configure(&s) == nil {
		t.Fatal("Expected error due to incorrect default value.")
	}
}

func TestMustConfigure(t *testing.T) {
	defer func() {
		// We are supposed to get a panic, so silence it.
		recover()
	}()
	var s = struct {
		X bool `flag:"x,test,test"`
	}{}
	MustConfigure(&s)
	t.FailNow()
}

func TestConfigureAndParse(t *testing.T) {
	var s = struct {
		X string `flag:"xx,test,Test 1."`
		Y bool   `flag:"y,f,Test 2."`
	}{}
	if err := ConfigureAndParse(&s); err != nil {
		t.Fatal("Did not expect error: " + err.Error())
	}
	if flag.Parsed() == false {
		t.Fatal("Expected command line flags to be parsed by now.")
	}
}

func TestConfigureAndParseFaulty(t *testing.T) {
	var s = struct {
		Y bool `flag:"y,bla,Test 2."`
	}{}
	if err := ConfigureAndParse(&s); err == nil {
		t.Fatal("Expected an error but got nothing.")
	}
}

func TestMustConfigureAndParseFailing(t *testing.T) {
	defer func() {
		// We are supposed to get a panic, so silence it.
		recover()
	}()
	var s = struct {
		X bool `flag:"xxx,test,test"`
	}{}
	MustConfigureAndParse(&s)
	t.FailNow()
}

func TestMustConfigureAndParseSuccessfully(t *testing.T) {
	var s = struct {
		X bool `flag:"xxxx,True,test"`
	}{}
	MustConfigureAndParse(&s)
	if !flag.Parsed() {
		t.Fatal("Expected an command line flags to be parsed by now.")
	}
}

func TestErrorOnInvalidDataType(t *testing.T) {
	var s = struct {
		Invalid uintptr `flag:"xxxxxx,,"`
	}{}
	if err := Configure(&s); err == nil {
		t.Fatal("Expected error because of unsupported data type.")
	}
}

func TestRecursiveStructProcessing(t *testing.T) {
	var outer = struct {
		Inner struct {
			V int `flag:"innerv,1"`
		}
	}{}
	err := Configure(&outer)
	if err != nil {
		t.Fatal("Unexpected error: " + err.Error())
	}
	f := flag.Lookup("innerv")
	if f == nil {
		t.Fatal("Could not find configured flag.")
	}
	if f.Name != "innerv" {
		t.Error("Configured flag has incorrect name.")
	}
	if f.DefValue != "1" {
		t.Error("Configured flag has incorrect default value.")
	}
	if f.Usage != "" {
		t.Error("Configured flag has incorrect usage description.")
	}
}

func TestBadInnerStruct(t *testing.T) {
	var outer = struct {
		Inner struct {
			V uint `flag:"innerv,-1"`
		}
	}{}
	err := Configure(&outer)
	if err == nil {
		t.Fatal("Expected error because of invalid default value.")
	}
}

func TestMixedInnerStructProcessing(t *testing.T) {
	var outer = struct {
		Before uint `flag:"outerBefore,3,some description"`
		Blank  uint
		Inner  struct {
			Dummy  int
			Inside string `flag:"innerInside,2,inside information"`
		}
		After int `flag:"outerAfter,1,final remark"`
	}{}
	err := Configure(&outer)
	if err != nil {
		t.Fatal("Unexpected error: " + err.Error())
	}
	flagBefore := flag.Lookup("outerBefore")
	if flagBefore.Name != "outerBefore" || flagBefore.DefValue != "3" || flagBefore.Usage != "some description" {
		t.Error("Flag outerBefore data is invalid.")
	}
	flagInside := flag.Lookup("innerInside")
	if flagInside.Name != "innerInside" || flagInside.DefValue != "2" || flagInside.Usage != "inside information" {
		t.Error("Flag innerInside data is invalid.")
	}
	flagAfter := flag.Lookup("outerAfter")
	if flagAfter.Name != "outerAfter" || flagAfter.DefValue != "1" || flagAfter.Usage != "final remark" {
		t.Error("Flag outerAfter data is invalid.")
	}
}

func TestRegisterTypeDerivedFromPrimitive(t *testing.T) {
	var s = struct {
		D aliasInt `flag:"flagValueAliasInt,-10,Alias of int, still works as primitive int flag."`
	}{}
	Configure(&s)
	flagAlias := flag.Lookup("flagValueAliasInt")
	if flagAlias == nil {
		t.Fatal("Could not find defined flagValueAliasInt.")
	}
	if flagAlias.Name != "flagValueAliasInt" || flagAlias.DefValue != "-10" || flagAlias.Usage != "Alias of int, still works as primitive int flag." {
		t.Error("Flag flagValueAliasInt data is invalid.")
	}
}

type aliasInt int

func TestRegisterValueInterfaceFlag(t *testing.T) {
	var s = struct {
		D dummyInt `flag:"flagValueDummyInt,,My first flag.Value implementation."`
	}{}
	err := Configure(&s)
	if err != nil {
		t.Fatal("Unexpected error: " + err.Error())
	}
	flagDummyInt := flag.Lookup("flagValueDummyInt")
	if flagDummyInt == nil {
		t.Fatal("Expected a flag, but got nil.")
	}
	if flagDummyInt.Name != "flagValueDummyInt" || flagDummyInt.DefValue != "0" || flagDummyInt.Usage != "My first flag.Value implementation." {
		t.Fatal("Flag data is invalid.")
	}
}

func TestRegisterValueInterfaceFlagNilPointer(t *testing.T) {
	var s = struct {
		D *dummyInt `flag:"flagValueDummyIntNilPointer,,My first flag.Value implementation."`
	}{}
	err := Configure(&s)
	if err == nil {
		t.Fatal("Expected an error since the pointer is nil, but didn't get anything.")
	}
}

func TestRegisterValueInterfaceFlagPointer(t *testing.T) {
	var s = struct {
		D *dummyInt `flag:"flagValueDummyIntPointer,,My first flag.Value implementation."`
	}{D: new(dummyInt)}
	err := Configure(&s)
	if err != nil {
		t.Fatal("Unexpected error: " + err.Error())
	}
	flagDummyInt := flag.Lookup("flagValueDummyIntPointer")
	if flagDummyInt == nil {
		t.Fatal("Expected a flag, but got nil.")
	}
	if flagDummyInt.Name != "flagValueDummyIntPointer" || flagDummyInt.DefValue != "0" || flagDummyInt.Usage != "My first flag.Value implementation." {
		t.Fatal("Flag data is invalid.")
	}
}

func TestRegisterPrimitiveFlagNilPointer(t *testing.T) {
	var s = struct {
		D *int `flag:"flagValueIntPointer,123,My first primitive pointer flag."`
	}{}
	err := Configure(&s)
	if err == nil {
		t.Fatal("Expected an error but got nothing.")
	}
}

func TestRegisterPrimitiveFlagPointer(t *testing.T) {
	var s = struct {
		D *int `flag:"flagValueIntPointer,123,My first primitive pointer flag."`
	}{D: new(int)}
	*s.D = 123
	err := Configure(&s)
	if err != nil {
		t.Fatal("Unexpected error: " + err.Error())
	}
	flagInt := flag.Lookup("flagValueIntPointer")
	if flagInt == nil {
		t.Fatal("Expected a flag, but got nil.")
	}
	if flagInt.Name != "flagValueIntPointer" || flagInt.DefValue != "123" || flagInt.Usage != "My first primitive pointer flag." {
		t.Fatal("Flag data is invalid.")
	}
}

func TestRegisterFlagValueThroughInterfaceNil(t *testing.T) {
	var s = struct {
		D interface{} `flag:"flagValueUnknownInterface,,Value which turns out is flag.Value compatible."`
	}{}
	err := Configure(&s)
	if err == nil {
		t.Fatal("Expected error because D is nil.")
	}
}

func TestRegisterFlagValueThroughInterfaceValueNilPointer(t *testing.T) {
	var v *dummyInt
	var s = struct {
		D interface{} `flag:"flagValueUnknownInterface,,Value which turns out is flag.Value compatible."`
	}{D: v}
	err := Configure(&s)
	if err == nil {
		t.Fatal("Expected error because D is nil.")
	}
}

func TestRegisterFlagValueThroughInterfaceValueZeroValue(t *testing.T) {
	var v dummyInt
	var s = struct {
		D interface{} `flag:"flagValueUnknownInterface,,Value which turns out is flag.Value compatible."`
	}{D: v}
	err := Configure(&s)
	if err == nil {
		t.Fatal("Expected error because D is nil.")
	}
}

func TestRegisterFlagValueThroughInterface(t *testing.T) {
	var s = struct {
		D interface{} `flag:"flagValueUnknownInterface,,Value which turns out is flag.Value compatible."`
	}{D: new(dummyInt)}
	err := Configure(&s)
	if err != nil {
		t.Fatal("Unexpected error: " + err.Error())
	}
	flagUnknown := flag.Lookup("flagValueUnknownInterface")
	if flagUnknown == nil {
		t.Fatal("Expected a flag, but got nil.")
	}
	if flagUnknown.Name != "flagValueUnknownInterface" || flagUnknown.DefValue != "0" || flagUnknown.Usage != "Value which turns out is flag.Value compatible." {
		t.Fatal("Flag data is invalid.")
	}
}

type dummyInt int

func (d *dummyInt) String() string {
	return strconv.Itoa(int(*d))
}

func (d *dummyInt) Set(value string) error {
	i, err := strconv.Atoi(value)
	if err == nil {
		*d = dummyInt(i)
	}
	return err
}

func TestRegisterDurationNilPointer(t *testing.T) {
	var s = struct {
		D *time.Duration `flag:"flagDuration,1h,Specify duration"`
	}{}
	err := Configure(&s)
	if err == nil {
		t.Fatal("Expected an error because of nil pointer.")
	}
}

func TestRegisterDuration(t *testing.T) {
	var s = struct {
		D time.Duration `flag:"flagDuration,1h,Specify duration"`
	}{}
	err := Configure(&s)
	if err != nil {
		t.Fatal("Unexpected error: " + err.Error())
	}
	flagInterface := flag.Lookup("flagDuration")
	if flagInterface == nil {
		t.Fatal("Expected a flag, but got nil.")
	}
	if flagInterface.Name != "flagDuration" || flagInterface.DefValue != "1h0m0s" || flagInterface.Usage != "Specify duration" {
		t.Fatal("Flag data is invalid.")
	}
}

func TestRegisterDurationBadDefault(t *testing.T) {
	var s = struct {
		D time.Duration `flag:"flagDuration,1abcde,Specify duration"`
	}{}
	err := Configure(&s)
	if err == nil {
		t.Fatal("Expected error because of bad defaults.")
	}
}

func TestRegisterDurationPointer(t *testing.T) {
	var s = struct {
		D *time.Duration `flag:"flagDurationPointer,1h,Specify duration"`
	}{D: new(time.Duration)}
	err := Configure(&s)
	if err != nil {
		t.Fatal("Unexpected error: " + err.Error())
	}
	flagInterface := flag.Lookup("flagDurationPointer")
	if flagInterface == nil {
		t.Fatal("Expected a flag, but got nil.")
	}
	if flagInterface.Name != "flagDurationPointer" || flagInterface.DefValue != "1h0m0s" || flagInterface.Usage != "Specify duration" {
		t.Fatal("Flag data is invalid.")
	}
}

func TestErrInvalidDefault(t *testing.T) {
	var s = struct {
		D int `flag:"flagInvalidDefault,abcde,Test invalid defaults..."`
	}{}
	err := Configure(&s)
	if err == nil {
		t.Fatal("Error was expected but got nothing.")
	}
	if !strings.HasPrefix(err.Error(), "invalid default value for field 'D' (tag 'flagInvalidDefault'): ") {
		t.Fatal("Expected a different error message than was provided.")
	}
}
