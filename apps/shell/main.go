package main

import (
	"github.com/desertbit/grumble"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

var grumbleShell = grumble.New(&grumble.Config{
	Name:                  "diligent",
	Description:           "Diligent: A SQL load runner",
	HistoryFile:           "/tmp/diligent.hist",
	Prompt:                "diligent » ",
	PromptColor:           color.New(color.FgCyan, color.Bold),
	HelpHeadlineColor:     color.New(color.FgCyan),
	HelpHeadlineUnderline: false,
	HelpSubCommands:       false,
})

func main() {
	log.SetLevel(log.WarnLevel)
	grumble.Main(grumbleShell)
}
