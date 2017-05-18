package processor

import (
	"time"
	r "gopkg.in/gorethink/gorethink.v3"
)


type RethinkDBProcessor struct {
	Processor
	pipelineId 	string
	session		*r.Session
}

type RethinkDataSchema struct {
	Id 			string 			`json:"_id"`
	UserId	 	string			`json:"_userId"`
	Interval	time.Duration	`json:"interval"`
	StartTime	string			`json:"starttime"`
	EndTime 	string			`json:"endtime"`
	GroupBy		[]string		`json:"groupBy"`
	Data		[]DataSchema	`json:"Data"`
}

type DataSchema struct {
	Group 		[]string		`json:"group"`
	Stat 		[]StatSchema	`json:"stat"`
}

type StatSchema struct {
	Column	 	string			`json:"column"`
	Operation 	string			`json:"operation"`
	Value 		string			`json:"value"`
}

func NewRethinkDBProcessor(pipelineId string) *RethinkDBProcessor {
	rp := new(RethinkDBProcessor)
	rp.pipelineId = pipelineId
	session, err := r.Connect(r.ConnectOpts{
		Address:    "0.0.0.0:32772",
		InitialCap: 10,
		MaxOpen:    10,
	})

	if err == nil && session != nil {
		rp.session = session
	} else {
		logger.Error(err)
	}

	return rp
}

func (p* RethinkDBProcessor) process(in chan interface{}) chan interface{} {

	go func(){
		for {
			select {
			case data := <-in:
				if data != nil {
					var dataMap = data.(map[string]interface{})
					p.insert(dataMap)
					logger.Debug(dataMap)
				}
			}
		}
	}()

	// end processor
	return nil;
}

func (p* RethinkDBProcessor) insert(dataMap map[string]interface{}) {

	r.DB("test").Table("test").Insert(dataMap).Exec(p.session)
}