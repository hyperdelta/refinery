package processor

type Processor struct {

}

type ProcessorInterface interface {
	process(in chan []byte) chan []byte
}

func ChainProcessors(processors []ProcessorInterface) (chan []byte, chan []byte) {
	var in = make(chan []byte)
	var out chan []byte = nil;

	var tmp_in chan []byte = in
	var tmp_out chan []byte

	for _, processor := range processors {
		tmp_out = processor.process(tmp_in)
		tmp_in = tmp_out
	}

	out = tmp_out

	return in, out
}