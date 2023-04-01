// Code generated by mockery v2.16.0. DO NOT EDIT.

package data

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockAppDataService is an autogenerated mock type for the AppDataService type
type MockAppDataService struct {
	mock.Mock
}

// AddFile provides a mock function with given fields: ctx, filepath, description
func (_m *MockAppDataService) AddFile(ctx context.Context, filepath string, description string) (File, error) {
	ret := _m.Called(ctx, filepath, description)

	var r0 File
	if rf, ok := ret.Get(0).(func(context.Context, string, string) File); ok {
		r0 = rf(ctx, filepath, description)
	} else {
		r0 = ret.Get(0).(File)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, filepath, description)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFile provides a mock function with given fields: ctx, id
func (_m *MockAppDataService) GetFile(ctx context.Context, id string) (File, error) {
	ret := _m.Called(ctx, id)

	var r0 File
	if rf, ok := ret.Get(0).(func(context.Context, string) File); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(File)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockAppDataService interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockAppDataService creates a new instance of MockAppDataService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockAppDataService(t mockConstructorTestingTNewMockAppDataService) *MockAppDataService {
	mock := &MockAppDataService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}