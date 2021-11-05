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
- the ` + "`sha512`" + ` function should return a hex value not the bytes

#### On app init

` + mdBlockCodeSample("js", `
apiKey = getApiKey();
apiKeyID = getApiKeyID();
key = "Basic " + apiKeyID + ":" + sha512(apiKey);
`) + `
`
