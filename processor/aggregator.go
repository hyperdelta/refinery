package processor


type Aggregator struct {

}

type AggregatorInfo struct {
	in <-chan interface{}
	lifetime int
}

func CreateAggregatorFromQuery(q* Query) (*Aggregator, error) {
	ag := new(Aggregator)
	err := ag.setupProcessors(q)

	return ag, err
}

func (a* Aggregator) setupProcessors(q* Query) error {

	var processors []ProcessorInterface
	var err error

	processors = append(processors, NewJsonParseProcessor())
	processors = append(processors, NewFilterProcessor(q.WhereQuery))

	// entry point
	//in := p.ChainProcessors(processors)

	// generate aggregator id
	//aggId := generateRandomID("agg-")

	return err
}
