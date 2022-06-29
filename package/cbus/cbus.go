package cbus

type Bus interface {
	RegisterHandler(definition HandlerDefinition)

	Handle(i Input)
}

type HandlerDefinition struct {
	handler  Handler
	resolver func(i Input) bool
}

func NewHandlerDefinition(handler Handler, resolver func(i Input) bool) HandlerDefinition {
	return HandlerDefinition{
		handler:  handler,
		resolver: resolver,
	}
}

type commandBus struct {
	output   Output
	handlers []HandlerDefinition
}

func NewCommandBus(output Output) Bus {
	return &commandBus{
		//todo
		handlers: make([]HandlerDefinition, 0),
		output:   output,
	}
}

func (cb commandBus) Handle(i Input) {
	for _, def := range cb.handlers {
		if def.resolver(i) {
			def.handler.Handle(i, cb.output)

			return
		}
	}
}

func (cb *commandBus) RegisterHandler(definition HandlerDefinition) {
	cb.handlers = append(cb.handlers, definition)
}
