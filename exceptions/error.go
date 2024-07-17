package exceptions

type MockamboError struct {
	EType   string   `json:"type"`
	Stack   []string `json:"stack,omitempty"`
	Source  string   `json:"source"`
	Message string   `json:"message"`
}

func Wrap(eType string, origin error) *MockamboError {
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
	return e.Message
}
