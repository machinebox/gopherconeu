package dataset

import (
	"testing"

	"github.com/matryer/is"
)

func TestDataset(t *testing.T) {
	is := is.New(t)
	// X are normally the samples
	x := []string{}
	// Y are the expected outputs
	y := []string{}

	categories, err := ReadDataset("../20_newsgroup", func(body, category string) {
		x = append(x, body)
		y = append(y, category)
	})
	is.NoErr(err)

	// print categories
	for k, v := range categories {
		t.Logf("%v -> %v", k, v)
	}

	is.Equal(len(categories), 20)

	// inspect some news
	//t.Log(y[20], x[20])

}
