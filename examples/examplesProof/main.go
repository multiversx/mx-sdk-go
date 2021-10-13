package main

import (
	"encoding/hex"

	"github.com/ElrondNetwork/elrond-go-core/hashing/blake2b"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/examples"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/interactors"
)

var log = logger.GetOrCreate("examples/examplesProof")

// TODO this will work only if the Merkle Proof routes are opened on the testnet
func main() {
	ep := blockchain.NewElrondProxy(examples.TestnetGateway, nil)

	mpv, err := interactors.NewMerkleProofVerifier(&marshal.GogoProtoMarshalizer{}, blake2b.NewBlake2b())
	if err != nil {
		log.Error("could not create Merkle Proof verifier", "error", err.Error())
		return
	}

	aliceAddress := "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th"
	proofResponse, err := ep.GetMerkleProof(aliceAddress)
	if err != nil {
		log.Error("could not get Merkle Proof", "error", err.Error())
		return
	}

	proof, err := hexArrayToBytes(proofResponse.Data.Proof)
	if err != nil {
		log.Error("could not convert hexProof to proof bytes", "error", err.Error())
		return
	}

	ok, err := mpv.VerifyProof(proofResponse.Data.RootHash, aliceAddress, proof)
	if err != nil {
		log.Error("could not verify Merkle Proof", "error", err.Error())
		return
	}

	log.Info("response", "isProofValid", ok)
}

func hexArrayToBytes(hexValue []string) ([][]byte, error) {
	bytesValue := make([][]byte, 0)
	for i := range hexValue {
		val, err := hex.DecodeString(hexValue[i])
		if err != nil {
			return nil, err
		}

		bytesValue = append(bytesValue, val)
	}

	return bytesValue, nil
}
