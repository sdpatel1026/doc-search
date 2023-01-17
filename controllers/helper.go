package controllers

import (
	"io/ioutil"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/sdpatel1026/doc-search/configs"
)

func processError(gContext *gin.Context, errMsg interface{}, requestID string, httpStatusCode int) {
	var response Response = make(Response)
	response[configs.KEY_ERROR] = errMsg
	response[configs.KEY_REQ_ID] = requestID
	gContext.AbortWithStatusJSON(httpStatusCode, response)
}

func readFile(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}
