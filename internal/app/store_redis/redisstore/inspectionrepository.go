package redisstore

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/eugenefoxx/http-rest-api-starline/internal/app/model"
)

type InspectionRepository struct {
	redis *Redis
}

// GetListShowDataByEO ...
func (r *InspectionRepository) GetListShowDataByEO(ctx context.Context, eo string) (*model.Inspections, error) {
	//getshowDataByDate := model.Inspection{}
	getshowDataByDateList := make(model.Inspections, 0)

	//cmd := r.client.Get(eo)
	cmd := r.redis.client.Get(ctx, "eos")

	cmdb, err := cmd.Bytes()
	if err != nil {
		return &getshowDataByDateList, err
	}

	b := bytes.NewReader(cmdb)

	//var res &getshowDataByDateList

	if err := gob.NewDecoder(b).Decode(&getshowDataByDateList); err != nil {
		return &getshowDataByDateList, err
	}

	return &getshowDataByDateList, nil

}

func (r *InspectionRepository) SetListShowDataByEO(ctx context.Context, n *model.Inspections) error {
	var b bytes.Buffer

	//strN := strconv.Itoa(n.IdRoll)
	//n = make(model.Inspections, 0)
	//var res *model.Inspections
	if err := gob.NewEncoder(&b).Encode(n); err != nil {
		return err
	}

	return r.redis.client.Set(ctx, "eos", b.Bytes(), 25*time.Second).Err()
}
