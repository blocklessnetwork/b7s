package codes

// Code represents the status of an action. It is a rough equivalent of the HTTP status code.
type Code string

// Response codes.
const (
	OK        Code = "200"
	Accepted  Code = "202"
	NoContent Code = "204"

	Invalid       Code = "400"
	NotAuthorized Code = "401"
	NotPermitted  Code = "403"
	NotFound      Code = "404"
	Timeout       Code = "408"

	Error          Code = "500"
	NotImplemented Code = "501"
	NotAvailable   Code = "503"
	NotSupported   Code = "505"
	Unknown        Code = "520"
)

func (c Code) String() string {
	return string(c)
}
