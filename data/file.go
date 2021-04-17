package data

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/xid"
	"github.com/tidwall/buntdb"
)

// File represents an uploaded file.
type File struct {
	ID          string    `json:"id"`       // Unique File ID
	Created     time.Time `json:"created"`  // File create time
	FilePath    string    `json:"filepath"` // Full filepath to the file
	Description string    `json:"details"`  // Additional information (like the files involved)
}

// AddEvent adds an event to the system
func (store Manager) AddFile(filepath, description string) (File, error) {
	//	Our return item
	retval := File{}

	newFile := File{
		ID:          xid.New().String(), // Generate a new id
		Created:     time.Now(),
		FilePath:    filepath,
		Description: description,
	}

	//	Serialize to JSON format
	encoded, err := json.Marshal(newFile)
	if err != nil {
		return retval, fmt.Errorf("problem serializing the data: %s", err)
	}

	//	Save it to the database:
	err = store.systemdb.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(GetKey("File", newFile.ID), string(encoded), &buntdb.SetOptions{})
		return err
	})

	//	If there was an error saving the data, report it:
	if err != nil {
		return retval, fmt.Errorf("problem saving the file: %s", err)
	}

	//	Set our retval:
	retval = newFile

	//	Return our data:
	return retval, nil
}

// GetAllFiles gets all files in the system
func (store Manager) GetAllFiles() ([]File, error) {
	//	Our return item
	retval := []File{}

	//	Set our prefix
	prefix := GetKey("File")

	//	Iterate over our values:
	err := store.systemdb.View(func(tx *buntdb.Tx) error {
		tx.Descend(prefix, func(key, val string) bool {

			if len(val) > 0 {
				//	Create our item:
				item := File{}

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
		return retval, fmt.Errorf("problem getting the list of files: %s", err)
	}

	//	Return our data:
	return retval, nil
}

// DeleteFile deletes a file from the system
func (store Manager) DeleteFile(id string) error {

	//	Remove it from the database:
	err := store.systemdb.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(GetKey("File", id))
		return err
	})

	//	If there was an error removing the data, report it:
	if err != nil {
		return fmt.Errorf("problem removing the file: %s", err)
	}

	//	Return our data:
	return nil
}
