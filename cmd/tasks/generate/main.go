package main

import (
	"DeutschBot/internal"
	"DeutschBot/internal/bot/learn"
	"DeutschBot/package/translator"
	"fmt"
	"io"
	"os"
	"strings"
)

type details struct {
	Meaning []string
}

// todo move from controller to service
func main() {
	projectDir, _ := os.Getwd()
	dictionaryFile, err := os.Open(projectDir + "/bin/top-words.json")
	if err != nil {
		panic(err)
	}
	defer dictionaryFile.Close()

	content, _ := io.ReadAll(dictionaryFile)

	var words map[string]details
	internal.Deserialize(content, &words)

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
	var examples []learn.Example
	if len(translation.Examples()) > 0 {
		for _, example := range translation.Examples() {
			examples = append(examples, learn.Example{
				Usage:   example.Meaning(),
				Meaning: example.Usage(),
			})

			if len(examples) > 2 {
				break
			}
		}
	}

	learn.CreateTask(
		fmt.Sprintf("What does \"%s\" mean?", translation.Word()),
		d.Meaning,
		examples,
	)
}

func createTaskType2(translation translator.Translation, d details) {
	var examples []learn.Example
	if len(translation.Examples()) > 0 {
		for _, example := range translation.Examples() {
			examples = append(examples, learn.Example{
				Usage:   example.Usage(),
				Meaning: example.Meaning(),
			})

			if len(examples) > 2 {
				break
			}
		}
	}

	learn.CreateTask(
		fmt.Sprintf("How do you say \"%s\"?", strings.Join(d.Meaning, "; ")),
		[]string{translation.Word()},
		examples,
	)
}
