package geo

import "alice088/sparser/internal/pkg/dto"

func Init() (*dto.GEO, error) {
	geo := &dto.GEO{}

	if err := Update(geo, 0); err != nil {
		return nil, err
	}

	return geo, nil
}
