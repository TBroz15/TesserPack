package types

type JSONConfig struct {
	Minify bool	`json:"minify"`
	Strict bool	`json:"strict"`
	Cache  bool	`json:"cache"`
}

type LANGConfig struct {
	Minify bool	`json:"minify"`
	Cache  bool	`json:"cache"`
}

type PNGConfig struct {
	Optimize 	  bool `json:"optimize"`
	Progressive   bool `json:"progressive"`
	CompressLevel byte `json:"compressLevel"`
	Quality		  byte `json:"quality"`
}

type JPGConfig struct {
	Optimize 	  bool `json:"optimize"`
	Progressive   bool `json:"progressive"`
	Quality		  byte `json:"quality"`
}

type CompilerConfig struct {
	JSON JSONConfig `json:"json"`
	LANG LANGConfig `json:"lang"`
	PNG  PNGConfig  `json:"png"`
	JPG  JPGConfig  `json:"jpg"`
}

type IgnoreConfig struct {
	Dirs  []string `json:"dirs"`
	Exts  []string `json:"exts"`
	Files []string `json:"files"`
}

type TesserPackConfig struct {
	Compiler CompilerConfig `json:"compiler"`
	Ignore   IgnoreConfig   `json:"ignore"`
}