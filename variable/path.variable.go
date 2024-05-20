package variable

import (
	"path/filepath"
	"project/env"
)

var DatabasePath = filepath.Join(env.GetPwd(), "database")
var TempPath = filepath.Join(env.GetPwd(), "temp")
