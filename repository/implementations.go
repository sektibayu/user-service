package repository

import (
	"context"
)

func (r *Repository) CreateNewUser(ctx context.Context, input User) (User, error) {
	var u User
	row := r.Db.QueryRowContext(ctx, 
			"INSERT INTO users (full_name, phone_number, hash_password, salt) values ($1, $2, $3, $4) RETURNING id, full_name, phone_number", 
			input.FullName, input.PhoneNumber, input.HashPassword, input.Salt)
	err := row.Scan(&u.Id, &u.FullName, &u.PhoneNumber)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func (r *Repository) GetUserById(ctx context.Context, id string) (User, error) {
	var u User
	row := r.Db.QueryRowContext(ctx, "SELECT id, full_name, phone_number, hash_password, salt FROM users WHERE id = $1", id)
	err := row.Scan(&u.Id, &u.FullName, &u.PhoneNumber, &u.HashPassword, &u.Salt)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (r *Repository) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (User, error) {
	var u User
	row := r.Db.QueryRowContext(ctx, "SELECT id, full_name, phone_number, hash_password, salt FROM users WHERE phone_number = $1", phoneNumber)
	err := row.Scan(&u.Id, &u.FullName, &u.PhoneNumber, &u.HashPassword, &u.Salt)
	if err != nil {
		return User{}, err
	}
	return u, nil
}



func (r *Repository) UpdateUserById(ctx context.Context, input User) (User, error) {
	var u User
	row := r.Db.QueryRowContext(ctx, "UPDATE users SET full_name = $1, phone_number = $2 WHERE id = $3 RETURNING id, full_name, phone_number", input.FullName, input.PhoneNumber, input.Id)
	err := row.Scan(&u.Id, &u.FullName, &u.PhoneNumber)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (r *Repository) UpdateSuccessLoginCountById(ctx context.Context, id string) (error) {
	_, err := r.Db.ExecContext(ctx, "UPDATE users SET login_success_count = login_success_count + 1 WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
