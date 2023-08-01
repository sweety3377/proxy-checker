package model

type CheckedProxyResponse []struct {
	Addr string `json:"addr,omitempty"`
}

type Response struct {
	CountryCode string `json:"countryCode,omitempty"`
}
