package main

import (
	"github.com/qaqcatz/impomysql/task"
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

func doTask(args []string) {
	if len(args) <= 2 {
		log.Fatal("len(args) <= 2")
	}
	taskConfig, err := task.NewTaskConfig(args[2])
	if err != nil {
		log.Fatal("new task config error: ", err)
	}
	_, err = task.RunTask(taskConfig, nil, nil)
	if err != nil {
		log.Fatal("task error: ", err)
	}
}

func doTaskPool(args []string) {
	if len(args) <= 2 {
		log.Fatal("len(args) <= 2")
	}
	taskPoolConfig, err := task.NewTaskPoolConfig(args[2])
	if err != nil {
		log.Fatal("new task pool config error: ", err)
	}
	_, err = task.RunTaskPool(taskPoolConfig)
	if err != nil {
		log.Fatal("task pool error: ", err)
	}
}