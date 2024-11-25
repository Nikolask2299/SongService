package models

// Song model info
// @Description Song information about the account
type Song struct {
	ID string `json:"id"`
	Group string `json:"group"` 
	Song string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Text string `json:"text"`
	Link string `json:"link"`
}

// Filter song model info
// @Description Filter song model info
type FilterSong struct {
	Group string `json:"group"` 
	Song string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Text string `json:"text"`
	Link string `json:"link"` 
}

// New song model info
// @Description Song information about user
type NewSong struct {
	Group string `json:"group"`
	Song string `json:"song"`
}

