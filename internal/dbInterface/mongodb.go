package dbInterface

import (
	"context"
	"log/slog"
	"sort"

	"github.com/shivam-cse/contextual-news-api/internal/models/newsArticle"
	"github.com/shivam-cse/contextual-news-api/pkg/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NewsDbInterface struct {
	DB     *mongo.Database
	Logger *slog.Logger
}

// NewNewsDbInterface is the constructor for NewsDbInterface.
func NewNewsDbInterface(db *mongo.Database, logger *slog.Logger) *NewsDbInterface {
	return &NewsDbInterface{
		DB:     db,
		Logger: logger,
	}
}

func (newsDbInterface *NewsDbInterface) FindAllArticles(
	ctx context.Context,
	collName string,
	maxSize int64,
) ([]newsArticle.NewsArticleDBResponse, error) {
	newsDbInterface.Logger.Debug("'Data Layer': Fetching latest news articles...")
	coll := newsDbInterface.DB.Collection(collName)

	// Create a filter with an opts
	// Find all articles and
	// Sort the result by publication_date in descending order
	filter := bson.M{}
	opts := options.Find().SetSort(bson.D{primitive.E{Key: "publication_date", Value: -1}})
	if maxSize > 0 {
		opts.SetLimit(maxSize)
	}

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var newsArticles []newsArticle.NewsArticleDBResponse
	if err = cursor.All(ctx, &newsArticles); err != nil {
		return nil, err
	}

	newsDbInterface.Logger.Debug("Successfully fetched latest news articles")
	return newsArticles, nil
}

func (newsDbInterface *NewsDbInterface) FindArticlesByCategory(
	ctx context.Context,
	collName string,
	maxSize int64,
	category string,
) ([]newsArticle.NewsArticleDBResponse, error) {
	newsDbInterface.Logger.Debug("'Data Layer': Fetching news articles by category...")

	coll := newsDbInterface.DB.Collection(collName)

	// Create a filter with an opts.
	// For the 'category' filter, we use a regex to match the category case-insensitively.
	// For example: if category is "Sports", it will match "Sports", "sports", "Sports Cricket", "cricket Sports", "sportnews" etc.
	// And finally, sort the result by publication_date in descending order.
	filter := bson.M{"category": bson.M{"$regex": primitive.Regex{Pattern: category, Options: "i"}}}
	opts := options.Find().SetSort(bson.D{primitive.E{Key: "publication_date", Value: -1}})
	if maxSize > 0 {
		opts.SetLimit(maxSize)
	}

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var newsArticles []newsArticle.NewsArticleDBResponse
	if err := cursor.All(ctx, &newsArticles); err != nil {
		return nil, err
	}

	return newsArticles, nil
}

func (newsDbInterface *NewsDbInterface) FindArticlesByScore(
	ctx context.Context,
	collName string,
	maxSize int64,
	threshold float64,
) ([]newsArticle.NewsArticleDBResponse, error) {
	newsDbInterface.Logger.Debug("'Data Layer': Fetching news articles by score...")
	coll := newsDbInterface.DB.Collection(collName)

	// Create a filter with an opts
	// For the 'threshold' filter, we use a range query to match scores greater than or equal to the specified threshold.
	// And finally, sort the result by relevance_score in descending order.
	filter := bson.M{"relevance_score": bson.M{"$gte": threshold}}
	opts := options.Find().SetSort(bson.D{primitive.E{Key: "relevance_score", Value: -1}})
	if maxSize > 0 {
		opts.SetLimit(maxSize)
	}

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var newsArticles []newsArticle.NewsArticleDBResponse
	if err := cursor.All(ctx, &newsArticles); err != nil {
		return nil, err
	}

	return newsArticles, nil
}

func (newsDbInterface *NewsDbInterface) FindArticlesBySearchQuery(
	ctx context.Context,
	collName string,
	maxSize int64,
	query string,
) ([]newsArticle.NewsArticleDBResponse, error) {
	newsDbInterface.Logger.Debug("'Data Layer': Searching news articles...")
	coll := newsDbInterface.DB.Collection(collName)

	// Create a filter with an opts
	filter := bson.M{"$text": bson.M{"$search": query}}
	opts := options.Find()
    opts.SetProjection(bson.M{"score": bson.M{"$meta": "textScore"}})

    // Sort by text matching score and relevance_score
    opts.SetSort(bson.D{
		primitive.E{Key: "score", Value: bson.M{"$meta": "textScore"}},
        primitive.E{Key: "relevance_score", Value: -1},
	})
	if maxSize > 0 {
		opts.SetLimit(maxSize)
	}

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var newsArticles []newsArticle.NewsArticleDBResponse
	if err := cursor.All(ctx, &newsArticles); err != nil {
		return nil, err
	}

	return newsArticles, nil
}

func (newsDbInterface *NewsDbInterface) FindArticlesBySource(
	ctx context.Context,
	collName string,
	maxSize int64,
	source string,
) ([]newsArticle.NewsArticleDBResponse, error) {
	newsDbInterface.Logger.Debug("'Data Layer': Fetching news articles by source...")
	coll := newsDbInterface.DB.Collection(collName)

	// Create a filter with an opts
	// For the 'source' filter, we use a regex to match the source case-insensitively.
	// for example: if source is "BBC", it will match "BBC", "bbc", "BBC News", "News BBC", "newsbbc" etc.
	// And finally, sort the result by publication_date in descending order
	filter := bson.M{"source_name": bson.M{"$regex": primitive.Regex{Pattern: source, Options: "i"}}}
	opts := options.Find().SetSort(bson.D{primitive.E{Key: "publication_date", Value: -1}})
	if maxSize > 0 {
		opts.SetLimit(maxSize)
	}

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var newsArticles []newsArticle.NewsArticleDBResponse
	if err := cursor.All(ctx, &newsArticles); err != nil {
		return nil, err
	}

	return newsArticles, nil
}

func (newsDbInterface *NewsDbInterface) FindArticlesNearby(
	ctx context.Context,
	collName string,
	maxSize int64,
	latitude float64,
	longitude float64,
	radius float64,
) ([]newsArticle.NewsArticleDBResponse, error) {
	newsDbInterface.Logger.Debug("'Data Layer': Fetching nearby news articles...")
	coll := newsDbInterface.DB.Collection(collName)

	// Create a filter with an opts
	// For the 'location' filter, we use a geospatial query to find articles near the specified coordinates.
	// And we use a 'maxDistance' to limit the search to a specific radius.
	// Finally, it sorts the results with nearest first
	filter := bson.M{
		"location": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": bson.A{longitude, latitude},
				},
				"$maxDistance": radius * 1000, // radius in meters
			},
		},
	}
	
	opts := options.Find()
	if maxSize > 0 {
		opts.SetLimit(maxSize)
	}

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var newsArticles []newsArticle.NewsArticleDBResponse
	if err := cursor.All(ctx, &newsArticles); err != nil {
		return nil, err
	}

	return newsArticles, nil
}

func (newsDbInterface *NewsDbInterface) InsertUserEvent(ctx context.Context, event newsArticle.UserEvent) error {
	newsDbInterface.Logger.Debug("'Data Layer': Inserting user event...")
	coll := newsDbInterface.DB.Collection(constants.USER_EVENT)

	_, err := coll.InsertOne(ctx, event)
	if err != nil {
		return err
	}

	return nil
}

func (newsDbInterface *NewsDbInterface) FindTrendingArticles(
	ctx context.Context,
	newsCollName string,
	userEventCollName string,
	maxSize int64, 
	latitude float64, 
	longitude float64, 
	radius float64,
) ([]newsArticle.NewsArticleDBResponse, error) {
	newsDbInterface.Logger.Debug("'Data Layer': Fetching trending news articles...")
	coll := newsDbInterface.DB.Collection(newsCollName)

	
	// get all the user-events 
	userEvents, err := newsDbInterface.GetAllUserEvents(ctx, userEventCollName)
	if err != nil {
		return nil, err
	}

	newsDbInterface.Logger.Debug("====> Fetched all user events for trending articles: ", "events", userEvents)

	// Calculate the user events scores
	trendingScores := make(map[string]float64)
	for _, event := range userEvents {
		score := 1.0 // for view
		if event.EventType == "click" {
			score = 3.0 // for click
		}
		trendingScores[event.ArticleID] += score
	}

	// sort the userEvent by score
	type kv struct {
		Key   string
		Value float64
	}
	var sortedScores []kv
	for k, v := range trendingScores {
		sortedScores = append(sortedScores, kv{k, v})
	}
	sort.Slice(sortedScores, func(i, j int) bool {
		return sortedScores[i].Value > sortedScores[j].Value
	})

	// get the top N trending articles ids
	var trendingArticleIDs []string
	for i := 0; i < int(maxSize) && i < len(sortedScores); i++ {
		trendingArticleIDs = append(trendingArticleIDs, sortedScores[i].Key)
	}

	// If no valid article IDs found, return empty slice
	if len(trendingArticleIDs) == 0 {
		return []newsArticle.NewsArticleDBResponse{}, nil
	}

	filter := bson.M{"_id": bson.M{"$in": trendingArticleIDs}}
	opts := options.Find().SetSort(bson.D{primitive.E{Key: "timestamp", Value: -1}})

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var newsArticles []newsArticle.NewsArticleDBResponse
	if err := cursor.All(ctx, &newsArticles); err != nil {
		return nil, err
	}

	return newsArticles, nil
}

func (newsDbInterface *NewsDbInterface) GetAllUserEvents(ctx context.Context, collName string) ([]newsArticle.UserEvent, error) {
	newsDbInterface.Logger.Debug("'Data Layer': Fetching all user events...")
	coll := newsDbInterface.DB.Collection(collName)

	filter := bson.M{}
	opts := options.Find().SetSort(bson.D{primitive.E{Key: "timestamp", Value: -1}})

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var userEvents []newsArticle.UserEvent
	if err := cursor.All(ctx, &userEvents); err != nil {
		return nil, err
	}

	return userEvents, nil
}