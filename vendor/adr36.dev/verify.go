// Copyright 2024 NASD Inc.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file or at
// https://opensource.org/licenses/MIT.

package adr36

import (
	"context"

	api "adr36.dev/api"
	tx "cosmossdk.io/api/cosmos/tx/v1beta1"
	"cosmossdk.io/x/tx/signing"
	"cosmossdk.io/x/tx/signing/aminojson"
	"github.com/cosmos/cosmos-proto/anyutil"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/protobuf/types/known/anypb"
)

func VerifySignature(pubkey cryptotypes.PubKey, data []byte, signature []byte) bool {
	handler := aminojson.NewSignModeHandler(aminojson.SignModeHandlerOptions{})
	rawPubkey, _ := codectypes.NewAnyWithValue(pubkey)

	msg, _ := anyutil.New(&api.MsgSignData{
		Signer: sdk.AccAddress(pubkey.Address()).String(),
		Data:   data,
	})

	bz, _ := handler.GetSignBytes(
		context.Background(),
		signing.SignerData{
			Address:       sdk.AccAddress(pubkey.Address()).String(),
			ChainID:       "",
			AccountNumber: 0,
			Sequence:      0,
			PubKey: &anypb.Any{
				TypeUrl: rawPubkey.TypeUrl,
				Value:   rawPubkey.Value,
			},
		},
		signing.TxData{
			Body: &tx.TxBody{
				Messages: []*anypb.Any{msg},
			},
			AuthInfo: &tx.AuthInfo{
				Fee: &tx.Fee{},
			},
		},
	)

	return pubkey.VerifySignature(bz, signature)
}
