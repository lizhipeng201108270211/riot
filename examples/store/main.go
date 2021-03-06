// Copyright 2016 ego authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-ego/riot/engine"
	"github.com/go-ego/riot/types"
)

var (
	// searcher is coroutine safe
	searcher = engine.Engine{}
)

func initEngine() {
	var path = "./index"

	searcher.Init(types.EngineInitOptions{
		Using: 1,
		IndexerInitOptions: &types.IndexerInitOptions{
			IndexType: types.DocIdsIndex,
		},
		UseStorage:    true,
		StorageFolder: path,
		StorageEngine: "bg", // bg: badger, lbd: leveldb, bolt: bolt
		SegmenterDict: "../../data/dict/dictionary.txt",
		StopTokenFile: "../../data/dict/stop_tokens.txt",
	})
	defer searcher.Close()
	os.MkdirAll(path, 0777)

	// Add the document to the index, docId starts at 1
	searcher.IndexDocument(1, types.DocIndexData{Content: "Google Is Experimenting With Virtual Reality Advertising"}, false)
	searcher.IndexDocument(2, types.DocIndexData{Content: "Google accidentally pushed Bluetooth update for Home speaker early"}, false)
	searcher.IndexDocument(3, types.DocIndexData{Content: "Google is testing another Search results layout with rounded cards, new colors, and the 4 mysterious colored dots again"}, false)

	// Wait for the index to refresh
	searcher.FlushIndex()

	log.Println("recover index number:", searcher.NumDocumentsIndexed())
}

func restoreIndex() {
	var path = "./index"

	searcher.Init(types.EngineInitOptions{
		Using: 1,
		IndexerInitOptions: &types.IndexerInitOptions{
			IndexType: types.DocIdsIndex,
		},
		UseStorage:    true,
		StorageFolder: path,
		StorageEngine: "bg", // bg: badger, lbd: leveldb, bolt: bolt
		SegmenterDict: "../../data/dict/dictionary.txt",
		StopTokenFile: "../../data/dict/stop_tokens.txt",
	})
	defer searcher.Close()
	os.MkdirAll(path, 0777)

	// Wait for the index to refresh
	searcher.FlushIndex()

	log.Println("recover index number:", searcher.NumDocumentsIndexed())
}

func main() {
	initEngine()

	sea := searcher.Search(types.SearchRequest{
		Text: "google testing",
		RankOptions: &types.RankOptions{
			OutputOffset: 0,
			MaxOutputs:   100,
		}})

	fmt.Println("search---------", sea, "; docs=", sea.Docs)
}
