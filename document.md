# BM-25
## Formula
 $$BM25 =  \sum IDF(Qi){{f(Qi,D)*(K+1)\over f(Qi,D)+K*(1-B+B*{fieldLen \over avgFieldLen})}}$$

 </br>

 $$IDF = ln(1+{N-Df(Qi)+0.5 \over Df(Qi)+0.5})$$
      Qi =          ith query term.
      IDF(Qi)=      inverse document frequency of the ith query term.
      f(Qi,D) =     how many times does the ith query term occur in  document D. 
      K =           constant which helps determine term frequency saturation characteristics.
      B =           constant that controls the effects of the length of the document compared to the average length.
      fieldLen =    no. of terms in document.
      avgFieldLen = average of fieldLen of documents in corpus.
      N =           total documents in corpus  
      Df(Qi) =      no. document in which ith query term appears.

</br>
</br>
</br>


# TF-IDF

## Formula
   $$tf.idf = tf * log ({1+N\over1+Df})$$

      tf = frequency of term in document  
      N = total documents in corpus  
      Df = no. document in which term appear


## Draw backs of TF-IDF
  1. ***Zero value issue:-*** 
     If word of interest appear in all document its tf-idf value reduce to zero,  even if   
     it only appears once in each of them. 


  2. ***Extensive margin issue:-***
     This arises because the inverse document frequency portion of tf-idf does not take into account how often a word appears in other documents (the intensive margin), just whether it appears at all (the extensive margin). This causes it to be overly sensitive when there are changes on the extensive margin and overly resistant to change when there are changes on the intensive margin, even though the latter is arguably more important.

     ***Example:***
     consider the following scenario. Imagine we have three documents: A, B, and C.
     Imagine apple making up 10% of words in document A. If apple did not appear at all in document B and C, the tf-idf value would be relatively high. However, if it were to appear just a single time in B, the value would plummet dramatically (though not quite to zero so long as it doesn’t appear in C). Despite the two scenarios being almost identical, with the only change being one more outside appearance of apple, the value changes dramatically, all because that change came on the extensive margin (i.e. apple went from being in no outside documents to one). Of course, in practical terms there is little difference between the word apple appearing just once or not at all outside of document A. Yet tf-idf would suggest that the two cases are radically different.

      On the flip side, the tf-idf would not change at all if apple were to go from appearing just once in document B to making up literally 100% of the document’s words. That is, it doesn’t matter how many times the word apple appears in document B, because tf-idf ignores the intensive margin. It only matters whether it appears at all. However, if our goal is to get a measure that indicates relative importance, it is quite critical to distinguish between these scenarios. In the former, apple is relatively important to document A compared to B (only showing up once in B), while in the latter apple is not at all uniquely important for document A (at 10% of its words) relative to how important it is for document B (100% of its words). Yet tf-idf sees no difference.
   3. It doesn't take document in account: 




