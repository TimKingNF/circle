package sqldriver

import (
	"errors"
	"fmt"
	mgo "gopkg.in/mgo.v2"
)

//	conn

type SQLConn interface {
	Collection(name string) (SQLCollection, error)
	Database() string
}

type mySQLConn struct {
	db     *mgo.Database
	dbname string
	driver sqlDriver
}

func NewSQLConn(addr, db string) (SQLConn, error) {
	sqlDriver := newSQLDriver(addr)
	err := sqlDriver.Dial()
	if err != nil {
		return nil, err
	}
	return sqlDriver.DB(db)
}

func (conn *mySQLConn) Collection(name string) (SQLCollection, error) {
	if conn == nil {
		return nil, errors.New("The sql conn is nil.")
	}
	if conn.db == nil {
		return nil, errors.New("The sql disconnected.")
	}
	return &mySQLCollection{collection: conn.db.C(name)}, nil
}

func (conn *mySQLConn) Database() string {
	return fmt.Sprintf("%s@%s", conn.driver.Addr(), conn.dbname)
}
