package data_test

import (
	"os"
	"testing"
	"time"

	"github.com/danesparza/fxaudio/data"
	"github.com/danesparza/fxaudio/event"
	"github.com/danesparza/fxaudio/mediatype"
)

func TestEvent_AddEvent_ValidEvent_Successful(t *testing.T) {

	//	Arrange
	systemdb := getTestFiles()

	db, err := data.NewManager(systemdb)
	if err != nil {
		t.Fatalf("NewManager failed: %s", err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(systemdb)
	}()

	testEvent := data.Event{EventType: event.SystemStartup, MediaType: "Test", Details: "Unit test details"}

	//	Act
	newEvent, err := db.AddEvent(testEvent.EventType, testEvent.MediaType, testEvent.Details, "127.0.0.1", 2*time.Hour)

	//	Assert
	if err != nil {
		t.Errorf("AddEvent - Should add event without error, but got: %s", err)
	}

	if newEvent.Created.IsZero() {
		t.Errorf("AddEvent failed: Should have set an item with the correct datetime: %+v", newEvent)
	}

}

func TestEvent_GetEvent_ValidEvent_Successful(t *testing.T) {

	//	Arrange
	systemdb := getTestFiles()

	db, err := data.NewManager(systemdb)
	if err != nil {
		t.Fatalf("NewManager failed: %s", err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(systemdb)
	}()

	testEvent := data.Event{EventType: event.SystemStartup, MediaType: mediatype.Unknown, Details: "Unit test details"}

	//	Act
	newEvent, _ := db.AddEvent(testEvent.EventType, testEvent.MediaType, testEvent.Details, "127.0.0.1", 2*time.Hour)
	gotEvent, err := db.GetEvent(newEvent.ID)

	//	Assert
	if err != nil {
		t.Errorf("GetEvent - Should get event without error, but got: %s", err)
	}

	if newEvent.Created.IsZero() {
		t.Errorf("GetEvent failed: Should get an item with the correct datetime: %+v", gotEvent)
	}

	if newEvent.Details != gotEvent.Details || gotEvent.Details == "" {
		t.Errorf("GetEvent failed: Should get an item with the correct details: %+v", gotEvent)
	}

}

func TestEvent_GetEvent_ExpiredEvent_ReturnsError(t *testing.T) {

	//	Arrange
	systemdb := getTestFiles()

	db, err := data.NewManager(systemdb)
	if err != nil {
		t.Fatalf("NewManager failed: %s", err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(systemdb)
	}()

	testEvent := data.Event{EventType: event.SystemStartup, MediaType: "Test", Details: "Unit test details"}

	//	Act
	newEvent, _ := db.AddEvent(testEvent.EventType, testEvent.MediaType, testEvent.Details, "127.0.0.1", 2*time.Second)

	//	-- Wait for 5 seconds -- TTL should expire and the event should no longer be available:
	time.Sleep(5 * time.Second)

	_, err = db.GetEvent(newEvent.ID)

	//	Assert
	if err == nil {
		t.Errorf("GetEvent - Should get error for expired event, but didn't get error")
	}

}

func TestEvent_GetAllEvents_ValidEvents_Successful(t *testing.T) {

	//	Arrange
	systemdb := getTestFiles()

	db, err := data.NewManager(systemdb)
	if err != nil {
		t.Fatalf("NewManager failed: %s", err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(systemdb)
	}()

	testEvent1 := data.Event{EventType: event.SystemStartup, MediaType: mediatype.Audio, Details: "Unit test 1 details"}
	testEvent2 := data.Event{EventType: event.SystemStartup, MediaType: mediatype.Audio, Details: "Unit test 2 details"}
	testEvent3 := data.Event{EventType: event.SystemStartup, MediaType: mediatype.Audio, Details: "Unit test 3 details"}

	//	Act
	db.AddEvent(testEvent1.EventType, testEvent1.MediaType, testEvent1.Details, "127.0.0.1", 2*time.Hour)
	newEvent2, _ := db.AddEvent(testEvent2.EventType, testEvent2.MediaType, testEvent2.Details, "127.0.0.1", 2*time.Hour)
	db.AddEvent(testEvent3.EventType, testEvent3.MediaType, testEvent3.Details, "127.0.0.1", 2*time.Hour)
	gotEvents, err := db.GetAllEvents()

	//	Assert
	if err != nil {
		t.Errorf("GetAllEvents - Should get all events without error, but got: %s", err)
	}

	if len(gotEvents) < 2 {
		t.Errorf("GetAllEvents failed: Should get all items but got: %v", len(gotEvents))
	}

	if gotEvents[1].Details != newEvent2.Details {
		t.Errorf("GetAllEvents failed: Should get an item with the correct details: %+v", gotEvents[1])
	}

}
