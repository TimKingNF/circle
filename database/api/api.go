package api

import (
	"bytes"
	cmn "circle/common"
	args "circle/database/args"
	sql "circle/database/sqldriver"
	"circle/logging"
	"fmt"
)

var (
	loggerArgs cmn.LoggerArgs = args.LoggerArgs

	logger logging.Logger = cmn.NewLogger(cmn.NewLoggerArgs(
		loggerArgs.ConsoleLog(),
		loggerArgs.OutputfileLog(),
		loggerArgs.OutputfilePath(),
		loggerArgs.OutputfilePrefix()+"_api",
	))
)

const (
	DataError = "          %s [ERROR:%s] \n"

	InsertResultTemplete = "Insert %s[%s]: \n" +
		" SUCCESS[%d]:\n%s \n FAILED[%d]: %s "

	RemoveOneTemplete = "RemoveOne %s[%s] %s: SELECTOR[%v]"
	RemoveAllTemplete = "RemoveAll %s[%s] %s: SELECTOR[%v]"

	UpdateDataTemplete = "SELECTOR[%v] DATA[%v]"
	UpdateOneTemplete  = "UpdateOne %s[%s] %s: %s"
	UpdateAllTemplete  = "UpdateAll %s[%s] %s: %s"

	API_FAILED  = "FAILED"
	API_SUCCESS = "SUCCESS"
)

type ApiLayout interface {
	Insert(collection string, doc ...interface{})
	Iter(collection string, doc interface{}) []sql.DataM
	RemoveOne(collection string, selector interface{})
	RemoveAll(collection string, selector interface{})
	SelectOne(collection string, selector, result interface{})
	SelectAll(collection string, selector, doc interface{}) []sql.DataM
	UpdateOne(collection string, selector, newData interface{})
	UpdateAll(collection string, selector, newData interface{})
}

type myApiLayout struct {
	dbConn sql.SQLConn
}

func NewApiLayout(dburl, db string) ApiLayout {
	dbConn, err := sql.NewSQLConn(dburl, db)
	if err != nil {
		return nil
	}
	return &myApiLayout{
		dbConn: dbConn,
	}
}

func (api *myApiLayout) Insert(collection string, docs ...interface{}) {
	c, err := api.dbConn.Collection(collection)
	if err != nil {
		logger.Errorln(err)
		return
	}

	var failData, successData bytes.Buffer
	var failNum, successNum int8

	for k, _ := range docs {
		var saveData string
		if collection == "page" {
			saveData = fmt.Sprintf("%v ", *(docs[k].(*cmn.Page)))
		} else {
			saveData = fmt.Sprintf("%v ", docs[k])
		}
		err = c.Insert(docs[k])
		if err != nil {
			failNum++
			failData.WriteString(fmt.Sprintf(DataError, saveData, err.Error()))
		} else {
			successNum++
			successData.WriteString("          ")
			successData.WriteString(saveData)
			successData.WriteString("\n")
		}
	}

	logger.Infoln(fmt.Sprintf(InsertResultTemplete,
		api.dbConn.Database(),
		collection,
		successNum,
		successData.String(),
		failNum,
		failData.String()))
	return
}

func (api *myApiLayout) Iter(collection string, doc interface{}) []sql.DataM {
	c, err := api.dbConn.Collection(collection)
	if err != nil {
		logger.Errorln(err)
		return nil
	}
	return c.Iter(doc)
}

func (api *myApiLayout) RemoveOne(collection string, selector interface{}) {
	c, err := api.dbConn.Collection(collection)
	if err != nil {
		logger.Errorln(err)
		return
	}
	if err = c.RemoveOne(selector); err != nil {
		logger.Errorln(fmt.Sprintf(RemoveOneTemplete,
			api.dbConn.Database(),
			collection,
			API_FAILED,
			selector))
	} else {
		logger.Infoln(fmt.Sprintf(RemoveOneTemplete,
			api.dbConn.Database(),
			collection,
			API_SUCCESS,
			selector))
	}
	return
}

func (api *myApiLayout) RemoveAll(collection string, selector interface{}) {
	c, err := api.dbConn.Collection(collection)
	if err != nil {
		logger.Errorln(err)
		return
	}
	if err = c.RemoveAll(selector); err != nil {
		logger.Errorln(fmt.Sprintf(RemoveAllTemplete,
			api.dbConn.Database(),
			collection,
			API_FAILED,
			selector))
	} else {
		logger.Infoln(fmt.Sprintf(RemoveAllTemplete,
			api.dbConn.Database(),
			collection,
			API_SUCCESS,
			selector))
	}
	return
}

func (api *myApiLayout) SelectOne(collection string, selector interface{}, result interface{}) {
	c, err := api.dbConn.Collection(collection)
	if err != nil {
		logger.Errorln(err)
		return
	}
	c.SelectOne(selector, result)
	return
}

func (api *myApiLayout) SelectAll(collection string, selector interface{}, doc interface{}) []sql.DataM {
	c, err := api.dbConn.Collection(collection)
	if err != nil {
		logger.Errorln(err)
		return nil
	}
	return c.SelectAll(selector, doc)
}

func (api *myApiLayout) UpdateOne(collection string, selector, newData interface{}) {
	c, err := api.dbConn.Collection(collection)
	if err != nil {
		logger.Errorln(err)
		return
	}
	err = c.UpdateOne(selector, newData)
	if err != nil {
		logger.Errorln(fmt.Sprintf(UpdateOneTemplete,
			api.dbConn.Database(),
			collection,
			API_FAILED,
			fmt.Sprintf(DataError,
				fmt.Sprintf(UpdateDataTemplete, selector, newData),
				err.Error())))
	} else {
		logger.Infoln(fmt.Sprintf(UpdateOneTemplete,
			api.dbConn.Database(),
			collection,
			API_SUCCESS,
			fmt.Sprintf(UpdateDataTemplete, selector, newData)))
	}
	return
}

func (api *myApiLayout) UpdateAll(collection string, selector, newData interface{}) {
	c, err := api.dbConn.Collection(collection)
	if err != nil {
		logger.Errorln(err)
		return
	}
	err = c.UpdateAll(selector, newData)
	if err != nil {
		logger.Errorln(fmt.Sprintf(UpdateAllTemplete,
			api.dbConn.Database(),
			collection,
			API_FAILED,
			fmt.Sprintf(DataError,
				fmt.Sprintf(UpdateDataTemplete, selector, newData),
				err.Error())))
	} else {
		logger.Infoln(fmt.Sprintf(UpdateAllTemplete,
			api.dbConn.Database(),
			collection,
			API_SUCCESS,
			fmt.Sprintf(UpdateDataTemplete, selector, newData)))
	}
	return
}
