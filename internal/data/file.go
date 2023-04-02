package data

import (
	"context"
	"fmt"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	"time"
)

// File represents an uploaded file.
type File struct {
	ID          string   `json:"id"`       // Unique File ID
	Created     int64    `json:"created"`  // File create time
	FilePath    string   `json:"filepath"` // Full filepath to the file
	Description string   `json:"details"`  // Additional information (like the files involved)
	Tags        []string `json:"tags"`     // Tag list to associate with this file
}

func (a appDataService) AddFile(ctx context.Context, filepath, description string) (File, error) {
	//	Our return item
	retval := File{
		ID:          xid.New().String(), // Generate a new id
		Created:     time.Now().Unix(),
		FilePath:    filepath,
		Description: description,
	}

	query := `insert into media(id, filepath, description, created) 
				values($1, $2, $3, $4);`

	stmt, err := a.DB.PrepareContext(ctx, query)
	if err != nil {
		return retval, err
	}

	_, err = stmt.ExecContext(ctx, retval.ID, retval.FilePath, retval.Description, retval.Created)
	if err != nil {
		return retval, fmt.Errorf("problem executing query: %v", err)
	}

	return retval, nil
}

// GetFile gets information about a single file in the system based on its id
func (a appDataService) GetFile(ctx context.Context, id string) (File, error) {
	//	Our return item
	retval := File{}

	query := `select id, filepath, description, created 
				from media 
				where id = $1;`

	stmt, err := a.DB.PreparexContext(ctx, query)
	if err != nil {
		return retval, err
	}

	rows, err := stmt.QueryxContext(ctx, id)
	if err != nil {
		return retval, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Err(closeErr).Msg("unable to close rows")
		}
	}()

	for rows.Next() {
		err = rows.StructScan(&retval)
		if err != nil {
			return retval, fmt.Errorf("problem reading into struct: %v", err)
		}
	}

	//	Return our data:
	return retval, nil
}

// GetAllFiles gets all files in the system
func (a appDataService) GetAllFiles(ctx context.Context) ([]File, error) {
	//	Our return item
	retval := []File{}

	query := `select id, filepath, description, created 
				from media;`

	stmt, err := a.DB.PreparexContext(ctx, query)
	if err != nil {
		return retval, err
	}

	rows, err := stmt.QueryxContext(ctx)
	if err != nil {
		return retval, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Err(closeErr).Msg("unable to close rows")
		}
	}()

	for rows.Next() {
		item := File{}
		err = rows.StructScan(&item)
		if err != nil {
			return retval, fmt.Errorf("problem reading into struct: %v", err)
		}

		retval = append(retval, item)
	}

	//	Return our data:
	return retval, nil
}

// DeleteFile deletes a file from the system
func (a appDataService) DeleteFile(ctx context.Context, id string) error {

	query := `delete from media where id = $1;`

	stmt, err := a.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("problem executing query: %v", err)
	}

	return nil
}
