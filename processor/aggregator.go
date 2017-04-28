package processor

import (
	"github.com/hyperdelta/refinery/config"
)

type Aggregator struct {
	In chan []byte
	Out chan []byte
	Info* AggregatorInfo
}

type AggregatorInfo struct {
	Lifetime int
	Endpoint string
	Id string
}

var AggregatorsMap map[string]*Aggregator = make(map[string]*Aggregator)

func CreateAggregatorFromQuery(q* Query) (*Aggregator, error) {
	ag := new(Aggregator)
	in, out, err := ag.setupProcessors(q)

	ag.In = in
	ag.Out = out
	ag.Info = new(AggregatorInfo)
	ag.Info.Id = config.GenerateRandomID("agg-");

	AggregatorsMap[ag.Info.Id] = ag

	return ag, err
}

func (a* Aggregator) setupProcessors(q* Query) (chan []byte, chan []byte, error) {

	var processors []ProcessorInterface
	var err error

	processors = append(processors, NewJsonParseProcessor())
	processors = append(processors, NewFilterProcessor(q.WhereQuery))
	processors = append(processors, NewSelectProcessor(q.SelectFields, q.GroupByQuery))

	// entry point
	in, out := ChainProcessors(processors)

	return in, out, err
}
