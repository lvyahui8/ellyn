package ellyn_agent

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func varDefTest(t *testing.T, list []*VarDef) {
	defList := NewVarDefList(list)
	encoded := defList.Encode()
	newDefList := decodeVarDef(encoded)
	require.Equal(t, encoded, newDefList.Encode())
}

func TestVarDefList(t *testing.T) {
	varDefTest(t, []*VarDef{
		&VarDef{
			Names: []string{"a", "b"},
			Type:  "string",
		},
		&VarDef{
			Names: []string{"c"},
			Type:  "int",
		},
	})
	varDefTest(t, []*VarDef{})
	varDefTest(t, nil)
}
