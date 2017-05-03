package processor

import "github.com/hyperdelta/refinery/query"

type SelectProcessor struct {
	Processor
	extractKeys []string
}


func NewSelectProcessor(s []query.SelectQueryItem, g []query.GroupByQueryItem) *SelectProcessor {
	sp := new(SelectProcessor)
	return sp
}

func (p* SelectProcessor) process(in chan []byte) chan []byte {
	out := make(chan []byte)

	go func(){
		select {
		case data := <-in:
			out <- data
		}
	}()

	return out
}
