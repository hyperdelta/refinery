package test

import (
	"testing"
	"github.com/hyperdelta/refinery/query"
	"github.com/hyperdelta/refinery/pipeline"
	"github.com/hyperdelta/refinery/trie"
	"time"
	"io/ioutil"
	"log"
)

func TestPipeline(t *testing.T) {

	var q, _ = query.Get(queryJson)
	pipeline.CreateFromQuery(q)

	for _, p := range pipeline.PipelineList {
		go func() {
			for {
				select {
				case data := <- p.Out:
					if data != nil {
						log.Print(data)
					}
					break
				}
			}
		}()
	}

	var dataJsonList [][]byte = [][]byte {
		getByteArray("log/log-1.json"),
		getByteArray("log/log-2.json"),
		getByteArray("log/log-3.json"),
		getByteArray("log/log-4.json"),
		getByteArray("log/log-5.json"),
	}

	go func() {
		var count = 0
		for {
			pipeline.SendDataToAllPipeline(dataJsonList[count % 5])
			count += 1
		}
	}()

	time.Sleep(time.Second * 3)	// 3초간 유지
}

func getByteArray(path string) []byte {
	d, _ := ioutil.ReadFile(path)
	return d
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

	if "data2" != retData2 {
		t.Fatal("TestTrie Fail, expected data = data2, retrieve data = " + retData2.(string))
	}

	if retData3 != nil {
		t.Fatal("TestTrie Fail, expected data = nil, retrieve data = " + retData3.(string))
	}
}

var queryJson []byte = []byte(`{
	"interval": 1,
	"select": [
	{
		"column": "member_id",
		"operation": "count",
		"as": "member_id_count"
	}
	],
	"where": {
		"and": [
		{
			"column": "payload.body.paymentData.NewSmilePay.TotalMoney",
			"operation": "gte",
			"value": "10000"
		},
		{
			"column": "payload.body.paymentData.NewSmilePay.TotalMoney",
			"operation": "lte",
			"value": "100000"
		}
		],
		"or": null
	},
	"groupBy": [
	{
		"column": "payload.body.shippingAddressList.[0].DeliveryAddr1",
		"depth": 2
	}
	]
}`)