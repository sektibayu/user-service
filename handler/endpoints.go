package handler

import (
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

type RegistrationValidate struct {
	FullName    string	`validate:"required"`
	Password 	string  `validate:"required"`
	PhoneNumber	string	`validate:"required"`
}

// (POST /registration)
func (s *Server) Registration(ctx echo.Context) error {
	var req generated.RegisterRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Bad Request: "+ err.Error(),
		})
	}

	var errors []generated.ErrorField
	errors = append(errors, validatePhoneNumber(req.PhoneNumber)...)
	errors = append(errors, validateFullname(req.FullName)...)
	errors = append(errors, validatePassword(req.Password)...)

	if len(errors) > 0 {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Bad Request",
			Errors: errors,
		})
	}

	salt := generateRandomSalt(SALT_SIZE)
	hashPass := hashPassword(req.Password, salt)
	user := repository.User{
		FullName: req.FullName,
		PhoneNumber: req.PhoneNumber,
		HashPassword: hashPass,
		Salt: string(salt),
	}

	u, err := s.Repository.CreateNewUser(ctx.Request().Context(), user)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
					Message: "Bad Request",
					Errors: []generated.ErrorField{
						{
							Field: "phone_number",
							Message: "phone_numer already exists",
						},
					},
				})
			}
		}
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "Internal Server error: "+ err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, generated.RegisterResponse{
		Id: u.Id,
	})
}

// (POST /login)
func (s *Server) Login(ctx echo.Context) error {
	var req generated.LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Bad Request: "+ err.Error(),
		})
	}

	u, err := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), req.PhoneNumber)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Bad Request",
		})
	}
	if !doPasswordsMatch(u.HashPassword, req.Password, []byte(u.Salt)) {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Bad Request",
		})
	}
	// update last login
	go s.Repository.UpdateSuccessLoginCountById(ctx.Request().Context(), u.Id)

	// generate token
	token, err := generateToken(u.Id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, generated.LoginResponse{
		Id: u.Id,
		Token: token,
	})
}

// (GET /my-profile)
func (s *Server) GetMyProfile(ctx echo.Context, params generated.GetMyProfileParams) error {
	token := ctx.Request().Header.Get("x-auth-token")
	id, err := verifyToken(token)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{
			Message: "Forbidden",
		})
	}

	u, err := s.Repository.GetUserById(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "Internal Server error: "+ err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, generated.ProfileResponse{
		FullName: u.FullName,
		PhoneNumber: u.PhoneNumber,
	})
}

// (PATCH /my-profile)
func (s *Server) UpdateMyProfile(ctx echo.Context, params generated.UpdateMyProfileParams) error {
	token := ctx.Request().Header.Get("x-auth-token")
	id, err := verifyToken(token)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{
			Message: "Forbidden",
		})
	}

	var req generated.ProfileRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Bad Request: "+ err.Error(),
		})
	}

	u, err := s.Repository.UpdateUserById(ctx.Request().Context(), repository.User{
		Id: id,
		FullName: req.FullName,
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return ctx.JSON(http.StatusConflict, generated.ErrorResponse{
					Message: "conflict",
					Errors: []generated.ErrorField{
						{
							Field: "phone_number",
							Message: "phone_numer already exists",
						},
					},
				})
			}
		}
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "Internal Server error: "+ err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, generated.ProfileResponse{
		FullName: u.FullName,
		PhoneNumber: u.PhoneNumber,
	})
}
