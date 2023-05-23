package test

import (
	"fmt"
	"testing"

	"github.com/Southclaws/enumerator/example"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerated(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	s, err := example.NewProjectStatus("invalid")
	r.Error(err)
	r.Empty(s)

	s, err = example.NewProjectStatus("success")
	r.NoError(err)

	a.Equal("success", s.String())
	a.Equal("success", fmt.Sprintf("%s", s))
	a.Equal(`"success"`, fmt.Sprintf("%q", s))
	a.Equal("Success", fmt.Sprintf("%v", s))
}
