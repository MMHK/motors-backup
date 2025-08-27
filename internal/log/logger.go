package log

import (
	"github.com/op/go-logging"
)

var Logger = logging.MustGetLogger("MOTORS_BACKUP")

func init() {
	format := logging.MustStringFormatter(
		`MOTORS_BACKUP %{shortfunc} %{level:.4s} %{shortfile}
%{id:03x} %{message}`,
	)
	logging.SetFormatter(format)
}
