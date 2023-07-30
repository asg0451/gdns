package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"strings"
)

type Header struct {
	Id             uint16
	Flags          Flag
	NumQuestions   uint16
	NumAnswers     uint16
	NumAuthorities uint16
	NumAdditionals uint16
}

func (h *Header) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, h)
	if err != nil {
		return []byte{}, fmt.Errorf("marshalling header: %w", err)
	}
	return buf.Bytes(), nil
}

type Question struct {
	Name  Name
	Type  RecordType
	Class uint16
}

func (q *Question) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	// write name, then rest
	_, err := buf.Write(q.Name)
	if err != nil {
		return []byte{}, fmt.Errorf("writing name: %w", err)
	}

	err = binary.Write(buf, binary.BigEndian, q.Type)
	if err != nil {
		return []byte{}, fmt.Errorf("marshalling question: %w", err)
	}

	err = binary.Write(buf, binary.BigEndian, q.Class)
	if err != nil {
		return []byte{}, fmt.Errorf("marshalling question: %w", err)
	}
	return buf.Bytes(), nil
}

// TODO: private
type Name []byte

func NewName(name string) Name {
	encoded := make([]byte, 0, len(name))
	for _, part := range strings.Split(name, ".") {
		encoded = append(encoded, byte(len(part)))
		encoded = append(encoded, []byte(part)...)
	}
	return append(encoded, 0)
}

const (
	ClassIn uint16 = 1
)

type RecordType uint16

const (
	RecordTypeA RecordType = 1
)

type Flag uint16

const (
	FlagRecursionDesired Flag = 1 << 8
)

type Query struct {
	Header   Header
	Question Question
}

func (q *Query) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	bs, err := q.Header.MarshalBinary()
	if err != nil {
		return []byte{}, fmt.Errorf("marshalling query.header: %w", err)
	}
	_, err = buf.Write(bs)
	if err != nil {
		return []byte{}, fmt.Errorf("writing query.header: %w", err)
	}
	bs, err = q.Question.MarshalBinary()
	if err != nil {
		return []byte{}, fmt.Errorf("marshalling query.question: %w", err)
	}
	_, err = buf.Write(bs)
	if err != nil {
		return []byte{}, fmt.Errorf("writing query.question: %w", err)
	}
	return buf.Bytes(), nil
}

func NewQuery(name string, receordType RecordType) Query {
	return Query{
		Header: Header{
			Id:           uint16(rand.Intn(65535)),
			Flags:        FlagRecursionDesired,
			NumQuestions: 1,
		},
		Question: Question{
			Name:  NewName(name),
			Type:  RecordTypeA,
			Class: ClassIn,
		},
	}
}

// H: 2 bytes (as an integer)
// I: 4 bytes (as an integer)
// 4s: 4 bytes (as a byte string)
