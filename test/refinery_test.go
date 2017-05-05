package test

import (
	"testing"
	"log"
	"github.com/hyperdelta/refinery/query"
	"github.com/hyperdelta/refinery/pipeline"
	"github.com/hyperdelta/refinery/trie"
	"time"
	"fmt"
)

func TestPipeline(t *testing.T) {

	var q, _ = query.Get(queryJson)
	pipeline.CreateFromQuery(q)

	for _, p := range pipeline.PipelineList {
		go func() {
			for {
				select {
				case data := <- p.Out:
					fmt.Printf("\n%s\n", &data)
				}
			}
		}()
	}

	go func() {
		for {
			for _, p := range pipeline.PipelineList {
				log.Print("send data")
				p.In <- dataJson
				time.Sleep(time.Second)
			}
		}
	}()

	time.Sleep(time.Second * 3)
}

func TestTrie(t *testing.T) {
	var tt = trie.NewTrie()

	tt.Add("data1", "a", "b", "c")
	tt.Add("data2", "a", "b", "d")

	var retData = tt.Retrieve("a", "b", "c")
	var retData2 = tt.Retrieve("a", "b", "d")
	var retData3 = tt.Retrieve("a", "b")

	if "data1" !=  retData {
		t.Fatal("TestTrie Fail, expected data = data1, retrieve data = " + retData.(string))
	}

	if  "data2" != retData2 {
		t.Fatal("TestTrie Fail, expected data = data2, retrieve data = " + retData2.(string))
	}

	if retData3 != nil {
		t.Fatal("TestTrie Fail, expected data = nil, retrieve data = " + retData3.(string))
	}
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
"operation": "eq",
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
"depth": 2
}
]


}`)