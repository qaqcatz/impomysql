package main

import (
	"github.com/qaqcatz/impomysql/task"
	"github.com/qaqcatz/impomysql/tasktool/affversion"
	"log"
	"os"
	"strconv"
)

// use task taskConfigPath
// or  taskpool taskConfigPoolPath
// or  affversion dbmsOutputPath version dsn threadNum [whereVersionEQ], see tasktool.MayAffect
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
	case "affversion":
		doAffVersion(args)
	default:
		log.Fatal("please use task, taskpool, affversion")
	}
}

func doTask(args []string) {
	if len(args) <= 2 {
		log.Fatal("[doTask]len(args) <= 2")
	}
	taskConfig, err := task.NewTaskConfig(args[2])
	if err != nil {
		log.Fatal("[doTask]new task config error: ", err)
	}
	_, err = task.RunTask(taskConfig, nil, nil)
	if err != nil {
		log.Fatal("[doTask]task error: ", err)
	}
}

func doTaskPool(args []string) {
	if len(args) <= 2 {
		log.Fatal("[doTaskPool]len(args) <= 2")
	}
	taskPoolConfig, err := task.NewTaskPoolConfig(args[2])
	if err != nil {
		log.Fatal("[doTaskPool]new task pool config error: ", err)
	}
	_, err = task.RunTaskPool(taskPoolConfig)
	if err != nil {
		log.Fatal("[doTaskPool]task pool error: ", err)
	}
}

func doAffVersion(args []string) {
	if len(args) <= 5 {
		log.Fatal("[doAffVersion]len(args) <= 5")
	}
	dbmsOutputPath := args[2]
	version := args[3]
	dsn := args[4]
	threadNumStr := args[5]
	threadNum, err := strconv.Atoi(threadNumStr)
	if err != nil {
		log.Fatal("[doAffVersion]parse threadNum error")
	}
	if threadNum <= 0 {
		log.Fatal("[doAffVersion]threadNum <= 0")
	}
	whereVersionEQ := ""
	if len(args) > 6 {
		whereVersionEQ = args[6]
	}
	err = affversion.AffVersion(dbmsOutputPath, version, dsn, threadNum, whereVersionEQ)
	if err != nil {
		log.Fatal("[doAffVersion]affect version error: ", err)
	}
}