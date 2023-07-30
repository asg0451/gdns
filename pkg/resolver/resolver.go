package resolver

import (
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
	resp := make([]byte, 512)
	_, err = conn.Read(resp)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}
	fmt.Printf("%#v\n", string(resp))
	return nil, nil
}
