# WIP

This project is still a WIP and should not be used in a production environment.

Documentation is also not yet available due to this

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

## Setup

Requirements:

- GoLang 1.15+
- Mongodb _(mongodb compass is a great db viewer)_

```sh
cp .env.example .env
vim .env

go run .
```

Before you start fiddling around with the code make sure to read [CONTRIBUTING.md](/CONTRIBUTING.md)

## Auth

How to generate a token:

_The functions below don't exsist they explain the kind of function that should be called_
_sha256 can also be used, use use replace sha512 with sha256 everywhere below_

#### On app init

```js
apiKey = getApiKey();
apiKeyID = getApiKeyID();
seed = fetchJson("/v1/auth/seed").seed;
salt = random(32)
key = sha512(seed + apiKey + salt);
```

#### For every request

```js
key = sha512(key + apiKey + salt);
return "Authorization: Basic " + base64(`sha512:${apiKeyID}:${salt}:${key}`);
```

#### If auth fails while having a theoretically valid key

1: just retry it (the server might be offline or whatever)

```js
return "Authorization: Basic " + base64(`sha512:${apiKeyID}:${salt}:${key}`);
```

2: Get a new salt and start over (basically going back to "on init")

```js
seed = fetchJson("/v1/auth/salt").seed;
salt = random(32);
key = sha512(apiKey + salt);
key = sha512(key + apiKey + salt);
return "Authorization: Basic " + base64(`sha512:${apiKeyID}:${salt}:${key}`);
```
