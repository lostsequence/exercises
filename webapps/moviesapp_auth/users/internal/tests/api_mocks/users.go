// Code generated by MockGen. DO NOT EDIT.
// Source: users.go
//
// Generated by this command:
//
//	mockgen -source users.go -destination ../../tests/api_mocks/users.go package apimocks
//
// Package mock_handlers is a generated GoMock package.
package mock_handlers

import (
	context "context"
	reflect "reflect"
	domain "movies-auth/users/internal/domain"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockUsersService is a mock of UsersService interface.
type MockUsersService struct {
	ctrl     *gomock.Controller
	recorder *MockUsersServiceMockRecorder
}

// MockUsersServiceMockRecorder is the mock recorder for MockUsersService.
type MockUsersServiceMockRecorder struct {
	mock *MockUsersService
}

// NewMockUsersService creates a new mock instance.
func NewMockUsersService(ctrl *gomock.Controller) *MockUsersService {
	mock := &MockUsersService{ctrl: ctrl}
	mock.recorder = &MockUsersServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsersService) EXPECT() *MockUsersServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUsersService) Create(ctx context.Context, user domain.User) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, user)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUsersServiceMockRecorder) Create(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUsersService)(nil).Create), ctx, user)
}

// Login mocks base method.
func (m *MockUsersService) Login(ctx context.Context, user domain.User) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, user)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUsersServiceMockRecorder) Login(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUsersService)(nil).Login), ctx, user)
}

// MockSessionService is a mock of SessionService interface.
type MockSessionService struct {
	ctrl     *gomock.Controller
	recorder *MockSessionServiceMockRecorder
}

// MockSessionServiceMockRecorder is the mock recorder for MockSessionService.
type MockSessionServiceMockRecorder struct {
	mock *MockSessionService
}

// NewMockSessionService creates a new mock instance.
func NewMockSessionService(ctrl *gomock.Controller) *MockSessionService {
	mock := &MockSessionService{ctrl: ctrl}
	mock.recorder = &MockSessionServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionService) EXPECT() *MockSessionServiceMockRecorder {
	return m.recorder
}

// CreateSession mocks base method.
func (m *MockSessionService) CreateSession(userId int) (domain.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", userId)
	ret0, _ := ret[0].(domain.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockSessionServiceMockRecorder) CreateSession(userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockSessionService)(nil).CreateSession), userId)
}

// DeleteSession mocks base method.
func (m *MockSessionService) DeleteSession(key uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSession", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSession indicates an expected call of DeleteSession.
func (mr *MockSessionServiceMockRecorder) DeleteSession(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSession", reflect.TypeOf((*MockSessionService)(nil).DeleteSession), key)
}

// ValidateSession mocks base method.
func (m *MockSessionService) ValidateSession(key uuid.UUID) (domain.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateSession", key)
	ret0, _ := ret[0].(domain.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateSession indicates an expected call of ValidateSession.
func (mr *MockSessionServiceMockRecorder) ValidateSession(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateSession", reflect.TypeOf((*MockSessionService)(nil).ValidateSession), key)
}
