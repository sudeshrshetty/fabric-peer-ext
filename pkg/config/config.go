/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package config

import (
	"path/filepath"

	"github.com/hyperledger/fabric/core/ledger/ledgerconfig"
	"github.com/spf13/viper"
)

const (
	confRoles                = "ledger.roles"
	confPvtDataCacheSize     = "ledger.blockchain.pvtDataStorage.cacheSize"
	confTransientDataLeveldb = "transientDataLeveldb"
)

// GetRoles returns the roles of the peer. Empty return value indicates that the peer has all roles.
func GetRoles() string {
	return viper.GetString(confRoles)
}

// GetPvtDataCacheSize returns the number of pvt data per block to keep the in the cache
func GetPvtDataCacheSize() int {
	pvtDataCacheSize := viper.GetInt(confPvtDataCacheSize)
	if !viper.IsSet(confPvtDataCacheSize) {
		pvtDataCacheSize = 10
	}
	return pvtDataCacheSize
}

// GetTransientDataLevelDBPath returns the filesystem path that is used to maintain the transient data level db
func GetTransientDataLevelDBPath() string {
	return filepath.Join(ledgerconfig.GetRootPath(), confTransientDataLeveldb)
}
