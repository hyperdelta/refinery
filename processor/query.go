package processor

import (
	"encoding/json"
	"github.com/hyperdelta/refinery/config"
	"log"
	"fmt"
	"errors"
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

type WhereQuery struct {
	And []WhereQueryItem	`json:"and"`
	Or []WhereQueryItem		`json:"or"`
}

/**
{
	"column": "ItemNo",
	"operation": "equal",
 	"value": "1234567"
}
 */
type WhereQueryItem struct {
	Column string		`json:"column"`
	Operation string	`json:"operation"`
	Value string		`json:"value"`
}

type GroupByQueryItem struct {
	Column string 	`json:"column"`
	Pattern string 	`json:"pattern"`
}

func GetQueryObject(body []byte) (*Query, error) {
	q := Query{}

	err := json.Unmarshal(body, &q)

	if err != nil {
		panic(err)
	}

	if config.Debug {
		log.Print("GetQueryObject() register")
		fmt.Println(q)
	}

	// query validation
	err = validate(q)

	if err != nil {
		return nil, err
	} else {
		return &q, nil
	}

}

func validate(q Query) error {
	if len(q.WhereQuery.And) > 0 && len(q.WhereQuery.Or) > 0 {
		// or 와 and 는 둘 다 존재할 수 없음
		return errors.New("Where Clause Error: Both Or and And cannot be existed")
	}

	return nil
}

