# `RT-CV` matcher

_Real time curriculum vitae matcher_

This is an api for CV scrapers to upload found CVs to where this tries to match it to a list of user defined search profiles _(hence the name realtime)_. The actions taken when a CV is matched is based on the search profile.

**What this isn't:**

- A CV scraper for x website (that's up to you)
- A gui where you can search for CVs (this is only an API nor a database)
- A CV database (it might cache some cv information but **CVs are never written to disk by this program**)

**Goals:**

- Easy to understand API for uploading scraped CVs to and to set search profiles
- GDPR compliant
- Fast

## Quickly start hacking:

Only requires GoLang:
_note that this uses an in memory database that resets every time the app restarts, the default contents is defined in mock/mock.go_

```bash
USE_TESTING_DB=true go run .
```

## Setup

Requirements:

- GoLang 1.14+
- nodejs 14+
- Mongodb _(mongodb compass is a great db viewer)_

API:

```bash
cp .env.example .env
vim .env

go run .
```

Dashboard:

```bash
cd dashboard

npm i
npm run build
```

Before you commit make sure to read [CONTRIBUTING.md](/CONTRIBUTING.md)

## Dashboard

for more information about the dashboard see [/dashboard](/dashboard)

## API Documentation

There is a auto generated OpenAPI v3 schema available at `GET /api/v1/schema/openAPI`
You can download that as a `.json` and convert it to [human readable documentation](https://openapi-generator.tech/docs/generators#documentation-generators) using [openapi-generator](https://openapi-generator.tech) _([how to](https://stackoverflow.com/questions/59727169/how-to-generate-api-documentation-using-openapi-generator))_

There is also a human readable documentation page in the dashboard on the `/docs` page
