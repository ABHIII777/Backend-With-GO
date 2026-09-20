package main

type UserInputStruct struct {
	Method string
	// Target string
	Path string
	Query map[string] string
	Version string
	Headers map[string] string
	Body string
}