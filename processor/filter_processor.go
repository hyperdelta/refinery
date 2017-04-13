package processor

type FilterProcessor struct {
	Processor
}

func NewFilterProcessor(w WhereQuery ) *FilterProcessor {
	fp := new(FilterProcessor)

	return fp
}

func (p* FilterProcessor) process(in <-chan interface{}) <-chan interface{} {
	out := make(<-chan interface{})


	return out
}
