package tfidf

import "github.com/dchest/stemmer/porter2"

//RanksBM25 calculate BM25 rank of the document for the given query
func (tfIDF *TFIDF) RanksBM25(query string) (docsRank map[int64]float64) {
	tokens := tfIDF.tokenizer.Tokens(query)
	docsRank = make(map[int64]float64)
	for _, token := range tokens {
		if _, ok := tfIDF.stopWords[token]; ok {
			continue
		}
		token = porter2.Stemmer.Stem(token)
		docList, isFound := tfIDF.termDocs[token]
		if isFound {
			for docPos := range docList {
				bm25 := tfIDF.BM25(docPos, token)
				docsRank[docPos] = docsRank[docPos] + bm25
			}

		}
	}
	return docsRank
}
