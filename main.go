package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// use task taskConfigPath
// or  taskpool taskConfigPoolPath
func main() {
	args := os.Args
	if len(args) <= 1 {
		log.Fatal("len(args) <= 1")
	}
	switch args[1] {
	case "task":
		doTask(args)
	case "taskpool":
		doTaskPool(args)
	default:
		log.Fatal("please use task or taskpool")
	}
}

func readTaskConfig(configPath string) *TaskConfig {
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal("read task config error: ", err)
	}
	var config TaskConfig
	err = json.Unmarshal(configData, &config)
	if err != nil {
		log.Fatal("unmarshal task config error: ", err)
	}
	return &config
}

func doTask(args []string) {
	if len(args) <= 2 {
		log.Fatal("len(args) <= 2")
	}
	err := RunTask(readTaskConfig(args[2]))
	if err != nil {
		log.Fatal("task error: ", err)
	}
}

func readTaskPoolConfig(configPath string) *TaskPoolConfig {
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal("read task pool config error: ", err)
	}
	var config TaskPoolConfig
	err = json.Unmarshal(configData, &config)
	if err != nil {
		log.Fatal("unmarshal task pool config error: ", err)
	}
	return &config
}

func doTaskPool(args []string) {
	if len(args) <= 2 {
		log.Fatal("len(args) <= 2")
	}
	err := RunTaskPool(readTaskPoolConfig(args[2]))
	if err != nil {
		log.Fatal("task error: ", err)
	}
}