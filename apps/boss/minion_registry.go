package main

import (
	log "github.com/sirupsen/logrus"
	"sync"
)

// MinionRegistry maintains a register of all registered minions and their MinionManagers
// It allows new minions to be registered and existing minions to be unregistered
// It is thread safe
type MinionRegistry struct {
	mut      sync.Mutex // Must be taken for all operations on the MinionRegistry
	registry map[string]*MinionProxy
}

func NewMinionRegistry() *MinionRegistry {
	return &MinionRegistry{
		registry: make(map[string]*MinionProxy),
	}
}

func (r *MinionRegistry) RegisterMinion(addr string) error {
	log.Infof("RegisterMinion(%s)", addr)
	r.mut.Lock()
	defer r.mut.Unlock()

	// Close any existing proxy for the same minion and delete the entry
	p := r.registry[addr]
	if p != nil {
		log.Infof("RegisterMinion(%s): Existing entry found. Closing existing proxy", addr)
		p.Close()
	}
	delete(r.registry, addr)

	// Create a new proxy and add the entry in the registry
	p, err := NewMinionProxy(addr)
	if err != nil {
		return err
	}
	r.registry[addr] = p
	return nil
}

func (r *MinionRegistry) UnregisterMinion(addr string) error {
	log.Infof("UnregisterMinion(%s)", addr)
	r.mut.Lock()
	defer r.mut.Unlock()

	p := r.registry[addr]
	if p != nil {
		log.Infof("UnregisterMinion(%s): Found registered minion - removing", addr)
		p.Close()
		delete(r.registry, addr)
	} else {
		log.Infof("UnregisterMinion(%s): No such minion. Ignoring request", addr)
	}
	return nil
}

func (r *MinionRegistry) Minions() map[string]*MinionProxy {
	log.Infof("Minions")
	r.mut.Lock()
	defer r.mut.Unlock()

	return r.registry
}

func (r *MinionRegistry) GetNumMinions() int {
	log.Infof("GetNumMinions")
	r.mut.Lock()
	defer r.mut.Unlock()

	return len(r.registry)
}

//func (r *MinionRegistry) GetMinionAddrs() []string {
//	log.Infof("GetMinionAddrs")
//	r.mut.Lock()
//	defer r.mut.Unlock()
//
//	addrs := make([]string, len(r.registry))
//
//	i := 0
//	for addr := range r.registry {
//		addrs[i] = addr
//		i++
//	}
//	return addrs
//}
//
//func (r *MinionRegistry) GetMinionManagers() []*MinionProxy {
//	log.Infof("GetMinionManagers")
//	r.mut.Lock()
//	defer r.mut.Unlock()
//
//	mms := make([]*MinionProxy, len(r.registry))
//
//	i := 0
//	for _, p := range r.registry {
//		mms[i] = p
//		i++
//	}
//	return mms
//}
