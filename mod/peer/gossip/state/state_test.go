/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package state

import (
	"testing"

	"github.com/hyperledger/fabric/gossip/protoext"

	"github.com/trustbloc/fabric-peer-ext/pkg/roles"

	"github.com/pkg/errors"

	"github.com/hyperledger/fabric/gossip/discovery"
	"github.com/hyperledger/fabric/gossip/util"
	"github.com/hyperledger/fabric/protos/common"
	proto "github.com/hyperledger/fabric/protos/gossip"
	"github.com/stretchr/testify/require"
)

func TestProviderExtension(t *testing.T) {

	//all roles
	rolesValue := make(map[roles.Role]struct{})
	roles.SetRoles(rolesValue)
	defer func() { roles.SetRoles(nil) }()

	sampleError := errors.New("not implemented")

	handleAddPayload := func(payload *proto.Payload, blockingMode bool) error {
		return sampleError
	}

	handleStoreBlock := func(block *common.Block, pvtData util.PvtDataCollections) error {
		return sampleError
	}

	extension := NewGossipStateProviderExtension("test", nil)

	//test extension.AddPayload
	require.Error(t, sampleError, extension.AddPayload(handleAddPayload)(nil, false))

	//test extension.StoreBlock
	require.Error(t, sampleError, extension.StoreBlock(handleStoreBlock)(nil, util.PvtDataCollections{}))

	//test extension.AntiEntropy
	handled := make(chan bool, 1)
	antiEntropy := func() {
		handled <- true
	}
	extension.AntiEntropy(antiEntropy)()

	select {
	case ok := <-handled:
		require.True(t, ok)
		break
	default:
		require.Fail(t, "anti entropy handle should get called in case of all roles")

	}

	//test extension.HandleStateRequest
	handleStateRequest := func(msg protoext.ReceivedMessage) {
		handled <- true
	}

	extension.HandleStateRequest(handleStateRequest)(nil)

	select {
	case ok := <-handled:
		require.True(t, ok)
		break
	default:
		require.Fail(t, "state request handle should get called in case of all roles")
	}
}

func TestProviderByEndorser(t *testing.T) {

	//make sure roles is endorser not committer
	if roles.IsCommitter() {
		rolesValue := make(map[roles.Role]struct{})
		rolesValue[roles.EndorserRole] = struct{}{}
		roles.SetRoles(rolesValue)
		defer func() { roles.SetRoles(nil) }()
	}
	require.False(t, roles.IsCommitter())
	require.True(t, roles.IsEndorser())

	sampleError := errors.New("not implemented")

	handleAddPayload := func(payload *proto.Payload, blockingMode bool) error {
		return sampleError
	}

	handleStoreBlock := func(block *common.Block, pvtData util.PvtDataCollections) error {
		return sampleError
	}

	extension := NewGossipStateProviderExtension("test", nil)

	//test extension.AddPayload
	require.Nil(t, extension.AddPayload(handleAddPayload)(nil, false))

	//test extension.StoreBlock
	require.Nil(t, extension.StoreBlock(handleStoreBlock)(nil, util.PvtDataCollections{}))

	//test extension.AntiEntropy
	handled := make(chan bool, 1)
	antiEntropy := func() {
		handled <- true
	}
	extension.AntiEntropy(antiEntropy)()

	select {
	case <-handled:
		require.Fail(t, "anti entropy handle shouldn't get called in case of endorsers")
	default:
		//do nothing
	}

	//test extension.HandleStateRequest
	handleStateRequest := func(msg protoext.ReceivedMessage) {
		handled <- true
	}

	extension.HandleStateRequest(handleStateRequest)(nil)

	select {
	case ok := <-handled:
		require.True(t, ok)
		break
	default:
		require.Fail(t, "state request handle should get called in case of endorsers")
	}
}

func TestPredicate(t *testing.T) {
	predicate := func(peer discovery.NetworkMember) bool {
		return true
	}
	extension := NewGossipStateProviderExtension("test", nil)
	require.True(t, extension.Predicate(predicate)(discovery.NetworkMember{Properties: &proto.Properties{Roles: []string{"endorser"}}}))
	require.False(t, extension.Predicate(predicate)(discovery.NetworkMember{Properties: &proto.Properties{Roles: []string{"committer"}}}))
	require.True(t, extension.Predicate(predicate)(discovery.NetworkMember{Properties: &proto.Properties{Roles: []string{}}}))
	require.False(t, extension.Predicate(predicate)(discovery.NetworkMember{Properties: &proto.Properties{Roles: []string{""}}}))
}

func TestProviderByCommitter(t *testing.T) {

	//make sure roles is committer not endorser
	if roles.IsEndorser() {
		rolesValue := make(map[roles.Role]struct{})
		rolesValue[roles.CommitterRole] = struct{}{}
		roles.SetRoles(rolesValue)
		defer func() { roles.SetRoles(nil) }()
	}
	require.True(t, roles.IsCommitter())
	require.False(t, roles.IsEndorser())

	sampleError := errors.New("not implemented")

	handleAddPayload := func(payload *proto.Payload, blockingMode bool) error {
		return sampleError
	}

	handleStoreBlock := func(block *common.Block, pvtData util.PvtDataCollections) error {
		return sampleError
	}

	extension := NewGossipStateProviderExtension("test", nil)

	//test extension.AddPayload
	require.Error(t, sampleError, extension.AddPayload(handleAddPayload)(nil, false))

	//test extension.StoreBlock
	require.Error(t, sampleError, extension.StoreBlock(handleStoreBlock)(nil, util.PvtDataCollections{}))

	//test extension.AntiEntropy
	handled := make(chan bool, 1)
	antiEntropy := func() {
		handled <- true
	}
	extension.AntiEntropy(antiEntropy)()

	select {
	case ok := <-handled:
		require.True(t, ok)
		break
	default:
		require.Fail(t, "anti entropy handle should get called in case of committer")

	}

	//test extension.HandleStateRequest
	handleStateRequest := func(msg protoext.ReceivedMessage) {
		handled <- true
	}

	extension.HandleStateRequest(handleStateRequest)(nil)

	select {
	case <-handled:
		require.Fail(t, "state request handle shouldn't get called in case of endorsers")
	default:
		//do nothing
	}
}
