package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/spire/pkg/common/nodeutil"
	"github.com/spiffe/spire/pkg/server/plugin/datastore"
	"github.com/spiffe/spire/proto/spire/common"
	"github.com/spiffe/spire/proto/spire/types"
)

// AuthorizedEntryFetcher is the interface to fetch authorized entries
type AuthorizedEntryFetcher interface {
	// FetchAuthorizedEntries fetches the entries that the specified
	// SPIFFE ID is authorized for
	FetchAuthorizedEntries(ctx context.Context, id spiffeid.ID) ([]*types.Entry, error)
}

// AuthorizedEntryFetcherFunc is an implementation of AuthorizedEntryFetcher
// using a function.
type AuthorizedEntryFetcherFunc func(ctx context.Context, id spiffeid.ID) ([]*types.Entry, error)

// FetchAuthorizedEntries fetches the entries that the specified
// SPIFFE ID is authorized for
func (fn AuthorizedEntryFetcherFunc) FetchAuthorizedEntries(ctx context.Context, id spiffeid.ID) ([]*types.Entry, error) {
	return fn(ctx, id)
}

// AttestedNodeToProto converts an agent from the given *common.AttestedNode with
// the provided selectors to *types.Agent
func AttestedNodeToProto(node *common.AttestedNode, selectors []*types.Selector) (*types.Agent, error) {
	if node == nil {
		return nil, errors.New("missing node")
	}

	spiffeID, err := spiffeid.FromString(node.SpiffeId)
	if err != nil {
		return nil, fmt.Errorf("node has malformed SPIFFE ID: %v", err)
	}

	return &types.Agent{
		Id:                   ProtoFromID(spiffeID),
		AttestationType:      node.AttestationDataType,
		X509SvidSerialNumber: node.CertSerialNumber,
		X509SvidExpiresAt:    node.CertNotAfter,
		Selectors:            selectors,
		Banned:               nodeutil.IsAgentBanned(node),
	}, nil
}

// NodeSelectorsToProto converts node selectors from the given
// *datastore.NodeSelectors to []*types.Selector
func NodeSelectorsToProto(nodeSelectors *datastore.NodeSelectors) ([]*types.Selector, error) {
	if nodeSelectors == nil {
		return nil, errors.New("missing node selectors")
	}

	var selectors []*types.Selector
	for _, s := range nodeSelectors.Selectors {
		selectors = append(selectors, &types.Selector{
			Type:  s.Type,
			Value: s.Value,
		})
	}

	return selectors, nil
}
