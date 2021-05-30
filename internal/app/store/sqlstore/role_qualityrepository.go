package sqlstore

import (
	"database/sql"

	"github.com/eugenefoxx/http-rest-api-starline/internal/app/model"
	"github.com/eugenefoxx/http-rest-api-starline/internal/app/store"
)

type RoleQualityRepository struct {
	store *Store
}

func (r *RoleQualityRepository) ListRoleQuality() (*model.RoleQualitys, error) {
	roleQuality := model.Role_Quality{}
	roleQualityList := make(model.RoleQualitys, 0)

	rows, err := r.store.db.Query(
		"SELECT role FROM role_quality",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&roleQuality.Role,
		)
		for err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		roleQualityList = append(roleQualityList, roleQuality)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &roleQualityList, nil
}
