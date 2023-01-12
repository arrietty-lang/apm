package apm

import "encoding/json"

type Pkg struct {
	Deps []*Dependencies `json:"deps,omitempty"`
}

type Dependencies struct {
	Url     string `json:"url"`
	Version string `json:"version"`
}

func UnmarshalPkgJson(bytes []byte) (*Pkg, error) {
	var pkgJson Pkg
	err := json.Unmarshal(bytes, &pkgJson)
	if err != nil {
		return nil, err
	}
	return &pkgJson, nil
}
