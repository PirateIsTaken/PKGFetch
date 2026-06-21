package main

type Repository struct {
	Name        string
	Description string
	Link        string

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
			Name:        "AppFlowy-IO/AppFlowy",
			Description: "Open source alternative to Notion.",
			Link:        "https://github.com/AppFlowy-IO/AppFlowy",

			Stars:       65000,
			Forks:       4200,
			IsFork:      false,
			IsArchived:  false,
			HasReleases: true,
		},
		{
			Name:        "zed-industries/zed",
			Description: "High-performance multiplayer code editor.",
			Link:        "https://github.com/zed-industries/zed",

			Stars:       62000,
			Forks:       1800,
			IsFork:      false,
			IsArchived:  false,
			HasReleases: true,
		},
		{
			Name:        "godotengine/godot",
			Description: "Cross-platform open source game engine.",
			Link:        "https://github.com/godotengine/godot",

			Stars:       102000,
			Forks:       23000,
			IsFork:      false,
			IsArchived:  false,
			HasReleases: true,
		},
		{
			Name:        "randomdev/AppFlowy",
			Description: "Personal fork of AppFlowy with custom changes.",
			Link:        "https://github.com/randomdev/AppFlowy",

			Stars:       12,
			Forks:       1,
			IsFork:      true,
			IsArchived:  true,
			HasReleases: false,
		},
		{
			Name:        "anotherRandom/AppFlowy",
			Description: "Personal fork of AppFlowy with custom changes.",
			Link:        "https://github.com/randomdev/AppFlowy",

			Stars:       10,
			Forks:       1,
			IsFork:      true,
			IsArchived:  true,
			HasReleases: true,
		},
		{
			Name:        "coolguy/zed-fork",
			Description: "Experimental fork of Zed.",
			Link:        "https://github.com/coolguy/zed-fork",

			Stars:       130,
			Forks:       22,
			IsFork:      true,
			IsArchived:  false,
			HasReleases: true,
		},
		{
			Name:        "oldtools/LegacyEditor",
			Description: "A discontinued text editor.",
			Link:        "https://github.com/oldtools/LegacyEditor",

			Stars:       8200,
			Forks:       600,
			IsFork:      false,
			IsArchived:  true,
			HasReleases: true,
		},
		{
			Name:        "abandoned/OldLauncher",
			Description: "Old game launcher no longer maintained.",
			Link:        "https://github.com/abandoned/OldLauncher",

			Stars:       2400,
			Forks:       120,
			IsFork:      false,
			IsArchived:  true,
			HasReleases: true,
		},
		{
			Name:        "somebody/NewCoolApp",
			Description: "Promising new application.",
			Link:        "https://github.com/somebody/NewCoolApp",

			Stars:       350,
			Forks:       17,
			IsFork:      false,
			IsArchived:  true,
			HasReleases: false,
		},
		{
			Name:        "docs/AppFlowy-Docs",
			Description: "Documentation repository for AppFlowy.",
			Link:        "https://github.com/docs/AppFlowy-Docs",

			Stars:       500,
			Forks:       90,
			IsFork:      false,
			IsArchived:  false,
			HasReleases: false,
		},
		{
			Name:        "example/TestProject",
			Description: "Test project with no releases.",
			Link:        "https://github.com/example/TestProject",

			Stars:       15000,
			Forks:       500,
			IsFork:      false,
			IsArchived:  false,
			HasReleases: false,
		},
		{
			Name:        "HelixEditor/helix",
			Description: "A post-modern text editor.",
			Link:        "https://github.com/HelixEditor/helix",

			Stars:       42000,
			Forks:       3200,
			IsFork:      false,
			IsArchived:  false,
			HasReleases: true,
		},
		{
			Name:        "obsidianmd/obsidian-releases",
			Description: "Release repository for Obsidian.",
			Link:        "https://github.com/obsidianmd/obsidian-releases",

			Stars:       18000,
			Forks:       900,
			IsFork:      false,
			IsArchived:  false,
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
