package processor

import (
	"github.com/hyperdelta/refinery/query"
	"github.com/hyperdelta/refinery/trie"
	"strings"
	"strconv"
	"time"
)

type StatisticProcessor struct {
	Processor

	SelectQueryList  []query.SelectQueryItem
	GroupByQueryList []query.GroupByQueryItem

	HasGroupBy       bool

	Interval         time.Duration
	tickerChannel 	 <-chan time.Time
	trie             *trie.Trie
}

type StatTrieData struct {
	DataMap map[string]*StatData
}

type StatData struct {
	Column	 	string			`json:"column"`
	Operation 	string			`json:"operation"`
	Value 		int64			`json:"value"`
}

const MIN_VALUE int64 =  -int64(^uint(0) >> 1) - 1
const WILDCARD string = "@"

func NewStatisticTrieData() *StatTrieData {
	data := new(StatTrieData)
	data.DataMap = make(map[string]*StatData)

	return data
}

func NewStatisticData(column string) *StatData {
	elem := new(StatData)
	elem.Column = column
	elem.Value = MIN_VALUE

	return elem
}

func NewStatisticProcessor(interval time.Duration, s []query.SelectQueryItem, g []query.GroupByQueryItem) *StatisticProcessor {
	sp := new(StatisticProcessor)

	sp.Interval = interval
	sp.SelectQueryList = s
	sp.GroupByQueryList = g

	sp.tickerChannel = time.NewTicker(time.Second * interval).C

	sp.buildStatisticPlan(s, g)

	return sp
}

func (p* StatisticProcessor) process(in chan interface{}) chan interface{} {
	out := make(chan interface{})

	go func(){
		for {
			select {
			case data := <-in:
				if data != nil {
					p.doOperation(data.(map[string]string))
				}
			case <- p.tickerChannel:
				out <- p.trie.ToDataMap()
				p.trie.Clear()
			}
		}
	}()

	return out
}

func (p* StatisticProcessor) buildStatisticPlan(s []query.SelectQueryItem, g []query.GroupByQueryItem) {

	if g != nil && len(g) > 0 {
		p.HasGroupBy = true
		p.trie = trie.NewTrie()
	}

	// select 에 As 가 비어있는 경우 column 으로 대체
	for _, v := range s {
		if v.As == "" {
			v.As = v.Column
		}
	}
}

func (p* StatisticProcessor) doOperation(data map[string]string) {

	var prefix []string

	if p.HasGroupBy {
		/// make prefix
		for _, groupby_item := range p.GroupByQueryList {
			var groupby_data = data[groupby_item.Column]
			var splited_groupby_data = strings.Split(groupby_data, " ")

			if len(splited_groupby_data) >= groupby_item.Depth {
				prefix = append(prefix, splited_groupby_data[:groupby_item.Depth]...)
			} else {
				// 요청한 group by depth 보다 모자람
				// 나머지는 wildcard 로 대체
				prefix = append(prefix, splited_groupby_data...)
				for i := 0; i < groupby_item.Depth - len(splited_groupby_data) ; i++ {
					prefix = append(prefix, WILDCARD)
				}
			}
		}

	} else {
		// groupby 없음 - prefix 는 wildcard 로..
		prefix = append(prefix, WILDCARD)
	}

	// find groupby prefix
	var trieData = p.trie.Retrieve(prefix ...)
	var statisticData *StatData
	//var addTrieData bool = false

	// 처음 생긴 groupby 정보라면 새로 생성
	if trieData == nil {
		trieData = NewStatisticTrieData()
		//addTrieData = true
	}

	// select 문을 순회하면서 통계 처리
	for _, select_item := range p.SelectQueryList {
		statisticData = (trieData.(*StatTrieData)).DataMap[select_item.As]

		if statisticData == nil {
			statisticData = NewStatisticData(select_item.As)
			(trieData.(*StatTrieData)).DataMap[select_item.As] = statisticData
		}

		statisticData.doOperation(data[select_item.Column], select_item.Operation)
	}

	p.trie.Add(trieData, prefix...)
}

func (d*StatData) doOperation(value string, op string) {

	switch op {
	case "count":
		if d.Value == MIN_VALUE {
			d.Value = 0
		}

		d.Value += 1
		return
	}

	v, err := strconv.ParseInt(value, 0, 64)

	if err == nil {

		switch op {
		case "sum":
			if d.Value == MIN_VALUE {
				d.Value = 0
			}

			d.Value += v
			return
		case "max":
			if d.Value < v {
				d.Value = v
			}
			return
		case "min":
			if d.Value == MIN_VALUE || d.Value > v {
				d.Value = v
			}
			return
		}
	}
}