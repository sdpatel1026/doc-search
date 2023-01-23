package tfidf

import (
	"crypto/md5"
	"encoding/hex"
	"math"

	"github.com/dchest/stemmer/porter2"
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
	docIndex      map[string]int64          // train document index
	indexDocName  map[int64]string          // train document name mapped with index
	termFreqs     []map[string]uint64       // terms frequency for each train document
	termDocsCount map[string]int64          // number of document in which term t appears in train data
	termDocs      map[string]map[int64]bool //list of doc in which term t present.
	docsTermLen   []uint64                  //len of each doc in words.
	termsLen      uint64                    //total terms in corpus
	n             int64                     // number of documents in train data
	stopWords     map[string]interface{}    // words to be remove.
	tokenizer     tokenize.Tokenizer
}

// New new model with default
func New() *TFIDF {

	if tfIdf != nil {
		return tfIdf
	}
	tfIdf = &TFIDF{
		docIndex:      make(map[string]int64),
		indexDocName:  make(map[int64]string),
		termFreqs:     make([]map[string]uint64, 0),
		termDocs:      make(map[string]map[int64]bool),
		termDocsCount: make(map[string]int64),
		docsTermLen:   make([]uint64, 0),
		termsLen:      0,
		n:             0,
		tokenizer:     &tokenize.EnTokenizer{},
	}
	stopwords := []string{`i`, `me`, `my`, `myself`, `we`, `our`, `ours`, `ourselves`, `you`, "you`re", "you`ve", "you`ll", "you`d", `your`, `yours`, `yourself`, `yourselves`, `he`, `him`, `his`, `himself`, `she`, "she`s", `her`, `hers`, `herself`, `it`, "it`s", `its`, `itself`, `they`, `them`, `their`, `theirs`, `themselves`, `what`, `which`, `who`, `whom`, `this`, `that`, "that`ll", `these`, `those`, `am`, `is`, `are`, `was`, `were`, `be`, `been`, `being`, `have`, `has`, `had`, `having`, `do`, `does`, `did`, `doing`, `a`, `an`, `the`, `and`, `but`, `if`, `or`, `because`, `as`, `until`, `while`, `of`, `at`, `by`, `for`, `with`, `about`, `against`, `between`, `into`, `through`, `during`, `before`, `after`, `above`, `below`, `to`, `from`, `up`, `down`, `in`, `out`, `on`, `off`, `over`, `under`, `again`, `further`, `then`, `once`, `here`, `there`, `when`, `where`, `why`, `how`, `all`, `any`, `both`, `each`, `few`, `more`, `most`, `other`, `some`, `such`, `no`, `nor`, `not`, `only`, `own`, `same`, `so`, `than`, `too`, `very`, `s`, `t`, `can`, `will`, `just`, `don`, "don`t", `should`, "should`ve", `now`, `d`, `ll`, `m`, `o`, `re`, `ve`, `y`, `ain`, `aren`, "aren`t", `couldn`, "couldn`t", `didn`, "didn`t", `doesn`, "doesn`t", `hadn`, "hadn`t", `hasn`, "hasn`t", `haven`, "haven`t", `isn`, "isn`t", `ma`, `mightn`, "mightn`t", `mustn`, "mustn`t", `needn`, "needn`t", `shan`, "shan`t", `shouldn`, "shouldn`t", `wasn`, "wasn`t", `weren`, "weren`t", `won`, "won`t", `wouldn`, "wouldn`t"}
	tfIdf.AddStopWords(stopwords...)
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
		tfIDF.n += 1
		var docTokenCount uint64 = 0
		for term, freq := range termFreq {
			termDocSet, isFound := tfIDF.termDocs[term]
			if !isFound {
				termDocSet = make(map[int64]bool)
				// tfIDF.termDocs[term] = termDocSet
			}
			docTokenCount += freq
			termDocSet[tfIDF.n] = true
			tfIDF.termDocs[term] = termDocSet
			tfIDF.termDocsCount[term] += 1
		}
		tfIDF.termFreqs = append(tfIDF.termFreqs, termFreq)
		tfIDF.docsTermLen = append(tfIDF.docsTermLen, docTokenCount)
		tfIDF.termsLen = tfIDF.termsLen + docTokenCount
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
//
// DocName return document name
func (tfIDF *TFIDF) DocName(docPos int64) string {
	return tfIDF.indexDocName[docPos]
}

// weight calculate tf-idf weight of doc in the corpus.
// func (tfIDF *TFIDF) weight(docPos int) map[string]float64 {

// 	weight := make(map[string]float64)
// 	termFreq := tfIDF.termFreqs[docPos-1]
// 	// fmt.Println("***************************************")
// 	for term, freq := range termFreq {
// 		// fmt.Printf("term: %v\n", term)
// 		// fmt.Printf("freq: %v\n", freq)
// 		// fmt.Printf("tfIDF.termDocs[term]: %v\n", tfIDF.termDocs[term])
// 		// fmt.Printf("tfIDF.n: %v\n", tfIDF.n)
// 		weight[term] = findTfIdf(freq, 1, tfIDF.termDocsCount[term], tfIDF.n)
// 		// weight[term] = findTfIdf(freq, docTerms, f.termDocs[term], f.n)
// 	}
// 	// fmt.Printf("termFreq: %v\n", termFreq)
// 	return weight

// }

// Cal calculate tf-idf weight for specified document
// func (tfIDF *TFIDF) Cal(doc string) (weight map[string]float64) {
// 	weight = make(map[string]float64)

// 	var termFreq map[string]int
// 	doc = strings.ToLower(doc)
// 	docPos := tfIDF.docPos(doc)
// 	if docPos < 0 {
// 		termFreq = tfIDF.termFreq(doc)
// 	} else {
// 		termFreq = tfIDF.termFreqs[docPos-1]
// 	}
// 	// docTerms := 0
// 	// for _, freq := range termFreq {
// 	// 	docTerms += freq
// 	// }
// 	for term, freq := range termFreq {
// 		weight[term] = findTfIdf(freq, 1, tfIDF.termDocsCount[term], tfIDF.n)
// 		// weight[term] = findTfIdf(freq, docTerms, f.termDocs[term], f.n)
// 	}

// 	return weight
// }

// termFreq calculate term-freq of each term in document.
func (tfIDF *TFIDF) termFreq(tokens []string) (m map[string]uint64) {
	m = make(map[string]uint64)
	for _, term := range tokens {
		if _, ok := tfIDF.stopWords[term]; ok {
			continue
		}
		term = porter2.Stemmer.Stem(term)
		freq, isFound := m[term]
		if !isFound {
			freq = 1
			m[term] = freq
		} else {
			m[term] = freq + 1
		}
	}
	return
}

// docHashPos return the position of doc in corpus.
func (tfIDF *TFIDF) docHashPos(hash string) int64 {
	if pos, ok := tfIDF.docIndex[hash]; ok {
		return pos
	}

	return -1
}

// docPos return the position of doc in corpus.
func (tfIDF *TFIDF) docPos(doc string) int64 {
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

// BM25 calculate BM25 score of the term.
func (tfIDF *TFIDF) BM25(docPos int64, term string) float64 {
	termFreq := tfIDF.termFreqs[docPos-1][term]
	termDocCount := tfIDF.termDocsCount[term]
	docLen := tfIDF.docsTermLen[docPos-1]
	docsLen := tfIDF.termsLen
	IDF := math.Log(1 + ((float64(tfIDF.n) - float64(termDocCount) + 0.5) / (float64(termDocCount) + 0.5)))

	avgDocsLen := float64(docsLen) / float64(tfIDF.n)
	numerator := float64(termFreq) * (K + 1)
	deno := float64(termFreq) + K*(1-B+B*(float64(docLen)/avgDocsLen))
	bm25 := IDF * (numerator / deno)
	return bm25

}
