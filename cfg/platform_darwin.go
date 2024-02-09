package cfg

import (
	"fmt"
	"os"
	"path/filepath"
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		AddInitErr(fmt.Errorf("cfg init() can't get homedir: %w", err))
		return
	}
	appDataPath = filepath.Join(homeDir, "Library", "Application Support", "Daytrader")
	err = os.Mkdir(appDataPath, 0750)
	if err != nil && !os.IsExist(err) {
		AddInitErr(err)
		return
	}
}
