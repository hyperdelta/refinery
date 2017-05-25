package processor

import (
	"time"
	r "gopkg.in/gorethink/gorethink.v3"
	"github.com/hyperdelta/refinery/config"
	"github.com/hyperdelta/refinery/query"
)


type RethinkDBProcessor struct {
	Processor

	pipelineId 	string
	session		*r.Session
	startTime 	time.Time
	endTime 	time.Time

	query 		query.Query
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
	Stat 		[]StatData		`json:"stat"`
}

var (
	conf = config.RefineryConfig
)

func NewRethinkDBProcessor(pipelineId string, query query.Query) *RethinkDBProcessor {
	rp := new(RethinkDBProcessor)

	rp.pipelineId = pipelineId
	rp.query = query

	session, err := r.Connect(r.ConnectOpts {
		Address:    conf.RethinkDBAddress + ":" + string(conf.RethinkDBPort),
		InitialCap: 10,
		MaxOpen:    10,
	})

	if err == nil && session != nil {
		rp.session = session
	} else {
		logger.Error(err)
	}

	rp.startTime = time.Now()
	return rp
}

func (p* RethinkDBProcessor) process(in chan interface{}) chan interface{} {

	go func(){
		for {
			select {
			case data := <-in:
				if data != nil {
					var dataMap = data.(map[string]interface{})
					p.endTime = time.Now()
					p.insert(dataMap)

					p.startTime = time.Now()
				}
			}
		}
	}()

	// end processor
	return nil;
}


func (p* RethinkDBProcessor) insert(dataMap map[string]interface{}) {

	var dataList [len(dataMap)]DataSchema;
	var i = 0;

	// Data
	for k, v := range dataMap {
		var group [1]string;

		if k == WILDCARD {	// wildcard 의 경우, group 이 없는 케이스이므로 nil 전달
			group = nil
		} else {
			group[0] = k
		}

		var data = DataSchema {
			Group: group,
			Stat: v,
		}

		dataList[i] = data
		i ++
	}

	// GroupBy list



	var resultData = RethinkDataSchema {
		StartTime: p.startTime.Format("2017-04-06 12:56:00"),
		EndTime: p.endTime.Format("2017-04-06 12:56:05"),
		GroupBy: len(p.query.GroupByQuery) > 0 ?  : nil,
		Data: dataList,
	};

	r.DB("test").Table("test").Insert(resultData).Exec(p.session)
}