package structs

type Translate struct {
	Text string `json:"text"`
	To   string `json:"to"`
}
type TranslationRes struct {
	Translation []Translate `json:"translations"`
}
