package request

type HTTPRequestStruct struct {
	Method  string
	Target  string
	Version string
	Headers map[string]string
	Body    string
}
