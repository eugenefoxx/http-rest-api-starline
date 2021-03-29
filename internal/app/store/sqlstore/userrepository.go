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
	tx, err := r.store.db.Begin()
	if err != nil {
		return err
	}
	if err := s.Validate(); err != nil {
		return err
	}

	if err := s.BeforCreate(); err != nil {
		return err
	}
	_, errTx := tx.Exec(
		"UPDATE users SET encrypted_password = $1 WHERE email = $2",
		s.EncryptedPassword,
		s.Email,
	//	s.Tabel,
	)
	if errTx != nil {
		tx.Rollback()
		panic(err)
	}
	return tx.Commit()
}

// SuperIngenerQuality
func (r *UserRepository) CreateUserByManager(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO users (email, encrypted_password, firstname, lastname, role, groups, tabel) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		u.Email,
		u.EncryptedPassword,
		u.FirstName,
		u.LastName,
		u.Role,
		u.Groups,
		u.Tabel,
	).Scan(&u.ID)
}

func (r *UserRepository) ListUsersQuality() (u *model.Users, err error) {
	showUsersQuality := model.User{}
	showUsersQualityList := make(model.Users, 0)

	selectUsers := `SELECT id, email, firstname, lastname, 
	role, groups, tabel FROM users WHERE groups = 'качество';`

	rows, err := r.store.db.Query(
		selectUsers,
	)

	if err != nil {
		fmt.Println(err, "error in func ListUsersQuality()")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&showUsersQuality.ID,
			&showUsersQuality.Email,
			&showUsersQuality.FirstName,
			&showUsersQuality.LastName,
			&showUsersQuality.Role,
			&showUsersQuality.Groups,
			&showUsersQuality.Tabel,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		showUsersQualityList = append(showUsersQualityList, showUsersQuality)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &showUsersQualityList, nil
}

func (r *UserRepository) ListUsersWarehouse() (u *model.Users, err error) {
	showUsersWarehouse := model.User{}
	showUsersWarehouseList := make(model.Users, 0)

	selectUsers := `SELECT id, email, firstname, lastname, 
	role, groups, tabel FROM users WHERE groups = 'склад';`

	rows, err := r.store.db.Query(
		selectUsers,
	)

	if err != nil {
		fmt.Println(err, "error in func ListUsersWarehouse()")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&showUsersWarehouse.ID,
			&showUsersWarehouse.Email,
			&showUsersWarehouse.FirstName,
			&showUsersWarehouse.LastName,
			&showUsersWarehouse.Role,
			&showUsersWarehouse.Groups,
			&showUsersWarehouse.Tabel,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		showUsersWarehouseList = append(showUsersWarehouseList, showUsersWarehouse)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &showUsersWarehouseList, nil
}

func (r *UserRepository) EditUserByManager(id int) (*model.User, error) {

	u := &model.User{}
	//	selectEdit := `SELECT email, firstname, lastname,
	//	role, tabel from users where id = $1;`
	if err := r.store.db.QueryRow(
		//		selectEdit,
		"SELECT id, email, firstname, lastname, role, tabel FROM users WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		&u.Email,
		&u.FirstName,
		&u.LastName,
		&u.Role,
		&u.Tabel,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) UpdateUserByManager(s *model.User) error {

	updateUser := `UPDATE users SET email = $1, firstname = $2,
	lastname = $3, role = $4, tabel = $5 WHERE id = $6;`
	_, err := r.store.db.Exec(
		updateUser,
		s.Email,
		s.FirstName,
		s.LastName,
		s.Role,
		s.Tabel,
		s.ID,
	)

	if err != nil {
		panic(err)
	}

	return nil
}

func (r *UserRepository) DeleteUserByManager(s *model.User) error {

	deleteUser := `DELETE FROM users WHERE id = $1;`
	_, err := r.store.db.Exec(
		deleteUser,
		s.ID,
	)
	if err != nil {
		panic(err)
	}

	return nil
}
