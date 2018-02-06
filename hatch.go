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
	v := viper.New()
	if h.configName != "" {
		v.SetConfigName(h.configName) // name of config file (without extension)
	}
	if h.configType != "" {
		v.SetConfigType(h.configType)
	}
	for _, p := range h.configPaths {
		v.AddConfigPath(p)
		pathProvided = true
	}
	if h.configName != "" && h.configType != "" && pathProvided {
		err := v.ReadInConfig() // Find and read the config file
		if err != nil {         // Handle errors reading the config file
			log.Printf("hatch: config file: %s\n", err)
		}
	}

	processStruct(v, i.GetType(), "")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))

	err := v.Unmarshal(i)
	if err != nil {
		log.Fatalf("hatch: unable to decode into struct, %v\n", err)
	}
}

func processStruct(v *viper.Viper, t reflect.Type, prefix string) {
	if t.Kind() != reflect.Struct {
		log.Fatalf("hatch: calling processStruct on: %s, %s which is not Kind of Struct\n", t.Name(), t.Kind())
	}
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		r, d := parseTag(sf)

		k := sf.Tag.Get("mapstructure")
		if k == "" {
			k = sf.Name
		}
		k = prefix + k
		if sf.Type.Kind() == reflect.Struct {
			processStruct(v, sf.Type, k+".")
		} else {
			if d != "" {
				v.SetDefault(k, d)
			}
			v.BindEnv(k)
		}
		if r && !v.IsSet(k) {
			log.Fatalf("hatch: required key: %s is not set\n", k)
		}
	}
}

func parseTag(sf reflect.StructField) (bool, string) {
	h := sf.Tag.Get("hatch")
	kv := strings.SplitN(h, ",", 2)
	t1, t2 := func() (string, string) {
		switch len(kv) {
		case 2:
			return kv[0], kv[1]
		case 1:
			return kv[0], ""
		default:
			return "", ""

		}
	}()
	r, _ := strconv.ParseBool(t1)
	return r, t2
}
