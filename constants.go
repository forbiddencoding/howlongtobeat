package howlongtobeat

import "time"

const (
	// hltbBaseURL is the base URL for the HowLongToBeat.
	hltbBaseURL = "https://howlongtobeat.com"
	// hltbSearchEndpoint is the default endpoint for the HowLongToBeat search API.
	hltbSearchEndpoint = "/api/finder"
	// hltbTokenURL is the URL to retrieve the token for the HowLongToBeat API.
	hltbTokenURL = "https://howlongtobeat.com/api/finder/init"
	// hltbGameURL is the base URL for the HowLongToBeat game API.
	hltbGameURL = "https://howlongtobeat.com/game"
	// defaultRequestTimeout is the default timeout for outgoing requests, we wait up to 30 seconds.
	defaultRequestTimeout = 30 * time.Second
)
