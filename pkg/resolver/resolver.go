package resolver

import (
	"bytes"
	"context"
	"fmt"
	"gdns/pkg/protocol"
	"net"
	"time"

	"github.com/davecgh/go-spew/spew"
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
	packet, err := protocol.ParsePacket(rdr)
	if err != nil {
		return nil, fmt.Errorf("parsing packet: %w", err)
	}
	spew.Dump(packet)
	return nil, nil
}
