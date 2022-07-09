package helper

import (
	"github.com/diogoalbuquerque/migration-customers/internal/entity"
)

func SplitBetweenAvailableBuckets(legacyPeople []entity.LegacyPerson, availableBuckets int) [][]entity.LegacyPerson {

	var separatedBuckets [][]entity.LegacyPerson
	defaultSize := len(legacyPeople) / availableBuckets
	maxBuckets := len(legacyPeople) - defaultSize*availableBuckets

	size := defaultSize + 1
	for i, idx := 0, 0; i < availableBuckets; i++ {
		if i == maxBuckets {
			size--
			if size == 0 {
				break
			}
		}
		separatedBuckets = append(separatedBuckets, legacyPeople[idx:idx+size])
		idx += size
	}
	return separatedBuckets
}
