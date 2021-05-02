package redisstore

import "github.com/eugenefoxx/http-rest-api-starline/internal/app/model"

func (c *Client) GetListShowDataByEO (eo string) (s *model.Inspections, err) {
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