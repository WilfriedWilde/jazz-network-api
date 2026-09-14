package countries

import "github.com/biter777/countries"

type Country struct {
	Name   string `json:"name"`
	Alpha2 string `json:"code"`
}

// TODO: build countries.json for backend and frontend
func main() {
	all := countries.All()
}
