package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	. "github.com/WebGou/baaplogger"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/mapping"
)

var batchSize = 200
var imgIndex bleve.Index

//InitSearch init search
func InitSearch() {

	indexPath := "index"
	// open the index
	done := make(chan struct{})
	var err error
	imgIndex, err = bleve.Open(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		Log.Informational("Creating new index...")
		// create a mapping
		indexMapping, err := BuildIndexMapping()
		if err != nil {
			log.Fatal(err)
		}
		imgIndex, err = bleve.New(indexPath, indexMapping)
		if err != nil {
			log.Fatal(err)
		}

		// index data in the background
		go func() {
			defer func() { done <- struct{}{} }()
			err = Indexkeyword(imgIndex)
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
	searchResults, err := imgIndex.Search(search)

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

	keywordMapping.AddFieldMappingsAt("Keywords", englishTextFieldMapping)

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
		idx := strings.Index(path, "images")
		batch.Index(path[idx:], jsonDoc)
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

//DoSearchP do the search
func DoSearchP(c *gin.Context) {
	SrcTx := c.PostForm("SrcTx")
	if SrcTx == "" {
		c.Redirect(http.StatusFound, "/search")
	}

	query := bleve.NewMatchQuery(SrcTx)
	search := bleve.NewSearchRequest(query)
	sr, err := imgIndex.Search(search)

	if err != nil {
		Log.Error(err.Error())
		return
	}

	ser := []string{}
	if sr.Total > 0 {
		if sr.Request.Size > 0 {
			//rv = fmt.Sprintf("%d matches, showing %d through %d, took %s\n", sr.Total, sr.Request.From+1, sr.Request.From+len(sr.Hits), sr.Took)
			for _, hit := range sr.Hits {
				//rv += fmt.Sprintf("%5d. %s (%f)\n", i+sr.Request.From+1, hit.ID, hit.Score)
				ser = append(ser, hit.ID)
			}
		}
	}

	s := sessions.Default(c)
	if len(ser) > 0 {
		s.AddFlash(ser, "result")

	} else {
		s.AddFlash(fmt.Sprintf(`sorry, cannot find information for "%s"`, SrcTx), "srmsg")
	}

	s.AddFlash(SrcTx, "lastTx")
	s.Save()
	c.Redirect(http.StatusFound, "/search")
	//Log.Informational("%v", searchResults)

}

//DoSearchG do search get
func DoSearchG(c *gin.Context) {
	s := sessions.Default(c)
	srmsg := s.Flashes("srmsg")

	msg := ""
	if len(srmsg) > 0 {
		msg = srmsg[0].(string)
	}

	result := []string{}
	resultS := s.Flashes("result")
	if len(resultS) > 0 {
		result = resultS[0].([]string)
	}

	lastTxS := s.Flashes("lastTx")

	lastTx := ""
	if len(lastTxS) > 0 {
		lastTx = lastTxS[0].(string)
	}

	s.Save()

	c.HTML(http.StatusOK, "globSearch.tmpl", gin.H{"msg": msg, "lastTx": lastTx, "result": result})
}
