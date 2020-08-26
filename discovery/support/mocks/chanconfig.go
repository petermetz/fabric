// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/petermetz/fabric/common/channelconfig"
	"github.com/petermetz/fabric/discovery/support/acl"
)

type ChanConfig struct {
	GetChannelConfigStub        func(cid string) channelconfig.Resources
	getChannelConfigMutex       sync.RWMutex
	getChannelConfigArgsForCall []struct {
		cid string
	}
	getChannelConfigReturns struct {
		result1 channelconfig.Resources
	}
	getChannelConfigReturnsOnCall map[int]struct {
		result1 channelconfig.Resources
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ChanConfig) GetChannelConfig(cid string) channelconfig.Resources {
	fake.getChannelConfigMutex.Lock()
	ret, specificReturn := fake.getChannelConfigReturnsOnCall[len(fake.getChannelConfigArgsForCall)]
	fake.getChannelConfigArgsForCall = append(fake.getChannelConfigArgsForCall, struct {
		cid string
	}{cid})
	fake.recordInvocation("GetChannelConfig", []interface{}{cid})
	fake.getChannelConfigMutex.Unlock()
	if fake.GetChannelConfigStub != nil {
		return fake.GetChannelConfigStub(cid)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getChannelConfigReturns.result1
}

func (fake *ChanConfig) GetChannelConfigCallCount() int {
	fake.getChannelConfigMutex.RLock()
	defer fake.getChannelConfigMutex.RUnlock()
	return len(fake.getChannelConfigArgsForCall)
}

func (fake *ChanConfig) GetChannelConfigArgsForCall(i int) string {
	fake.getChannelConfigMutex.RLock()
	defer fake.getChannelConfigMutex.RUnlock()
	return fake.getChannelConfigArgsForCall[i].cid
}

func (fake *ChanConfig) GetChannelConfigReturns(result1 channelconfig.Resources) {
	fake.GetChannelConfigStub = nil
	fake.getChannelConfigReturns = struct {
		result1 channelconfig.Resources
	}{result1}
}

func (fake *ChanConfig) GetChannelConfigReturnsOnCall(i int, result1 channelconfig.Resources) {
	fake.GetChannelConfigStub = nil
	if fake.getChannelConfigReturnsOnCall == nil {
		fake.getChannelConfigReturnsOnCall = make(map[int]struct {
			result1 channelconfig.Resources
		})
	}
	fake.getChannelConfigReturnsOnCall[i] = struct {
		result1 channelconfig.Resources
	}{result1}
}

func (fake *ChanConfig) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getChannelConfigMutex.RLock()
	defer fake.getChannelConfigMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ChanConfig) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ acl.ChannelConfigGetter = new(ChanConfig)
