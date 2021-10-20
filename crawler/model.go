package crawler

// SearchResult modela os dados de retorno
type SearchResult struct {
	ResultRank  int    `json:"rank,omitempty"`
	ResultURL   string `json:"url,omitempty"`
	ResultTitle string `json:"title,omitempty"`
	ResultDesc  string `json:"description,omitempty"`
}
