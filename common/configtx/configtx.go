/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package configtx

import (
	cb "github.com/petermetz/fabric/protos/common"
)

// Validator provides a mechanism to propose config updates, see the config update results
// and validate the results of a config update.
type Validator interface {
	// Validate attempts to apply a configtx to become the new config
	Validate(configEnv *cb.ConfigEnvelope) error

	// Validate attempts to validate a new configtx against the current config state
	ProposeConfigUpdate(configtx *cb.Envelope) (*cb.ConfigEnvelope, error)

	// ChainID retrieves the chain ID associated with this manager
	ChainID() string

	// ConfigProto returns the current config as a proto
	ConfigProto() *cb.Config

	// Sequence returns the current sequence number of the config
	Sequence() uint64
}
