package data

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	_ "modernc.org/sqlite"
	"os"
	"path"
)

type AppDataService interface {
	AddFile(ctx context.Context, filepath, description string) (File, error)
	GetFile(ctx context.Context, id string) (File, error)
	GetAllFiles(ctx context.Context) ([]File, error)
	DeleteFile(ctx context.Context, id string) error
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
	err := os.MkdirAll(datasource, 0777)
	if err != nil {
		log.Fatal().Err(err).Str("datasource", datasource).Msg("Problem creating datasource directory")
	}

	//	Connect to the datasource
	dbname := path.Join(datasource, "fxaudio.db")
	db, err := sqlx.Open("sqlite", dbname)
	if err != nil {
		log.Fatal().Err(err).Str("dbname", dbname).Msg("Problem connecting to SQLite database")
	}

	//	Run migrations
	driver, err := sqlite.WithInstance(db.DB, &sqlite.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("problem setting up driver for migrations")
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		viper.GetString("datastore.migrationsource"),
		dbname, driver)
	if err != nil {
		log.Fatal().Err(err).Msg("problem creating migrator config")
	}

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
