/*
	Note: This file is autogenerated! Do not edit it manually!
	Edit client_image_template.go instead, and run
	hack/generate-client.sh afterwards.
*/

package client

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/gitops-toolkit/pkg/runtime"
	"github.com/weaveworks/gitops-toolkit/pkg/storage"
	"github.com/weaveworks/gitops-toolkit/pkg/storage/filterer"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ImageClient is an interface for accessing Image-specific API objects
type ImageClient interface {
	// New returns a new Image
	New() *api.Image
	// Get returns the Image matching given UID from the storage
	Get(runtime.UID) (*api.Image, error)
	// Set saves the given Image into persistent storage
	Set(*api.Image) error
	// Patch performs a strategic merge patch on the object with
	// the given UID, using the byte-encoded patch given
	Patch(runtime.UID, []byte) error
	// Find returns the Image matching the given filter, filters can
	// match e.g. the Object's Name, UID or a specific property
	Find(filter filterer.BaseFilter) (*api.Image, error)
	// FindAll returns multiple Images matching the given filter, filters can
	// match e.g. the Object's Name, UID or a specific property
	FindAll(filter filterer.BaseFilter) ([]*api.Image, error)
	// Delete deletes the Image with the given UID from the storage
	Delete(uid runtime.UID) error
	// List returns a list of all Images available
	List() ([]*api.Image, error)
}

// Images returns the ImageClient for the IgniteInternalClient instance
func (c *IgniteInternalClient) Images() ImageClient {
	if c.imageClient == nil {
		c.imageClient = newImageClient(c.storage, c.gv)
	}

	return c.imageClient
}

// imageClient is a struct implementing the ImageClient interface
// It uses a shared storage instance passed from the Client together with its own Filterer
type imageClient struct {
	storage  storage.Storage
	filterer *filterer.Filterer
	gvk      schema.GroupVersionKind
}

// newImageClient builds the imageClient struct using the storage implementation and a new Filterer
func newImageClient(s storage.Storage, gv schema.GroupVersion) ImageClient {
	return &imageClient{
		storage:  s,
		filterer: filterer.NewFilterer(s),
		gvk:      gv.WithKind(api.KindImage.Title()),
	}
}

// New returns a new Object of its kind
func (c *imageClient) New() *api.Image {
	log.Tracef("Client.New; GVK: %v", c.gvk)
	obj, err := c.storage.New(c.gvk)
	if err != nil {
		panic(fmt.Sprintf("Client.New must not return an error: %v", err))
	}
	return obj.(*api.Image)
}

// Find returns a single Image based on the given Filter
func (c *imageClient) Find(filter filterer.BaseFilter) (*api.Image, error) {
	log.Tracef("Client.Find; GVK: %v", c.gvk)
	object, err := c.filterer.Find(c.gvk, filter)
	if err != nil {
		return nil, err
	}

	return object.(*api.Image), nil
}

// FindAll returns multiple Images based on the given Filter
func (c *imageClient) FindAll(filter filterer.BaseFilter) ([]*api.Image, error) {
	log.Tracef("Client.FindAll; GVK: %v", c.gvk)
	matches, err := c.filterer.FindAll(c.gvk, filter)
	if err != nil {
		return nil, err
	}

	results := make([]*api.Image, 0, len(matches))
	for _, item := range matches {
		results = append(results, item.(*api.Image))
	}

	return results, nil
}

// Get returns the Image matching given UID from the storage
func (c *imageClient) Get(uid runtime.UID) (*api.Image, error) {
	log.Tracef("Client.Get; UID: %q, GVK: %v", uid, c.gvk)
	object, err := c.storage.Get(c.gvk, uid)
	if err != nil {
		return nil, err
	}

	return object.(*api.Image), nil
}

// Set saves the given Image into the persistent storage
func (c *imageClient) Set(image *api.Image) error {
	log.Tracef("Client.Set; UID: %q, GVK: %v", image.GetUID(), c.gvk)
	return c.storage.Set(c.gvk, image)
}

// Patch performs a strategic merge patch on the object with
// the given UID, using the byte-encoded patch given
func (c *imageClient) Patch(uid runtime.UID, patch []byte) error {
	return c.storage.Patch(c.gvk, uid, patch)
}

// Delete deletes the Image from the storage
func (c *imageClient) Delete(uid runtime.UID) error {
	log.Tracef("Client.Delete; UID: %q, GVK: %v", uid, c.gvk)
	return c.storage.Delete(c.gvk, uid)
}

// List returns a list of all Images available
func (c *imageClient) List() ([]*api.Image, error) {
	log.Tracef("Client.List; GVK: %v", c.gvk)
	list, err := c.storage.List(c.gvk)
	if err != nil {
		return nil, err
	}

	results := make([]*api.Image, 0, len(list))
	for _, item := range list {
		results = append(results, item.(*api.Image))
	}

	return results, nil
}
