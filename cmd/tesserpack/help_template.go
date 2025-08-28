package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var titleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("2"))

var usageStyle = titleStyle.SetString("-Usage-").Render()
var optionsStyle = titleStyle.SetString("-Options-").Render()
var commandsStyle = titleStyle.SetString("-Commands-").Render()
var globalOptionsStyle = titleStyle.SetString("-Global Options-").Render()

var SubCommandHelpTemplate = fmt.Sprintf(`{{template "helpNameTemplate" .}}

%v
   {{template "usageTemplate" .}}{{if .Category}}

CATEGORY:
   {{.Category}}{{end}}{{if .Description}}

DESCRIPTION:
   {{template "descriptionTemplate" .}}{{end}}{{if .VisibleFlagCategories}}

%v{{template "visibleFlagCategoryTemplate" .}}{{else if .VisibleFlags}}

%v{{template "visibleFlagTemplate" .}}{{end}}{{if .VisiblePersistentFlags}}

%v{{template "visiblePersistentFlagTemplate" .}}{{end}}
`,
usageStyle,
optionsStyle,
optionsStyle,
globalOptionsStyle,
)

var RootHelpTemplate = fmt.Sprintf(`A build tool that compiles and optimize Minecraft packs for easier download. You know why would you download this right?
https://github.com/TBroz15/TesserPack

%v
   {{if .UsageText}}{{wrap .UsageText 3}}{{else}}{{.FullName}} {{if .VisibleFlags}}[global options]{{end}}{{if .VisibleCommands}} [command [command options]]{{end}}{{if .ArgsUsage}} {{.ArgsUsage}}{{else}}{{if .Arguments}} [arguments...]{{end}}{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

VERSION:
   {{.Version}}{{end}}{{end}}{{if .Description}}

DESCRIPTION:
   {{template "descriptionTemplate" .}}{{end}}
{{- if len .Authors}}

AUTHOR{{template "authorsTemplate" .}}{{end}}{{if .VisibleCommands}}

%v{{template "visibleCommandCategoryTemplate" .}}{{end}}{{if .VisibleFlagCategories}}

%v:{{template "visibleFlagCategoryTemplate" .}}{{else if .VisibleFlags}}

%v:{{template "visibleFlagTemplate" .}}{{end}}{{if .Copyright}}

COPYRIGHT:
   {{template "copyrightTemplate" .}}{{end}}
`, 
usageStyle,
commandsStyle,
globalOptionsStyle,
globalOptionsStyle)