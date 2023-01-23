package configs

const (
	THRESHOLD_SIMILARITIES float64 = 0.40
	OUTPUT_LEN             int     = 10 // no. of documents given in output

	TECHNICAL_ERROR       = "Something went wrong. Please try again later..."
	FILE_MISSING_ERROR    = "docs required."
	READING_ERROR         = "error in reading the document."
	TEXT_MISSING_ERROR    = "text required."
	DOC_TRAINED           = "document successfully trained."
	DOC_ALEREADY_TRAINED  = "document is already trained."
	TRAINING_NOT_REQUIRED = "This document does not needs to be trained, as it does not contains any useful information."
	KEY_ERROR             = "error"
	KEY_MSG               = "result"
	KEY_RANK              = "rank"
	KEY_REQ_ID            = "req_id"
	KEY_DOCS              = "docs"
	KEY_DOC_ID            = "doc_id"
	KEY_DOC_NAME          = "doc_name"
	KEY_TEXT              = "text"
	KEY_RESULT            = "result"
)
