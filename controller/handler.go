package controller

import (
	"fmt"
)

// Handler is implemented by any handler.
// The Handle method is used to process event
type Handler interface {
	Init(c *Config) error
	ObjectCreated(obj interface{})
	ObjectDeleted(event KubeEvent)
	ObjectUpdated(oldObj interface{}, event KubeEvent)
}

// Map maps each event handler function to a name for easily lookup
var Map = map[string]interface{}{
	"default": &DefaultHandler{},
}

// DefaultHandler handler implements Handler interface,
// print each event with JSON format
type DefaultHandler struct {
}

// Init initializes handler configuration
// Do nothing for default handler
func (d *DefaultHandler) Init(c *Config) error {
	fmt.Println("DefaultHandler init ")

	return nil
}

func (d *DefaultHandler) ObjectCreated(obj interface{}) {
	fmt.Println("DefaultHandler ObjectCreated ", obj)

}

func (d *DefaultHandler) ObjectDeleted(event KubeEvent) {
	fmt.Println("DefaultHandler ObjectDeleted ", event)

}

func (d *DefaultHandler) ObjectUpdated(oldObj interface{}, event KubeEvent) {
	fmt.Println("DefaultHandler ObjectUpdated ", event)

}
