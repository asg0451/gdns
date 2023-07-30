package resolver

import (
	"context"
	"gdns/pkg/protocol"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	ctx := context.Background()
	_, err := Query(ctx, protocol.NewQuery("google.com", protocol.RecordTypeA))
	require.NoError(t, err)
}
