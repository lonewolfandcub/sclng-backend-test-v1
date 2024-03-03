package github

type (
	License struct {
		Name string `json:"name"`
	}

	Language struct {
		Name  string `json:"name"`
		Bytes uint   `json:"bytes"`
	}

	Repository struct {
		ID        uint    `json:"id"`
		FullName  string  `json:"full_name"`
		Language  string  `json:"language"`
		License   License `json:"license"`
		URL       string  `json:"url"`
		CreatedAt string  `json:"created_at"`
	}

	StatsRepository struct {
		ID        uint       `json:"id"`
		FullName  string     `json:"full_name"`
		Languages []Language `json:"languages"`
		License   License    `json:"license"`
		URL       string     `json:"url"`
	}

	GithubRepositoriesReponse struct {
		Items []Repository `json:"items"`
	}
)
