# `RT-CV` matcher

_Real time curriculum vitae matcher_

This is an api for CV scrapers to upload found CVs to where this tries to match it to a list of user defined search profiles at the moment of a CV upload _(hence the name realtime)_. The actions taken when a CV is matched is based on the search profile.
This tool also provides helper tools for managing secrets and keys.

**What this isn't:**

- A CV scraper for x website (that's up to you)
- A gui where you can search for CVs (this is only an API nor a database)
- A CV database (it might cache some cv information but **CVs are never written to any disk by this program**)

**What this is:**

An API that matches defined search profiles to scraped CVs.
On one side you have your CV scraper that scraps CVs and uploads them to this API.
On the other side you have a controller program build by you that defines the search profiles where scraped CVs should be compared to.
The API can also store scraper secrets like site login credentials.

We intent to make this API GDPR compliant.

## Quickly start hacking:

Requirements:

- GoLang 1.15+
- Nodejs 14+

```bash
cd dashboard
npm install
npm run build
cd ..
echo USE_TESTING_DB=true > .env

go run .
```

Copy the printed out dashboard token en head over to [localhost:4000](http://localhost:4000)

_note that this uses an in memory database that resets every time the app restarts, the default contents is defined in [mock/mock.go](./mock/mock.go)_

## Full Setup

Extra requirements:

- Mongodb _(mongodb compass is a great db viewer)_

Make sure to also create a new mongodb database, the collections are created automatically by this program

```bash
cp .env.example .env
vim .env

go run .
```

## API Docs

Head over to [localhost:4000/docs](http://localhost:4000/docs) to get the api docs

## Contribute

Make sure to read

- [/dashboard](/dashboard)
- [/CONTRIBUTING.md](/CONTRIBUTING.md)

## Docker

Build the project using

```sh
docker build -t rtcv:latest .
```

Run the project using
_Note that you probably want to change the environment variables_

```sh
docker run -it --rm -e USE_TESTING_DB=true -p 4000:4000 rtcv:latest
```
