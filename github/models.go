package github

type (
	License struct {
		Name string `json:"name"`
	}

	Repository struct {
		ID         uint    `json:"id"`
		FullName   string  `json:"full_name"`
		Language   string  `json:"language"`
		License    License `json:"license"`
		Repository string  `json:"repository"`
		URL        string  `json:"url"`
		CreatedAt  string  `json:"created_at"`
	}
	RepositoriesReponse struct {
		Items []Repository `json:"items"`
	}
)
