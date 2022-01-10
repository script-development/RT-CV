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

- GoLang 1.16+
- Nodejs 15+

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
- Dart 2+

Make sure to also create a new mongodb database, the collections are created automatically by this program

```bash
cp .env.example .env
vim .env

cd pdf_generator
dart pub get
dart compile exe bin/pdf_generator.dart
cd ..

go run .
```

## API Docs

Head over to [localhost:4000/docs](http://localhost:4000/docs) to get the api docs

## Contribute

Make sure to read

- [/dashboard](/dashboard)
- [/pdf_generator](/pdf_generator)
- [/CONTRIBUTING.md](/CONTRIBUTING.md)

## Docker

Build the project using

```sh
docker build -t rtcv:latest .
```

Quickly run the project using docker

```sh
docker run -it --rm -e USE_TESTING_DB=true -p 4000:4000 rtcv:latest
```

<details><summary>Run the full project in docker</summary><br/>

```sh
# create a docker network so RT-CV and mongodb can communicate without exposing ports
docker network create f2f

# run the mongodb database
docker run \
    -d \
    -v /data/db:/data/db \
    --network f2f \
    mongo:5.0


# create an env file for the RT-CV app
# you can also use -e for every env variable but there might be a lot so this is easier
cp .env.example .env
vim .env

# run RT-CV
docker run \
  -d \
  --network f2f \
  --env-file $(pwd)/.env \
  -p 127.0.0.1:4000:4000 \
  rtcv:latest
```

</details>
