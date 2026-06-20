package main

import "strings"

func SearchForPkg(pkgName string, repos []Repository) []Repository {
	var searchResult []Repository
	pkgName = strings.ToLower(pkgName)

	for _, repo := range repos {
		if strings.Contains(
			strings.ToLower(repo.Name),
			pkgName,
		) {
			searchResult = append(searchResult, repo)
		}
	}

	return searchResult
}
