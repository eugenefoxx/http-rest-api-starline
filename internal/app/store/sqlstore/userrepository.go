package sqlstore

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/eugenefoxx/http-rest-api-starline/internal/app/model"
	"github.com/eugenefoxx/http-rest-api-starline/internal/app/store"
)

// UserRepository ...
type UserRepository struct {
	store *Store
}

// Create ...
func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO users (email, encrypted_password, firstname, lastname) VALUES ($1, $2, $3, $4) RETURNING id",
		u.Email,
		u.EncryptedPassword,
		u.FirstName,
		u.LastName,
	).Scan(&u.ID)
}

// Find ...
func (r *UserRepository) Find(id int) (*model.User, error) {

	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password, firstname, lastname, Coalesce (role, ''), Coalesce (groups, ''), Coalesce (tabel, '') FROM users WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
		&u.FirstName,
		&u.LastName,
		&u.Role,
		&u.Groups,
		&u.Tabel,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

// FindByEmail ...
func (r *UserRepository) FindByEmail(email string, tabel string) (*model.User, error) {
	fmt.Println("nikname - ", email, tabel)
	u := &model.User{}
	err := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password, firstname, lastname, Coalesce(role, ''), Coalesce (groups, ''), Coalesce (tabel, '') FROM users WHERE email = $1 OR tabel = $2",
		email, tabel,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
		&u.FirstName,
		&u.LastName,
		&u.Role,
		&u.Groups,
		&u.Tabel,
	) /*; err != nil {
		fmt.Println("Ошибка на 76 - ")
		if err != sql.ErrNoRows {
			fmt.Println("Ошибка на 77 - ")
			return nil, store.ErrRecordNotFound
		}
		fmt.Println("Ошибка на 80 - ")
		return nil, err
	}*/
	if err != nil {
		if err == sql.ErrNoRows {
			// there were no rows, but otherwise no error occurred
		} else {
			log.Fatal(err)
		}
	}
	//	email2 := u.Email
	//	tabel2 := u.Tabel
	//	fmt.Println("FindByEmail email - ", email2)
	//	fmt.Println("FindByEmail tabel - ", tabel2)
	return u, nil
}

func (r *UserRepository) UpdatePass(s *model.User) error {
	if err := s.Validate(); err != nil {
		return err
	}

	if err := s.BeforCreate(); err != nil {
		return err
	}
	_, err := r.store.db.Exec(
		"UPDATE users SET encrypted_password = $1 WHERE email = $2",
		s.EncryptedPassword,
		s.Email,
	//	s.Tabel,
	)
	if err != nil {
		panic(err)
	}
	return nil
}
