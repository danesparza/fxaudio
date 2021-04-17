package data_test

import (
	"os"
	"testing"

	"github.com/danesparza/fxaudio/data"
)

func TestFile_AddFile_ValidFile_Successful(t *testing.T) {

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

	testFile := data.File{FilePath: "crossbones1.mp3", Description: "Unit test file"}

	//	Act
	newFile, err := db.AddFile(testFile.FilePath, testFile.Description)

	//	Assert
	if err != nil {
		t.Errorf("AddFile - Should add file without error, but got: %s", err)
	}

	if newFile.Created.IsZero() {
		t.Errorf("AddFile failed: Should have set an item with the correct datetime: %+v", newFile)
	}

}
