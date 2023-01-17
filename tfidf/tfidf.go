package tfidf

import (
	"crypto/md5"
	"encoding/hex"
	"math"

	"github.com/sdpatel1026/doc-search/configs"
	"github.com/sdpatel1026/doc-search/tfidf/tokenize"
)

var tfIdf *TFIDF

// TFIDF tfidf model
type TFIDF struct {
	docIndex     map[string]int         // train document index in TermFreqs
	indexDocName map[int]string         // train document name mapped with index in TermFreqs
	termFreqs    []map[string]int       // terms frequency for each train document
	weights      []map[string]float64   //tf-idf of weight for each train document
	termDocs     map[string]int         // number of documents for each term in train data
	n            int                    // number of documents in train data
	stopWords    map[string]interface{} // words to be remove.
	tokenizer    tokenize.Tokenizer
}

// New new model with default
func New() *TFIDF {

	if tfIdf != nil {
		return tfIdf
	}
	tfIdf = &TFIDF{
		docIndex:     make(map[string]int),
		indexDocName: make(map[int]string),
		termFreqs:    make([]map[string]int, 0),
		weights:      make([]map[string]float64, 0),
		termDocs:     make(map[string]int),
		n:            0,
		tokenizer:    &tokenize.EnTokenizer{},
	}
	return tfIdf
}

// AddStopWords add stop words to be remove
func (tfIDF *TFIDF) AddStopWords(words ...string) {
	if tfIDF.stopWords == nil {
		tfIDF.stopWords = make(map[string]interface{})
	}

	for _, word := range words {
		tfIDF.stopWords[word] = nil
	}
}

// TrainDocs train documents
func (tfIDF *TFIDF) TrainDocs(docs map[string][]byte) []map[string]interface{} {
	results := make([]map[string]interface{}, len(docs))
	for docName, content := range docs {
		h := hash(content)
		docPos := tfIDF.docHashPos(h)
		if docPos >= 0 {
			result := make(map[string]interface{})
			result[configs.KEY_DOC_ID] = docPos
			result[configs.KEY_DOC_NAME] = docName
			result[configs.KEY_MSG] = configs.DOC_TRAINED
			results = append(results, result)
			continue
		}

		termFreq := tfIDF.termFreq(string(content))

		//not required to train doc as it does not contain useful information.
		if len(termFreq) == 0 {
			result := make(map[string]interface{})
			result[configs.KEY_DOC_ID] = -1
			result[configs.KEY_DOC_NAME] = docName
			result[configs.KEY_MSG] = configs.TRAINING_NOT_REQUIRED
			results = append(results, result)
			continue
		}
		tfIDF.termFreqs = append(tfIDF.termFreqs, termFreq)
		tfIDF.docIndex[h] = tfIDF.n
		tfIDF.indexDocName[tfIDF.n] = docName
		docID := tfIDF.n
		tfIDF.n++
		for term := range termFreq {
			tfIDF.termDocs[term]++
			tfIDF.UpdateWeights(term)
		}
		weight := tfIDF.weight(docID)
		tfIDF.weights = append(tfIDF.weights, weight)
		result := make(map[string]interface{})
		result[configs.KEY_DOC_ID] = docID
		result[configs.KEY_DOC_NAME] = docName
		result[configs.KEY_MSG] = configs.DOC_TRAINED
		results = append(results, result)
	}
	return results
}

// UpdateWieghts update weight of all docs present in corpus.
func (tfIDF *TFIDF) UpdateWeights(term string) {
	for pos, weight := range tfIDF.weights {
		_, isTermPresent := weight[term]
		if isTermPresent {
			weight[term] = findTfIdf(tfIDF.termFreqs[pos][term], 1, tfIDF.termDocs[term], tfIDF.n)
		}
	}
}

// weight calculate weight of doc.
func (tfIDF *TFIDF) weight(docPos int) (weight map[string]float64) {

	weight = make(map[string]float64)
	termFreq := tfIDF.termFreqs[docPos]
	for term, freq := range termFreq {
		weight[term] = findTfIdf(freq, 1, tfIDF.termDocs[term], tfIDF.n)
		// weight[term] = findTfIdf(freq, docTerms, f.termDocs[term], f.n)
	}
	return weight

}

// Cal calculate tf-idf weight for specified document
func (tfIDF *TFIDF) Cal(doc string) (weight map[string]float64) {
	weight = make(map[string]float64)

	var termFreq map[string]int

	docPos := tfIDF.docPos(doc)
	if docPos < 0 {
		termFreq = tfIDF.termFreq(doc)
	} else {
		termFreq = tfIDF.termFreqs[docPos]
	}
	// docTerms := 0
	// for _, freq := range termFreq {
	// 	docTerms += freq
	// }
	for term, freq := range termFreq {
		weight[term] = findTfIdf(freq, 1, tfIDF.termDocs[term], tfIDF.n)
		// weight[term] = findTfIdf(freq, docTerms, f.termDocs[term], f.n)
	}

	return weight
}

// termFreq calculate term-freq of each term in document.
func (tfIDF *TFIDF) termFreq(doc string) (m map[string]int) {
	m = make(map[string]int)

	tokens := tfIDF.tokenizer.Tokens(doc)
	if len(tokens) == 0 {
		return
	}

	for _, term := range tokens {
		if _, ok := tfIDF.stopWords[term]; ok {
			continue
		}

		m[term]++
	}

	return
}

// docHashPos return the position of doc in corpus.
func (tfIDF *TFIDF) docHashPos(hash string) int {
	if pos, ok := tfIDF.docIndex[hash]; ok {
		return pos
	}

	return -1
}

// docPos return the position of doc in corpus.
func (tfIDF *TFIDF) docPos(doc string) int {
	return tfIDF.docHashPos(hash([]byte(doc)))
}

// hash return hash of the doc content.
func hash(text []byte) string {
	h := md5.New()
	h.Write(text)
	return hex.EncodeToString(h.Sum(nil))
}

// findTfIdf calculate tf-idf.
func findTfIdf(termFreq, docTerms, termDocs, N int) float64 {
	tf := float64(termFreq) / float64(docTerms)
	idf := math.Log(float64(1+N) / (1 + float64(termDocs)))
	return tf * idf
}
