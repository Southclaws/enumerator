package example

//go:generate go run -mod=mod github.com/Southclaws/enumerator

import "fmt"

type statusEnum string

const (
	success   statusEnum = "success"
	failure   statusEnum = "failure"
	inBetween statusEnum = "inbetween"
	notSure   statusEnum = "notsure"
)

func Hi(in string) {
	status, err := NewStatus(in)
	if err != nil {
		panic(err)
	}

	fmt.Println("Status is:", status.String())
}
