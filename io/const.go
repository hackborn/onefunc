package io

import (
	"errors"

	"github.com/hackborn/onefunc/sys"
)

// Variables used in strings.
const (
	AppPathVar      = "$" + sys.AppPath + "$"
	AppDocPathVar   = "$" + sys.AppDocumentsPath + "$"
	AppCachePathVar = "$" + sys.AppCachePath + "$"
)

var SkipFile = errors.New("skip this file")

var sysvars = []string{AppPathVar, AppDocPathVar, AppCachePathVar}
