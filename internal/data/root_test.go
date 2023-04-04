package data_test

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func init() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.SetDefault("datastore.system", filepath.Join(home, "fxaudio", "db", "fxaudio.db"))
	viper.SetDefault("datastore.migrationsource", "../../scripts/sqlite/migrations")
}
