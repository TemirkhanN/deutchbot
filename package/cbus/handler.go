package cbus

type Input string

type Output interface {
	Write(output string)
	Writeln(output string)
}

type Handler interface {
	Handle(input Input, o Output)
}
