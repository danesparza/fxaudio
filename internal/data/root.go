package data

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Add the file migrations source to golang-migrate
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	_ "modernc.org/sqlite" // Register relevant drivers.
	"os"
	"path/filepath"
)

type AppDataService interface {
	AddFile(ctx context.Context, filepath, description string) (File, error)
	GetFile(ctx context.Context, id string) (File, error)
	GetAllFiles(ctx context.Context) ([]File, error)
	DeleteFile(ctx context.Context, id string) error
	UpdateTags(ctx context.Context, id string, tags []string) error
}

type appDataService struct {
	*sqlx.DB
}

func NewAppDataService(db *sqlx.DB) AppDataService {
	return &appDataService{db}
}

// InitSqlite initializes SQLite and returns a pointer to the db
// object, or an error
func InitSqlite(datasource string) (*sqlx.DB, error) {
	log.Info().Msg("InitSqlite")

	//	Make sure the path is created:
	err := os.MkdirAll(filepath.Dir(datasource), 0777)
	if err != nil {
		log.Fatal().Err(err).Str("datasource", datasource).Msg("Problem creating datasource directory")
	}

	//	Connect to the datasource
	db, err := sqlx.Open("sqlite", datasource)
	if err != nil {
		log.Fatal().Err(err).Str("datasource", datasource).Msg("Problem connecting to SQLite database")
	}

	//	Create a 'driver' object from the existing connection
	driver, err := sqlite.WithInstance(db.DB, &sqlite.Config{})
	if err != nil {
		log.Fatal().Str("datasource", datasource).Err(err).Msg("problem setting up driver for migrations")
	}

	//	Format migration source:
	migrationSource := fmt.Sprintf("file://%s", viper.GetString("datastore.migrationsource"))

	//	Create a new migrator with the driver (and existing connection)
	migrator, err := migrate.NewWithDatabaseInstance(
		migrationSource,
		datasource, driver)
	if err != nil {
		log.Fatal().
			Str("datasource", datasource).
			Str("migrationSource", migrationSource).
			Err(err).Msg("problem creating migrator config")
	}

	//	Run the migrations
	err = migrator.Up()
	switch err {
	case migrate.ErrNoChange:
		log.Info().Msg("sqlite schema is up-to-date")
	case nil:
		log.Info().Msg("sqlite schema was updated successfully")
	default:
		log.Err(err).Msg("problem running migrations")
		return db, err
	}

	return db, nil
}
