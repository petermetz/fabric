// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import common "github.com/petermetz/fabric/protos/common"
import mock "github.com/stretchr/testify/mock"

// BlockVerifier is an autogenerated mock type for the BlockVerifier type
type BlockVerifier struct {
	mock.Mock
}

// VerifyBlockSignature provides a mock function with given fields: sd, config
func (_m *BlockVerifier) VerifyBlockSignature(sd []*common.SignedData, config *common.ConfigEnvelope) error {
	ret := _m.Called(sd, config)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*common.SignedData, *common.ConfigEnvelope) error); ok {
		r0 = rf(sd, config)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
