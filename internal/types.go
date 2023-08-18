package internal

type GetSidResponse struct {
	Sid          string   `json:"sid"`
	Upgrades     []string `json:"upgrades"`
	PingTimeout  int      `json:"pingTimeout"`
	PingInterval int      `json:"pingInterval"`
}

type AnswerResponse struct {
	Answer   AnswerDetails
	Response AskResponse
}

type AskRequest struct {
	Source      SearchSource `json:"source"`
	Language    string       `json:"language,omitempty"`
	Timezone    string       `json:"timezone,omitempty"`
	SearchFocus SearchFocus  `json:"search_focus,omitempty"`
	Gpt4        bool         `json:"gpt4,omitempty"`
	Mode        SearchMode   `json:"mode,omitempty"`
}

type AskResponse struct {
	UUID           string   `json:"uuid"`
	Mode           string   `json:"mode"`
	SearchFocus    string   `json:"search_focus"`
	RelatedQueries []string `json:"related_queries"`
	Gpt4           bool     `json:"gpt4"`
	QueryStr       string   `json:"query_str"`
	Text           string   `json:"text"`
}

type AnswerDetails struct {
	Text       string `json:"answer"`
	WebResults []struct {
		Name     string `json:"name"`
		URL      string `json:"url"`
		Snippet  string `json:"snippet"`
		Client   string `json:"client"`
		MetaData any    `json:"meta_data"`
	} `json:"web_results"`
	Chunks          []string `json:"chunks"`
	EntityLinks     any      `json:"entity_links"`
	ExtraWebResults []struct {
		Name     string `json:"name"`
		URL      string `json:"url"`
		Snippet  string `json:"snippet"`
		Client   string `json:"client"`
		MetaData any    `json:"meta_data"`
	} `json:"extra_web_results"`
	DeletedUrls   []any `json:"deleted_urls"`
	ImageMetadata []any `json:"image_metadata"`
}
