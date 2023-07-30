package resolver

import (
	"bytes"
	"context"
	"fmt"
	"gdns/pkg/protocol"
	"net"
	"time"
)

func Query(ctx context.Context, q protocol.Query) (*int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	dlr := &net.Dialer{}
	conn, err := dlr.DialContext(ctx, "udp", "8.8.8.8:53")
	if err != nil {
		return nil, fmt.Errorf("dialing: %w", err)
	}
	defer conn.Close()
	bs, err := q.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshalling query: %w", err)
	}
	_, err = conn.Write(bs)
	if err != nil {
		return nil, fmt.Errorf("writing query: %w", err)
	}
	resp := make([]byte, 512) // TODO: read more than length, multiple etc
	_, err = conn.Read(resp)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	rdr := bytes.NewReader(resp)
	header, err := protocol.ParseHeader(rdr)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response.header: %w", err)
	}
	fmt.Printf("header: %+v\n", header)
	question, err := protocol.ParseQuestion(rdr)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response.question: %w", err)
	}
	fmt.Printf("question: %+v\n", question)
	record, err := protocol.ParseRecord(rdr)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response.record: %w", err)
	}
	fmt.Printf("record: %+v\n", record)

	return nil, nil
}
