package example

//go:generate go run -mod=mod github.com/Southclaws/enumerator

import "fmt"

type projectStatusEnum string

const (
	success   projectStatusEnum = "success"
	failure   projectStatusEnum = "failure"
	inBetween projectStatusEnum = "inbetween"
	notSure   projectStatusEnum = "notsure"
)

type secondStatusEnum string

const (
	firstValue  secondStatusEnum = "first"
	secondValue secondStatusEnum = "second"
	thirdValue  secondStatusEnum = "third"
)

func Hi(in string) {
	status, err := NewProjectStatus(in)
	if err != nil {
		panic(err)
	}

	fmt.Println("Status is:", status.String())
}
