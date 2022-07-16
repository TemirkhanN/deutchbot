package internal

import (
	"encoding/json"
	"log"
)

func Serialize(data interface{}) []byte {
	result, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func Deserialize(serialized []byte, into interface{}) {
	err := json.Unmarshal(serialized, into)
	if err != nil {
		log.Fatal(err)
	}
}
