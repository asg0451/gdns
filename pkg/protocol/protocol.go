package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
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

// TODO: api?
func ParseHeader(rdr io.Reader) (*Header, error) {
	h := &Header{}
	err := binary.Read(rdr, binary.BigEndian, h)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling header: %w", err)
	}
	return h, nil
}

type Question struct {
	Name  string
	Type  RecordType
	Class uint16
}

func (q *Question) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	// write name, then rest
	if _, err := buf.Write(encodeName(q.Name)); err != nil {
		return []byte{}, fmt.Errorf("writing name: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, q.Type); err != nil {
		return []byte{}, fmt.Errorf("marshalling question: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, q.Class); err != nil {
		return []byte{}, fmt.Errorf("marshalling question: %w", err)
	}
	return buf.Bytes(), nil
}

func ParseQuestion(rdr io.Reader) (*Question, error) {
	var err error
	q := &Question{}
	if q.Name, err = decodeName(rdr); err != nil {
		return nil, fmt.Errorf("reading name: %w", err)
	}
	if err = binary.Read(rdr, binary.BigEndian, &q.Type); err != nil {
		return nil, fmt.Errorf("reading type: %w", err)
	}
	if err = binary.Read(rdr, binary.BigEndian, &q.Class); err != nil {
		return nil, fmt.Errorf("reading class: %w", err)
	}

	return q, nil
}

func encodeName(name string) []byte {
	encoded := make([]byte, 0, len(name))
	for _, part := range strings.Split(name, ".") {
		encoded = append(encoded, byte(len(part)))
		encoded = append(encoded, []byte(part)...)
	}
	return append(encoded, 0)
}

// TODO: api? also is this right
func decodeName(rdr io.Reader) (string, error) {
	ln := make([]byte, 1)
	part := make([]byte, 64)
	parts := make([]string, 0, 2)
	for {
		// read length byte
		_, err := rdr.Read(ln)
		if err != nil {
			return "", fmt.Errorf("reading length: %w", err)
		}
		length := ln[0]
		if length == 0 {
			break
		}
		// read that many bytes
		n, err := rdr.Read(part[:length])
		if err != nil {
			return "", fmt.Errorf("reading part: %w", err)
		}
		if n != int(length) {
			return "", fmt.Errorf("reading part: expected %d bytes, got %d", length, n)
		}

		parts = append(parts, string(part[:length]))
	}
	return strings.Join(parts, "."), nil
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

func NewQuery(name string, recordType RecordType) Query {
	return Query{
		Header: Header{
			Id:           uint16(rand.Intn(65535)),
			Flags:        FlagRecursionDesired,
			NumQuestions: 1,
		},
		Question: Question{
			Name:  name,
			Type:  recordType,
			Class: ClassIn,
		},
	}
}

type Record struct {
	Name  string
	Type  uint16
	Class uint16
	TTL   uint16
	Data  []byte
}

func (r *Record) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.Write(encodeName(r.Name)); err != nil {
		return []byte{}, fmt.Errorf("writing record.name: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, r.Type); err != nil {
		return []byte{}, fmt.Errorf("marshalling record.type: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, r.Class); err != nil {
		return []byte{}, fmt.Errorf("marshalling record.class: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, r.TTL); err != nil {
		return []byte{}, fmt.Errorf("marshalling record.ttl: %w", err)
	}
	if err := buf.WriteByte(byte(len(r.Data))); err != nil {
		return []byte{}, fmt.Errorf("writing record.data.length: %w", err)
	}
	if _, err := buf.Write(r.Data); err != nil {
		return []byte{}, fmt.Errorf("writing record.data: %w", err)
	}
	return buf.Bytes(), nil
}

func ParseRecord(rdr io.Reader) (*Record, error) {
	var err error
	r := &Record{}
	if r.Name, err = decodeName(rdr); err != nil {
		return nil, fmt.Errorf("reading name: %w", err)
	}
	if err = binary.Read(rdr, binary.BigEndian, &r.Type); err != nil {
		return nil, fmt.Errorf("reading type: %w", err)
	}
	if err = binary.Read(rdr, binary.BigEndian, &r.Class); err != nil {
		return nil, fmt.Errorf("reading class: %w", err)
	}
	if err = binary.Read(rdr, binary.BigEndian, &r.TTL); err != nil {
		return nil, fmt.Errorf("reading ttl: %w", err)
	}
	var dataLength byte
	if err = binary.Read(rdr, binary.BigEndian, &dataLength); err != nil {
		return nil, fmt.Errorf("reading data length: %w", err)
	}
	r.Data = make([]byte, dataLength)
	if _, err = rdr.Read(r.Data); err != nil {
		return nil, fmt.Errorf("reading data: %w", err)
	}

	return r, nil
}

// H: 2 bytes (as an integer)
// I: 4 bytes (as an integer)
// 4s: 4 bytes (as a byte string)
