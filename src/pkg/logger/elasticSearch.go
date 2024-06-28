package logger

import (
	"bytes"
	"compress/gzip"
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog"
	"log"
	"os"
	"time"
)

type ESLogger struct {
	zerologger    *zerolog.Logger
	c             *elasticsearch.Client
	esBulkIndexer esutil.BulkIndexer
}

var _ Logger = (*ESLogger)(nil)

func NewElasticSearchLogger(addresses []string, user, pass string) Logger {
	logger := zerolog.New(os.Stdout)

	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses:                addresses,
		Username:                 user,
		Password:                 pass,
		CompressRequestBody:      true,
		CompressRequestBodyLevel: gzip.DefaultCompression,
		PoolCompressor:           true,
		EnableMetrics:            true,
	})

	if err != nil {
		log.Panicf("Error creating the elastic search client: %s", err)
	}

	_, err = esClient.Indices.Create("logs")
	if err != nil {
		log.Panicf("Error creating the index: %s", err)
	}

	esLogger := &ESLogger{
		zerologger: &logger,
	}

	esLogger.c = esClient

	esBulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:     esClient,
		FlushBytes: 1 << 20,
		Index:      "logs",
		OnError: func(ctx context.Context, err error) {
			esLogger.zerologger.Error().Err(err).Msg("Error in the bulk indexer")
		},
		OnFlushStart: func(ctx context.Context) context.Context {
			esLogger.zerologger.Info().Msg("Start flushing the bulk indexer")
			return ctx
		},
		OnFlushEnd: func(ctx context.Context) {
			esLogger.zerologger.Info().Msg("End flushing the bulk indexer")
		},
		FlushInterval: time.Second * 90,
	})

	if err != nil {
		log.Panicf("Error creating the bulk indexer: %s", err)
	}

	esLogger.esBulkIndexer = esBulkIndexer

	return esLogger
}

type record struct {
	Level string    `json:"level"`
	Time  time.Time `json:"time"`
	Msg   string    `json:"msg"`
}

func (l *ESLogger) Info(msg string) {
	data, _ := json.Marshal(record{Level: "info", Time: time.Now(), Msg: msg})

	err := l.esBulkIndexer.Add(context.Background(), esutil.BulkIndexerItem{
		Action: "index",
		Body:   bytes.NewReader(data),
		OnFailure: func(ctx context.Context, _ esutil.BulkIndexerItem, _ esutil.BulkIndexerResponseItem, err error) {
			l.zerologger.Error().Err(err).Msg("Error has been occurred while sending item from the bulk indexer")
		},
	})
	if err != nil {
		l.zerologger.Error().Err(err).Msg("Error has been occurred while adding item to the bulk indexer")
	}
	l.zerologger.Info().Msg(msg)
}

func (l *ESLogger) Error(msg string) {
	data, _ := json.Marshal(record{Level: "error", Time: time.Now(), Msg: msg})

	err := l.esBulkIndexer.Add(context.Background(), esutil.BulkIndexerItem{
		Action: "index",
		Body:   bytes.NewReader(data),
		OnFailure: func(ctx context.Context, _ esutil.BulkIndexerItem, _ esutil.BulkIndexerResponseItem, err error) {
			l.zerologger.Error().Err(err).Msg("Error has been occurred while sending item from the bulk indexer")
		},
	})
	if err != nil {
		l.zerologger.Error().Err(err).Msg("Error has been occurred while adding item to the bulk indexer")
	}
	l.zerologger.Error().Msg(msg)
}

func (l *ESLogger) Fatal(msg string) {
	data, _ := json.Marshal(record{Level: "fatal", Time: time.Now(), Msg: msg})

	err := l.esBulkIndexer.Add(context.Background(), esutil.BulkIndexerItem{
		Action: "index",
		Body:   bytes.NewReader(data),
		OnFailure: func(ctx context.Context, _ esutil.BulkIndexerItem, _ esutil.BulkIndexerResponseItem, err error) {
			l.zerologger.Error().Err(err).Msg("Error has been occurred while sending item from the bulk indexer")
		},
	})
	if err != nil {
		l.zerologger.Error().Err(err).Msg("Error has been occurred while adding item to the bulk indexer")
	}

	l.zerologger.Fatal().Msg(msg)
}
