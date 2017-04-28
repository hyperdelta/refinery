package processor

type SelectProcessor struct {
	Processor
	extractKeys []string
}


func NewSelectProcessor(s []SelectQueryItem, g []GroupByQueryItem) *SelectProcessor {
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
