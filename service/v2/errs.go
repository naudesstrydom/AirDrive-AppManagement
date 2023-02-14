package v2

import (
	"fmt"

	"github.com/IceWhaleTech/CasaOS-AppManagement/common"
)

var (
	ErrComposeExtensionNameXCasaOSNotFound = fmt.Errorf("extension `%s` not found", common.ComposeExtensionNameXCasaOS)
	ErrMainAppNotFound                     = fmt.Errorf("main app not found")
	ErrComposeAppNotFound                  = fmt.Errorf("compose app not found")
	ErrComposeAppNotMatch                  = fmt.Errorf("compose app not match")
	ErrComposeFileNotFound                 = fmt.Errorf("compose file not found")
)
