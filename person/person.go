package person

import "errors"

type Person struct {
	name string
	age  int
}

func New(name string, age int) (*Person, error) {
	p := &Person{}
	p.name = name
	if age > 120 {
		return p, errors.New("the age is too large")
	}
	p.age = age
	return p, nil

}
