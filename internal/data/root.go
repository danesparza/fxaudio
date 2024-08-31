package data

import (
	"context"
	"github.com/danesparza/fxaudio/scripts"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Add the file migrations source to golang-migrate
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	_ "modernc.org/sqlite" // Register relevant drivers.
	"os"
	"path/filepath"
)

type AppDataService interface {
	AddFile(ctx context.Context, filepath, description string) (File, error)
	GetFile(ctx context.Context, id string) (File, error)
	GetAllFiles(ctx context.Context) ([]File, error)
	GetAllFilesWithTag(ctx context.Context, tag string) ([]File, error)
	DeleteFile(ctx context.Context, id string) error
	UpdateTags(ctx context.Context, id string, tags []string) error
	GetConfig(ctx context.Context) (SystemConfig, error)
	SetConfig(ctx context.Context, config SystemConfig) error
}

type appDataService struct {
	*sqlx.DB
}

// InitConfig performs runtime config
func (a appDataService) InitRuntimeConfig() error {

	/* Turn on foreign key support (I can't believe we have to do this) */
	/* More information: https://www.sqlite.org/foreignkeys.html */
	_, err := a.DB.Exec(`PRAGMA foreign_keys = ON;`)
	if err != nil {
		return err
	}

	return nil
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

	//	Create a 'databaseDriver' object from the existing connection
	databaseDriver, err := sqlite.WithInstance(db.DB, &sqlite.Config{})
	if err != nil {
		log.Fatal().Str("datasource", datasource).Err(err).Msg("problem setting up databaseDriver for migrations")
	}

	//	Create the new migration source
	sourceDriver, err := iofs.New(scripts.FS, "sqlite/migrations")
	if err != nil {
		log.Fatal().Err(err).
			Str("datasource", datasource).
			Str("migrationsource", "(iofs)sqlite/migrations").
			Msg("problem setting up migration source")
	}

	//	Create a new migrator with the databaseDriver (and existing connection)
	migrator, err := migrate.NewWithInstance("iofs", sourceDriver, datasource, databaseDriver)
	if err != nil {
		log.Fatal().
			Str("datasource", datasource).
			Str("migrationsource", "(iofs)sqlite/migrations").
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
