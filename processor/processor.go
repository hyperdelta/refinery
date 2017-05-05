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

func ChainProcessors(processors []ProcessorInterface) (chan interface{}, chan interface{}) {
	var in = make(chan interface{})
	var out chan interface{} = nil;

	var tmp_in chan interface{} = in
	var tmp_out chan interface{}

	for _, processor := range processors {
		tmp_out = processor.process(tmp_in)
		tmp_in = tmp_out
	}

	out = tmp_out

	return in, out
}