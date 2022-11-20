package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

func main() {
	fname := "data.json"
	file, err := os.Create(fname)
	if err != nil {
		log.Fatalf("Cannot create json file: %q", err)
		return
	}
	defer file.Close()

	writer := json.NewEncoder(file)

	results := []map[string]string{}

	c := colly.NewCollector()
	c.OnHTML("table#customers", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, el *colly.HTMLElement) {
			row := map[string]string{
				"candidate": el.ChildText("td:nth-child(1)"),
				"votes":     el.ChildText("td:nth-child(2)"),
				"party":     el.ChildText("td:nth-child(3)"),
			}
			if i != 0 {
				results = append(results, row)
			}
		})
	})
	c.Visit("https://www.w3schools.com/html/html_tables.asp")
	j, err := json.Marshal(results)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		fmt.Println(string(j))
	}
	writer.Encode(j)
}
