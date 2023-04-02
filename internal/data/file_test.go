package data_test

import (
	"context"
	"github.com/danesparza/fxaudio/internal/data"
	"github.com/sanity-io/litter"
	"github.com/spf13/viper"
	"testing"
)

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
