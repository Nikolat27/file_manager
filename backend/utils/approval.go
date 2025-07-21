package utils

type ApprovalError struct {
	Type    string
	Message string
}

func (e *ApprovalError) Error() string {
	return e.Message
}
