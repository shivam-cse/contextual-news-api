package newsResponse

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shivam-cse/contextual-news-api/pkg/constants"
)

// APIResponse defines the standard JSON response structure.
type APIResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Articles  interface{} `json:"articles,omitempty"`
    Metadata interface{} `json:"metadata,omitempty"`
    ErrorDetails   interface{} `json:"errorDetails,omitempty"`
}

// Success sends a standardized successful response.
func Success(
    c *gin.Context,
    logger *slog.Logger,
    statusCode int,
    message string,
    articles interface{},
    length int,
) {
    logger.Info("API Success",
        slog.String("message", message),
        slog.String("path", c.Request.URL.Path),
    )

	c.JSON(statusCode, APIResponse{
        Status: constants.SUCCESS,
        Message: message,
        Articles: articles,
        Metadata: map[string]interface{}{
            "count": length,
            "query": c.Request.URL.Query(),
            "path": c.Request.URL.Path,
        },
    })
}

func Error(c *gin.Context, logger *slog.Logger, statusCode int, message string, err error) {
    errorDetails := ""
    if err != nil {
        errorDetails = err.Error()
    }

    logger.Error(
        "API Error",
        slog.String("message", message),
        slog.String("internal_error", errorDetails),
        slog.String("path", c.Request.URL.Path),
    )

    // Return a generic error message to the client.
    c.AbortWithStatusJSON(statusCode, APIResponse{
        Status: constants.FAILED,
        Message: message,
        Metadata: map[string]interface{}{
            "count": 0,
            "query": c.Request.URL.Query(),
            "path": c.Request.URL.Path,
        },
        ErrorDetails: errorDetails, // Or a more generic message
    })
}

func TimeOut(c *gin.Context) {
    c.AbortWithStatusJSON(http.StatusGatewayTimeout, APIResponse{
        Status: constants.FAILED,
        Message: "Request timeout - please try again later",
        Metadata: map[string]interface{}{
            "count": 0,
            "query": c.Request.URL.Query(),
            "path": c.Request.URL.Path,
        },
        ErrorDetails: "The request could not be completed within the allowed time limit",
    })
}