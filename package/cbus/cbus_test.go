package cbus_test

import (
	"DeutschBot/package/cbus"

	"github.com/stretchr/testify/assert"

	"testing"
)

type simpleHandler struct {
}

func (sh simpleHandler) Handle(i cbus.Input, o cbus.Output) {
	o.Write("World")
}

type outputSpy struct {
	lastOutput string
}

func (o *outputSpy) Write(text string) {
	o.lastOutput = text
}

func (o *outputSpy) Writeln(text string) {
	o.lastOutput = text + "\n"
}

func TestCommandBus_HandleWithoutHandlers(t *testing.T) {
	output := &outputSpy{}
	commandBus := cbus.NewCommandBus(output)

	commandBus.Handle("Hello")

	assert.Empty(t, output.lastOutput)
}

func TestCommandBus_HandleUnknownCommand(t *testing.T) {
	output := &outputSpy{}
	commandBus := cbus.NewCommandBus(output)
	commandBus.RegisterHandler(cbus.NewHandlerDefinition(simpleHandler{}, func(i cbus.Input) bool {
		return false
	}))

	commandBus.Handle("Hello")

	assert.Empty(t, output.lastOutput)
}

func TestCommandBus_Handle(t *testing.T) {
	output := &outputSpy{}
	commandBus := cbus.NewCommandBus(output)
	commandBus.RegisterHandler(cbus.NewHandlerDefinition(simpleHandler{}, func(i cbus.Input) bool {
		return true
	}))

	commandBus.Handle("Hello")

	assert.Equal(t, "World", output.lastOutput)
}
