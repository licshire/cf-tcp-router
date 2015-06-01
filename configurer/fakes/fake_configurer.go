// This file was generated by counterfeiter
package fakes

import (
	"sync"

	cf_tcp_router "github.com/GESoftware-CF/cf-tcp-router"
	"github.com/GESoftware-CF/cf-tcp-router/configurer"
)

type FakeRouterConfigurer struct {
	MapBackendHostsToAvailablePortStub        func(backendHostInfos cf_tcp_router.BackendHostInfos) (cf_tcp_router.RouterHostInfo, error)
	mapBackendHostsToAvailablePortMutex       sync.RWMutex
	mapBackendHostsToAvailablePortArgsForCall []struct {
		backendHostInfos cf_tcp_router.BackendHostInfos
	}
	mapBackendHostsToAvailablePortReturns struct {
		result1 cf_tcp_router.RouterHostInfo
		result2 error
	}
}

func (fake *FakeRouterConfigurer) MapBackendHostsToAvailablePort(backendHostInfos cf_tcp_router.BackendHostInfos) (cf_tcp_router.RouterHostInfo, error) {
	fake.mapBackendHostsToAvailablePortMutex.Lock()
	fake.mapBackendHostsToAvailablePortArgsForCall = append(fake.mapBackendHostsToAvailablePortArgsForCall, struct {
		backendHostInfos cf_tcp_router.BackendHostInfos
	}{backendHostInfos})
	fake.mapBackendHostsToAvailablePortMutex.Unlock()
	if fake.MapBackendHostsToAvailablePortStub != nil {
		return fake.MapBackendHostsToAvailablePortStub(backendHostInfos)
	} else {
		return fake.mapBackendHostsToAvailablePortReturns.result1, fake.mapBackendHostsToAvailablePortReturns.result2
	}
}

func (fake *FakeRouterConfigurer) MapBackendHostsToAvailablePortCallCount() int {
	fake.mapBackendHostsToAvailablePortMutex.RLock()
	defer fake.mapBackendHostsToAvailablePortMutex.RUnlock()
	return len(fake.mapBackendHostsToAvailablePortArgsForCall)
}

func (fake *FakeRouterConfigurer) MapBackendHostsToAvailablePortArgsForCall(i int) cf_tcp_router.BackendHostInfos {
	fake.mapBackendHostsToAvailablePortMutex.RLock()
	defer fake.mapBackendHostsToAvailablePortMutex.RUnlock()
	return fake.mapBackendHostsToAvailablePortArgsForCall[i].backendHostInfos
}

func (fake *FakeRouterConfigurer) MapBackendHostsToAvailablePortReturns(result1 cf_tcp_router.RouterHostInfo, result2 error) {
	fake.MapBackendHostsToAvailablePortStub = nil
	fake.mapBackendHostsToAvailablePortReturns = struct {
		result1 cf_tcp_router.RouterHostInfo
		result2 error
	}{result1, result2}
}

var _ configurer.RouterConfigurer = new(FakeRouterConfigurer)