package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/", usage)
	http.HandleFunc("/area", areaHandler)
	http.HandleFunc("/url", urlHandler)
	http.HandleFunc("/bulk", bulkHandler)
	http.HandleFunc("/summary", summaryHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	fmt.Println(":: Listening on port http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}

func usage(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(`
1. Usage

1.1. AreaName

Requests at

/area?name=pradesh-1/district-jhapa

where name is valid kantipur url part representing an electoral area.
This is supposed to be extracted from a kantipur url.

Example: https://electionapi.osac.org.np/area?name=pradesh-1/district-jhapa

1.2. URL

Requests at

/url?url=https://election.ekantipur.com/pradesh-1/district-jhapa?lng=eng

where url must be valid kantipur url in format similar to url in above example.

Example: https://electionapi.osac.org.np/url?url=https://election.ekantipur.com/pradesh-1/district-jhapa?lng=eng

1.3 Bulk List
Requests at

/bulk?list=pradesh-1/district-jhapa,pradesh-3/district-kathmandu

Where list must be list of valid AreaNames sepearated by commas.

Example: https://electionapi.osac.org.np/bulk?list=pradesh-1/district-jhapa,pradesh-3/district-kathmandu

1.4. Summary

Requests at

/summary

Gives all party names, their wins and leads count in Federal and provincial category.

Example: https://electionapi.osac.org.np/summary
`))

}

func areaHandler(w http.ResponseWriter, r *http.Request) {
	areaName := r.URL.Query().Get("name")
	url := fmt.Sprintf("https://election.ekantipur.com/%s?lng=eng", areaName)
	results := fetchArea(url)

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(results)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Errored out : %q", err)
	}
}

func urlHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	results := fetchArea(url)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(results)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Errored out : %q", err)
	}
}

func summaryHandler(w http.ResponseWriter, r *http.Request) {
	results := fetchSummary()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(results)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Errored out : %q", err)
	}
}

type result struct {
	district string
	data     map[string][]map[string]interface{}
}

func bulkHandler(w http.ResponseWriter, r *http.Request) {
	areaNameSlug := r.URL.Query().Get("list")
	areaNames := strings.Split(areaNameSlug, ",")
	results := map[string]map[string][]map[string]interface{}{}
	resultChannel := make(chan result)

	for _, areaName := range areaNames {
		url := urlFromAreaName(areaName)
		district := getDistrictName(url)
		go func(url, district string) {
			resultChannel <- result{district: district, data: fetchArea(url)}
		}(url, district)
		time.Sleep(2 * time.Second)
	}
	for i := 0; i < len(areaNames); i++ {
		r := <-resultChannel
		results[r.district] = r.data
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(results)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Errored out : %q", err)
	}
}

func getDistrictName(url string) string {
	areaName := strings.TrimSuffix(url, "?lng=eng")
	areaParts := strings.Split(areaName, "-")
	districtName := areaParts[len(areaParts)-1]
	return districtName
}

func urlFromAreaName(areaName string) string {
	url := fmt.Sprintf("https://election.ekantipur.com/%s?lng=eng", areaName)
	return url
}
