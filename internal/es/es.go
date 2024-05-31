// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-05-31, by liasica

package es

import (
	"github.com/elastic/go-elasticsearch/v8"
)

type Elastic struct {
	client     *elasticsearch.TypedClient
	datastream string
	index      string
}

func New(apiKey, datastream string, addresses []string) (*Elastic, error) {
	c, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		APIKey:    apiKey,
		Addresses: addresses,
	})
	if err != nil {
		return nil, err
	}
	return &Elastic{
		client:     c,
		datastream: datastream,
		index:      datastream + "*",
	}, err
}
