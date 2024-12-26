package github

type Contributor struct {
	Username string `json:"username"`
	Href     string `json:"href"`
	Avatar   string `json:"avatar"`
}

type Repository struct {
	Author             string        `json:"author"`
	Name               string        `json:"name"`
	Avatar             string        `json:"avatar"`
	URL                string        `json:"url"`
	Description        string        `json:"description"`
	Language           string        `json:"language"`
	LanguageColor      string        `json:"languageColor"`
	Stars              int           `json:"stars"`
	Forks              int           `json:"forks"`
	CurrentPeriodStars int           `json:"currentPeriodStars"`
	BuiltBy            []Contributor `json:"builtBy"`
}
