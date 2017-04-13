package processor

//import "github.com/buger/jsonparser"

type JsonParseProcessor struct {
	Processor
	extractKeys []string
}


func NewJsonParseProcessor() *JsonParseProcessor {
	jpp := new(JsonParseProcessor)
	return jpp
}


func (p* JsonParseProcessor) process(in <-chan interface{}) <-chan interface{} {
	out := make(<-chan interface{})



	return out
}
