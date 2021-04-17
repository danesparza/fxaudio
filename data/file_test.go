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

func TestFile_GetAllFiles_ValidFiles_Successful(t *testing.T) {

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

	testFile1 := data.File{FilePath: "crossbones1.mp3", Description: "Unit test 1 file"}
	testFile2 := data.File{FilePath: "crossbones2.mp3", Description: "Unit test 2 file"}
	testFile3 := data.File{FilePath: "crossbones3.mp3", Description: "Unit test 3 file"}

	//	Act
	db.AddFile(testFile1.FilePath, testFile1.Description)
	newFile2, _ := db.AddFile(testFile2.FilePath, testFile2.Description)
	db.AddFile(testFile3.FilePath, testFile3.Description)

	gotFiles, err := db.GetAllFiles()

	//	Assert
	if err != nil {
		t.Errorf("GetAllFiles - Should get all files without error, but got: %s", err)
	}

	if len(gotFiles) < 2 {
		t.Errorf("GetAllFiles failed: Should get all items but got: %v", len(gotFiles))
	}

	if gotFiles[1].Description != newFile2.Description {
		t.Errorf("GetAllFiles failed: Should get an item with the correct details: %+v", gotFiles[1])
	}

}
