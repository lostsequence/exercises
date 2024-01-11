package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"movies-auth/users/internal/domain"
	mock_api "movies-auth/users/internal/tests/api_mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/mock/gomock"
)

// не хватает обработки ошибок от SessionService
func TestCreate(t *testing.T) {
	type fields struct {
		login    string
		password string
	}

	testCases := []struct {
		name                string
		fields              fields
		mockUserServiceInit func(s *mock_api.MockUsersService)
		header              http.Header
		wantStatusCode      int
		wantErrMessage      string
		wantError           bool
	}{
		{
			name:                "fail_content_type",
			fields:              fields{},
			mockUserServiceInit: func(s *mock_api.MockUsersService) {},
			header: http.Header{
				"Content-Type": []string{
					"text/plain",
				},
			},
			wantStatusCode: http.StatusUnsupportedMediaType,
			wantErrMessage: "not supported content type\n",
			wantError:      true,
		},
		{
			name:   "fail_conflict",
			fields: fields{},
			mockUserServiceInit: func(s *mock_api.MockUsersService) {
				s.EXPECT().Create(gomock.Any(), gomock.Any()).Return(domain.User{}, domain.ErrConflict)
			},
			header: http.Header{
				"Content-Type": []string{
					"application/json",
				},
			},
			wantStatusCode: http.StatusConflict,
			wantErrMessage: "user already exist\n",
			wantError:      true,
		},
		{
			name:   "fail_conflict",
			fields: fields{},
			mockUserServiceInit: func(s *mock_api.MockUsersService) {
				s.EXPECT().Create(gomock.Any(), gomock.Any()).Return(domain.User{}, domain.ErrConflict)
			},
			header: http.Header{
				"Content-Type": []string{
					"application/json",
				},
			},
			wantStatusCode: http.StatusConflict,
			wantErrMessage: "user already exist\n",
			wantError:      true,
		},
		{
			name:   "fail_internal_error",
			fields: fields{},
			mockUserServiceInit: func(s *mock_api.MockUsersService) {
				s.EXPECT().Create(gomock.Any(), gomock.Any()).Return(domain.User{}, errors.New("unexpected error"))
			},
			header: http.Header{
				"Content-Type": []string{
					"application/json",
				},
			},
			wantStatusCode: http.StatusInternalServerError,
			wantErrMessage: "unexpected error\n",
			wantError:      true,
		},
		{
			name: "success_user_created",
			fields: fields{
				login:    "user1",
				password: "12345678",
			},
			mockUserServiceInit: func(s *mock_api.MockUsersService) {
				s.EXPECT().Create(gomock.Any(), gomock.Any()).Return(domain.User{
					Login:    "user1",
					Password: "12345678",
				}, nil)
			},
			header: http.Header{
				"Content-Type": []string{
					"application/json",
				},
			},
			wantStatusCode: http.StatusCreated,
			wantErrMessage: "unexpected error\n",
			wantError:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			su := mock_api.NewMockUsersService(ctrl)
			if tc.mockUserServiceInit != nil {
				tc.mockUserServiceInit(su)
			}

			ss := mock_api.NewMockSessionService(ctrl)

			h := NewUsersHandler(su, ss)
			payload := domain.User{
				Login:    tc.fields.login,
				Password: tc.fields.password,
			}

			body, _ := json.Marshal(payload)
			req := httptest.NewRequest(http.MethodPost, "/api/", bytes.NewReader(body))
			req.Header = tc.header
			recorder := httptest.NewRecorder()

			h.Register(recorder, req)
			if recorder.Result().StatusCode != tc.wantStatusCode {
				t.Errorf("expected status code: %d, got: %d", tc.wantStatusCode, recorder.Result().StatusCode)
			}

			if tc.wantError {
				respBody := recorder.Body.Bytes()
				msg := string(respBody)

				if msg != tc.wantErrMessage {
					t.Errorf("expected error message: %s, got: %s", tc.wantErrMessage, msg)
				}
				return
			}

			var user domain.User
			_ = json.NewDecoder(recorder.Result().Body).Decode(&user)

			if tc.fields.login != user.Login && tc.fields.password != user.Password {
				t.Errorf("expected user: %v, got: %v", tc.fields, user)
			}
		})
	}
}
