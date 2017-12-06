package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	. "github.com/WebGou/baaplogger"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/mapping"
)

var batchSize = 200

//InitSearch init search
func InitSearch() {

	indexPath := "index"
	// open the index
	done := make(chan struct{})
	beerIndex, err := bleve.Open(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		Log.Informational("Creating new index...")
		// create a mapping
		indexMapping, err := BuildIndexMapping()
		if err != nil {
			log.Fatal(err)
		}
		beerIndex, err = bleve.New(indexPath, indexMapping)
		if err != nil {
			log.Fatal(err)
		}

		// index data in the background
		go func() {
			defer func() { done <- struct{}{} }()
			err = Indexkeyword(beerIndex)
			if err != nil {
				log.Fatal(err)
			}
		}()
		<-done
	} else if err != nil {
		Log.Error(err.Error())
	} else {
		Log.Informational("Opening existing index...")
	}

	query := bleve.NewMatchQuery("20171108175602")
	search := bleve.NewSearchRequest(query)
	searchResults, err := beerIndex.Search(search)

	if err != nil {
		Log.Error(err.Error())
		return
	}
	/* 	beerIndex.Index("111", `{"keywords":"test"}`)

	query = bleve.NewMatchQuery("test")
	search = bleve.NewSearchRequest(query)
	searchResults, err = beerIndex.Search(search) */

	Log.Informational("first serach result,%v", searchResults)
}

//BuildIndexMapping build the index mapping
func BuildIndexMapping() (mapping.IndexMapping, error) {

	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	keywordMapping := bleve.NewDocumentMapping()

	keywordMapping.AddFieldMappingsAt("keywords", englishTextFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("config", keywordMapping)

	indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping, nil
}

//Indexkeyword index keyword
func Indexkeyword(i bleve.Index) error {

	// open the directory
	jsonDir := filepath.Join(exePath, "images")

	// walk the directory entries for indexing
	Log.Informational("Indexing...")
	count := 0
	startTime := time.Now()
	batch := i.NewBatch()
	batchCount := 0

	filepath.Walk(jsonDir, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(info.Name(), "config.json") == false {
			return nil
		}
		// read the bytes
		jsonBytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		// parse bytes as json
		var jsonDoc interface{}
		err = json.Unmarshal(jsonBytes, &jsonDoc)
		if err != nil {
			return err
		}
		//ext := filepath.Ext(filename)
		//docID := filename[:(len(filename) - len(ext))]
		batch.Index(path, jsonDoc)
		batchCount++

		if batchCount >= batchSize {
			err = i.Batch(batch)
			if err != nil {
				return err
			}
			batch = i.NewBatch()
			batchCount = 0
		}
		count++
		if count%1000 == 0 {
			indexDuration := time.Since(startTime)
			indexDurationSeconds := float64(indexDuration) / float64(time.Second)
			timePerDoc := float64(indexDuration) / float64(count)
			Log.Informational("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
		}

		return nil

	})

	// flush the last batch
	if batchCount > 0 {
		err := i.Batch(batch)
		if err != nil {
			Log.Error(err.Error())
		}
	}
	indexDuration := time.Since(startTime)
	indexDurationSeconds := float64(indexDuration) / float64(time.Second)
	timePerDoc := float64(indexDuration) / float64(count)
	Log.Informational("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
	return nil
}
