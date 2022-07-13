package main

import (
	"DeutschBot/internal"
	"DeutschBot/internal/bot/learn"
	"DeutschBot/package/translator"
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
		translation, err := internal.Translator.Translate(word, translator.DE, translator.EN)
		if err != nil {
			fmt.Println(err)
			continue
		}

		createTaskType1(translation, d)
		createTaskType2(translation, d)
	}

	fmt.Println("Tasks successfully regenerated")
}

func createTaskType1(translation translator.Translation, d details) {
	answers, err := json.Marshal(d.Meaning)
	if err != nil {
		log.Fatal(err)
	}

	task := &learn.Task{
		Question: fmt.Sprintf("What does \"%s\" mean?", translation.Word()),
		Answers:  answers,
	}

	if len(translation.Examples()) > 0 {
		var examples []learn.Example
		for _, example := range translation.Examples() {
			examples = append(examples, learn.Example{
				Usage:   example.Meaning(),
				Meaning: example.Usage(),
			})

			if len(examples) > 2 {
				break
			}
		}
		task.SetExamples(examples)
	}

	learn.TaskRepository.Save(task)
}

// will be buggy for multivalue words
func createTaskType2(translation translator.Translation, d details) {
	for _, meaning := range d.Meaning {
		answers, err := json.Marshal([]string{translation.Word()})
		if err != nil {
			log.Fatal(err)
		}
		task := learn.Task{
			Question: fmt.Sprintf("How do you say \"%s\"?", meaning),
			Answers:  answers,
		}

		if len(translation.Examples()) > 0 {
			var examples []learn.Example
			for _, example := range translation.Examples() {
				examples = append(examples, learn.Example{
					Usage:   example.Usage(),
					Meaning: example.Meaning(),
				})

				if len(examples) > 2 {
					break
				}
			}
			task.SetExamples(examples)
		}

		learn.TaskRepository.Save(&task)
	}
}
