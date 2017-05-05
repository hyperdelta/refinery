package query

import (
	"errors"
	"strings"
	"strconv"
)

// Where -> T | And | Or
//	 And -> Where*
//   Or -> Where*
type WhereQuery struct {
	TerminalQuery
	AndQuery
	OrQuery
}

type AndQuery struct {
	And []WhereQuery 	`json:"and"`
}

type OrQuery struct {
	Or []WhereQuery		`json:"or"`
}

type TerminalQuery struct {
	Column string			`json:"column"`
	Operation string		`json:"operation"`
	Value string			`json:"value"`
}

func (w *WhereQuery) Eval(data map[string]string) bool {
	if w.And != nil {
		for _, q := range w.And {
			if q.Eval(data) == false {
				return false
			}
		}
		return true
	} else if w.Or != nil {

		for _, q := range w.Or {
			if (q.Eval(data)) {
				return true
			}
		}
		return false
	} else {
		// find lval
		var lval = data[w.Column]
		result := evalTerm(lval, w.Operation, w.Value)

		logger.Debug("lval = " + lval + " op = " + w.Operation +
			" rval = " + w.Value + ", result = " + strconv.FormatBool(result))

		return result
	}
}

func (w *WhereQuery) GetColumnListFromQuery(list *[][]string) {
	if w.And != nil {
		for _, q := range w.And {
			q.GetColumnListFromQuery(list)
		}
	} else if w.Or != nil {
		for _, q := range w.Or {
			q.GetColumnListFromQuery(list)
		}
	} else {
		// terminal
		*list = append(*list, strings.Split(w.Column, "."))
	}
}

func evalTerm(lval string, op string, rval string) bool {

	var result bool = false

	// string operation
	switch op {
	case "match":
		return strings.Contains(lval, rval);
	}

	// number operation
	f_lval, lval_err := strconv.ParseFloat(lval, 64);
	f_rval, rval_err := strconv.ParseFloat(rval, 64);

	if lval_err != nil || rval_err != nil {
		result = false
	}

	switch op {
	case "eq":
		return f_lval == f_rval
	case "neq":
		return f_lval != f_rval
	case "gte":
		return f_lval >= f_rval
	case "gt":
		return f_lval > f_rval
	case "lte":
		return f_lval <= f_rval
	case "lt":
		return f_lval < f_rval
	}

	return result
}

func (w *WhereQuery) Validate() error {

	var count = 0

	if w.Column != "" || w.Value != "" || w.Operation != "" {
		count ++
	}

	if w.And != nil {
		count ++
	}

	if w.Or != nil {
		count ++
	}

	if count > 1 {
		return errors.New("Where Clause Error: Root Where node should have only one child")
	}

	return nil
}

