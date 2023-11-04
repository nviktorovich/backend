package client

//go:generate mockgen -source=./scouter.go -destination=./testdata/scouter.go --package=testdata
type Scouter interface {
	GetAll(titles []string, in string) (map[string]float64, error)
	GetSpecial(title string, in string) (map[string]float64, error)
}
