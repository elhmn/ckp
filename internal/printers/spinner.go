package printers

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

//ISpinner defines a spinner interface
type ISpinner interface {
	Message(m string)
	Clear()
	Stop()
	Start()
}

//Spinner is a spinner object
type Spinner struct {
	s *spinner.Spinner
}

//NewSpinner returns a new Spinner
func NewSpinner() *Spinner {
	spin := Spinner{s: spinner.New(spinner.CharSets[11], 100*time.Millisecond)}
	return &spin
}

//Message will set the spinner suffix message
func (s Spinner) Message(m string) {
	s.s.Suffix = fmt.Sprintf(" %s", m)
}

//Clear will clear the spinner suffix
func (s Spinner) Clear() {
	s.s.Suffix = ""
}

//Stop stops the spinner
func (s Spinner) Stop() {
	s.s.Stop()
}

//Start starts the spinner
func (s Spinner) Start() {
	s.s.Start()
}
