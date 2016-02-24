package sqldiriver

import (
	"errors"
	mgo "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

// database: [ mongo db ]

//	diriver

type sqlDiriver interface {
	//	connected SQL
	Dial() error

	//	use db
	DB(db string) (SQLConn, error)

	//	address
	Addr() string
}

type mySqlDiriver struct {
	address string
	session *mgo.Session
}

func newSQLDiriver(addr string) sqlDiriver {
	return &mySqlDiriver{
		address: addr,
	}
}

func (diriver *mySqlDiriver) Dial() error {
	session, err := mgo.Dial(diriver.address)
	if err != nil {
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	diriver.session = session
	return nil
}

func (diriver *mySqlDiriver) DB(db string) (SQLConn, error) {
	if diriver.session == nil {
		return nil, errors.New("Can't use database before connected SQL.")
	}
	dbConn := diriver.session.DB(db)
	return &mySQLConn{db: dbConn, dbname: db, diriver: diriver}, nil
}

func (diriver *mySqlDiriver) Addr() string {
	return diriver.address
}
