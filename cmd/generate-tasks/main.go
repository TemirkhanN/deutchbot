package main

import (
	"DeutchBot/internal/database"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type details struct {
	Meaning []string
}

// todo move from controller to service
func main() {
	projectDir, _ := os.Getwd()
	dictionaryFile, err := os.Open(projectDir + "/top500words.json")
	if err != nil {
		panic(err)
	}
	defer dictionaryFile.Close()

	content, _ := ioutil.ReadAll(dictionaryFile)

	var words map[string]details

	json.Unmarshal(content, &words)

	tasks := make([]database.Task, 0)

	taskId := 0
	for word, d := range words {
		taskId++
		tasks = append(tasks, createTaskType1(taskId, word, d))
		taskId++
		reverseTasks := createTaskType2(taskId, word, d)
		tasks = append(tasks, reverseTasks...)
		//todo messy
		taskId += len(reverseTasks) - 1
	}

	rawContent, err := yaml.Marshal(tasks)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(projectDir+"/internal/database/tasks_tmp.yaml", rawContent, 0)
	if err != nil {
		panic(err)
	}

	fmt.Println("Tasks successfully regenerated")
}

func createTaskType1(id int, word string, d details) database.Task {
	return database.Task{
		Id:       id,
		Question: fmt.Sprintf("What does \"%s\" mean?", word),
		Answers:  d.Meaning,
	}
}

// will be buggy for multivalue words
func createTaskType2(id int, word string, d details) []database.Task {
	tasks := make([]database.Task, len(d.Meaning))

	for i, meaning := range d.Meaning {
		tasks[i] = database.Task{
			Id:       id,
			Question: fmt.Sprintf("How do you say \"%s\"?", meaning),
			Answers:  []string{word},
		}
		id++
	}

	return tasks
}
