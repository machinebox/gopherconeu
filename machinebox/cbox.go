package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/machinebox/sdk-go/boxutil"
	"github.com/machinebox/sdk-go/classificationbox"
)


// Using classification box to create and validate a classifier for the news20 problem
//
// Get MB_KEY from https://machinebox.io/
// $ export MB_KEY=your_key
// $ docker-compose up
//
// to tear down the box
// $ docker-compose down

func main() {
	dataset := "../20_newsgroup/"
	classificationbox := classificationbox.New("http://localhost:8080")
	source := rand.NewSource(time.Now().UnixNano())
	fmt.Println("Waiting for classificationbox to be ready...")
	boxutil.WaitForReady(context.Background(), classificationbox)
	fmt.Println("Done!")


	// Load all the examples in memory


	// print how many examples and categories


	// shuffle examples


	// split 80% trainning 20% validation


	// train


	// validate

	
	// print precision 


}

func shuffle(array []???????, source rand.Source) {
	random := rand.New(source)
	for i := len(array) - 1; i > 0; i-- {
		j := random.Intn(i + 1)
		array[i], array[j] = array[j], array[i]
	}
}
