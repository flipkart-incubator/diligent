package main

import (
	"github.com/desertbit/grumble"
	"github.com/flipkart-incubator/diligent/pkg/datagen"
	"strings"
)

func init() {
	dsCmd := &grumble.Command{
		Name:    "dataspec",
		Help:    "work with dataspecs",
		Aliases: []string{"ds"},
	}
	grumbleApp.AddCommand(dsCmd)

	dsCreateCmd := &grumble.Command{
		Name:    "create",
		Help:    "create a data spec",
		Aliases: []string{"cr"},
		Flags: func(f *grumble.Flags) {
			f.Int("n", "num-recs", 10, "Number of records")
			f.Int("s", "rec-size", 1024, "Approx size of each record")
		},
		Args: func(a *grumble.Args) {
			a.String("name", "a name for the dataspec")
		},
		Run: dsCreate,
	}
	dsCmd.AddCommand(dsCreateCmd)

	dsDescCmd := &grumble.Command{
		Name: "desc",
		Help: "describe a data spec",
		Flags: func(f *grumble.Flags) {
			f.Bool("v", "verbose", false, "show verbose information")
		},
		Args: func(a *grumble.Args) {
			a.String("name", "name of the data spec")
		},
		Run: dsDesc,
	}
	dsCmd.AddCommand(dsDescCmd)
}

func dsCreate(c *grumble.Context) error {
	name := c.Args.String("name")
	numRecs := c.Flags.Int("num-recs")
	recSize := c.Flags.Int("rec-size")

	c.App.Println("Creating dataspec...")
	c.App.Println("name:", name)
	c.App.Println("num-recs:", numRecs)
	c.App.Println("rec-size:", recSize)

	spec := datagen.NewSpec(numRecs, recSize)
	if !strings.HasSuffix(name, ".json") {
		name = name + ".json"
	}
	err := spec.SaveToFile(name)
	if err != nil {
		return err
	} else {
		c.App.Println("Done!")
	}
	return nil
}

func dsDesc(c *grumble.Context) error {
	name := c.Args.String("name")
	if !strings.HasSuffix(name, ".json") {
		name = name + ".json"
	}
	spec, err := datagen.LoadSpecFromFile(name)
	if err != nil {
		return err
	}
	if c.Flags.Bool("verbose") {
		jsonString := string(spec.Json())
		c.App.Println(jsonString)
	} else {
		c.App.Println("Type: ", spec.SpecType)
		c.App.Println("Version: ", spec.Version)
		c.App.Println("NumRecs: ", spec.KeyGenSpec.NumKeys())
		c.App.Println("RecSize: ", spec.RecordSize)
	}
	return nil
}
