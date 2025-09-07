package types

type JSONConfig struct {
	Strict bool	`json:"strict"`
}

type PNGConfig struct {
	Compression   byte `json:"compressLevel"`
	Q			  byte `json:"quality"`
	Effort        byte `json:"effort"`
}

type JPGConfig struct {
	Q			  byte `json:"quality"`
}

type CompilerConfig struct {
	JSON JSONConfig `json:"json"`
	PNG  PNGConfig  `json:"png"`
	JPG  JPGConfig  `json:"jpg"`

	Cache bool `json:"cache"`
}

type TesserPackConfig struct {
	Compiler   CompilerConfig `json:"compiler"`
	IgnoreGlob []string		  `json:"ignoreGlob"`
}