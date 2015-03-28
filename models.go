package main

import (
	"encoding/json"
	"errors"

	"appengine/datastore"
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

type Member struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	URL       string `json:"url"`
}
