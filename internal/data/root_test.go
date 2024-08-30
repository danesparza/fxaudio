package data_test

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"testing"
)

func init() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.SetDefault("datastore.system", filepath.Join(home, "fxaudio", "db", "fxaudio.db"))
}

func TestSliceToJSONArray(t *testing.T) {
	//	Arrange
	testSlice := []string{"one", "two", "three"}
	expectedString := `["one","two","three"]`

	//	Act
	jsonArray, _ := json.Marshal(testSlice)
	jsonString := string(jsonArray)

	//	Assert
	if jsonString != expectedString {
		t.Errorf("Problem serializing to JSON.  Got %s", jsonString)
	}

}
