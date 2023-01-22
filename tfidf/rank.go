package tfidf

//Ranks finds a rank of all docs in corpus w.r.t queryWeight.
// func (tfIDF *TFIDF) Ranks(queryWeight map[string]float64) (docsRank map[int]float64) {
// 	docsRank = make(map[int]float64)
// 	for pos, weight := range tfIDF.weights {
// 		rank := Cosine(queryWeight, weight)
// 		docsRank[pos+1] = rank
// 	}
// 	return
// }

//Ranks finds a rank of all docs in corpus w.r.t queryWeight.
func (tfIDF *TFIDF) Ranks(queryWeight map[string]float64) (docsRank map[int]float64) {
	docsRank = make(map[int]float64)
	var pos = 1
	for pos <= tfIDF.n {
		weight := tfIDF.weight(pos)
		// fmt.Printf("\n-----------%s---------\n", tfIDF.indexDocName[pos])
		// fmt.Printf("weight: %v\n", weight)
		rank := Cosine(queryWeight, weight)
		docsRank[pos] = rank
		pos++
	}
	return
}
