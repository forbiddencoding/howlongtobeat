package howlongtobeat

import (
	"math"
	"strings"
)

/*
calculateJaccardSimilarity calculates the Jaccard Similarity between two titles.

The Jaccard Coefficient is a measure of the similarity between the two titles.
It is defined as the size of the intersection divided by the size of the union of the sample sets.

The function first converts titles to lowercase and then splits them into words (tokenizes the titles).
An intersection of these sets of words is found.
The sizes of the intersection and the unique words across both titles are calculated.
The ratio of the length of intersection to the total unique words gives the Jaccard Similarity.
This value is rounded to two decimal places for precision.

Parameters:
title1 -- First title (type: string)
title2 -- Second title (type: string)

Returns:
Jaccard Similarity of the inputted titles (type: float64).
It would return 0 if there are no unique words in both titles.

Note:
The function assumes the words in the titles are separated by spaces.
It may see words as unique that aren't really unique due to attached punctuation.
For example, "word." and "word" will be seen as different words.
*/
func calculateJaccardSimilarity(title1, title2 string) float64 {
	// Convert to lowercase
	title1 = strings.ToLower(title1)
	title2 = strings.ToLower(title2)

	var words1 = strings.Fields(title1)
	var words2 = strings.Fields(title2)

	intersection := make(map[string]bool)
	unique := make(map[string]bool)

	for _, word := range words1 {
		unique[word] = true
	}

	for _, word := range words2 {
		if _, value := unique[word]; value {
			intersection[word] = true
		}
		unique[word] = true
	}

	intersectionSize := float64(len(intersection))
	uniqueSize := float64(len(unique))

	if uniqueSize > 0 {
		result := intersectionSize / uniqueSize
		return math.Round(result*100) / 100 // Round the value to two decimal places
	}

	return 0
}
