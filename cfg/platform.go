package cfg

func ApplicationDataPath() string {
	return appDataPath
}

var (
	// Local path to the app data folder. Each supported platform needs to set
	// this in an init()
	appDataPath = ""
)
