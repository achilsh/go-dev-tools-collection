package example

import (
	"fmt"
	"log"
	"testing"
)

type MyKeyer interface {
	Key() string
}

type myKey struct {
	id   int64
	name string
	addr *int
}

func (m myKey) Key() string {
	return fmt.Sprintf("%v:%v", m.id, m.name)
}

func TestMyKeyInMap(t *testing.T) {
	var xyz = map[MyKeyer]int{
		myKey{
			id:   10,
			name: "abc",
		}: 11,
		myKey{
			id:   101,
			name: "abc",
		}: 1000,
	}

	for k, v := range xyz {
		fmt.Printf("key: %v, value: %v\n", k, v)
	}

	r := xyz[myKey{
		id:   10,
		name: "abc",
	}]
	fmt.Printf("r value: %v\n", r)
}
func TestTaskWrapper(t *testing.T) {
	log.Printf("this is demo.")
}
