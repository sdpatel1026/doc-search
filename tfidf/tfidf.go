package tfidf

import (
	"crypto/md5"
	"encoding/hex"
	"math"
	"math/big"
	"strings"

	"github.com/sdpatel1026/doc-search/configs"
	"github.com/sdpatel1026/doc-search/tfidf/tokenize"
)

const (
	K float64 = 1.2
	B float64 = 0.75
)

var tfIdf *TFIDF

// TFIDF tfidf model
type TFIDF struct {
	docIndex      map[string]int          // train document index
	indexDocName  map[int]string          // train document name mapped with index
	termFreqs     []map[string]*big.Int   // terms frequency for each train document
	termDocsCount map[string]int          // number of document in which term t appears in train data
	termDocs      map[string]map[int]bool //list of doc in which term t present.
	docsTermLen   []*big.Int              //len of each doc in words.
	termsLen      *big.Int                //total terms in corpus
	n             int                     // number of documents in train data
	stopWords     map[string]interface{}  // words to be remove.
	tokenizer     tokenize.Tokenizer
}

// New new model with default
func New() *TFIDF {

	if tfIdf != nil {
		return tfIdf
	}
	tfIdf = &TFIDF{
		docIndex:      make(map[string]int),
		indexDocName:  make(map[int]string),
		termFreqs:     make([]map[string]*big.Int, 0),
		termDocs:      make(map[string]map[int]bool),
		termDocsCount: make(map[string]int),
		docsTermLen:   make([]*big.Int, 0),
		termsLen:      big.NewInt(0),
		n:             0,
		tokenizer:     &tokenize.EnTokenizer{},
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
	results := make([]map[string]interface{}, 0)
	for docName, content := range docs {
		docHash := hash(content)
		docPos := tfIDF.docHashPos(docHash)
		if docPos >= 1 {
			result := make(map[string]interface{})
			result[configs.KEY_DOC_ID] = docPos
			result[configs.KEY_DOC_NAME] = docName
			result[configs.KEY_MSG] = configs.DOC_ALEREADY_TRAINED
			results = append(results, result)
			continue
		}
		tokens := tfIDF.tokenizer.Tokens(string(content))
		tfIDF.n += 1
		termFreq := tfIDF.termFreq(tokens)
		//not required to train doc as it does not contain useful information.
		if len(termFreq) == 0 {
			result := make(map[string]interface{})
			result[configs.KEY_DOC_ID] = -1
			result[configs.KEY_DOC_NAME] = docName
			result[configs.KEY_MSG] = configs.TRAINING_NOT_REQUIRED
			results = append(results, result)
			continue
		}
		for term := range termFreq {
			termDocSet, isFound := tfIDF.termDocs[term]
			if !isFound {
				termDocSet = make(map[int]bool)
				tfIDF.termDocs[term] = termDocSet
			}
			termDocSet[tfIDF.n] = true
			tfIDF.termDocsCount[term] += 1
		}
		tfIDF.termFreqs = append(tfIDF.termFreqs, termFreq)
		docTokenCount := big.NewInt(int64(len(tokens)))
		tfIDF.docsTermLen = append(tfIDF.docsTermLen, docTokenCount)
		tfIDF.termsLen = tfIDF.termsLen.Add(tfIDF.termsLen, docTokenCount)
		tfIDF.docIndex[docHash] = tfIDF.n
		tfIDF.indexDocName[tfIDF.n] = docName

		result := make(map[string]interface{})
		result[configs.KEY_DOC_ID] = tfIDF.n
		result[configs.KEY_DOC_NAME] = docName
		result[configs.KEY_MSG] = configs.DOC_TRAINED
		results = append(results, result)
	}
	// var pos = 1
	// for pos <= tfIDF.n {
	// 	fmt.Println("---------------------------")
	// 	fmt.Printf("\n%s", tfIdf.indexDocName[pos])
	// 	fmt.Printf("\ntermFreqs --- %v\n", tfIDF.termFreqs[pos-1])
	// 	pos++
	// }
	return results
}

// UpdateWieghts update weight of all docs present in corpus.
//
//	func (tfIDF *TFIDF) updateWeights(term string) {
//		for pos, weight := range tfIDF.weights {
//			_, isTermPresent := weight[term]
//			if isTermPresent {
//				weight[term] = findTfIdf(tfIDF.termFreqs[pos][term], 1, tfIDF.termDocs[term], tfIDF.n)
//				tfIDF.weights[pos] = weight
//			}
//		}
//	}
func (tfIDF *TFIDF) DocName(docPos int) string {
	return tfIDF.indexDocName[docPos]
}

// weight calculate tf-idf weight of doc in the corpus.
func (tfIDF *TFIDF) weight(docPos int) map[string]float64 {

	weight := make(map[string]float64)
	termFreq := tfIDF.termFreqs[docPos-1]
	// fmt.Println("***************************************")
	for term, freq := range termFreq {
		// fmt.Printf("term: %v\n", term)
		// fmt.Printf("freq: %v\n", freq)
		// fmt.Printf("tfIDF.termDocs[term]: %v\n", tfIDF.termDocs[term])
		// fmt.Printf("tfIDF.n: %v\n", tfIDF.n)
		weight[term] = findTfIdf(freq, 1, tfIDF.termDocsCount[term], tfIDF.n)
		// weight[term] = findTfIdf(freq, docTerms, f.termDocs[term], f.n)
	}
	// fmt.Printf("termFreq: %v\n", termFreq)
	return weight

}

// Cal calculate tf-idf weight for specified document
func (tfIDF *TFIDF) Cal(doc string) (weight map[string]float64) {
	weight = make(map[string]float64)

	var termFreq map[string]int
	doc = strings.ToLower(doc)
	docPos := tfIDF.docPos(doc)
	if docPos < 0 {
		termFreq = tfIDF.termFreq(doc)
	} else {
		termFreq = tfIDF.termFreqs[docPos-1]
	}
	// docTerms := 0
	// for _, freq := range termFreq {
	// 	docTerms += freq
	// }
	for term, freq := range termFreq {
		weight[term] = findTfIdf(freq, 1, tfIDF.termDocsCount[term], tfIDF.n)
		// weight[term] = findTfIdf(freq, docTerms, f.termDocs[term], f.n)
	}

	return weight
}

// termFreq calculate term-freq of each term in document.
func (tfIDF *TFIDF) termFreq(tokens []string) (m map[string]*big.Int) {
	m = make(map[string]*big.Int)
	for _, term := range tokens {
		if _, ok := tfIDF.stopWords[term]; ok {
			continue
		}
		freq, isFound := m[term]
		if !isFound {
			freq = big.NewInt(1)
			m[term] = freq
		} else {
			m[term] = freq.Add(freq, big.NewInt(1))
		}
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
	idf := math.Log(float64(1+N) / (float64(1 + termDocs)))
	// idf := (float64(1+N) / float64(1+termDocs))
	return tf * idf
}

func (tfIDF *TFIDF) BM25(docPos int, term string) *big.Float {
	termFreq := tfIDF.termFreqs[docPos-1][term]
	termDocCount := tfIDF.termDocsCount[term]
	docLen := tfIDF.docsTermLen[docPos-1]
	docsLen := tfIDF.termsLen
	IDF := math.Log(1 + ((float64(tfIDF.n) - float64(termDocCount) + 0.5) / (float64(termDocCount) + 0.5)))

	// avgDocsLen := docsLen.Q(docsLen, big.NewInt(int64(tfIDF.n)))
	f := new(big.Float)
	fTermFreq := new(big.Float)
	fTermFreq = fTermFreq.SetInt(termFreq)
	fNumerator := f.Mul(big.NewFloat(K+1), fTermFreq)
	fDocLen := big.Float{}
	fDocsLen := big.Float{}
	avgDocsLen := f.Quo(fDocLen.SetInt(docLen), fDocsLen.SetInt(docsLen))

	f1bb := f.Add(big.NewFloat(1-B), f.Mul(big.NewFloat(B), avgDocsLen))
	fk1bb := f.Mul(big.NewFloat(K), f1bb)

	fDeno := f.Add(fTermFreq, fk1bb)

	fBM25 := f.Mul(big.NewFloat(IDF), f.Quo(fNumerator, fDeno))
	return fBM25

}
