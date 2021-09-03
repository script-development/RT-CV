package controller

import "strings"

func mdBlockCodeSample(lang, sample string) string {
	return "```" + lang + "\n" + strings.Trim(sample, "\n") + "\n```"
}

var schemaDescription = `
# RT-CV matcher

_Real time curriculum vitae matcher_

## auth

Every authenticated route is tagged with one base tag.
In the case of routes with a required auth role, the role(s) is/are also added to the tags.
Note that in the case of routes with multiple auth roles you only need one of those roles to access the route.

## generate auth token

Notes:

- The functions below don't exist they explain the kind of function that should be called
- sha256 can also be used, use use replace sha512 with sha256 everywhere below
- the ` + "`sha512`" + ` function should return a hex value not the bytes

#### On app init

` + mdBlockCodeSample("js", `
apiKey = getApiKey();
apiKeyID = getApiKeyID();
seed = fetchJson("/api/v1/auth/seed").seed;
salt = random(32);
key = sha512(seed + apiKey + salt);
`) + `

#### For every request

` + mdBlockCodeSample("js", `
key = sha512(key + apiKey + salt);
return (
  "Authorization: Basic " + base64(`+"`sha512:${apiKeyID}:${salt}:${key}`"+`)
);
`) + `

#### If auth fails while having a theoretically valid key

1: just retry it (the server might be offline or whatever)

` + mdBlockCodeSample("js", `
return (
  "Authorization: Basic " + base64(`+"`sha512:${apiKeyID}:${salt}:${key}`"+`)
);
`) + `

2: Get a new salt and start over (basically going back to "on init")

` + mdBlockCodeSample("js", `
seed = fetchJson("/api/v1/auth/seed").seed;
salt = random(32);
key = sha512(apiKey + salt);
key = sha512(key + apiKey + salt);
return "Authorization: Basic " + base64(`+"`sha512:${apiKeyID}:${salt}:${key}`"+`);
`)
