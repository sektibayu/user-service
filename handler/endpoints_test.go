package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRegistration(t *testing.T) {
	// setup mock
	ctrl := gomock.NewController(t)
	repo := repository.NewMockRepositoryInterface(ctrl)
	h := handler.NewServer(handler.NewServerOptions{
		Repository: repo,
	})

	var testcases = []struct {
		name        string
		params 		*generated.RegisterRequest
		mockFunc  func (ctx context.Context)
		expect      interface{}
	}{
		{
			name: "when no params",
			params: &generated.RegisterRequest{
				FullName: "",
				PhoneNumber: "",
				Password: "",
			},
			mockFunc: func(ctx context.Context) {},
			expect: &generated.ErrorResponse{
				Message: "Bad Request",
				Errors: []generated.ErrorField{
					{ Field: "phone_number", Message: "phone_number must be at minimum 10 char and maximum 13 char" },
					{ Field: "full_name", Message: "full_name must be at minimum 3 char and maximum 60 char" },
					{ Field: "password", Message: "password must be at minimum 6 char and maximum 64 char" },
				},
			},
		},
		{
			name: "when phone number not using indonesia code ",
			params: &generated.RegisterRequest{
				FullName: "bayu sektiaji",
				PhoneNumber: "+61899999999",
				Password: "Kupu=Kupu4",
			},
			mockFunc: func(ctx context.Context) {},
			expect: &generated.ErrorResponse{
				Message: "Bad Request",
				Errors: []generated.ErrorField{
					{ Field: "phone_number", Message: "phone_number must start indonesia country code: +62" },
				},
			},
		},
		{
			name: "when success",
			params: &generated.RegisterRequest{
				FullName: "bayu sektiaji",
				PhoneNumber: "+62899999999",
				Password: "Kupu=Kupu4",
			},
			mockFunc: func(ctx context.Context) {
				repo.EXPECT().CreateNewUser(gomock.Any(), gomock.Any()).Return(repository.User{
					Id: "asdf-adsf-adf-afd",
					FullName: "bayu sektiaji",
					PhoneNumber: "+62899999999",
					HashPassword: "lalalili",
					Salt: "lalalili",
				}, nil)
			},
			expect: &generated.RegisterResponse{
				Id: "asdf-adsf-adf-afd",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.params)
			req := httptest.NewRequest(http.MethodPost, "/registration", strings.NewReader(string(body)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := echo.New().NewContext(req, rec)
			tc.mockFunc(ctx.Request().Context())
			h.Registration(ctx)
			var result interface{}
			_, isError := tc.expect.(*generated.ErrorResponse)
			if isError {
				result = &generated.ErrorResponse{}
			} else {
				result = &generated.RegisterResponse{}
			}
			json.Unmarshal(rec.Body.Bytes(), result)
			rec.Body.Reset()

			assert.Equal(t, tc.expect, result)
		})
	}
}
