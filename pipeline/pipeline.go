package pipeline

import (
	"github.com/hyperdelta/refinery/config"
	"github.com/hyperdelta/refinery/query"
	"github.com/hyperdelta/refinery/processor"
)

var (
	PipelineList map[string]*Pipeline = make(map[string]*Pipeline)
)

type State uint8

const (
	INITIALIZE State = iota
	RUNNING
	PAUSE
	CLOSING
	CLOSED
)

type Pipeline struct {
	In chan interface{}	// input channel
	//Out chan interface{}	// output channel

	Lifetime int
	Endpoint string
	Id string

	State State
	processors []processor.ProcessorInterface
}

func CreateFromQuery(q *query.Query) (*Pipeline, error) {
	var p *Pipeline = &Pipeline {
	}
	var err error

	p.Initialize(q)
	PipelineList[p.Id] = p

	return p, err
}

func Close(id string) {
	var p *Pipeline = Retrieve(id)

	if p != nil {
		p.Close()
	}

	// remove pipeline
	delete(PipelineList, id)
}

func CloseAll() {
	for _, p := range PipelineList {
		p.Close()
	}
}

func Retrieve(id string) (*Pipeline) {
	v, ok := PipelineList[id]

	if ok {
		return v
	} else {
		return nil
	}
}

func Pause(id string) {
	var p *Pipeline = Retrieve(id)

	if p != nil {
		p.Pause()
	}
}

func PauseAll()  {
	for _, p := range PipelineList {
		p.Pause()
	}
}

func Resume (id string) {
	var p *Pipeline = Retrieve(id)

	if p != nil {
		p.Resume()
	}
}

func ResumeAll() {
	for _, p := range PipelineList {
		p.Resume()
	}
}

func SendDataToPipeline(id string, data []byte) {
	var p *Pipeline = Retrieve(id);

	if p != nil && p.State == RUNNING {
		p.In <- data
	}
}

func SendDataToAllPipeline(data []byte) {
	for _, p := range PipelineList {
		if p.State == RUNNING {
			p.In <- data
		}
	}
}

func (p* Pipeline) Resume() {
	p.State = RUNNING
}

func (p* Pipeline) Pause() {
	p.State = PAUSE
}

func (p* Pipeline) Close() {
	p.State = CLOSING

	close(p.In)

	//TODO: remove endpoint

	p.State = CLOSED
}

func (p* Pipeline) Initialize(q* query.Query) error {
	var err error

	p.SetupProcessors(q)
	p.Id = config.GenerateRandomID("pipeline-")

	p.State = INITIALIZE
	//TODO: generate Rethink endpoint

	//TODO: save Query into Redis
	p.State = RUNNING

	return err
}

func (p* Pipeline) SetupProcessors(q *query.Query) error {

	var processors []processor.ProcessorInterface
	var err error

	processors = append(processors, processor.NewJsonParseProcessor(q))
	processors = append(processors, processor.NewFilterProcessor(q.WhereQuery))
	processors = append(processors, processor.NewStatisticProcessor(q.Interval, q.SelectFields, q.GroupByQuery))
	processors = append(processors, processor.NewRethinkDBProcessor(p.Id, q))

	// entry point
	in := processor.ChainProcessors(processors)

	p.processors = processors
	p.In = in
	//p.Out = out

	return err
}
