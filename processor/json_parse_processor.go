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


func (p* JsonParseProcessor) process(in chan []byte) chan []byte {
	out := make(chan []byte)

	go func(){
		select {
		case data := <-in:
			out <- data
		}
	}()

	return out
}
