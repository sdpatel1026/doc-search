package tfidf

import "github.com/dchest/stemmer/porter2"

//Ranks finds a rank of all docs in corpus w.r.t queryWeight.
// func (tfIDF *TFIDF) Ranks(queryWeight map[string]float64) (docsRank map[int]float64) {
// 	docsRank = make(map[int]float64)
// 	for pos, weight := range tfIDF.weights {
// 		rank := Cosine(queryWeight, weight)
// 		docsRank[pos+1] = rank
// 	}
// 	return
// }

//Ranks finds a rank of all docs in corpus w.r.t queryWeight using TF-IDF and cosine similarities.
// func (tfIDF *TFIDF) Ranks(queryWeight map[string]float64) (docsRank map[int]float64) {
// 	docsRank = make(map[int]float64)
// 	var pos = 1
// 	for pos <= tfIDF.n {
// 		weight := tfIDF.weight(pos)
// 		// fmt.Printf("\n-----------%s---------\n", tfIDF.indexDocName[pos])
// 		// fmt.Printf("weight: %v\n", weight)
// 		rank := Cosine(queryWeight, weight)
// 		docsRank[pos] = rank
// 		pos++
// 	}
// 	return
// }

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
