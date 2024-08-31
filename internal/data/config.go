package data

import (
	"context"
	"fmt"
)

// SystemConfig represents the system configuration information
type SystemConfig struct {
	AlsaDevice string `db:"alsa_device" json:"alsa_device"`
}

// GetConfig gets the system configuration from the database
func (a appDataService) GetConfig(ctx context.Context) (SystemConfig, error) {
	retval := SystemConfig{}
	query := `select alsa_device from system_config limit 1;`

	err := a.DB.GetContext(ctx, &retval, query)
	if err != nil {
		return retval, fmt.Errorf("getConfig: %w", err)
	}

	return retval, nil
}

// SetConfig updates the system configuration
func (a appDataService) SetConfig(ctx context.Context, config SystemConfig) error {
	query := `update system_config set alsa_device = $1;`

	_, err := a.DB.ExecContext(ctx, query, config.AlsaDevice)
	if err != nil {
		return fmt.Errorf("setConfig: %w", err)
	}

	return nil
}
