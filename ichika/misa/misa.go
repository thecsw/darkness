package misa

import (
	"github.com/charmbracelet/log"
	"github.com/thecsw/darkness/v3/emilia/puck"
	"github.com/thecsw/darkness/v3/ichika/kuroko"
)

// Logger is the logger for Akane.
var logger = puck.NewLogger("Misa 🍎", log.InfoLevel)

func initLog() {
	logger = puck.NewLogger("Misa 🍎", kuroko.LogLevel())
}
