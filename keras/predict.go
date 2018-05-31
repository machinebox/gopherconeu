package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/machinebox/gopherconeu/dataset"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

const (
	maxlen   = 1000  // maximum len of the sequence of words
	maxWords = 20000 // maximum number of words that the model handles
)

// function to load the word index

// bodyToVector translates a text body into a vector.
func bodyToVector(body string, wordIndex map[string]int32, dim int) []int32 {
	words := strings.Fields(body)
	vector := make([]int32, dim)
	//vector := [1000]int32{}
	if len(words) > dim {
		words = words[:dim]
	}
	offset := dim - len(words)
	for pos, w := range words {
		idx, ok := wordIndex[strings.ToLower(w)]
		if !ok {
			continue
		}
		vector[pos+offset] = int32(idx)
	}
	return vector
}

func loadIndex(indexFile string, maxIndex int) (map[string]int32, error) {
	index := map[string]int32{}
	file, err := os.Open(indexFile)
	if err != nil {
		return index, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	r.Comma = ';'
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return index, err
		}
		if len(record) != 2 {
			return index, errors.New("column mismatch on the index")
		}
		i, err := strconv.Atoi(record[1])
		if err != nil {
			return index, err
		}
		if maxIndex != 0 && i >= maxIndex {
			continue
		}
		index[record[0]] = int32(i)
	}
	return index, nil
}

func invertIndex(index map[string]int32) ([]string, error) {
	invert := make([]string, len(index))
	for k, v := range index {
		// normally the values are 0..len
		// to make it more robush it should validate that
		if int(v) >= len(invert) {
			return invert, errors.New("the index is malformed")
		}
		invert[v] = k
	}
	return invert, nil
}

func vectorToLabel(y []float32, labels []string) (string, error) {
	// get the max, activation, and map to the index
	max := float32(0.0)
	idx := 0
	for i := range y {
		if y[i] > max {
			max = y[i]
			idx = i
		}
	}
	if idx > len(labels) {
		return "", errors.New("the index for the label is out of bounds")
	}
	return labels[idx], nil

}

// function to load label index
// function to transform a string to a vector
// function to transform a vector of labels to a label probabilities

func main() {

	wordIndex, err := loadIndex("./words.csv", maxlen)
	if err != nil {
		log.Fatal("can not load the index", err)
	}

	body := "Computers are good, and Go is awesome"

	vector := bodyToVector(body, wordIndex, maxWords)

	fmt.Printf("vector %+v", vector)

	labelIndex, err := loadIndex("./labels.csv", 0)
	if err != nil {
		log.Fatal("can not load the index for labels", err)
	}

	labels, err := invertIndex(labelIndex)
	if err != nil {
		log.Fatal("can not load the invert the index for labels", err)
	}

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
	x := vector

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

	category, err := vectorToLabel(y[0], labels)
	if err != nil {
		log.Fatal("can not transform the output", err)
		return
	}

	fmt.Println("Result: ", y)

	fmt.Println("Category: ", category)

}
