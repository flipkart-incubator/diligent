package main

import (
	"database/sql"
	"github.com/flipkart-incubator/diligent/pkg/datagen"
	"github.com/flipkart-incubator/diligent/pkg/metrics"
	"github.com/flipkart-incubator/diligent/pkg/proto"
)

type DataContext struct {
	dataSpecName string
	dataSpec     *datagen.Spec
	dataGen      *datagen.DataGen
}

type DBContext struct {
	driver string
	url    string
	db     *sql.DB
}

type MinionContext struct {
	minions map[string]proto.MinionClient
}

type ControllerApp struct {
	data    DataContext
	db      DBContext
	minions MinionContext
	metrics *metrics.DiligentMetrics
}

var controllerApp *ControllerApp

func init() {
	controllerApp = &ControllerApp{
		data: DataContext{
			dataSpecName: "",
			dataSpec:     nil,
			dataGen:      nil,
		},
		db: DBContext{
			driver: "",
			url:    "",
		},
		minions: MinionContext{
			minions: make(map[string]proto.MinionClient),
		},
	}
}