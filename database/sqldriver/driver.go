package sqldriver

import (
	"errors"
	mgo "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

// database: [ mongo db ]

//	driver

type sqlDriver interface {
	//	connected SQL
	Dial() error

	//	use db
	DB(db string) (SQLConn, error)

	//	address
	Addr() string
}

type mySqlDriver struct {
	address string
	session *mgo.Session
}

func newSQLDriver(addr string) sqlDriver {
	return &mySqlDriver{
		address: addr,
	}
}

func (driver *mySqlDriver) Dial() error {
	session, err := mgo.Dial(driver.address)
	if err != nil {
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	driver.session = session
	return nil
}

func (driver *mySqlDriver) DB(db string) (SQLConn, error) {
	if driver.session == nil {
		return nil, errors.New("Can't use database before connected SQL.")
	}
	dbConn := driver.session.DB(db)
	return &mySQLConn{db: dbConn, dbname: db, driver: driver}, nil
}

func (driver *mySqlDriver) Addr() string {
	return driver.address
}
