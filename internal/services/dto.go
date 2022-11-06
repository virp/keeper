package services

type Metadata struct {
	Key   string
	Value string
}

type Item struct {
	Name     string
	Type     string
	Data     []byte
	Metadata []Metadata
}
