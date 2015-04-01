package main

import (
	"encoding/json"
	"errors"

	"appengine"
	"appengine/datastore"
)

const (
	ConfigurationEntityKind = "Configuration"
)

type Organization struct {
	Name    string
	Members []Member `datastore:"-"`
}

func (o *Organization) Load(c <-chan datastore.Property) error {
	loadChan := make(chan datastore.Property, 100)

	for prop := range c {
		if prop.Name == "Members" {

			jsonBytes, ok := prop.Value.([]byte)
			if !ok {
				return errors.New("Could not convert Members to JSON []byte")
			}

			err := json.Unmarshal(jsonBytes, &o.Members)
			if err != nil {
				return err
			}

		} else {
			loadChan <- prop
		}
	}

	// If I don't close it the process got stuck
	close(loadChan)

	return datastore.LoadStruct(o, loadChan)
}

func (o *Organization) Save(c chan<- datastore.Property) error {
	json, err := json.Marshal(o.Members)

	if err != nil {
		return err
	}

	c <- datastore.Property{
		Name:    "Members",
		Value:   json,
		NoIndex: true,
	}

	return datastore.SaveStruct(o, c)
}

func UpdateOrganizationMembers(c appengine.Context, org *Organization) error {
	c.Infof("Updating %#v", org)

	configurationID := "organization"
	key := datastore.NewKey(c, ConfigurationEntityKind, configurationID, 0, nil)
	_, err := datastore.Put(c, key, org)

	return err
}

func LoadOrganization(c appengine.Context) (*Organization, error) {
	key := datastore.NewKey(c, ConfigurationEntityKind, "organization", 0, nil)

	var org Organization
	err := datastore.Get(c, key, &org)

	return &org, err

}
