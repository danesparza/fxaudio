package data

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tidwall/buntdb"
)

// Manager is the data manager
type Manager struct {
	systemdb *buntdb.DB
}

// NewManager creates a new instance of a Manager and returns it
func NewManager(systemdbpath string) (*Manager, error) {
	retval := new(Manager)

	//	Make sure the path already exists:
	if err := os.MkdirAll(filepath.Dir(systemdbpath), os.FileMode(0664)); err != nil {
		return nil, err
	}

	sysdb, err := buntdb.Open(systemdbpath)
	if err != nil {
		return retval, fmt.Errorf("problem opening the systemDB: %s", err)
	}
	retval.systemdb = sysdb

	//	Create our indexes
	sysdb.CreateIndex("Event", "Event:*", buntdb.IndexString)

	//	Return our Manager reference
	return retval, nil
}

// Close closes the data Manager
func (store Manager) Close() error {
	syserr := store.systemdb.Close()

	if syserr != nil {
		return fmt.Errorf("an error occurred closing the manager.  Syserr: %s ", syserr)
	}

	return nil
}

// GetKey returns a key to be used in the storage system
func GetKey(entityType string, keyPart ...string) string {
	allparts := []string{}
	allparts = append(allparts, entityType)
	allparts = append(allparts, keyPart...)
	return strings.Join(allparts, ":")
}
