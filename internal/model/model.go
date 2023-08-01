package model

type CheckedProxyResponse []struct {
	Addr string `json:"addr,omitempty"`
}

type Response struct {
	RegionName string `json:"regionName,omitempty"`
}
