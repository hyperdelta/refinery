package processor

import (
	"github.com/hyperdelta/refinery/query"
	"strings"
	"github.com/buger/jsonparser"
	"encoding/json"
	"reflect"
)

type JsonParseProcessor struct {
	Processor
	columnPathList [][]string
}


func NewJsonParseProcessor(q *query.Query) *JsonParseProcessor {
	jpp := new(JsonParseProcessor)
	jpp.getColumnListFromQuery(q)

	return jpp
}

func appendDistinct(slice [][]string, item []string) [][]string {
	for _, e := range slice {
		if reflect.DeepEqual(e, item) {
			return slice
		}
	}

	return append(slice, item)
}

/*
 * Query 로 부터 필요한 ColumnPath 추출
 */
func (p* JsonParseProcessor) getColumnListFromQuery(q *query.Query) {

	if q.SelectFields != nil {
		for _, field := range q.SelectFields {
			p.columnPathList = appendDistinct(p.columnPathList, (strings.Split(field.Column, ".")))
		}
	}

	if &q.WhereQuery != nil {
		var list [][]string
		q.WhereQuery.GetColumnListFromQuery(&list)

		for _, field := range list {
			p.columnPathList = appendDistinct(p.columnPathList, field)
		}
	}

	if q.GroupByQuery != nil {
		for _, field := range q.GroupByQuery {
			p.columnPathList = appendDistinct(p.columnPathList, (strings.Split(field.Column, ".")))
		}
	}

	out, _ := json.Marshal(p.columnPathList)
	logger.Debug(string(out))
}

func (p* JsonParseProcessor) process(in chan interface{}) chan interface{} {
	out := make(chan interface{})

	go func(){
		for {
			select {
			case data := <-in:
				logger.Debug("RUN JsonParseProcessor")

				var m = p.getDataMap(data.([]byte))
				out <- m
			}
		}
	}()

	return out
}

func (p* JsonParseProcessor) getDataMap(data []byte) map[string]string {
	result := make(map[string]string)

	jsonparser.EachKey(data, func(idx int, value []byte, vt jsonparser.ValueType, err error){
		var key string = strings.Join(p.columnPathList[idx], ".")

		if err == nil {
			logger.Debug("key = " + key + " value = " + string(value))

			result[key] = string(value)
		} else {
			logger.Error("cannot find column " + key + " in data" )
			logger.Error(err)
		}
	}, p.columnPathList...)

	return result
}