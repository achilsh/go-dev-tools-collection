package balanceupstreams

type HttpElem struct {
	Targets string `json:"Target"`
	Weight  int    `json:"Weight"`
}
type UpstreamsBizElem struct {
	Name string
	Http []HttpElem
}
type UpstreamsBizConfig struct {
	UpstreamsBiz []UpstreamsBizElem
}
