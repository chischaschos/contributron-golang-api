package main

import (
	"encoding/json"
	"errors"

	"appengine"
	"appengine/datastore"
)

const (
	EventEntityKind = "Event"
)

type User struct {
	Login string `json:"login"`
}

type PullRequest struct {
	Merged   bool   `json:"merged"`
	MergedBy User   `json:"merged_by"`
	MergedAt string `json:"merged_at"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	User     `json:"user"`
}

type Payload struct {
	Action      string `json:"action"`
	PullRequest `json:"pull_request"`
}

type Event struct {
	ID   string
	Type string
	Payload
}

func (e *Event) Load(c <-chan datastore.Property) error {
	loadChan := make(chan datastore.Property, 100)

	for prop := range c {
		if prop.Name == "Payload" {

			jsonBytes, ok := prop.Value.([]byte)
			if !ok {
				return errors.New("Could not convert Payload to JSON []byte")
			}

			err := json.Unmarshal(jsonBytes, &e.Payload)
			if err != nil {
				return err
			}

		} else {
			loadChan <- prop
		}
	}

	// If I don't close it the process got stuck
	close(loadChan)

	return datastore.LoadStruct(e, loadChan)
}

func (e *Event) Save(c chan<- datastore.Property) error {
	json, err := json.Marshal(e.Payload)

	if err != nil {
		return err
	}

	c <- datastore.Property{
		Name:    "Payload",
		Value:   json,
		NoIndex: true,
	}

	return datastore.SaveStruct(e, c)
}

func UpdateEvents(c appengine.Context, events []Event) error {
	for _, event := range events {
		key := datastore.NewKey(c, EventEntityKind, event.ID, 0, nil)
		_, err := datastore.Put(c, key, &event)

		if err != nil {
			c.Infof("Error %#v, %#v updating: %#v", key, event, err)
			return err
		}
	}
	return nil
}
