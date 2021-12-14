package data_test

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	math_rand "math/rand"
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

func TestRoot_Rand_Test(t *testing.T) {
	//	Arrange
	maxNumber := 2
	numTests := 100

	//	Spit out several tests
	for j := 0; j <= numTests; j++ {
		testnum := math_rand.Intn(maxNumber)
		t.Logf("Random number between 0 and %v: %v", maxNumber, testnum)
		if testnum > 0 {
			t.Logf("WOW!  It's greater than zero -- that's amazing ***** ")
		}
	}
}

func init() {
	var b [8]byte
	crypto_rand.Read(b[:])
	math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}
