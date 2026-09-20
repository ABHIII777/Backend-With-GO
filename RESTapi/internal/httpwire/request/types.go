package request

type HTTPRequestStruct struct {
	Method  string
	Path string
	Query map[string] string
	Version string
	Headers map[string]string
	Body    string
}
