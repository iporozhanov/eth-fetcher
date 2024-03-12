// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"

	common "github.com/ethereum/go-ethereum/common"

	mock "github.com/stretchr/testify/mock"

	types "github.com/ethereum/go-ethereum/core/types"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

type Client_Expecter struct {
	mock *mock.Mock
}

func (_m *Client) EXPECT() *Client_Expecter {
	return &Client_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with given fields:
func (_m *Client) Close() {
	_m.Called()
}

// Client_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type Client_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *Client_Expecter) Close() *Client_Close_Call {
	return &Client_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *Client_Close_Call) Run(run func()) *Client_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Client_Close_Call) Return() *Client_Close_Call {
	_c.Call.Return()
	return _c
}

func (_c *Client_Close_Call) RunAndReturn(run func()) *Client_Close_Call {
	_c.Call.Return(run)
	return _c
}

// TransactionByHash provides a mock function with given fields: ctx, txHash
func (_m *Client) TransactionByHash(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error) {
	ret := _m.Called(ctx, txHash)

	if len(ret) == 0 {
		panic("no return value specified for TransactionByHash")
	}

	var r0 *types.Transaction
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) (*types.Transaction, bool, error)); ok {
		return rf(ctx, txHash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) *types.Transaction); ok {
		r0 = rf(ctx, txHash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Hash) bool); ok {
		r1 = rf(ctx, txHash)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(context.Context, common.Hash) error); ok {
		r2 = rf(ctx, txHash)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Client_TransactionByHash_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TransactionByHash'
type Client_TransactionByHash_Call struct {
	*mock.Call
}

// TransactionByHash is a helper method to define mock.On call
//   - ctx context.Context
//   - txHash common.Hash
func (_e *Client_Expecter) TransactionByHash(ctx interface{}, txHash interface{}) *Client_TransactionByHash_Call {
	return &Client_TransactionByHash_Call{Call: _e.mock.On("TransactionByHash", ctx, txHash)}
}

func (_c *Client_TransactionByHash_Call) Run(run func(ctx context.Context, txHash common.Hash)) *Client_TransactionByHash_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.Hash))
	})
	return _c
}

func (_c *Client_TransactionByHash_Call) Return(tx *types.Transaction, isPending bool, err error) *Client_TransactionByHash_Call {
	_c.Call.Return(tx, isPending, err)
	return _c
}

func (_c *Client_TransactionByHash_Call) RunAndReturn(run func(context.Context, common.Hash) (*types.Transaction, bool, error)) *Client_TransactionByHash_Call {
	_c.Call.Return(run)
	return _c
}

// TransactionReceipt provides a mock function with given fields: ctx, txHash
func (_m *Client) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	ret := _m.Called(ctx, txHash)

	if len(ret) == 0 {
		panic("no return value specified for TransactionReceipt")
	}

	var r0 *types.Receipt
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) (*types.Receipt, error)); ok {
		return rf(ctx, txHash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) *types.Receipt); ok {
		r0 = rf(ctx, txHash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Receipt)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Hash) error); ok {
		r1 = rf(ctx, txHash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_TransactionReceipt_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TransactionReceipt'
type Client_TransactionReceipt_Call struct {
	*mock.Call
}

// TransactionReceipt is a helper method to define mock.On call
//   - ctx context.Context
//   - txHash common.Hash
func (_e *Client_Expecter) TransactionReceipt(ctx interface{}, txHash interface{}) *Client_TransactionReceipt_Call {
	return &Client_TransactionReceipt_Call{Call: _e.mock.On("TransactionReceipt", ctx, txHash)}
}

func (_c *Client_TransactionReceipt_Call) Run(run func(ctx context.Context, txHash common.Hash)) *Client_TransactionReceipt_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.Hash))
	})
	return _c
}

func (_c *Client_TransactionReceipt_Call) Return(_a0 *types.Receipt, _a1 error) *Client_TransactionReceipt_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Client_TransactionReceipt_Call) RunAndReturn(run func(context.Context, common.Hash) (*types.Receipt, error)) *Client_TransactionReceipt_Call {
	_c.Call.Return(run)
	return _c
}

// NewClient creates a new instance of Client. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *Client {
	mock := &Client{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
