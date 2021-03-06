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

func TestFile_GetFile_ValidFile_Successful(t *testing.T) {

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

	gotFile, err := db.GetFile(newFile2.ID)

	//	Log the file details:
	t.Logf("File: %+v", gotFile)

	//	Assert
	if err != nil {
		t.Errorf("GetFile - Should get file without error, but got: %s", err)
	}

	if len(gotFile.ID) < 2 {
		t.Errorf("GetFile failed: Should get valid id but got: %v", gotFile.ID)
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

func TestFile_DeleteFile_ValidFiles_Successful(t *testing.T) {

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
	err = db.DeleteFile(newFile2.ID) //	Delete the 2nd file

	gotFiles, _ := db.GetAllFiles()

	//	Assert
	if err != nil {
		t.Errorf("DeleteFile - Should delete file without error, but got: %s", err)
	}

	if len(gotFiles) != 2 {
		t.Errorf("DeleteFile failed: Should remove an item but got: %v", len(gotFiles))
	}

	if gotFiles[1].Description == newFile2.Description {
		t.Errorf("DeleteFile failed: Should get an item with different details than the removed item but got: %+v", gotFiles[1])
	}

}
