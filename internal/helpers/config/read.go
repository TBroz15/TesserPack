package config

import (
	"os"
	"path"
	"reflect"
	"tesserpack/internal/types"

	"dario.cat/mergo"
	"github.com/charmbracelet/log"
	"github.com/goccy/go-json"
	"github.com/titanous/json5"
)

// dear self do not optimize this further. i know the code is trash but this is not the main focus -tuxebro
func customMerge(dst, src interface{}, useMergo bool) error {
	if (useMergo) {
		err := mergo.Merge(dst, src, mergo.WithOverride, mergo.WithTypeCheck)
		if err != nil {
			return err
		}
	}
	
	dstVal := reflect.ValueOf(dst).Elem()
	srcVal := reflect.ValueOf(src)

	for i := 0; i < dstVal.NumField(); i++ {
		
		dstField := dstVal.Field(i)
		srcField := srcVal.Field(i)

		if dstField.Kind() == reflect.Struct {
			dstInter := dstField.Addr().Interface()
			srcInter := srcField.Interface()

			customMerge(dstInter, srcInter, false)
		}

		if dstField.Kind() == reflect.Bool {
			dstField.SetBool(srcField.Bool())
		}
	}

	return nil
}

func ReadConf(inPath, confPath string) types.TesserPackConfig {
	conf := NewDefault()

	// if user did not define confPath via CLI, check inside of pack
	if confPath == "" {
		confPath = path.Join(inPath, ".tesserpackrc.json5")
	}

	confFileCont, err := os.ReadFile(confPath)
	if err != nil && confPath != "" {
		return conf
	} else if err != nil {
		log.Fatal("Failed to read config.", "err", err)
	}

	messageFromTuxeBro := "That should not happen. If it does, report to me!"

	var confFileJSON types.TesserPackConfig
	err = json5.Unmarshal(confFileCont, &confFileJSON)
	if err != nil {
		log.Fatal("Failed to parse config.", "err", err, "messageFromTuxeBro", messageFromTuxeBro)
	}

	err = customMerge(&conf, confFileJSON, true)
	if err != nil {
		log.Fatal("Failed to merge default and user defined config.", "err", err, "messageFromTuxeBro", messageFromTuxeBro)
	}

	if (log.GetLevel() == log.DebugLevel) {
		confByte, _ := json.MarshalIndent(conf, "", "  ")

		log.Debug("Successfully parsed config.", "conf", string(confByte))
	}


	return conf
}