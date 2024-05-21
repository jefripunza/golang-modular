package variable

import (
	"core/env"
	"path/filepath"
)

var DatabasePath = filepath.Join(env.GetPwd(), "database")
var TempPath = filepath.Join(env.GetPwd(), "temp")
