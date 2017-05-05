package processor

import (
	"github.com/hyperdelta/refinery/query"
	"github.com/hyperdelta/refinery/trie"
	"strings"
	"strconv"
)

type StatisticProcessor struct {
	Processor

	selectQueryList  []query.SelectQueryItem
	groupByQueryList []query.GroupByQueryItem

	hasGroupBy       bool

	interval         int
	trie             *trie.Trie
}

type StatisticTrieData struct {
	dataMap map[string]*StatisticData
}

type StatisticData struct {
	column string
	value int64
}

func NewStatisticTrieData() *StatisticTrieData {
	data := new(StatisticTrieData)
	data.dataMap = make(map[string]*StatisticData)

	return data
}

func NewStatisticData(column string) *StatisticData {
	elem := new(StatisticData)
	elem.column = column
	elem.value = 0

	return elem
}

func NewStatisticProcessor(interval int, s []query.SelectQueryItem, g []query.GroupByQueryItem) *StatisticProcessor {
	sp := new(StatisticProcessor)

	sp.interval = interval
	sp.selectQueryList = s
	sp.groupByQueryList = g

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
					out <- p.trie
				} else {
					// bypass
					out <- nil
				}
			}
		}
	}()

	return out
}

func (p* StatisticProcessor) buildStatisticPlan(s []query.SelectQueryItem, g []query.GroupByQueryItem) {

	if g != nil && len(g) > 0 {
		p.hasGroupBy = true
		p.trie = trie.NewTrie()
	}

}

func (p* StatisticProcessor) doOperation(data map[string]string) {

	var prefix []string

	if p.hasGroupBy {
		/// make prefix

		for _, groupby_item := range p.groupByQueryList {
			var groupby_data = data[groupby_item.Column]
			var splited_groupby_data = strings.Split(groupby_data, " ")

			prefix = append(prefix, splited_groupby_data[:groupby_item.Depth]...)
		}

	} else {
		// groupby 없음
	}

	// find groupby prefix
	var trieData = p.trie.Retrieve(prefix ...)
	var statisticData *StatisticData

	// 처음 생긴 groupby 정보라면 새로 생성
	if trieData == nil {
		trieData = NewStatisticTrieData()
		p.trie.Add(trieData, prefix...)
		p.trie.Print()
	}

	// select 문을 순회하면서 통계 처리
	for _, select_item := range p.selectQueryList {
		statisticData = (trieData.(*StatisticTrieData)).dataMap[select_item.As]

		if statisticData == nil {
			statisticData = NewStatisticData(select_item.As)
			(trieData.(*StatisticTrieData)).dataMap[select_item.As] = statisticData
		}

		statisticData.doOperation(data[select_item.Column], select_item.Operation)
	}

}

func (d* StatisticData) doOperation(value string, op string) {
	switch op {
	case "count":
		d.value += 1
		return
	}

	v, err := strconv.ParseInt(value, 0, 64)

	if err == nil {

		switch op {
		case "sum":
			d.value += v
			return
		case "max":
			if d.value < v {
				d.value = v
			}
			return
		case "min":
			if d.value > v {
				d.value = v
			}
			return
		}
	}
}