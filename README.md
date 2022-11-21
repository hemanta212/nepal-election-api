# Nepal Election API 2079

Developed using the fantastic [gocolly](https://github.com/gocolly/colly) scraping library, data source is [ekantipur](https://election.ekantipur.com)

This repo was developed as backend for election-bot for reddit hosted at https://github.com/pykancha/reddit-bots

## Installation
- Install go
- Clone the repo
```
git clone https://github.com/hemanta212/nepal-election-api
cd nepal-election-api
```
- Install dependencies
```
go mod tidy
```
- Run the server
```
go run .
```
- For building an executable, use;
```
go build .
./nepal-election-api
```


## Usage
###### AreaName

Requests at

```
/area?name=pradesh-1/district-jhapa
```

where name is valid kantipur url part representing an electoral area.
This is supposed to be extracted from a kantipur url.

###### URL

Requests at

```
/url?url=https://election.ekantipur.com/pradesh-1/district-jhapa?lng=eng
```

where url must be valid kantipur url in format similar to url in above example.

###### Bulk List
Requests at

```
/bulk?list=pradesh-1/district-jhapa,pradesh-3/district-kathmandu
```

Where list must be list of valid AreaNames sepearated by commas.
