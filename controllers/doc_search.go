package controllers

import (
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sdpatel1026/doc-search/configs"
	"github.com/sdpatel1026/doc-search/tfidf"
)

// DocSearch search important document for given text query.
func DocSearch(gContext *gin.Context) {
	log.Println("got an request to search the documents")
	requestID := uuid.New().String()
	text := gContext.Param(configs.KEY_TEXT)
	text = strings.TrimSpace(text)
	if text == "" {
		processError(gContext, configs.TEXT_MISSING_ERROR, requestID, http.StatusBadRequest)
		return
	}
	tfIDF := tfidf.New()
	textWeight := tfIDF.Cal(text)
	docsRank := tfIDF.Ranks(textWeight)
	docIDs := make([]int, len(docsRank))
	for docID := range docsRank {
		docIDs = append(docIDs, docID)
	}
	sort.SliceStable(docIDs, func(i, j int) bool {
		return docsRank[docIDs[i]] > docsRank[docIDs[j]]
	})
	response := make([]map[string]interface{}, configs.OUTPUT_LEN)
	var count = 0
	for _, docID := range docIDs {
		result := make(map[string]interface{})
		result[configs.KEY_DOC_ID] = docID
		result[configs.KEY_RANK] = docsRank[docID]
		response = append(response, result)
		count++
		if count == configs.OUTPUT_LEN {
			break
		}
	}
	gContext.JSON(http.StatusOK, gin.H{configs.KEY_RESULT: response})
}
