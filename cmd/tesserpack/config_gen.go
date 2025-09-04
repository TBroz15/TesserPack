package main

import (
	"tesserpack/internal/types"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

func ConfigGen() {
	_ = &types.CompilerConfig{}


	err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Enable strict JSON?").
				Description("When enabled, TesserPack will always assume all .json files are pure JSON without comments.\nIt will potentially increase compilation performance."),

		).Title("JSON"),
	).WithShowHelp(true).WithHeight(10).Run()

	if err != nil {
		log.Fatal(err)
	}
}