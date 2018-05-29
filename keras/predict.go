package main

import (
	"fmt"
	"log"

	dataset "github.com/machinebox/gopherconeu"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

const (
	maxlen   = 1000  // maximum len of the sequence of words
	maxWords = 20000 // maximum number of words that the model handles
)

// function to load the word index
// function to load label index
// function to transform a string to a vector
// function to transform a vector of labels to a label probabilities

func main() {

	// load the model with the name and the tags
	model, err := tf.LoadSavedModel("newsmodelkeras", []string{"newsmodelkerasTag"}, nil)

	if err != nil {
		fmt.Printf("Error loading saved model: %s\n", err.Error())
		return
	}

	defer model.Session.Close()

	_, err = dataset.ReadDataset("../20_newsgroup", func(body, label string) {
		// read some examples for dataset
	})
	if err != nil {
		log.Fatal("can not read the dataset", err)
	}

	// dummy input
	x := []int32{0, 1, 0, 1}

	tensor, _ := tf.NewTensor([1][]int32{0: x})

	// fmt.Println(tensor.Shape()) // prints the shape of the tensor

	result, err := model.Session.Run(
		map[tf.Output]*tf.Tensor{
			// Use the input layer that we named
			model.Graph.Operation("news_input_layer").Output(0): tensor,
		},
		[]tf.Output{
			// Use the output layer that we named
			model.Graph.Operation("news_output_layer/Softmax").Output(0),
		},
		nil,
	)
	if err != nil {
		fmt.Printf("Error running the session with input, err: %s\n", err.Error())
		return
	}

	y := result[0].Value().([][]float32)

	fmt.Println("Result: ", y)

}
