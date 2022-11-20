# Nepal Election API 2079

Developed using the fantastic [gocolly](https://github.com/gocolly/colly) scraping library, data source is [ekantipur](https://election.ekantipur.com)

This repo was developed as backend for election-bot for reddit hosted at https://github.com/pykancha/reddit-bots

## Installation
- Install go
- Clone the repo
```
git clone https://github.com/hemanta212/nepal-election-api
```

## Usage
1. AreaName

Requests at

/area?name=pradesh-1/district-jhapa

for more cities or general usecase see url method

2. URL

Requests at

/url?url=https://election.ekantipur.com/pradesh-1/district-jhapa?lng=eng

where url must be valid kantipur url in format similar to url in above example.

3. Bulk List
Requests at

/bulk?list=pradesh-1/district-jhapa,pradesh-3/district-kathmandu

Where list= must be valid AreaName sepearated by commas
