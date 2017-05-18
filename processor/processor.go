package processor

import "github.com/hyperdelta/refinery/log"

type Processor struct {

}

var (
	logger *log.Logger = log.Get()
)

type ProcessorInterface interface {
	process(in chan interface{}) chan interface{}
}

type DBProcessorInterface interface {
	ProcessorInterface
}

func ChainProcessors(processors []ProcessorInterface) (chan interface{}) {
	var in = make(chan interface{})

	var tmp_in chan interface{} = in
	var tmp_out chan interface{}

	for _, processor := range processors {
		tmp_out = processor.process(tmp_in)
		tmp_in = tmp_out
	}

	return in
}