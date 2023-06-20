// This file contains types that are used in the repository layer.
package repository


type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type User struct {
	Id 				string `json:"id", db:"id"`
	FullName 		string `json:"full_name", db:"full_name"`
	PhoneNumber 	string `json:"phone_number", db:"phone_number"`
	HashPassword 	string `json:"hash_password", db:"hash_password"`
	Salt			string `json:"salt", db:"salt"`
}
