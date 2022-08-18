package main

import (
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

var grumbleApp = grumble.New(&grumble.Config{
	Name:                  "diligent",
	Description:           "Diligent: A SQL load runner",
	HistoryFile:           "/tmp/diligent.hist",
	Prompt:                "diligent » ",
	PromptColor:           color.New(color.FgCyan, color.Bold),
	HelpHeadlineColor:     color.New(color.FgCyan),
	HelpHeadlineUnderline: false,
	HelpSubCommands:       false,

	Flags: func(f *grumble.Flags) {
		f.String("b", "boss", "", "host[:port] of boss server")
	},
})

func onGrumbleInit(a *grumble.App, flags grumble.FlagMap) error {
	bossAddr := flags.String("boss")
	if bossAddr == "" {
		return fmt.Errorf("please provide a valid address for the boss server using -b or --boss")
	}
	return nil
}

func main() {
	log.SetLevel(log.WarnLevel)
	grumbleApp.OnInit(onGrumbleInit)
	grumble.Main(grumbleApp)
}
