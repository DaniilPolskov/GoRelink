package types

type ShortenRequest struct {
	URL string `json:"url"`
}

type BatchShortenRequest struct {
	URL   string `json:"url"`
	Count int    `json:"count"`
}
