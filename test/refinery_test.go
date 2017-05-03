package test

import (
	"testing"
	"log"
	"github.com/hyperdelta/refinery/query"
	"github.com/hyperdelta/refinery/pipeline"
)

func TestAggregator(t *testing.T) {

	var q, _ = query.Get(queryJson)
	var p, _ = pipeline.CreateFromQuery(q)
	p.In <- dataJson

	log.Print(string(<- p.Out))
}

var dataJson []byte = []byte(`
{
"ItemNo": "1234567",
"PaymentAmount": 20000,
"ShippingAddress": "서울시 강남구.."
}

`)
var queryJson []byte = []byte(`{
"interval": 10,
"select": [
{
"column": "ItemNo",
"operation": "sum",
"as": "PaymentAmountSum"
},
{
"column": "ItemNo",
"operation": "count",
"as": "ItemNoCount"
}

],
"where": {
"and": [
{
"column": "ItemNo",
"operation": "equal",
"value": "1234567"
},
{
"column": "PaymentAmount",
"operation": "gte",
"value": "20000"
},
{
"column": "ShippingAddress",
"operation": "match",
"value": "서울시"
}

],
"or": null
},
"groupBy": [
{
"column": "ShippingAddress",
"pattern": "구"
}
]


}`)