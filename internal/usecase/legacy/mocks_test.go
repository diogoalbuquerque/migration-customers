// Code generated by MockGen. DO NOT EDIT.
// Source: internal/usecase/legacy/interfaces.go

package usecase_test

import (
	context "context"
	reflect "reflect"

	entity "github.com/diogoalbuquerque/migration-customers/internal/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockLegacyPersonDB2Repo is a mock of LegacyPersonDB2Repo interface.
type MockLegacyPersonDB2Repo struct {
	ctrl     *gomock.Controller
	recorder *MockLegacyPersonDB2RepoMockRecorder
}

// MockLegacyPersonDB2RepoMockRecorder is the mock recorder for MockLegacyPersonDB2Repo.
type MockLegacyPersonDB2RepoMockRecorder struct {
	mock *MockLegacyPersonDB2Repo
}

// NewMockLegacyPersonDB2Repo creates a new mock instance.
func NewMockLegacyPersonDB2Repo(ctrl *gomock.Controller) *MockLegacyPersonDB2Repo {
	mock := &MockLegacyPersonDB2Repo{ctrl: ctrl}
	mock.recorder = &MockLegacyPersonDB2RepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLegacyPersonDB2Repo) EXPECT() *MockLegacyPersonDB2RepoMockRecorder {
	return m.recorder
}

// GetBucketsAvailable mocks base method.
func (m *MockLegacyPersonDB2Repo) GetBucketsAvailable() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBucketsAvailable")
	ret0, _ := ret[0].(int)
	return ret0
}

// GetBucketsAvailable indicates an expected call of GetBucketsAvailable.
func (mr *MockLegacyPersonDB2RepoMockRecorder) GetBucketsAvailable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBucketsAvailable", reflect.TypeOf((*MockLegacyPersonDB2Repo)(nil).GetBucketsAvailable))
}

// GetLegacyPeople mocks base method.
func (m *MockLegacyPersonDB2Repo) GetLegacyPeople(ctx context.Context, limit int) ([]entity.LegacyPerson, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLegacyPeople", ctx, limit)
	ret0, _ := ret[0].([]entity.LegacyPerson)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLegacyPeople indicates an expected call of GetLegacyPeople.
func (mr *MockLegacyPersonDB2RepoMockRecorder) GetLegacyPeople(ctx, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLegacyPeople", reflect.TypeOf((*MockLegacyPersonDB2Repo)(nil).GetLegacyPeople), ctx, limit)
}

// UpdateLegacyPeople mocks base method.
func (m *MockLegacyPersonDB2Repo) UpdateLegacyPeople(ctx context.Context, legacyPeople []entity.LegacyPerson) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateLegacyPeople", ctx, legacyPeople)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateLegacyPeople indicates an expected call of UpdateLegacyPeople.
func (mr *MockLegacyPersonDB2RepoMockRecorder) UpdateLegacyPeople(ctx, legacyPeople interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLegacyPeople", reflect.TypeOf((*MockLegacyPersonDB2Repo)(nil).UpdateLegacyPeople), ctx, legacyPeople)
}
