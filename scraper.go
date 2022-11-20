package main

import (
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func fetchArea(url string) map[string][]map[string]string {
	results := map[string][]map[string]string{}
	lastLoggedConstituency := 0
	c := colly.NewCollector()
	c.OnHTML(".col-md-6", func(e *colly.HTMLElement) {
		constituencyName := e.ChildText("h3")

		constituencyNo := validateConstituencyNo(constituencyName, lastLoggedConstituency)
		if constituencyNo == -1 {
			return
		} else {
			lastLoggedConstituency = constituencyNo
		}
		constituencyName = strings.ToLower(constituencyName)

		nomineeData := []map[string]string{}
		e.ForEach("div.candidate-meta", func(i int, el *colly.HTMLElement) {
			if i%2 != 0 {
				return
			}
			nomineeName := el.ChildText("div.nominee-name")
			nomineeParty := el.ChildText("div.candidate-party-name")
			nomineeData = append(nomineeData, map[string]string{
				"name":  nomineeName,
				"party": nomineeParty,
				"votes": "0",
			})
		})
		results[constituencyName] = nomineeData
	})
	c.Visit(url)
	return results
}

func validateConstituencyNo(constituencyName string, lastLogged int) int {
	if strings.TrimSpace(constituencyName) == "" {
		return -1
	}
	constituencyNoStr := strings.TrimPrefix(constituencyName, "Constituency :")
	constituencyNoStr = strings.TrimSpace(constituencyNoStr)
	constituencyNo, err := strconv.Atoi(constituencyNoStr)
	if err != nil {
		return -1
	}
	if constituencyNo <= lastLogged {
		return -1
	} else {
		return constituencyNo
	}
}