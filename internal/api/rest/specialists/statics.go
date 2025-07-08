package specialists

const (
	Type = "type"
	City = "City"
)

type AreaType string

const (
	Makeup AreaType = "makeup"
	Sport  AreaType = "sport"
)

func StringToAreaType(s string) (AreaType, bool) {
	switch {
	case s == string(Makeup):
		return Makeup, true
	case s == string(Sport):
		return Sport, true
	default:
		return "", false
	}
}
