package main

import (
	"context"
	"fmt"
	"github.com/xulehexuwei/scikits"
	"reflect"
)

type Employee struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Age       int      `json:"age"`
	About     string   `json:"about"`
	Interests []string `json:"interests"`
}

func create() {
	client, err := scikits.NewEsClient("es")
	if err != nil {
		fmt.Println(err)
	}

	//使用结构体
	e1 := Employee{"Jane", "Smith", 32, "I like to collect rock albums", []string{"music"}}
	put1, err := client.Index().
		Index("megacorp").
		Type("employee").
		Id("1").
		BodyJson(e1).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put1.Id, put1.Index, put1.Type)

	//使用字符串
	e2 := `{"first_name":"John","last_name":"Smith","age":25,"about":"I love to go rock climbing","interests":["sports","music"]}`
	put2, err := client.Index().
		Index("megacorp").
		Type("employee").
		Id("2").
		BodyJson(e2).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put2.Id, put2.Index, put2.Type)

	e3 := `{"first_name":"Douglas","last_name":"Fir","age":35,"about":"I like to build cabinets","interests":["forestry"]}`
	put3, err := client.Index().
		Index("megacorp").
		Type("employee").
		Id("3").
		BodyJson(e3).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put3.Id, put3.Index, put3.Type)
}

func matchQuery() {
	searchResult := scikits.EsQueryByMatch("es", "megacorp", "first_name", "Douglas")
	querySet := make([]Employee, 0)
	if searchResult.TotalHits() > 0 {
		// 查询结果不为空，则遍历结果
		var b1 Employee
		// 通过Each方法，将es结果的json结构转换成struct对象
		for _, item := range searchResult.Each(reflect.TypeOf(b1)) {
			// 转换成Article对象
			if t, ok := item.(Employee); ok {
				querySet = append(querySet, t)
			}
		}
	}

	fmt.Println(querySet)
}

func main() {
	create()

	//_, err := scikits.NewEsClient("es")
	//if err != nil {
	//	fmt.Println(err)
	//}

	matchQuery()
}
