package protocol

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {

	h := Header{
		Id:             0x1314,
		Flags:          0,
		NumQuestions:   1,
		NumAdditionals: 0,
		NumAuthorities: 0,
		NumAnswers:     0,
	}
	bs, err := h.MarshalBinary()
	require.NoError(t, err)
	printBytes(bs)

	nm := NewName("google.com")

	q := Question{
		Name:  nm,
		Type:  1,
		Class: 1,
	}
	bs, err = q.MarshalBinary()
	require.NoError(t, err)
	printBytes(bs)

	printBytes(nm)

	query := NewQuery("google.com", RecordTypeA)
	fmt.Printf("%+#v\n", query)
	bs, err = query.MarshalBinary()
	require.NoError(t, err)
	printBytes(bs)
}

func printBytes(bs []byte) {
	for _, b := range bs {
		fmt.Printf("%#x ", b)
	}
	fmt.Println()
}
