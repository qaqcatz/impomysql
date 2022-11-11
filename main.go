package main

import (
	"encoding/json"
	"github.com/qaqcatz/impomysql/task"
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

func readTaskConfig(configPath string) *task.TaskConfig {
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal("read task config error: ", err)
	}
	var config task.TaskConfig
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
	err := task.RunTask(readTaskConfig(args[2]), nil)
	if err != nil {
		log.Fatal("task error: ", err)
	}
}

func readTaskPoolConfig(configPath string) *task.TaskPoolConfig {
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal("read task pool config error: ", err)
	}
	var config task.TaskPoolConfig
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
	err := task.RunTaskPool(readTaskPoolConfig(args[2]))
	if err != nil {
		log.Fatal("task error: ", err)
	}
}