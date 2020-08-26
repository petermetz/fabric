// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import transientstore "github.com/petermetz/fabric/core/transientstore"

// RWSetScanner is an autogenerated mock type for the RWSetScanner type
type RWSetScanner struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *RWSetScanner) Close() {
	_m.Called()
}

// Next provides a mock function with given fields:
func (_m *RWSetScanner) Next() (*transientstore.EndorserPvtSimulationResults, error) {
	ret := _m.Called()

	var r0 *transientstore.EndorserPvtSimulationResults
	if rf, ok := ret.Get(0).(func() *transientstore.EndorserPvtSimulationResults); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transientstore.EndorserPvtSimulationResults)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NextWithConfig provides a mock function with given fields:
func (_m *RWSetScanner) NextWithConfig() (*transientstore.EndorserPvtSimulationResultsWithConfig, error) {
	ret := _m.Called()

	var r0 *transientstore.EndorserPvtSimulationResultsWithConfig
	if rf, ok := ret.Get(0).(func() *transientstore.EndorserPvtSimulationResultsWithConfig); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transientstore.EndorserPvtSimulationResultsWithConfig)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
