package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	fmt.Println(":: Listening on port http://localhost:8090")
	http.HandleFunc("/", usage)
	http.HandleFunc("/area", areaHandler)
	http.HandleFunc("/url", urlHandler)
	http.HandleFunc("/bulk", bulkHandler)

	http.ListenAndServe(":8090", nil)
}

func usage(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(`
1. Usage
1.1. AreaName

Requests at

/area?name=pradesh-1/district-jhapa

for more cities or general usecase see url method

1.2. URL

Requests at

/url?url=https://election.ekantipur.com/pradesh-1/district-jhapa?lng=eng

where url must be valid kantipur url in format similar to url in above example.

1.3 Bulk List
Requests at

/bulk?list=pradesh-1/district-jhapa,pradesh-3/district-kathmandu

Where list= must be valid AreaName sepearated by commas
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

type result struct {
	district string
	data     map[string][]map[string]string
}

func bulkHandler(w http.ResponseWriter, r *http.Request) {
	areaNameSlug := r.URL.Query().Get("list")
	areaNames := strings.Split(areaNameSlug, ",")
	results := map[string]map[string][]map[string]string{}
	resultChannel := make(chan result)

	for _, areaName := range areaNames {
		url := urlFromAreaName(areaName)
		district := getDistrictName(url)
		go func(url, district string) {
			resultChannel <- result{district: district, data: fetchArea(url)}
		}(url, district)
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
