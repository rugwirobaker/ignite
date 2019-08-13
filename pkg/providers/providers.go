package providers

import (
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/network"
	"github.com/weaveworks/ignite/pkg/runtime"
	"github.com/weaveworks/ignite/pkg/storage"
)

// NetworkPlugins provides the initialized network plugins indexed by their name
var NetworkPlugins = make(map[network.PluginName]network.Plugin)

// NetworkPlugin provides the chosen network plugin that should be used
// This should be set after parsing user input on what network mode to use
var NetworkPlugin network.Plugin

// Runtime provides the default container runtime
var Runtime runtime.Interface

// Client is the default client that can be easily used
var Client *client.Client

// Storage is the default storage implementation
var Storage storage.Storage

type ProviderInitFunc func() error

// Populate initializes all providers
func Populate(providers []ProviderInitFunc) error {
	log.Trace("Populating providers...")
	for i, init := range providers {
		log.Tracef("Provider %d...", i)
		if err := init(); err != nil {
			return err
		}
	}

	return nil
}
