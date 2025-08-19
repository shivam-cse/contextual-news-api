package newsArticle


type LLMEntitiesAndIntentOutput struct {
	Intent   string        `json:"intent,omitempty"`
	Entities []string      `json:"entities,omitempty"`
	Keywords []string      `json:"keywords,omitempty"` // Important searchable terms
}