package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"tesserpack/internal/helpers/config"
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

func configGenPrompt(conf *types.TesserPackConfig) {
	var err error;
	logFatalOnErr := func() {
		if (err != nil) {log.Fatal(err)}
	}

	err = huh.NewConfirm().
		Title("Enable strict JSON?").
		Description("When enabled, TesserPack will always assume all .json files are pure JSON without comments.\nThis is optional since Minecraft Bedrock Edition can read JSON with comments.\nBut, it will potentially make png JSON optimization faster.").
		Value(&conf.Compiler.JSON.Strict).
		Run()
	logFatalOnErr()

	err = safeIntInput(
		"PNG Compression Level", 
		"Higher levels will result in better optimization in images.", 
		&conf.Compiler.PNG.Compression, 0, 9)
	logFatalOnErr()

	err = safeIntInput(
		"PNG Quality Level", 
		"Higher levels will make images more accurate but higher file size.", 
		&conf.Compiler.PNG.Q, 1, 100)
	logFatalOnErr()

	err = safeIntInput(
		"PNG Effort Level", 
		"Higher levels will result in better optimization images but higher CPU usage.", 
		&conf.Compiler.PNG.Effort, 1, 10)
	logFatalOnErr()

	err = safeIntInput(
		"JPG Quality Level", 
		"Higher levels will make images more accurate but higher file size.", 
		&conf.Compiler.JPG.Q, 1, 100)
	logFatalOnErr()

	err = huh.NewConfirm().
		Title("Enable caching?").
		Description("When enabled, TesserPack won't waste CPU usage by storing already optimized assets.\nThis will result into faster overall optimization.\nYou can run 'tesserpack clear-cache' since it will use up your storage.").
		Value(&conf.Compiler.Cache).
		Run()
	logFatalOnErr()

	ignoreGlobPatterns := 
`node_modules/*
node_modules/**/*
.git/*
.git/**/*
.vscode/*
.vscode/**/*
.github/*
.github/**/*
`
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Ignore List").
				Description("TesserPack will ignore files and directories via Glob patterns.\nGlob patterns are seperated by every new line.").
				Value(&ignoreGlobPatterns).
				ShowLineNumbers(true),		
		),
	).
	WithShowHelp(true).
	Run()
	logFatalOnErr()
	conf.IgnoreGlob = strings.Split(ignoreGlobPatterns, "\n")
}

func ConfigGen(doCreateRecommended bool) {
	conf := config.NewDefault()

	if _, err := os.Stat(".tesserpackrc.json5"); !os.IsNotExist(err) {
		log.Fatalf(".tesserpackrc.json5 already exists in your working directory.")
	}

	if (!doCreateRecommended) {
		configGenPrompt(&conf)
	}

	conf.IgnoreGlob = slices.DeleteFunc(conf.IgnoreGlob, func(elm string) bool {
		return elm == ""
	})

	// holy crap, minecraft movie reference
	confJsonMomoa, err := json.MarshalIndent(conf, "", "  ")
	if (err != nil) {
		log.Fatal(err)
	}

	log.Info("", "content", string(confJsonMomoa))

	os.WriteFile(".tesserpackrc.json5", confJsonMomoa, 0777)

	log.Info(".tesserpackrc.json5 is now created in your working directory!")
}
