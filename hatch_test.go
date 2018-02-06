package hatch

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type ignored struct {
	A string
}
type grandchild struct {
	A string
	B string `required:"true"`
	C string `required:"true"`
}
type child struct {
	A int
	B float64
	C string
	D []string
	E time.Duration
	F time.Duration
	G bool
	H bool
	I grandchild
}
type parent struct {
	Child       child
	Default     int    `default:"1234"`
	Required    string `required:"true"`
	CamelCase   string `mapstructure:"camel_case"`
	Ignored     *ignored
	Prefilled   string
	Overwritten string
}

func (c *parent) GetType() reflect.Type {
	return reflect.TypeOf(*c)
}

func TestHatch(t *testing.T) {
	os.Setenv("REQUIRED", "required")
	os.Setenv("CAMEL_CASE", "camel_case")
	os.Setenv("OVERWRITTEN", "overwritten")
	os.Setenv("CHILD__A", "1")
	os.Setenv("CHILD__B", "1.23")
	os.Setenv("CHILD__C", "string")
	os.Setenv("CHILD__D", "string1,string2,string3")
	os.Setenv("CHILD__E", "1h2m3s4ms5us6ns")
	os.Setenv("CHILD__F", "3723.004005006s")
	os.Setenv("CHILD__G", "true")
	os.Setenv("CHILD__H", "1")
	os.Setenv("CHILD__I__A", "grandchild")

	dur, _ := time.ParseDuration("1h2m3s4ms5us6ns")
	expected := &parent{
		Default:   1234,
		Required:  "required",
		CamelCase: "camel_case",
		Child: child{
			A: 1,
			B: 1.23,
			C: "string",
			D: []string{"string1", "string2", "string3"},
			E: dur,
			F: dur,
			G: true,
			H: true,
			I: grandchild{A: "grandchild", B: "b_from_yaml", C: "c_from_yaml"},
		},
		Ignored:     &ignored{A: "ignored"},
		Prefilled:   "prefilled",
		Overwritten: "overwritten",
	}

	actual := &parent{Ignored: &ignored{A: "ignored"},
		Prefilled:   "prefilled",
		Overwritten: "should_be_overwritten",
	}
	New().SetName("test").SetType("yaml").AddPath(".").Unmarshal(actual)
	assert.Equal(t, expected, actual)

	actual = &parent{Ignored: &ignored{A: "ignored"},
		Prefilled:   "prefilled",
		Overwritten: "should_be_overwritten",
	}
	NewWithConfig("test", "yaml", ".").Unmarshal(actual)
	assert.Equal(t, expected, actual)
}
