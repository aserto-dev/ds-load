package app

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/common/msg"
	"github.com/aserto-dev/ds-load/common/version"
	"github.com/aserto-dev/go-grpc/aserto/common/info/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

type InfoCmd struct{}

func (cmd *InfoCmd) Run(context *kong.Context) error {

	v := reflect.ValueOf(&ExecCmd{})

	configs, err := getConfigForStruct(v.Type())
	if err != nil {
		return err
	}

	ver := version.GetInfo()
	infoMsg := &msg.Info{
		Build: &info.BuildInfo{
			Version: ver.Version,
			Commit:  ver.Commit,
			Date:    ver.Date,
			Os:      ver.OS,
			Arch:    ver.Arch,
		},
		Description: AppDescription,
		Configs:     configs,
	}

	message, err := protojson.Marshal(infoMsg)
	if err != nil {
		return err
	}
	os.Stdout.Write(message)
	return nil
}

func getConfigForStruct(v reflect.Type) ([]*msg.ConfigElement, error) {
	var configs []*msg.ConfigElement

	// v := reflect.ValueOf(i)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Type.Kind() == reflect.Struct {
			cfgs, err := getConfigForStruct(field.Type)
			if err != nil {
				return nil, err
			}
			configs = append(configs, cfgs...)
			continue
		}
		tag := field.Tag
		optionalTag := tag.Get("optional")
		var optional bool
		var err error
		if optionalTag == "" {
			optional = false
		} else {
			optional, err = strconv.ParseBool(tag.Get("readonly"))
			if err != nil {
				return nil, err
			}
		}

		fieldName := tag.Get("name")
		help := tag.Get("help")

		var usage string
		var flagType msg.ConfigElementType
		switch field.Type.Name() {
		case "string":
			flagType = msg.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING
			usage = fmt.Sprintf("--%s=STRING", fieldName)
		case "int":
			flagType = msg.ConfigElementType_CONFIG_ELEMENT_TYPE_INTEGER
			usage = fmt.Sprintf("--%s=INT", fieldName)
		case "bool":
			flagType = msg.ConfigElementType_CONFIG_ELEMENT_TYPE_BOOLEAN
			usage = fmt.Sprintf("--%s", fieldName)
		default:
			flagType = msg.ConfigElementType_CONFIG_ELEMENT_TYPE_UNKNOWN
		}

		c := &msg.ConfigElement{
			Name:        fieldName,
			Type:        flagType,
			Description: help,
			Usage:       usage,
			Optional:    optional,
		}
		configs = append(configs, c)
	}
	return configs, nil
}
