package tfidf

//Ranks finds a rank of all docs in corpus w.r.t queryWeight.
func (tfIDF *TFIDF) Ranks(queryWeight map[string]float64) (docsRank map[int]float64) {
	docsRank = make(map[int]float64)
	for pos, weight := range tfIDF.weights {
		rank := Cosine(queryWeight, weight)
		docsRank[pos+1] = rank
	}
	return
}
