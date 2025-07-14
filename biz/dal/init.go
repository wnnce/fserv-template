package dal

import (
	"github.com/wnnce/fserv-template/config"
)

// init registers MongoDB and Postgresql initialization functions with the
// global configuration manager. This ensures that the corresponding database
// connections are established when the application starts.
func init() {
	config.RegisterConfigureReaders()
}
