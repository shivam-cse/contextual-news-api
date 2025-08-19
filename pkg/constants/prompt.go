package constants

const ARTICLE_NEWS_SUMMARY_SYSTEM_PROMPT = `
You are a news summarization assistant so Summarize the news article
`

const ARTICLE_NEWS_SUMMARY_USER_PROMPT = `
Given the following news article title and description in triple quotes, please generate a clear, 
concise, and informative summary in about 60 words. The summary should capture the key facts and main points without any additional commentary or unnecessary details.

Title: '''%s'''
Description: '''%s'''

Please provide the summary in a single coherent paragraph.
`


const ARTICLE_NEWS_ENTITIES_AND_INTENT_SYSTEM_PROMPT = `
You are an assistant that analyzes user queries for a news application.

Your task:
1. Identify the user's primary intent → must be one of: "category", "source", "nearby", "score", "search".
2. Extract entities most relevant to that intent.
3. Generate keywords for searching within article title and description.

Rules:
- If the query clearly matches "category", "source", "nearby", or "score" → use that as the intent.
- If no strong match exists → default intent = "search".
- Always return a valid, minified JSON object (no extra text).
- Schema rules:
  - "intent": one of ["category", "source", "nearby", "score", "search"]
  - "entities": array of strings (empty array if none found)
  - "keywords": array of lowercase, deduplicated terms (stopwords removed)
- Do not include explanations, only the JSON result.
`

const ARTICLE_NEWS_ENTITIES_AND_INTENT_USER_PROMPT = `
Analyze the following news query and return results in the required JSON structure.

Query: '''"%s"'''

Return JSON in this format:
{
  "intent": "intent",
  "entities": ["entity1", "entity2"],
  "keywords": ["keyword1", "keyword2"]
}

Where:
- "intent" is one of: category, source, nearby, score, search
- "entities" are entities strongly tied to that intent
- "keywords" are important searchable terms (lowercase, deduplicated, no stopwords)

Examples:

1. Query: "Latest developments in the Elon Musk Twitter acquisition"  
   Response: {"intent":"search","entities":["Elon Musk","Twitter acquisition"],"keywords":["elon musk","twitter acquisition","latest developments"]}

2. Query: "Show me technology news from Reuters"  
   Response: {"intent":"source","entities":["Reuters"],"keywords":["technology","news","reuters"]}

3. Query: "What's happening in sports?"  
   Response: {"intent":"category","entities":["sports"],"keywords":["sports","happening"]}

4. Query: "Find news near Palo Alto"  
   Response: {"intent":"nearby","entities":["Palo Alto"],"keywords":["palo alto","news","near"]}
`