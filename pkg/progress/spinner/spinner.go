package spinner

import (
	"io"

	"github.com/Phillezi/kthcloud-cli/pkg/scheduler"
	"github.com/pterm/pterm"
)

type SpinnerVariant int64

const (
	Spin SpinnerVariant = iota
	Waiting
	Done
	Errored
)

type Spinner struct {
	Printer   *pterm.SpinnerPrinter
	PrevState scheduler.JobState
}

func New(writer io.Writer, state scheduler.JobState, colors ...pterm.Color) *Spinner {
	return &Spinner{
		Printer: func() *pterm.SpinnerPrinter {
			switch state {
			case scheduler.Created:
				return waitingSpinner(writer, colors...)
			case scheduler.Done:
				return doneCheckmark(writer, colors...)
			case scheduler.Cancelled, scheduler.Errored:
				return errorCross(writer, colors...)
			default:
				return NewSpinner(writer, colors...)
			}
		}(),
	}
}

func NewSpinner(writer io.Writer, colors ...pterm.Color) *pterm.SpinnerPrinter {
	return pterm.DefaultSpinner.WithWriter(writer).WithSequence("⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏").WithStyle(pterm.NewStyle(colors...))
}

func doneCheckmark(writer io.Writer, colors ...pterm.Color) *pterm.SpinnerPrinter {
	return pterm.DefaultSpinner.WithWriter(writer).WithSequence("DONE").WithStyle(pterm.NewStyle(colors...))
}

func errorCross(writer io.Writer, colors ...pterm.Color) *pterm.SpinnerPrinter {
	return pterm.DefaultSpinner.WithWriter(writer).WithSequence("X").WithStyle(pterm.NewStyle(colors...))
}

func waitingSpinner(writer io.Writer, colors ...pterm.Color) *pterm.SpinnerPrinter {
	return pterm.DefaultSpinner.WithWriter(writer).WithSequence(".", "..", "...").WithStyle(pterm.NewStyle(colors...))
}
