package io

import (
	"github.com/hackborn/onefunc/sys"
)

// Variables used in strings.
const (
	AppPathVar      = "$" + sys.AppPath + "$"
	AppDocPathVar   = "$" + sys.AppDocumentsPath + "$"
	AppCachePathVar = "$" + sys.AppCachePath + "$"
)

var sysvars = []string{AppPathVar, AppDocPathVar, AppCachePathVar}
