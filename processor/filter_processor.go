package processor

import (
	"github.com/buger/jsonparser"
	"strings"
	"strconv"
	"github.com/hyperdelta/refinery/query"
)

const (
	AND = iota + 1
	OR
)

type FilterProcessor struct {
	Processor
	whereQuery query.WhereQuery
	operation int // and, or
	compareList []CompareItem
}

type CompareItem struct {
	path []string
	op string
	value string
}

func NewFilterProcessor(w query.WhereQuery) *FilterProcessor {
	fp := new(FilterProcessor)
	fp.whereQuery = w

	var list []query.WhereQuery

	if (len(w.And) > 0) {
		fp.operation = AND
		list = w.And

	} else {
		fp.operation = OR
		list = w.Or
	}

	fp.compareList = make([]CompareItem, len(list))

	for i, item := range list {
		var compare = CompareItem{
			path: strings.Split(item.Column, "."),
			op: item.Operation,
			value: item.Value,
		}

		fp.compareList[i] = compare
	}

	return fp
}

func (p* FilterProcessor) process(in chan []byte) chan []byte {
	out := make(chan []byte)

	go func(){
		select {
		case data := <-in:
			out <- p.filter(data)
		}
	}()

	return out
}

func (p* FilterProcessor) filter(data []byte) []byte {
	if (p.operation == OR ) {
		// or
		// 하나라도 True 가 나오면 바로 반환
		for _, compare := range p.compareList {
			value, _, _, _:= jsonparser.Get(data, compare.path...)

			if (p.doOperation(string(value), compare.op, compare.value)) {
				return data;
			}
		}

		return nil

	} else {
		// and
		// 하나라도 False 면 바로 반환
		for _, compare := range p.compareList {
			value, _, _, _:= jsonparser.Get(data, compare.path...)

			if (!p.doOperation(string(value), compare.op, compare.value)) {
				return nil;
			}
		}

		return data;
	}

	return nil;
}

func (p* FilterProcessor) doOperation(lval string, op string, rval string) bool {

	var ret bool = false

	f_lval, lval_err := strconv.ParseFloat(lval, 64);
	f_rval, rval_err := strconv.ParseFloat(rval, 64);

	if (lval_err == nil && rval_err == nil) {
		if (op == "equal") {
			ret = f_lval == f_rval;
		} else if (op == "gte") {
			ret = f_lval >= f_rval;
		} else if (op == "gt") {
			ret = f_lval > f_rval;
		} else if (op == "lte") {
			ret = lval <= rval;
		} else if (op == "lt") {
			ret = lval < rval;
		}
	} else {
		// 진짜 string 인가봄..
		// string 관련 연산만 하자
		if (op == "match") {
			ret = strings.Contains(lval, rval);
		}
	}

	return ret;
}