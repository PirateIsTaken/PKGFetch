package main

type Repository struct {
	Name        string
	Description string

	Stars int
	Forks int

	IsFork      bool
	IsArchived  bool
	HasReleases bool

	Score int
}

func SearchGithub(searchName string) []Repository {
	// @TEMP
	// FAKE REPOS
	return []Repository{
		{
			Name:        "JJJ/WorkingPackage",
			Description: "This is a test package 1",
			Stars:       696969,
			Forks:       100,
			IsFork:      false,
			IsArchived:  false,
			HasReleases: true,
		},
		{
			Name:        "HHH/WorkingPackage",
			Description: "This is a test package 1",
			Stars:       69697000,
			Forks:       100,
			IsFork:      false,
			IsArchived:  false,
			HasReleases: false,
		},
		{
			Name:        "Another/Package",
			Description: "Just Another Package",
			Stars:       1000,
			Forks:       10000,
			IsFork:      false,
			IsArchived:  false,
			HasReleases: true,
		},
		{
			Name:        "OneMore/Thing",
			Description: "Thing One more",
			Stars:       991928,
			Forks:       199991,
			IsFork:      true,
			IsArchived:  false,
			HasReleases: true,
		},
		{
			Name:        "Archi/One",
			Description: "The Archi",
			Stars:       991920008,
			Forks:       19999001,
			IsFork:      false,
			IsArchived:  true,
			HasReleases: true,
		},
	}
}

// Mechanics
func CalculateScore(repo *Repository) {
	score := 0

	score += repo.Stars / 100
	score += repo.Forks

	if !repo.IsArchived {
		score += 500
	}
	if !repo.IsFork {
		score += 500
	}
	if repo.HasReleases {
		score += 1000
	}

	repo.Score = score
}
