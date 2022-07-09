package helper_test

import (
	"github.com/diogoalbuquerque/migration-customers/internal/entity"
	"github.com/diogoalbuquerque/migration-customers/internal/helper"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_SplitBetweenAvailableBuckets_SameSize(t *testing.T) {
	lp := createLegacyPerson()
	lc := createLegacyCompany()
	buckets := helper.SplitBetweenAvailableBuckets([]entity.LegacyPerson{lp, lc}, 2)
	assert.Contains(t, buckets[0], lp, "This object should be contained in the slice.")
	assert.Contains(t, buckets[1], lc, "This object should be contained in the slice.")
}

func Test_SplitBetweenAvailableBuckets_DifferentSize(t *testing.T) {
	lp := createLegacyPerson()
	lc := createLegacyCompany()
	buckets := helper.SplitBetweenAvailableBuckets([]entity.LegacyPerson{lp, lc}, 3)
	assert.Contains(t, buckets[0], lp, "This object should be contained in the slice.")
	assert.Contains(t, buckets[1], lc, "This object should be contained in the slice.")
}

func createLegacyPerson() entity.LegacyPerson {
	birthDate := time.Now()
	return entity.LegacyPerson{NI: 88626159226, Name: "User One", BirthDate: &birthDate}
}

func createLegacyCompany() entity.LegacyPerson {
	return entity.LegacyPerson{NI: 87511532000171, Name: "Company One"}
}
