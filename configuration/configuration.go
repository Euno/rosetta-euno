// Copyright 2020 Coinbase, Inc.
// Copyright 2022 EUNO
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package configuration

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/Euno/rosetta-euno/bitcoin"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/coinbase/rosetta-sdk-go/storage/encoder"
	"github.com/coinbase/rosetta-sdk-go/types"
)

// Mode is the setting that determines if
// the implementation is "online" or "offline".
type Mode string

const (
	// Online is when the implementation is permitted
	// to make outbound connections.
	Online Mode = "ONLINE"

	// Offline is when the implementation is not permitted
	// to make outbound connections.
	Offline Mode = "OFFLINE"

	// Mainnet is the Euno Mainnet.
	Mainnet string = "MAINNET"

	// Testnet is Euno Testnet3.
	Testnet string = "TESTNET"

	// mainnetConfigPath is the path of the Bitcoin
	// configuration file for mainnet.
	mainnetConfigPath = "/app/euno-mainnet.conf"

	// testnetConfigPath is the path of the Bitcoin
	// configuration file for testnet.
	testnetConfigPath = "/app/euno-testnet.conf"

	mainnetRPCPort = 46461
	testnetRPCPort = 46465

	pruneDepth = int64(10000) //nolint

	minPruneHeight = int64(100000) //nolint

	// attempt to prune once an hour
	pruneFrequency = 60 * time.Minute

	// DataDirectory is the default location for all
	// persistent data.
	DataDirectory = "/data"

	bitcoindPath = "eunod"
	indexerPath  = "indexer"

	// allFilePermissions specifies anyone can do anything
	// to the file.
	allFilePermissions = 0777

	// ModeEnv is the environment variable read
	// to determine mode.
	ModeEnv = "MODE"

	// NetworkEnv is the environment variable
	// read to determine network.
	NetworkEnv = "NETWORK"

	// PortEnv is the environment variable
	// read to determine the port for the Rosetta
	// implementation.
	PortEnv = "PORT"
)

// PruningConfiguration is the configuration to
// use for pruning in the indexer.
type PruningConfiguration struct {
	Frequency time.Duration
	Depth     int64
	MinHeight int64
}

// Configuration determines how
type Configuration struct {
	Mode                   Mode
	Network                *types.NetworkIdentifier
	Params                 *chaincfg.Params
	Currency               *types.Currency
	GenesisBlockIdentifier *types.BlockIdentifier
	Port                   int
	RPCPort                int
	ConfigPath             string
	Pruning                *PruningConfiguration
	IndexerPath            string
	BitcoindPath           string
	Compressors            []*encoder.CompressorEntry
}

// LoadConfiguration attempts to create a new Configuration
// using the ENVs in the environment.
func LoadConfiguration(baseDirectory string) (*Configuration, error) {
	config := &Configuration{}
	config.Pruning = &PruningConfiguration{
		Frequency: pruneFrequency,
		Depth:     pruneDepth,
		MinHeight: minPruneHeight,
	}

	modeValue := Mode(os.Getenv(ModeEnv))
	switch modeValue {
	case Online:
		config.Mode = Online
		config.IndexerPath = path.Join(baseDirectory, indexerPath)
		if err := ensurePathExists(config.IndexerPath); err != nil {
			return nil, fmt.Errorf("%w: unable to create indexer path", err)
		}

		config.BitcoindPath = path.Join(baseDirectory, bitcoindPath)
		if err := ensurePathExists(config.BitcoindPath); err != nil {
			return nil, fmt.Errorf("%w: unable to create bitcoind path", err)
		}
	case Offline:
		config.Mode = Offline
	case "":
		return nil, errors.New("MODE must be populated")
	default:
		return nil, fmt.Errorf("%s is not a valid mode", modeValue)
	}

	networkValue := os.Getenv(NetworkEnv)
	switch networkValue {
	case Mainnet:
		config.Network = &types.NetworkIdentifier{
			Blockchain: bitcoin.Blockchain,
			Network:    bitcoin.MainnetNetwork,
		}
		config.GenesisBlockIdentifier = bitcoin.MainnetGenesisBlockIdentifier
		config.Params = bitcoin.MainnetParams
		config.Currency = bitcoin.MainnetCurrency
		config.ConfigPath = mainnetConfigPath
		config.RPCPort = mainnetRPCPort
	case Testnet:
		config.Network = &types.NetworkIdentifier{
			Blockchain: bitcoin.Blockchain,
			Network:    bitcoin.TestnetNetwork,
		}
		config.GenesisBlockIdentifier = bitcoin.TestnetGenesisBlockIdentifier
		config.Params = bitcoin.TestnetParams
		config.Currency = bitcoin.TestnetCurrency
		config.ConfigPath = testnetConfigPath
		config.RPCPort = testnetRPCPort
	case "":
		return nil, errors.New("NETWORK must be populated")
	default:
		return nil, fmt.Errorf("%s is not a valid network", networkValue)
	}

	portValue := os.Getenv(PortEnv)
	if len(portValue) == 0 {
		return nil, errors.New("PORT must be populated")
	}

	port, err := strconv.Atoi(portValue)
	if err != nil || len(portValue) == 0 || port <= 0 {
		return nil, fmt.Errorf("%w: unable to parse port %s", err, portValue)
	}
	config.Port = port

	return config, nil
}

// ensurePathsExist directories along
// a path if they do not exist.
func ensurePathExists(path string) error {
	if err := os.MkdirAll(path, os.FileMode(allFilePermissions)); err != nil {
		return fmt.Errorf("%w: unable to create %s directory", err, path)
	}

	return nil
}
