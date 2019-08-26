package cni

import (
	"context"
	"fmt"
	"sync"

	gocni "github.com/containerd/go-cni"
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/network"
	"github.com/weaveworks/ignite/pkg/runtime"
)

const (
	// TODO: CNIBinDir and CNIConfDir should maybe be globally configurable?

	// CNIBinDir describes the directory where the CNI binaries are stored
	CNIBinDir = "/opt/cni/bin"
	// CNIConfDir describes the directory where the CNI plugin's configuration is stored
	CNIConfDir = "/etc/cni/net.d"
	// netNSPathFmt gives the path to the a process network namespace, given the pid
	netNSPathFmt = "/proc/%d/ns/net"
)

type cniNetworkPlugin struct {
	cni     gocni.CNI
	runtime runtime.Interface
	once    *sync.Once
}

func GetCNINetworkPlugin(runtime runtime.Interface) (network.Plugin, error) {
	binDirs := []string{CNIBinDir}
	cniInstance, err := gocni.New(gocni.WithMinNetworkCount(2),
		gocni.WithPluginConfDir(CNIConfDir),
		gocni.WithPluginDir(binDirs))
	if err != nil {
		return nil, err
	}

	return &cniNetworkPlugin{
		runtime: runtime,
		cni:     cniInstance,
		once:    &sync.Once{},
	}, nil
}

func (plugin *cniNetworkPlugin) Name() network.PluginName {
	return network.PluginCNI
}

func (plugin *cniNetworkPlugin) PrepareContainerSpec(container *runtime.ContainerConfig) error {
	// No need for the container runtime to set up networking, as this plugin will do it
	container.NetworkMode = "none"
	return nil
}

func (plugin *cniNetworkPlugin) SetupContainerNetwork(containerid string) (*network.Result, error) {
	if err := plugin.initialize(); err != nil {
		return nil, err
	}

	c, err := plugin.runtime.InspectContainer(containerid)
	if err != nil {
		return nil, fmt.Errorf("CNI failed to retrieve network namespace path: %v", err)
	}

	netnsPath := fmt.Sprintf(netNSPathFmt, c.PID)
	result, err := plugin.cni.Setup(context.Background(), containerid, netnsPath)
	if err != nil {
		log.Errorf("failed to setup network for namespace %q: %v", containerid, err)
		return nil, err
	}

	return cniToIgniteResult(result), nil
}

func (plugin *cniNetworkPlugin) initialize() (err error) {
	plugin.once.Do(func() {
		if err = plugin.cni.Load(gocni.WithLoNetwork, gocni.WithDefaultConf); err != nil {
			log.Errorf("failed to load cni configuration: %v", err)
		}
	})
	return
}

func cniToIgniteResult(r *gocni.CNIResult) *network.Result {
	result := &network.Result{}
	for _, iface := range r.Interfaces {
		for _, ip := range iface.IPConfigs {
			result.Addresses = append(result.Addresses, network.Address{
				IP:      ip.IP,
				Gateway: ip.Gateway,
			})
		}
	}
	return result
}

func (plugin *cniNetworkPlugin) RemoveContainerNetwork(containerid string) error {
	if err := plugin.initialize(); err != nil {
		return err
	}

	// Lack of namespace should not be fatal on teardown
	c, err := plugin.runtime.InspectContainer(containerid)
	if err != nil {
		log.Infof("CNI failed to retrieve network namespace path: %v", err)
		return nil
	}

	netnsPath := fmt.Sprintf(netNSPathFmt, c.PID)
	if c.PID == 0 {
		log.Info("CNI failed to retrieve network namespace path, PID was 0")
		return nil
	}

	return plugin.cni.Remove(context.Background(), containerid, netnsPath)
}
