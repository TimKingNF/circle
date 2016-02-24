/* item */
package base

type Item map[string]interface{}

func (item Item) Valid() bool {
	return item != nil
}
