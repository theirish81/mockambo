package exceptions

// MockamboError is an error meant to track a dev-defined stack to simplify reporting to the user
type MockamboError struct {
	EType   string   `json:"type"`
	Stack   []string `json:"stack,omitempty"`
	Source  string   `json:"source"`
	Message string   `json:"message"`
}

// Wrap wraps an error, whatever it may be, within a MockamboError. If, however, the passed error is
// a MockamboError itself, then the previous EType gets pushed to the stack and gets replaced with the new one, while
// all other fields remain untouched
func Wrap(eType string, origin error) error {
	if origin == nil {
		return nil
	}
	if e, ok := origin.(*MockamboError); ok {
		e.Stack = append(e.Stack, e.EType)
		e.EType = eType
		return e
	}
	return &MockamboError{
		Message: origin.Error(),
		EType:   eType,
		Stack:   make([]string, 0),
		Source:  "mockambo",
	}
}

func (e MockamboError) Error() string {
	return "(" + e.EType + ") " + e.Message
}
