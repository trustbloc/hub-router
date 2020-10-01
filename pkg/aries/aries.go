/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package aries

import (
	"fmt"

	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/client/mediator"
	"github.com/hyperledger/aries-framework-go/pkg/client/outofband"
	ariescrypto "github.com/hyperledger/aries-framework-go/pkg/crypto"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"
	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	vdriapi "github.com/hyperledger/aries-framework-go/pkg/framework/aries/api/vdri"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
	"github.com/hyperledger/aries-framework-go/pkg/storage"
)

// Ctx framework context provider.
type Ctx interface {
	Service(id string) (interface{}, error)
	ServiceEndpoint() string
	StorageProvider() storage.Provider
	ProtocolStateStorageProvider() storage.Provider
	KMS() kms.KeyManager
	VDRIRegistry() vdriapi.Registry
	Crypto() ariescrypto.Crypto
}

// OutOfBand client.
type OutOfBand interface {
	CreateInvitation(protocols []string, opts ...outofband.MessageOption) (*outofband.Invitation, error)
}

// DIDExchange client.
type DIDExchange interface {
	CreateConnection(myDID string, theirDID *did.Doc, options ...didexchange.ConnectionOption) (string, error)
	RegisterActionEvent(chan<- service.DIDCommAction) error
}

// Mediator client.
type Mediator interface {
	RegisterActionEvent(chan<- service.DIDCommAction) error
}

// CreateOutofbandClient util function to create oob client.
func CreateOutofbandClient(ariesCtx outofband.Provider) (*outofband.Client, error) {
	oobClient, err := outofband.New(ariesCtx)
	if err != nil {
		return nil, fmt.Errorf("create out-of-band client : %w", err)
	}

	return oobClient, err
}

// CreateDIDExchangeClient util function to create did exchange client and registers for action event.
func CreateDIDExchangeClient(ctx Ctx, actionCh chan service.DIDCommAction) (DIDExchange, error) {
	didExClient, err := didexchange.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("create didexchange client : %w", err)
	}

	err = didExClient.RegisterActionEvent(actionCh)
	if err != nil {
		return nil, fmt.Errorf("register didexchange action event : %w", err)
	}

	return didExClient, nil
}

// CreateMediatorClient util function to create mediator client and registers for action event.
func CreateMediatorClient(ctx Ctx, actionCh chan service.DIDCommAction) (Mediator, error) {
	mediatorClient, err := mediator.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("create mediator client : %w", err)
	}

	err = mediatorClient.RegisterActionEvent(actionCh)
	if err != nil {
		return nil, fmt.Errorf("register mediator action event : %w", err)
	}

	return mediatorClient, nil
}
