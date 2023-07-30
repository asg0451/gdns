package protocol

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

	nm := "google.com"

	q := Question{
		Name:  name(nm),
		Type:  1,
		Class: 1,
	}
	bs, err = q.MarshalBinary()
	require.NoError(t, err)
	printBytes(bs)

	query := NewQuery("www.example.com", RecordTypeA)
	fmt.Printf("%+#v\n", query)
	bs, err = query.MarshalBinary()
	require.NoError(t, err)
	printBytes(bs)
	fmt.Printf("%+#v\n", bs)

	exp := []byte{
		// first 2 bytes are random id; remove them
		// 0x3c, 0x5f,
		0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x03, 0x77, 0x77, 0x77, 0x07, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x03,
		0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01}
	assert.Equal(t, exp, bs[2:])
}

func printBytes(bs []byte) {
	for _, b := range bs {
		fmt.Printf("%#x ", b)
	}
	fmt.Println()
}
