/* args */
package base

type Args interface {
	Check() error

	String() string
}

type Data interface {
	Valid() bool
}

type Entity interface {
	Id() uint32
}
