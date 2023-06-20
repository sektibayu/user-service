// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	CreateNewUser(ctx context.Context, input User) (User, error)
	GetUserById(ctx context.Context, id string) (User, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (User, error)
	UpdateUserById(ctx context.Context, input User) (User, error)
	UpdateSuccessLoginCountById(ctx context.Context, id string) (error)
}
