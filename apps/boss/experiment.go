package main

import (
	"fmt"
	"github.com/flipkart-incubator/diligent/pkg/proto"
	"sync"
	"time"
)

type ExperimentState int

const (
	_ ExperimentState = iota
	New
	Started
	Stopped
)

func (s ExperimentState) String() string {
	switch s {
	case New:
		return "New"
	case Started:
		return "Started"
	case Stopped:
		return "Stopped"
	default:
		panic(fmt.Sprintf("unknown experiment state %d", s))
	}
}

func (s ExperimentState) ToProto() proto.ExperimentState {
	switch s {
	case New:
		return proto.ExperimentState_NEW_EXPERIMENT
	case Started:
		return proto.ExperimentState_STARTED
	case Stopped:
		return proto.ExperimentState_STOPPED
	default:
		panic(fmt.Sprintf("unknown experiment state %d", s))
	}
}

type Experiment struct {
	mut       sync.Mutex
	name      string
	state     ExperimentState
	startTime time.Time
	stopTime  time.Time
}

func NewExperiment(name string) *Experiment {
	return &Experiment{
		name:  name,
		state: New,
	}
}

func (e *Experiment) Name() string {
	e.mut.Lock()
	defer e.mut.Unlock()
	return e.name
}

func (e *Experiment) State() ExperimentState {
	e.mut.Lock()
	defer e.mut.Unlock()
	return e.state
}

func (e *Experiment) HasEnded() bool {
	e.mut.Lock()
	defer e.mut.Unlock()
	return e.state == Stopped
}

func (e *Experiment) Start() error {
	e.mut.Lock()
	defer e.mut.Unlock()

	if e.state != New {
		return fmt.Errorf("experiment %s cannot be started. state=%s", e.name, e.state.String())
	}
	e.state = Started
	e.startTime = time.Now()
	return nil
}

func (e *Experiment) Stop() error {
	e.mut.Lock()
	defer e.mut.Unlock()

	if e.state != Started {
		return fmt.Errorf("experiment %s cannot be stopped. state=%s", e.name, e.state.String())
	}
	e.state = Stopped
	e.stopTime = time.Now()
	return nil
}

func (e *Experiment) ToProto() *proto.ExperimentInfo {
	e.mut.Lock()
	defer e.mut.Unlock()
	return &proto.ExperimentInfo{
		Name:      e.name,
		State:     e.state.ToProto(),
		StartTime: e.startTime.UnixMilli(),
		StopTime:  e.stopTime.UnixMilli(),
	}
}
