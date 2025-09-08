package config

import (
	"os"
	"path"
	"tesserpack/internal/types"

	"dario.cat/mergo"
	"github.com/charmbracelet/log"
	"github.com/titanous/json5"
)

func ReadConf(inPath, confPath string) types.TesserPackConfig {
	conf := NewDefault()

	// if user did not define confPath via CLI
	if confPath == "" {
		confPath = path.Join(inPath, ".tesserpackrc.json5")
	}

	confFileCont, err := os.ReadFile(confPath)
	if err != nil {
		log.Warn("Can't read config, using default config.", "err", err)
		return conf
	}

	var confFileJSON types.TesserPackConfig
	err = json5.Unmarshal(confFileCont, confFileJSON)
	if err != nil {
		log.Fatal("Failed to parse config.", "err", err)
	}

	err = mergo.Merge(conf, confFileJSON)
	if err != nil {
		log.Fatal("Failed to merge default and user defined config.", "err", err, "messageFromTuxeBro", "That should not happen.")
	}

	return conf
}