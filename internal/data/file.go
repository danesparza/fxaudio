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

/*

// GetAllFiles gets all files in the system
func (store Manager) GetAllFiles() ([]File, error) {
	//	Our return item
	retval := []File{}

	//	Set our prefix
	prefix := GetKey("File")

	//	Iterate over our values:
	err := store.systemdb.View(func(tx *buntdb.Tx) error {
		tx.Descend(prefix, func(key, val string) bool {

			if len(val) > 0 {
				//	Create our item:
				item := File{}

				//	Unmarshal data into our item
				bval := []byte(val)
				if err := json.Unmarshal(bval, &item); err != nil {
					return false
				}

				//	Add to the array of returned users:
				retval = append(retval, item)
			}

			return true
		})
		return nil
	})

	//	If there was an error, report it:
	if err != nil {
		return retval, fmt.Errorf("problem getting the list of files: %s", err)
	}

	//	Return our data:
	return retval, nil
}

// DeleteFile deletes a file from the system
func (store Manager) DeleteFile(id string) error {

	//	Remove it from the database:
	err := store.systemdb.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(GetKey("File", id))
		return err
	})

	//	If there was an error removing the data, report it:
	if err != nil {
		return fmt.Errorf("problem removing the file: %s", err)
	}

	//	Return our data:
	return nil
}
*/
