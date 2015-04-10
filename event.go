package main

import (
	"encoding/json"
	"errors"

	"code.google.com/p/appengine-go/appengine"

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

func LoadEvents(c appengine.Context) ([]Event, error) {
	var events []Event

	q := datastore.NewQuery(EventEntityKind)

	_, err := q.GetAll(c, &events)

	return events, err
}

func UpdateEvents(c appengine.Context, events []Event) error {
	var kBatch []*datastore.Key
	var eventBatch []Event
	batchSize := 0

	for _, event := range events {
		key := datastore.NewKey(c, EventEntityKind, event.ID, 0, nil)

		if batchSize < 500 {
			kBatch = append(kBatch, key)
			eventBatch = append(eventBatch, event)
			batchSize++
		} else {
			err := UpdateEventBatch(c, kBatch, eventBatch)

			if err != nil {
				c.Infof("Could not update historic archive: %#v", err)
				return err
			}

			kBatch = []*datastore.Key{}
			eventBatch = []Event{}
			batchSize = 0
		}
	}

	if batchSize > 0 {
		err := UpdateEventBatch(c, kBatch, eventBatch)

		if err != nil {
			c.Infof("Could not update historic archive: %#v", err)
			return err
		}
	}

	return nil

}

func UpdateEventBatch(c appengine.Context, kBatch []*datastore.Key, eventBatch []Event) error {
	_, err := datastore.PutMulti(c, kBatch, eventBatch)

	if err != nil {
		return err
	}

	return nil
}
