package errors

type ErrSessionDataMissing struct{}

func (e ErrSessionDataMissing) Error() string {
	return "session data missing"
}
