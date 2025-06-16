package internal

import (
	"os"
)

type Instance struct {
	*Logger
	*Shutdown
	Info *Info
}

var program *Instance = &Instance{
	Logger: NewLogger(os.Stdout),
	Info:   NewInfo(),
}

// New executes only one time and then returns the same Instance everytime.
// On execution parses all flags and stores them
func NewInstance() *Instance {

	return program
}
