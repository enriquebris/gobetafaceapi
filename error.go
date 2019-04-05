package gobetafaceapi

const (
	ErrorCodeNotImplemented = "missingImplementor"
	ErrorCodeJSONMarshal    = "json.marshal"
	ErrorCodeJSONUnmarshal  = "json.unmarshal"
	ErrorCodeHTTPRequest    = "http.request"
	ErrorCodeArguments      = "arguments"
)

func NewError(code string, message string) *Error {
	return &Error{
		code:    code,
		message: message,
	}
}

type Error struct {
	code    string
	message string
}

func (st *Error) Error() string {
	return st.message
}

func (st *Error) Code() string {
	return st.code
}
