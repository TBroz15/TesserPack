package config

import (
	"os"
	"path"
	"tesserpack/internal/types"

	"dario.cat/mergo"
	"github.com/charmbracelet/log"
	"github.com/titanous/json5"
)

func ReadConf(inPath string) types.TesserPackConfig {
	conf := NewDefault()

	confFilePath := path.Join(inPath, ".tesserpackrc")
	confFileCont, err := os.ReadFile(confFilePath)
	if err != nil {
		log.Warn("Can't read config, using default config.", "err", err)
		return conf
	}

	var confFileJSON types.TesserPackConfig
	err = json5.Unmarshal(confFileCont, confFileJSON)
	if err != nil {
		log.Error("Failed to parse config, using default config.", "err", err)
		return conf
	}

	err = mergo.Merge(conf, confFileJSON)
	if err != nil {
		log.Error("Failed to merge user and default config, using default config.", "err", err, "messageFromTuxeBro", "That should not happen.")
		return conf
	}

	return conf
}