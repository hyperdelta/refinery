package processor

type Processor struct {

}

type ProcessorInterface interface {
	process(in <-chan interface{}) <-chan interface{}
}

func ChainProcessors(processors []ProcessorInterface) <-chan interface{} {
	in := make(<-chan interface{})
	var out <-chan interface{} = nil

	for _, processor := range processors {
		if out != nil {
			in = out
		}

		out = processor.process(in)
	}

	return in
}