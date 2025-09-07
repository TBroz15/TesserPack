package config

import "tesserpack/internal/types"

func NewDefault() types.TesserPackConfig {
	return types.TesserPackConfig{
		Compiler: types.CompilerConfig{
			JSON: types.JSONConfig{
				Strict: false,
			},
			PNG: types.PNGConfig{
				Q: 100,
				Compression: 9,
				Effort: 10,
			},
			JPG: types.JPGConfig{
				Q: 100,
			},
			Cache: true,
		},
		IgnoreGlob: []string{
			"node_modules/",
			".git/",
			".vscode/",
			".github/",
		},
	}
}