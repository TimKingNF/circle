package main

import (
	api "circle/database/api"
	base "circle/database/base"
	dbcm "circle/database/databasectrlm"
	diriver "circle/database/sqldiriver"
	"fmt"
	"sync"
	// "time"
)

type Person struct {
	Name  string
	Phone string
}

var (
	dividerAddr = "127.0.0.1:8086"

	sqlAddr = "127.0.0.1:27017"

	dbname = "circle"

	databaseArgs base.DatabaseArgs = base.NewDatabaseArgs(
		dividerAddr,
		sqlAddr,
		dbname)
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	TestSocket()
}

func TestSocket() {
	dbCtrlM := dbcm.GenDatabaseControlModel()
	err := dbCtrlM.Init(databaseArgs)
	if err != nil {
		panic(fmt.Sprintf("database control model initialize failed, Error: %s", err))
	}

	//	socket test
	// for i := 0; i < 5; i++ {
	// 	dbCtrlM.Send(fmt.Sprintf("send the number[%d]", i))
	// 	time.Sleep(time.Second)
	// }

	var waitgroup sync.WaitGroup
	waitgroup.Add(1)
	waitgroup.Wait()
}

func TestAPI() {
	apiLayout := api.NewApiLayout(databaseArgs.SqlAddr(), databaseArgs.Dbname())

	//	insert api
	// apiLayout.Insert("col", &Person{"Ale", "+55 53 8402 9639"}, &Person{"Cla", "+55 53 8402 8510"})

	//	remove all api
	// apiLayout.RemoveAll("col", diriver.DataM{"name": "Cla"})
	// apiLayout.RemoveAll("col", nil)

	//	remove one api
	// apiLayout.RemoveOne("col", diriver.DataM{"name": "Ale"})

	//	update one api
	/*	apiLayout.UpdateOne("col",
		diriver.DataM{"name": "Ale"},
		diriver.DataM{
			"$set": diriver.DataM{"name": "Cla"},
		})*/

	//	update all api
	/*apiLayout.UpdateAll("col",
	diriver.DataM{"name": "Cla"},
	diriver.DataM{
		"$set": diriver.DataM{"addr": "广州"},
	})*/

	//	iter api
	person := Person{}
	ret := apiLayout.Iter("col", &person)
	for _, v := range ret {
		fmt.Println(v)
	}

	//	select one api
	/*person := Person{}
	apiLayout.SelectOne("col", diriver.DataM{"name": "Ale"}, &person)
	fmt.Println(person)*/

	//	select all api
	/*	person := Person{}
		ret := apiLayout.SelectAll("col", nil, &person)
		for _, v := range ret {
			fmt.Println(v)
		}*/
}

func TestSQL() {
	conn, _ := diriver.NewSQLConn(sqlAddr, "test1")
	collection, _ := conn.Collection("col")

	//	insert data
	err := collection.Insert(&Person{"Ale", "+55 53 8402 9639"}, &Person{"Cla", "+55 53 8402 8510"})
	fmt.Println(err)

	//	update data
	/*err := collection.UpdateAll(
		diriver.DataM{"name": "Cla"},
		diriver.DataM{
			"$set": diriver.DataM{"phone": "321"},
		})
	fmt.Println(err)*/

	//	select muti data
	/*person := Person{}
	ret := collection.SelectAll(diriver.DataM{"name": "Cla"}, &person)
	for _, v := range ret {
		fmt.Println(v)
	}*/

	//	remove data
	/*	err := collection.RemoveAll(nil)
		fmt.Println(err)*/

	//	select all data
	person := Person{}
	ret := collection.Iter(&person)
	for _, v := range ret {
		fmt.Println(v)
	}
}
