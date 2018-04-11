package nameday

import (
	"testing"

	"github.com/stretchr/testify/assert"

	natural "github.com/skrawler/go-natural"
)

func TestStrBetween(t *testing.T) {
	x, err := StrBetween(`$START$hej$END$`, `$START$`, `$END$`)
	assert.Equal(t, nil, err)
	assert.Equal(t, "hej", x)
}

func TestSweNamesOnDate(t *testing.T) {
	list := SweNamesOnDate(12, 15)
	assert.Equal(t, []Nameday{
		{Name: "Gottfrid", Date: natural.NewMonthDay("12-15"), Official: true},
		{Name: "Leon", Date: natural.NewMonthDay("12-15"), Official: false},
		{Name: "Leona", Date: natural.NewMonthDay("12-15"), Official: false},
		{Name: "Levina", Date: natural.NewMonthDay("12-15"), Official: false},
	}, list)
}

func TestSweNamedayFor(t *testing.T) {
	nd, err := SweNamedayFor("Martin")
	assert.Equal(t, nil, err)
	assert.Equal(t, Nameday{Name: "Martin", Date: natural.NewMonthDay("11-10"), Official: true}, nd)
}
