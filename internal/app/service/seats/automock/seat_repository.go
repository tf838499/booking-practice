// Code generated by MockGen. DO NOT EDIT.
// Source: ticket-booking/internal/app/service/seats (interfaces: SeatRepository)

// Package automock is a generated GoMock package.
package automock

import (
	context "context"
	reflect "reflect"
	PostgresDB "ticket-booking/internal/adapter/repository/pgDB"

	gomock "github.com/golang/mock/gomock"
)

// MockSeatRepository is a mock of SeatRepository interface.
type MockSeatRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSeatRepositoryMockRecorder
}

// MockSeatRepositoryMockRecorder is the mock recorder for MockSeatRepository.
type MockSeatRepositoryMockRecorder struct {
	mock *MockSeatRepository
}

// NewMockSeatRepository creates a new mock instance.
func NewMockSeatRepository(ctrl *gomock.Controller) *MockSeatRepository {
	mock := &MockSeatRepository{ctrl: ctrl}
	mock.recorder = &MockSeatRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSeatRepository) EXPECT() *MockSeatRepositoryMockRecorder {
	return m.recorder
}

// SelectSeatsList mocks base method.
func (m *MockSeatRepository) SelectSeatsList(arg0 context.Context, arg1 int32) ([]PostgresDB.Seat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectSeatsList", arg0, arg1)
	ret0, _ := ret[0].([]PostgresDB.Seat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectSeatsList indicates an expected call of SelectSeatsList.
func (mr *MockSeatRepositoryMockRecorder) SelectSeatsList(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectSeatsList", reflect.TypeOf((*MockSeatRepository)(nil).SelectSeatsList), arg0, arg1)
}

// SelectSectionsAndPricesByEventID mocks base method.
func (m *MockSeatRepository) SelectSectionsAndPricesByEventID(arg0 context.Context, arg1 int32) ([]PostgresDB.SelectSectionsAndPricesByEventIDRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectSectionsAndPricesByEventID", arg0, arg1)
	ret0, _ := ret[0].([]PostgresDB.SelectSectionsAndPricesByEventIDRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectSectionsAndPricesByEventID indicates an expected call of SelectSectionsAndPricesByEventID.
func (mr *MockSeatRepositoryMockRecorder) SelectSectionsAndPricesByEventID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectSectionsAndPricesByEventID", reflect.TypeOf((*MockSeatRepository)(nil).SelectSectionsAndPricesByEventID), arg0, arg1)
}
