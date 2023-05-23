package example

//go:generate go run -mod=mod github.com/Southclaws/enumerator

type projectStatusEnum string

const (
	success   projectStatusEnum = "success"   // Success
	failure   projectStatusEnum = "failure"   // Failure
	inBetween projectStatusEnum = "inbetween" // In between
	notSure   projectStatusEnum = "notsure"   // Not sure?
)

type secondStatusEnum string

const (
	firstValue  secondStatusEnum = "first"
	secondValue secondStatusEnum = "second"
	thirdValue  secondStatusEnum = "third"
)

func Status(in string) ProjectStatus {
	status, err := NewProjectStatus(in)
	if err != nil {
		panic(err)
	}

	return status
}
