package data_test

import (
	"context"
	"github.com/danesparza/fxaudio/internal/data"
	"github.com/sanity-io/litter"
	"github.com/spf13/viper"
	"reflect"
	"testing"
)

/*
func TestFile_AddFile_ValidFile_Successful(t *testing.T) {

	//	Arrange
	systemdb := getTestFiles()

	db, err := data2.NewManager(systemdb)
	if err != nil {
		t.Fatalf("NewManager failed: %s", err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(systemdb)
	}()

	testFile := data2.File{FilePath: "crossbones1.mp3", Description: "Unit test file"}

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

	db, err := data2.NewManager(systemdb)
	if err != nil {
		t.Fatalf("NewManager failed: %s", err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(systemdb)
	}()

	testFile1 := data2.File{FilePath: "crossbones1.mp3", Description: "Unit test 1 file"}
	testFile2 := data2.File{FilePath: "crossbones2.mp3", Description: "Unit test 2 file"}
	testFile3 := data2.File{FilePath: "crossbones3.mp3", Description: "Unit test 3 file"}

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

	db, err := data2.NewManager(systemdb)
	if err != nil {
		t.Fatalf("NewManager failed: %s", err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(systemdb)
	}()

	testFile1 := data2.File{FilePath: "crossbones1.mp3", Description: "Unit test 1 file"}
	testFile2 := data2.File{FilePath: "crossbones2.mp3", Description: "Unit test 2 file"}
	testFile3 := data2.File{FilePath: "crossbones3.mp3", Description: "Unit test 3 file"}

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

	db, err := data2.NewManager(systemdb)
	if err != nil {
		t.Fatalf("NewManager failed: %s", err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(systemdb)
	}()

	testFile1 := data2.File{FilePath: "crossbones1.mp3", Description: "Unit test 1 file"}
	testFile2 := data2.File{FilePath: "crossbones2.mp3", Description: "Unit test 2 file"}
	testFile3 := data2.File{FilePath: "crossbones3.mp3", Description: "Unit test 3 file"}

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
*/

func Test_appDataService_AddFile(t *testing.T) {
	type args struct {
		ctx         context.Context
		filepath    string
		description string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid arguments - saves to database",
			args: args{
				ctx:         context.TODO(),
				filepath:    "testfile1.mp3",
				description: "This is a test file",
			},
			wantErr: false,
		},
	}

	//	Initialize sqlite
	db, err := data.InitSqlite(viper.GetString("datastore.system"))
	if err != nil {
		t.Errorf("Problem initializing database: %v", err)
	}

	a := data.NewAppDataService(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := a.AddFile(tt.args.ctx, tt.args.filepath, tt.args.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			litter.Dump(got)
		})
	}
}

func Test_appDataService_GetFile(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    data.File
		wantErr bool
	}{
		{
			name: "Valid arguments - select from db",
			args: args{
				ctx: context.TODO(),
				id:  "cgk8ckj511ion26hceh0",
			},
			want: data.File{
				ID:          "cgk8ckj511ion26hceh0",
				Created:     1680377426,
				FilePath:    "testfile1.mp3",
				Description: "This is a test file",
				Tags:        nil,
			},
			wantErr: false,
		},
	}

	//	Initialize sqlite
	db, err := data.InitSqlite(viper.GetString("datastore.system"))
	if err != nil {
		t.Errorf("Problem initializing database: %v", err)
	}

	a := data.NewAppDataService(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := a.GetFile(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
