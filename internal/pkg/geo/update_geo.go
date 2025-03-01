package geo

import (
	"alice088/sparser/internal/pkg/dto"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

func Update(geo *dto.GEO, pos int) error {
	var geos []*dto.GEO

	jsonFile, err := os.Open("./configs/geo_conf.json")
	if err != nil {
		return err
	}

	jsonByteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(jsonByteValue, &geos)
	if err != nil {
		return err
	}

	if len(geos) == 0 {
		return errors.New("no geo objects")
	}

	if len(geos) < pos {
		return errors.New("pos out of range")
	}

	*geo = *geos[pos]

	return nil
}
