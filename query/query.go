package query

import (
	"encoding/json"
	"errors"
	"github.com/hyperdelta/refinery/log"
)

var (
	logger *log.Logger = log.Get()
	resolver ValueResolver
)

type Query struct {
	Interval     int				`json:"interval"`
	SelectFields []SelectQueryItem	`json:"select"`
	WhereQuery   WhereQuery			`json:"where"`
	GroupByQuery []GroupByQueryItem		`json:"groupBy"`
}

/**
{
	"column": "PaymentAmount",
	"operation": "sum",
	"as": "PaymentAmountSum"
}
 */
type SelectQueryItem struct {
	Column string		`json:"column"`
	Operation string	`json:"operation"`
	As string			`json:"as"`
}

type GroupByQueryItem struct {
	Column string 	`json:"column"`
	Pattern string 	`json:"pattern"`
}

type ValueResolver interface {
	Get(column string) string
}

func SetValueResolver(resolver ValueResolver) {
	resolver = resolver
}

func Get(body []byte) (*Query, error) {
	q := new(Query)

	err := json.Unmarshal(body, &q)

	if err != nil {
		panic(err)
	}

	logger.Debug("Get() register")

	// query validation
	err = validate(*q)

	if err != nil {
		return nil, err
	} else {
		return q, nil
	}
}

func validate(q Query) error {
	if len(q.WhereQuery.And) > 0 && len(q.WhereQuery.Or) > 0 {
		// or 와 and 는 둘 다 존재할 수 없음
		return errors.New("Where Clause Error: Both Or and And cannot be existed")
	}

	return nil
}

func (q *Query) GetBytes() ([]byte, error) {
	return json.Marshal(q)
}