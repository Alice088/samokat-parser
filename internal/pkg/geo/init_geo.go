package geo

import (
	"alice088/sparser/internal/pkg/dto"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

func Init() ([]*dto.GEO, error) {
	var geos []*dto.GEO

	jsonFile, err := os.Open("./configs/geo_conf.json")
	if err != nil {
		return nil, err
	}

	jsonByteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(jsonByteValue, &geos)
	if err != nil {
		return nil, err
	}

	if len(geos) == 0 {
		return nil, errors.New("no geo objects")
	}

	return geos, nil
}
