package types

type JSONConfig struct {
	Strict bool	`json:"strict"`
}

type PNGConfig struct {
	Progressive   bool `json:"progressive"`
	CompressLevel byte `json:"compressLevel"`
	Quality		  byte `json:"quality"`
}

type JPGConfig struct {
	Progressive   bool `json:"progressive"`
	Quality		  byte `json:"quality"`
}

type CompilerConfig struct {
	JSON JSONConfig `json:"json"`
	PNG  PNGConfig  `json:"png"`
	JPG  JPGConfig  `json:"jpg"`

	Cache bool `json:"cache"`
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