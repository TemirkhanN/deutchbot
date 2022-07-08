package main

import (
	"DeutchBot/internal/learn"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type details struct {
	Meaning []string
}

// todo move from controller to service
func main() {
	projectDir, _ := os.Getwd()
	dictionaryFile, err := os.Open(projectDir + "/bin/top500words.json")
	if err != nil {
		panic(err)
	}
	defer dictionaryFile.Close()

	content, _ := io.ReadAll(dictionaryFile)

	var words map[string]details

	json.Unmarshal(content, &words)

	for word, d := range words {
		createTaskType1(word, d)
		createTaskType2(word, d)
	}

	fmt.Println("Tasks successfully regenerated")
}

func createTaskType1(word string, d details) {
	answers, err := json.Marshal(d.Meaning)
	if err != nil {
		log.Fatal(err)
	}

	task := &learn.Task{
		Question: fmt.Sprintf("What does \"%s\" mean?", word),
		Answers:  answers,
	}

	learn.TaskRepository.Save(task)
}

// will be buggy for multivalue words
func createTaskType2(word string, d details) {
	for _, meaning := range d.Meaning {
		answers, err := json.Marshal([]string{word})
		if err != nil {
			log.Fatal(err)
		}
		task := learn.Task{
			Question: fmt.Sprintf("How do you say \"%s\"?", meaning),
			Answers:  answers,
		}

		learn.TaskRepository.Save(&task)
	}
}
