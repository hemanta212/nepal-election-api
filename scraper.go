package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func fetchArea(url string) map[string][]map[string]interface{} {
	results := map[string][]map[string]interface{}{}
	lastLoggedConstituency := 0

	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit 537.36(KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36"
	c.IgnoreRobotsTxt = true
	c.OnHTML(".col-md-6", func(e *colly.HTMLElement) {
		constituencyName := e.ChildText("h3")

		constituencyNo := validateConstituencyNo(constituencyName, lastLoggedConstituency)
		if constituencyNo == -1 {
			return
		} else {
			lastLoggedConstituency = constituencyNo
		}
		constituencyName = strings.ToLower(constituencyName)

		fmt.Println("Fetching nominee data...")
		results[constituencyName] = fetchNomineeData(e)
	})

	c.Visit(url)
	return results
}

func fetchSummary() map[string][]map[string]string {
	url := "https://election.ekantipur.com/?lng=eng"
	results := map[string][]map[string]string{}
	// lastLoggedConstituency := 0

	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit 537.36(KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36"
	c.IgnoreRobotsTxt = true
	c.OnHTML("div.parties", func(e *colly.HTMLElement) {
		e.ForEach("div.col-md-6", func(_ int, el *colly.HTMLElement) {
			levelName := strings.ToLower(el.ChildText("h2.title"))
			levelName = strings.TrimSpace(strings.Split(levelName, " ")[0])
			levelData := []map[string]string{}
			el.ForEach("div.row.gx-1", func(_ int, ele *colly.HTMLElement) {
				party := ele.ChildText("div.party-name")
				wins := ele.ChildText("div:nth-child(3)")
				leads := ele.ChildText("div:nth-child(4)")
				partyData := map[string]string{
					"name":  party,
					"wins":  wins,
					"leads": leads,
				}
				fmt.Println(e.DOM)
				levelData = append(levelData, partyData)
			})
			results[levelName] = levelData
		})
	})

	c.Visit(url)
	return results
}

func fetchNomineeData(e *colly.HTMLElement) []map[string]interface{} {
	fmt.Println("Starting fetching..")
	nomineeData := []map[string]interface{}{}

	card_text := e.ChildText("div.win-lead-text")
	has_won := strings.Contains(card_text, "won with")

	e.ForEach("div.candidate-wrappper", func(_ int, el *colly.HTMLElement) {
		nomineeName := el.ChildText("div.nominee-name")
		nomineeParty := el.ChildText("div.candidate-party-name")
		votes := el.ChildText("div.vote-count")
		fmt.Println(el.DOM, nomineeName, votes)
		if strings.TrimSpace(votes) == "" {
			votes = "0"
		}

		nomineeData = append(nomineeData, map[string]interface{}{
			"name":            nomineeName,
			"party":           nomineeParty,
			"votes":           votes,
			"winner_declared": has_won,
		})
	})

	return nomineeData
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
