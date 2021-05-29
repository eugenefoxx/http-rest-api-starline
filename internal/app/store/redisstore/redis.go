package redisstore

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/eugenefoxx/http-rest-api-starline/internal/app/model"
)

// GetListShowDataByEO ...
func (c *Client) GetListShowDataByEO(eo string) (*model.Inspections, error) {
	getshowDataByDate := model.Inspection{}
	//	getshowDataByDateList := make(model.Inspections, 0)

	cmd := c.client.Get(eo)

	cmdb, err := cmd.Bytes()
	if err != nil {
		return getshowDataByDate, err
	}

	b := bytes.NewReader(cmdb)

	var res getshowDataByDate

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return getshowDataByDate, err
	}

	return res, nil

}

func (c *Client) SetListShowDataByEO(n *model.Inspections) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(n); err != nil {
		return err
	}

	return c.client.Set(n.ID, b.Bytes(), 25*time.Second).Err()
}
