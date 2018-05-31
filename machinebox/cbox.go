package main

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/machinebox/gopherconeu/dataset"
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

	fmt.Println("Loading newsgroup on memory")
	examples, categories, err := readNewsDataset(dataset)
	if err != nil {
		fmt.Println("[ERROR] reading the dataset", err)
		return
	}

	fmt.Printf("Examples %v, Categories %v\n", len(examples), len(categories))

	model, err := createModel(classificationbox, categories)
	if err != nil {
		fmt.Println("[ERROR] creating model", err)
		return
	}
	fmt.Printf("Model ready %v:%v\n", model.ID, model.Name)

	// shuffle examples
	shuffle(examples, source)

	// split examples into train 80%, validation 20%
	split := 0.2
	v := int(split * float64(len(examples)))
	t := len(examples) - v

	train := examples[0:t]
	validation := examples[t:len(examples)]
	fmt.Printf("Train %v, Validation %v\n", len(train), len(validation))

	fmt.Printf("Training...\n")
	model.train(train)
	fmt.Printf("Finished Training\n")

	fmt.Printf("Validation...\n")
	correct, incorrect, _ := model.validate(validation)
	fmt.Printf("Finished Validation\n\n")
	acc := float64(correct) / float64(v)
	fmt.Printf("Correct %d Incorrect %d, Acc %f\n\n", correct, incorrect, acc)

	//fmt.Printf("Confusion Matrix\n")

	//printMatrix(matrix)

}

func readNewsDataset(dir string) ([]classificationbox.Example, map[string]int, error) {
	examples := []classificationbox.Example{}
	categories, err := dataset.ReadDataset(dir, func(body, class string) {
		ex := classificationbox.Example{
			Inputs: []classificationbox.Feature{{Key: "body", Type: "text", Value: body}},
			Class:  class,
		}
		examples = append(examples, ex)
	})
	return examples, categories, err
}

type Model struct {
	ID   string
	Name string
	// classificationbox client
	client *classificationbox.Client
}

func createModel(cbox *classificationbox.Client, categories map[string]int) (*Model, error) {
	classes := []string{}
	for category := range categories {
		classes = append(classes, category)
	}
	m, err := cbox.CreateModel(context.Background(), classificationbox.Model{
		Name:    "newsgroup",
		Classes: classes,
		Options: &classificationbox.ModelOptions{
			Skipgrams: 0,
			Ngrams:    0,
		},
	})
	if err != nil {
		return nil, err
	}
	return &Model{
		ID:     m.ID,
		Name:   m.Name,
		client: cbox,
	}, nil
}

func (m *Model) train(examples []classificationbox.Example) {
	fmt.Println("")
	for i, ex := range examples {
		err := m.client.Teach(context.Background(), m.ID, ex)
		if err != nil {
			fmt.Println("[WARN] Error teaching example", err, ex)
		}
		fmt.Printf("\rTeach [%v/%v]", i, len(examples))
	}
}

func (m *Model) validate(examples []classificationbox.Example) (correct int, incorrect int, matrix map[mkey]int) {
	fmt.Println("")
	correct = 0
	incorrect = 0
	// row predicted column actual
	matrix = map[mkey]int{}
	for i, ex := range examples {
		pred, err := m.client.Predict(context.Background(), m.ID, classificationbox.PredictRequest{
			Inputs: ex.Inputs,
		})
		if err != nil {
			fmt.Println("[WARN] Error teaching example", err, ex)
		}
		matrix[mkey{Predicted: pred.Classes[0].ID, Actual: ex.Class}]++
		if pred.Classes[0].ID == ex.Class {
			fmt.Printf("\rPredict [%v/%v] ok: %s", i, len(examples), ex.Class)
			correct++
		} else {
			fmt.Printf("\rPredict [%v/%v] f: Expected %s  Actual %s", i, len(examples), ex.Class, pred.Classes[0].ID)
			incorrect++
		}
	}
	return
}

func shuffle(array []classificationbox.Example, source rand.Source) {
	random := rand.New(source)
	for i := len(array) - 1; i > 0; i-- {
		j := random.Intn(i + 1)
		array[i], array[j] = array[j], array[i]
	}
}

type mkey struct {
	Predicted string
	Actual    string
}

func printMatrix(matrix map[mkey]int) {
	keys := []mkey{}
	for k, _ := range matrix {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Predicted < keys[j].Predicted
	})

	for _, k := range keys {
		if k.Predicted == k.Actual {
			fmt.Printf("OK %s -> %v\n", k.Predicted, matrix[k])
		} else {
			fmt.Printf("P:%s A:%s -> %v\n", k.Predicted, k.Actual, matrix[k])
		}

	}
}

func sanitize(s string) string {
	r := strings.NewReplacer("\n", "", "\r", "", "\t", " ", "|", " ", ":", " ", "_", " ")
	return r.Replace(s)
}
