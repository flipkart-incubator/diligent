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
	grumbleShell.AddCommand(dsCmd)

	dsCreateCmd := &grumble.Command{
		Name:    "create",
		Help:    "create a data spec",
		Aliases: []string{"cr"},
		Flags: func(f *grumble.Flags) {
			f.Int("n", "num-recs", 10, "Number of records")
			f.Int("s", "rec-size", 1024, "Approx size of each record")
		},
		Args: func(a *grumble.Args) {
			a.String("name", "a name for the dataspec", grumble.Default(""))
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
			a.String("name", "name of the data spec", grumble.Default(""))
		},
		Run: dsDesc,
	}
	dsCmd.AddCommand(dsDescCmd)

	//dsLoadCmd := &grumble.Command{
	//	Name:    "load",
	//	Help:    "load a data spec for use",
	//	Aliases: []string{"lo"},
	//	Args: func(a *grumble.Args) {
	//		a.String("name", "name of the data spec", grumble.Default(""))
	//	},
	//	Run: dsLoad,
	//}
	//dsCmd.AddCommand(dsLoadCmd)
	//
	//dsShowCmd := &grumble.Command{
	//	Name:    "show",
	//	Help:    "show the currently loaded dataspec",
	//	Aliases: []string{"sh"},
	//	Run: dsShow,
	//}
	//dsCmd.AddCommand(dsShowCmd)
}

func dsCreate(c *grumble.Context) error {
	name := c.Args.String("name")
	numRecs := c.Flags.Int("num-recs")
	recSize := c.Flags.Int("rec-size")
	if name == "" {
		return fmt.Errorf("please specify a name for the dataspec")
	}

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
	if name == "" {
		return fmt.Errorf("please specify a name for the dataspec")
	}
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

//func dsLoad(c *grumble.Context) error {
//	// File name validation and enhancement
//	name := c.Args.String("name")
//	if name == "" {
//		return fmt.Errorf("please specify a name for the dataspec")
//	}
//	if !strings.HasSuffix(name, ".json") {
//		name = name + ".json"
//	}
//
//	// Load dataSpec
//	c.App.Println("Loading data spec from file:", name)
//	spec, err := datagen.LoadSpecFromFile(name)
//	if err != nil {
//		return err
//	}
//	controllerApp.data.dataSpec = spec
//	controllerApp.data.dataSpecName = name
//	controllerApp.data.dataGen = datagen.NewDataGen(spec)
//
//	// Transfer dataSpec
//	success := 0
//	for minAddr, minClient := range controllerApp.minions.minions {
//		request := &proto.LoadDataSpecRequest {
//			SpecName: name,
//			DataSpec: proto.DataSpecToProto(spec),
//			Hash:     0,
//		}
//		response, err := minClient.LoadDataSpec(context.Background(), request)
//		if err != nil {
//			c.App.Printf("%s: Request failed (%v)\n", minAddr, err)
//			continue
//		}
//		if !response.GetIsOk() {
//			c.App.Printf("%s: Failed (%s)\n", minAddr, response.GetFailureReason())
//			continue
//		}
//		c.App.Printf("%s: OK\n", minAddr)
//		success++
//	}
//	c.App.Printf("Successful on %d/%d minions\n", success, len(controllerApp.minions.minions))
//	return nil
//}
//
//func dsShow(c *grumble.Context) error {
//	if controllerApp.data.dataSpec == nil {
//		return fmt.Errorf("no dataspec is currently loaded. Use the 'dataspec load' command to load one")
//	}
//
//	localInfo := fmt.Sprintf("Name=%s, Type=%s, Version=%d, NumRecs=%d, RecSize=%d",
//		controllerApp.data.dataSpecName,
//		controllerApp.data.dataSpec.SpecType,
//		controllerApp.data.dataSpec.Version,
//		controllerApp.data.dataSpec.KeyGenSpec.NumKeys(),
//		controllerApp.data.dataSpec.RecordSize)
//	c.App.Printf("Expected: %s\n", localInfo)
//	success := 0
//	for minAddr, minClient := range controllerApp.minions.minions {
//		// TODO: Add validations on returned info
//		request := &proto.GetDataSpecInfoRequest{}
//		response, err := minClient.GetDataSpecInfo(context.Background(), request)
//		if err != nil {
//			c.App.Printf("%s: Request failed (%v)\n", minAddr, err)
//			continue
//		}
//		if !response.GetIsOk() {
//			c.App.Printf("%s: Failed! (%s)\n", minAddr, response.GetFailureReason())
//			continue
//		}
//		info := fmt.Sprintf("Name=%s, Type=%s, Version=%d, NumRecs=%d, RecSize=%d",
//			response.GetDataSpecInfo().GetSpecName(),
//			response.GetDataSpecInfo().GetSpecType(),
//			response.GetDataSpecInfo().GetVersion(),
//			response.GetDataSpecInfo().GetNumRecs(),
//			response.GetDataSpecInfo().GetRecordSize())
//		c.App.Printf("%s: %s\n", minAddr, info)
//		success++
//	}
//	c.App.Printf("Successful on %d/%d minions\n", success, len(controllerApp.minions.minions))
//	return nil
//}
