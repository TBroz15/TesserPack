// my main.go is to blow up

package main

import (
	"context"
	"os"

	"tesserpack/internal/compiler"
	"tesserpack/internal/helpers"
	"tesserpack/internal/types"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

func main() {
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
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				inPath       := cmd.String("in")
				outPath      := cmd.String("out")
				isStrictJSON := cmd.Bool("strict-json")

				conf := types.Config{
					InPath:       inPath,
					OutPath:      outPath,
					IsStrictJSON: isStrictJSON,
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
	}

	cmd := &cli.Command{
		EnableShellCompletion: true,
        Name:  "tesserpack",
        Usage: "A build tool that compiles and optimize Minecraft packs for easier download. You know why you download this right?\nhttps://github.com/TBroz15/TesserPack",
		Commands: subCommands,
    }

    if err := cmd.Run(context.Background(), os.Args); err != nil {
        log.Fatal(err)
    }
}