package main

import (
	"errors"
	"fmt"
	"strconv"
	"tesserpack/internal/types"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/goccy/go-json"
)

func safeIntInput[T int | uint | byte](title, description string, value *T, min, max int) error {
	val := strconv.Itoa(int(*value))

	err := huh.NewInput().
		Title(title).
		Description(fmt.Sprintf("%v\nNumbers from %v to %v is only allowed. ", description, min, max)).
		Value(&val).
		Validate(func(s string) error {
			num, err := strconv.Atoi(s)
			if (err != nil) {
				return errors.New("invalid number")
			}

			if (num < min || num > max) {
				return errors.New("number is outside the range")
			}

			return nil
		}).
		Run()

	if (err != nil) {return err}

	num, err := strconv.Atoi(val)
	if (err != nil) {return err}

	*value = T(num)

	return nil
}

func ConfigGen() {
	conf := &types.TesserPackConfig{
		Compiler: types.CompilerConfig{
			JSON: types.JSONConfig{
				Strict: true,
			},
			PNG: types.PNGConfig{
				Quality: 100,
				CompressLevel: 9,
				Effort: 10,
			},
			JPG: types.JPGConfig{
				Quality: 100,
			},
			Cache: true,
		},
	}

	var err error;

	logFatalOnErr := func() {
		if (err != nil) {log.Fatal(err)}
	}

	err = huh.NewConfirm().
		Title("Enable strict JSON?").
		Description("When enabled, TesserPack will always assume all .json files are pure JSON without comments.\nIt will potentially increase optimization performance.").
		Value(&conf.Compiler.JSON.Strict).
		Run()
	logFatalOnErr()

	err = huh.NewConfirm().
		Title("Enable caching?").
		Description("When enabled, TesserPack won't re-optimize assets, resulting in faster overall optimization.\nYou can run 'tesserpack clear-cache' since it will use up your storage.").
		Value(&conf.Compiler.Cache).
		Run()
	logFatalOnErr()

	err = safeIntInput(
		"PNG Compression Level", 
		"Higher levels will result in better optimization in images.", 
		&conf.Compiler.PNG.CompressLevel, 0, 9)
	logFatalOnErr()

	err = safeIntInput(
		"PNG Quality Level", 
		"Higher levels will make images more accurate but higher file size.", 
		&conf.Compiler.PNG.Quality, 1, 100)
	logFatalOnErr()

	err = safeIntInput(
		"PNG Effort Level", 
		"Higher levels will result in better optimization images but higher CPU usage.", 
		&conf.Compiler.PNG.Effort, 1, 10)
	logFatalOnErr()

	err = safeIntInput(
		"JPG Quality Level", 
		"Higher levels will make images more accurate but higher file size.", 
		&conf.Compiler.JPG.Quality, 1, 100)
	logFatalOnErr()

	err = huh.NewText().
		Title("John balal?").
		Description("When enabled, TesserPack will always assume all .json files are pure JSON without comments.\nIt will potentially increase optimization performance.").
		
		Run()
	logFatalOnErr()

	// holy crap, minecraft movie reference
	jsonMomoa, err := json.MarshalIndent(conf, "", "  ")
	logFatalOnErr()

	log.Info(".tesserpackrc is created in your working directory!", "content", string(jsonMomoa))
}