# WIP

This project is still a WIP and should not be used in a production environment.

Documentation is also not yet available due to this

# `RT-CV` matcher

_Real time curriculum vitae matcher_

This is an api for CV scrapers to upload found CVs to where this tries to match it to a list of user defined search profiles _(hence the name realtime)_. The actions taken when a CV is matched is based on the search profile.

**What this isn't:**

- A CV scraper for x website (that's up to you)
- A gui where you can search for CVs (this is only an API)
- A CV database (it might cache some cv information but **CVs are never written to disk**)

**Goals:**

- Easy to understand API for uploading scraped CVs to and to set search profiles
- GDPR compliant
- Fast

# Auth

How to generate a token:

```js
// On init
apiKey = fetchApiKey();
salt = random(32);
key = sha512(apiKey + salt);

// For every message
key = sha512(key);
return "Authorization: Basic " + base64(`sha512:${keyID}:${salt}:${key}`);
```

_sha256 can also be used, to use make sure replace sha512 with sha256 everywhere above_
