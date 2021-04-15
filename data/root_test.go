package data_test

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
)

//	Gets the database path for this environment:
func getTestFiles() string {
	systemdb := os.Getenv("FXAUDIO_TEST_ROOT")

	if systemdb == "" {
		home, _ := homedir.Dir()
		if home != "" {
			systemdb = path.Join(home, "fxaudio", "db", "system.db")
		}
	}
	return systemdb
}

func TestRoot_GetTestDBPaths_Successful(t *testing.T) {

	systemdb := getTestFiles()

	if systemdb == "" {
		t.Fatal("The required FXAUDIO_TEST_ROOT environment variable is not set to the test database root path.  It should probably be $HOME/fxaudio/db/system.db")
	}

	t.Logf("System db path: %s", systemdb)
	t.Logf("System db folder: %s", filepath.Dir(systemdb))
}

func TestRoot_Databases_ShouldNotExistYet(t *testing.T) {
	//	Arrange
	systemdb := getTestFiles()

	//	Act

	//	Assert
	if _, err := os.Stat(systemdb); err == nil {
		t.Errorf("System database check failed: System db %s already exists, and shouldn't", systemdb)
	}
}
