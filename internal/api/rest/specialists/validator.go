package specialists

import (
	"errors"
	"net/url"
)

func validateGetSpecialistsQueries(queries url.Values) (AreaType, string, error) {
	areaType := queries.Get(Type)
	city := queries.Get(City)

	area, ok := StringToAreaType(areaType)
	if !ok {
		return "", "", errors.New("invalid specialist's area")
	}

	return area, city, nil
}
