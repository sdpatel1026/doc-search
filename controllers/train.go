package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sdpatel1026/doc-search/configs"
	"github.com/sdpatel1026/doc-search/tfidf"
)

// Train train documents
func Train(gContext *gin.Context) {
	log.Println("got an request to train the documents")
	requestID := uuid.New().String()
	form, err := gContext.MultipartForm()
	if err != nil {
		log.Printf("error in reading files from request body for a req_id %s is :%s\n", requestID, err.Error())
		processError(gContext, configs.TECHNICAL_ERROR, requestID, http.StatusInternalServerError)
		return
	}
	files := form.File[configs.KEY_DOCS]
	if len(files) < 1 {
		processError(gContext, configs.FILE_MISSING_ERROR, requestID, http.StatusBadRequest)
		return
	}
	responses := make([]map[string]interface{}, len(files))
	fileNameContent := make(map[string][]byte)
	for _, file := range files {
		fileContent, err := readFile(file)
		if err != nil {
			log.Printf("error in reading file for a req_id %s is: %s\n", requestID, err.Error())
			response := make(map[string]interface{})
			response[configs.KEY_DOC_NAME] = file.Filename
			response[configs.KEY_ERROR] = configs.READING_ERROR
			responses = append(responses, response)
			continue
		}
		fileNameContent[file.Filename] = fileContent
	}
	tfIDF := tfidf.New()
	results := tfIDF.TrainDocs(fileNameContent)
	responses = append(responses, results...)
	gContext.JSON(http.StatusOK, gin.H{configs.KEY_RESULT: responses})
}
