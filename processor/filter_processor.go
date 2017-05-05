package processor

import (
	"github.com/hyperdelta/refinery/query"
)

type FilterProcessor struct {
	Processor
	whereQuery query.WhereQuery
}

func NewFilterProcessor(w query.WhereQuery) *FilterProcessor {
	fp := new(FilterProcessor)
	fp.whereQuery = w

	return fp
}

func (p* FilterProcessor) process(in chan interface{}) chan interface{} {
	out := make(chan interface{})

	go func(){
		for {
			select {
			case data := <-in:
				logger.Debug("RUN FilterProcessor")

				okay := p.whereQuery.Eval(data.(map[string]string))

				if okay {
					out <- data
				} else {
					out <- nil
				}
			}
		}
	}()

	return out
}