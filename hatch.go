package hatch

import (
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type hatch struct {
	configName  string
	configType  string
	configPaths []string
}

type Hatchling interface {
	GetType() reflect.Type
}

func New() *hatch {
	return &hatch{}
}

func NewWithConfig(n string, t string, paths ...string) *hatch {
	return &hatch{
		configName:  n,
		configType:  t,
		configPaths: paths,
	}
}

func (h *hatch) SetName(s string) *hatch {
	h.configName = s
	return h
}

func (h *hatch) SetType(s string) *hatch {
	h.configType = s
	return h
}

func (h *hatch) AddPath(s string) *hatch {
	h.configPaths = append(h.configPaths, s)
	return h
}

func (h *hatch) AddPaths(s []string) *hatch {
	h.configPaths = append(h.configPaths, s...)
	return h
}

func (h *hatch) Unmarshal(i Hatchling) {
	var pathProvided bool
	if h.configName != "" {
		viper.SetConfigName(h.configName) // name of config file (without extension)
	}
	if h.configType != "" {
		viper.SetConfigType(h.configType)
	}
	for _, p := range h.configPaths {
		viper.AddConfigPath(p)
		pathProvided = true
	}
	if h.configName != "" && h.configType != "" && pathProvided {
		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			log.Printf("Hatch: config file: %s\n", err)
		}
	}

	processStruct(i.GetType(), "")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))

	err := viper.Unmarshal(i)
	if err != nil {
		log.Fatalf("Hatch: unable to decode into struct, %v\n", err)
	}
}

func processStruct(t reflect.Type, prefix string) {
	if t.Kind() != reflect.Struct {
		log.Fatalf("Hatch: calling processStruct on: %s, %s which is not Kind of Struct", t.Name(), t.Kind())
	}
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)

		k := sf.Tag.Get("mapstructure")
		if k == "" {
			k = sf.Name
		}
		k = prefix + k
		if sf.Type.Kind() == reflect.Struct {
			processStruct(sf.Type, k+".")
		} else {
			d := sf.Tag.Get("default")
			if d != "" {
				viper.SetDefault(k, d)
			}
			viper.BindEnv(k)
		}
		r, _ := strconv.ParseBool(sf.Tag.Get("required"))
		if r && !viper.IsSet(k) {
			log.Fatalf("Hatch: key: %s is not set\n", k)
		}
	}
}
