// my main.go is to blow up

package main

import (
	"context"
	"fmt"
	"os"

	"tesserpack/internal/compiler"
	"tesserpack/internal/helpers"
	"tesserpack/internal/helpers/cache"
	"tesserpack/internal/types"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

func main() {
	cli.SubcommandHelpTemplate = SubCommandHelpTemplate
	
	grassBlock := lipgloss.NewStyle().
    	SetString("⬒").
    	Foreground(lipgloss.Color("#2f9e44")).
		Render()

	tesserpackTitle := lipgloss.NewStyle().
		SetString("TesserPack").
		Underline(true).
		Bold(true).
		Render()

	tesserpackVersion := lipgloss.NewStyle().
		SetString("v0.4").
		Italic(true).
		Render()

	fmt.Printf("\n %v  %v %v\n\n", grassBlock, tesserpackTitle,tesserpackVersion)

	subCommands := []*cli.Command{
		{
			Name:    "compile",
			Aliases: []string{"c"},
			Usage:   "Starts the compilation process.",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "in",
					Aliases:  []string{"i"},
					Usage:    "Specify the input pack to be compiled.",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "out",
					Aliases:  []string{"o"},
					Usage:    "Specify where the optimized pack will be.", 
					Required: false,
				},
				&cli.BoolFlag{
					Name:     "strict-json",
					Aliases:  []string{"sj"},
					Usage:    "TesserPack will assume every .json file has no comments.", 
					Required: false,
				},
				&cli.BoolFlag{
					Name:     "disable-cache",
					Aliases:  []string{"dc"},
					Usage:    "TesserPack will complie everything from scratch.",
					Required: false,
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				inPath       := cmd.String("in")
				outPath      := cmd.String("out")
				isStrictJSON := cmd.Bool("strict-json")
				isCached	 := !cmd.Bool("disable-cache")

				conf := types.Config{
					InPath:       inPath,
					OutPath:      outPath,
					IsStrictJSON: isStrictJSON,
					IsCached:     isCached,
				} 
	
				err := compiler.StartCompile(&conf)
	
				return err
			},
		},
		{
			Name:    "help",
			Aliases: []string{"h"},
			Usage:   "???????",
		},
		{
			Name: 	 "clear-temp",
			Aliases: []string{"ct"},
			Usage:   "Clears all temporary directories.\nJust in case if TesserPack fails on compilation and hasn't deleted the temporary files.",
			Action: func(ctx context.Context, c *cli.Command) error {
				log.Info("Clearing temporary directories...")

				err := helpers.ClearTemp()
				if (err != nil) {return err}

				log.Info("Successfully cleared temporary directories.")

				return nil
			},
		},
		{
			Name: 	 "clear-cache",
			Aliases: []string{"cc"},
			Usage:   "Clears all cache files.",
			Action: func(ctx context.Context, c *cli.Command) error {
				log.Info("Clearing cache files...")

				err := cache.ClearCacheDir()
				if (err != nil) {return err}

				log.Info("Successfully cleared cache files.")

				return nil
			},
		},
	}

	cmd := &cli.Command{
		EnableShellCompletion: true,
        Name:  "tesserpack",
		Commands: subCommands,
		CustomRootCommandHelpTemplate: RootHelpTemplate,
		Suggest: true,
	}

    if err := cmd.Run(context.Background(), os.Args); err != nil {
        log.Fatal(err)
    }
}