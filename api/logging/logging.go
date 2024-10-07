package logging

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

const (
	Start      = "start"
	Index      = "minerva"
	Level      = "level"
	StatusCode = "statusCode"
	Message    = "message"
)

const (
	Error = "error"
	Info  = "info"
)

type Logging struct {
	client *elasticsearch.TypedClient
}

type LogDocument struct {
	Timestamp  time.Time `json:"timestamp"`
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
	Latency    int64     `json:"latency"`
	Path       string    `json:"path"`
	Level      string    `json:"level"`
	StatusCode int       `json:"statusCode"`
	Message    string    `json:"message,omitempty"`
	Method     string    `json:"method"`
}

func New(host string) (*Logging, error) {
	typedClient, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{host},
	})

	if err != nil {
		return nil, err
	}

	return &Logging{typedClient}, nil
}

func (receiver Logging) Set(c *gin.Context, level string, statusCode int, message string) {
	c.Set(Level, level)
	c.Set(StatusCode, statusCode)
	c.Set(Message, message)
}

func (receiver Logging) Start(c *gin.Context) {
	c.Set(Start, time.Now())
	c.Next()
}

func (receiver Logging) End(c *gin.Context) {
	c.Next()

	end := time.Now()
	start := c.GetTime(Start)
	latency := end.Sub(start)

	doc := LogDocument{
		Timestamp:  end,
		Start:      start,
		End:        end,
		Latency:    latency.Milliseconds(),
		Path:       c.FullPath(),
		Level:      c.GetString(Level),
		StatusCode: c.GetInt(StatusCode),
		Message:    c.GetString(Message),
		Method:     c.Request.Method,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := receiver.client.Index(Index).Request(doc).Do(ctx)
	if err != nil {
		log.Printf("error sending logs: %v", err)
	}
}
