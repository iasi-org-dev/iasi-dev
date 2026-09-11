// Package RC contains the cumulative return-code bitmask used by iasi-dev.
package RC

const (
	OK          = 0x00
	NothingToDo = 0x01
	Info        = 0x02
	Warning     = 0x04
	Attention   = 0x08

	Error    = 0x10
	Severe   = 0x20
	Fatal    = 0x40
	Reserved = 0x80

	NoticeMask = 0x0F
	ErrorMask  = 0xF0

	// Internal workflow marker. It is never exposed as the process exit code.
	Skip = 0x100
)

// Stop is used internally to unwind execution while preserving the accumulated RC.
type Stop struct {
	Code int
}

func IsErroneous(rc int) bool {
	return rc&ErrorMask != 0
}

func Add(current *int, rc int) int {
	if current == nil {
		return rc
	}

	*current |= rc
	return *current
}

func Value(current *int) int {
	if current == nil {
		return OK
	}

	return *current
}
