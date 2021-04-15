package data

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/xid"
	"github.com/tidwall/buntdb"
)

// Event represents a system in the system.  These events
// can be logged or passed (as meta information) to other systems
type Event struct {
	ID        string    `json:"id"`        // Unique Event ID
	Created   time.Time `json:"created"`   // Event creation time
	EventType string    `json:"eventtype"` // One of: System startup, File found, File processed, File error, System shutdown
	MediaType string    `json:"mediatype"` // The type of media involved: TV / Movie / Music / System
	Details   string    `json:"details"`   // Additional information (like the files involved)
}

// AddEvent adds an event to the system
func (store Manager) AddEvent(eventtype, mediatype, details string, expiresafter time.Duration) (Event, error) {
	//	Our return item
	retval := Event{}

	newEvent := Event{
		ID:        xid.New().String(), // Generate a new id
		Created:   time.Now(),
		EventType: eventtype,
		MediaType: mediatype,
		Details:   details,
	}

	//	Serialize to JSON format
	encoded, err := json.Marshal(newEvent)
	if err != nil {
		return retval, fmt.Errorf("problem serializing the data: %s", err)
	}

	//	Save it to the database:
	err = store.systemdb.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(GetKey("Event", newEvent.ID), string(encoded), &buntdb.SetOptions{Expires: true, TTL: expiresafter})
		return err
	})

	//	If there was an error saving the data, report it:
	if err != nil {
		return retval, fmt.Errorf("problem saving the event: %s", err)
	}

	//	Set our retval:
	retval = newEvent

	//	Return our data:
	return retval, nil
}

// GetEvent gets an event from the system
func (store Manager) GetEvent(id string) (Event, error) {
	//	Our return item
	retval := Event{}

	err := store.systemdb.View(func(tx *buntdb.Tx) error {
		item, err := tx.Get(GetKey("Event", id))
		if err != nil {
			return err
		}

		if len(item) > 0 {
			//	Unmarshal data into our item
			val := []byte(item)
			if err := json.Unmarshal(val, &retval); err != nil {
				return err
			}
		}

		return nil
	})

	//	If there was an error, report it:
	if err != nil {
		return retval, fmt.Errorf("problem getting the event: %s", err)
	}

	//	Return our data:
	return retval, nil
}

// GetAllEvents gets all events in the system
func (store Manager) GetAllEvents() ([]Event, error) {
	//	Our return item
	retval := []Event{}

	//	Set our prefix
	prefix := GetKey("Event")

	//	Iterate over our values:
	err := store.systemdb.View(func(tx *buntdb.Tx) error {
		tx.Descend(prefix, func(key, val string) bool {

			if len(val) > 0 {
				//	Create our item:
				item := Event{}

				//	Unmarshal data into our item
				bval := []byte(val)
				if err := json.Unmarshal(bval, &item); err != nil {
					return false
				}

				//	Add to the array of returned users:
				retval = append(retval, item)
			}

			return true
		})
		return nil
	})

	//	If there was an error, report it:
	if err != nil {
		return retval, fmt.Errorf("problem getting the list of events: %s", err)
	}

	//	Return our data:
	return retval, nil
}
