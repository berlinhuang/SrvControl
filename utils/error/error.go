package error

import "fmt"

type ErrCode uint32

const (
	ErrCodeNo       ErrCode = 0x0
	ErrCodeProtocol ErrCode = 0x1
)

var errCodeName = map[ErrCode]string{
	ErrCodeNo:       "NO_ERROR",
	ErrCodeProtocol: "PROTOCOL_ERROR",
}

func (e ErrCode) String() string {
	if s, ok := errCodeName[e]; ok {
		return s
	}
	return fmt.Sprintf("unknown error code 0x%x", uint32(e))
}

////////////////////////////////////////////////////////
type bufferError struct {
	Code   ErrCode
	Reason string
}

func (e bufferError) Error() string {
	return fmt.Sprintf("buffer error: %v, %v", e.Code, e.Reason)
}

////////////////////////////

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func New(text string) error {
	return &errorString{text}
}

var ErrOverBuffer = New("over buffer")
var ErrParamNotExist = New("param not exists ")
