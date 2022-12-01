package main

import (
	"fmt"
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
			f.Bool("", "reuse-payload", false, "Reuse the same payload for every record")
			f.Bool("", "skip-if-exists", false, "Reuse if a data spec file with same name and params exists")
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
	if !strings.HasSuffix(name, ".json") {
		name = name + ".json"
	}
	numRecs := c.Flags.Int("num-recs")
	recSize := c.Flags.Int("rec-size")
	reusePayload := c.Flags.Bool("reuse-payload")
	skipIfExists := c.Flags.Bool("skip-if-exists")
	tryCreate := true

	if skipIfExists {
		c.App.Println("Checking for existing dataspec file...")
		spec, err := datagen.LoadSpecFromFile(name)
		if err == nil {
			c.App.Println("Found and loaded existing dataspec file...")
			if spec.KeyGenSpec.NumKeys() == numRecs && spec.RecordSize == recSize {
				c.App.Println("Parameters match, will skip creating new dataspec file")
				tryCreate = false
			} else {
				c.App.Printf("num-recs: Expected=%d, Found=%d\n", numRecs, spec.KeyGenSpec.NumKeys())
				c.App.Printf("rec-size: Expected=%d, Found=%d\n", recSize, spec.RecordSize)
				return fmt.Errorf("existing dataspec found, but parameters do not match")
			}
		} else {
			c.App.Println("Was unable to find or load any existing dataspec file, will create new file")
		}
	}

	if tryCreate {
		c.App.Println("Creating dataspec...")
		c.App.Println("name:", name)
		c.App.Println("num-recs:", numRecs)
		c.App.Println("rec-size:", recSize)

		spec := datagen.NewSpec(numRecs, recSize, reusePayload)
		err := spec.SaveToFile(name)
		if err != nil {
			return err
		} else {
			c.App.Println("Done!")
		}
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
