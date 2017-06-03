package processor

import (
	"time"
	r "gopkg.in/gorethink/gorethink.v3"
	"github.com/hyperdelta/refinery/config"
	"github.com/hyperdelta/refinery/query"
	"strconv"
)


type RethinkDBProcessor struct {
	Processor

	pipelineId 	string
	session		*r.Session
	startTime 	time.Time
	endTime 	time.Time

	query 		*query.Query
}

type RethinkDataSchema struct {
	UserId	 	string			`json:"_userId"`
	Interval	time.Duration	`json:"interval"`
	StartTime	string			`json:"starttime"`
	EndTime 	string			`json:"endtime"`
	GroupBy		[]string		`json:"groupBy"`
	Data		[]DataSchema	`json:"Data"`
}

type DataSchema struct {
	Group 		[]string		`json:"group"`
	Stat 		[]*StatData		`json:"stat"`
}

var (
	conf = config.RefineryConfig
)

func NewRethinkDBProcessor(pipelineId string, query *query.Query) *RethinkDBProcessor {
	rp := new(RethinkDBProcessor)

	rp.pipelineId = pipelineId
	rp.query = query

	logger.Error(conf.RethinkDBAddress + ":" + strconv.Itoa(conf.RethinkDBPort))

	session, err := r.Connect(r.ConnectOpts {
		Address:    conf.RethinkDBAddress + ":" + strconv.Itoa(conf.RethinkDBPort),
		InitialCap: 10,
		MaxOpen:    10,
	})

	if err == nil && session != nil {
		rp.session = session
	} else {
		logger.Error(err)
	}

	r.DB(conf.RethinkDBName).TableCreate(query.UserId).Exec(rp.session)

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

	var dataList = make([]DataSchema, len(dataMap))
	var i = 0;

	// Data
	for k, v := range dataMap {
		var group = make([]string, 1);

		if k == WILDCARD {	// wildcard 의 경우, group 이 없는 케이스이므로 nil 전달
			group = nil
		} else {
			group[0] = k
		}

		var statList = make([]*StatData, len(v.(*StatTrieData).DataMap))

		var j = 0
		for _, stat_val := range v.(*StatTrieData).DataMap {
			statList[j] = stat_val
			j++
		}

		var data = DataSchema {
			Group: group,
			Stat: statList,
		}

		dataList[i] = data
		i ++
	}

	// GroupBy column name
	var groupByColumnNameList = make([]string, len(p.query.GroupByQuery))

	for i, v := range p.query.GroupByQuery {
		groupByColumnNameList[i] = v.Column;
	};

	var resultData = RethinkDataSchema {
		UserId: p.query.UserId,
		Interval: p.query.Interval,
		StartTime: p.startTime.Format("2006-01-02 15:04:05"),
		EndTime: p.endTime.Format("2006-01-02 15:04:05"),
		GroupBy: groupByColumnNameList,
		Data: dataList,
	};

	r.DB(conf.RethinkDBName).
		Table(p.query.UserId).
		Insert(resultData).
		Exec(p.session)
}