package main

import (
	"github.com/qaqcatz/impomysql/task"
	"github.com/qaqcatz/impomysql/tasktool/affversion"
	"github.com/qaqcatz/impomysql/tasktool/sqlsim"
	"log"
	"os"
)

// todo: use urfave/cli
// use task taskConfigPath
// or  taskpool taskPoolConfigPath
// or  sqlsim task taskConfigPath
// or  sqlsim taskpool taskPoolConfigPath
// or  affversion task taskConfigPath version [whereVersionEQ]
// or  affversion taskpool taskPoolConfigPath version [whereVersionEQ]
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
	case "sqlsim":
		doSqlSim(args)
	default:
		log.Fatal("[main]please use task, taskpool, affversion, sqlsim")
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
	if len(args) <= 4 {
		log.Fatal("[doAffVersion]len(args) <= 4")
	}
	version := args[4]
	whereVersionEQ := ""
	if len(args) > 5 {
		whereVersionEQ = args[5]
	}
	switch args[2] {
	case "task":
		taskConfig, err := task.NewTaskConfig(args[3])
		if err != nil {
			log.Fatal("[doAffVersion]new task config error: ", err)
		}
		err = affversion.AffVersionTask(taskConfig, nil, version, whereVersionEQ)
		if err != nil {
			log.Fatal("[doAffVersion]affversion task error: ", err)
		}
	case "taskpool":
		taskPoolConfig, err := task.NewTaskPoolConfig(args[3])
		if err != nil {
			log.Fatal("[doSqlSim]new task pool config error: ", err)
		}
		err = affversion.AffVersionTaskPool(taskPoolConfig, version, whereVersionEQ)
		if err != nil {
			log.Fatal("[doAffVersion]affversion task pool error: ", err)
		}
	default:
		log.Fatal("[doAffVersion]please use task, taskpool")
	}
}

func doSqlSim(args []string) {
	if len(args) <= 3 {
		log.Fatal("[doSqlSim]len(args) <= 3")
	}
	switch args[2] {
	case "task":
		taskConfig, err := task.NewTaskConfig(args[3])
		if err != nil {
			log.Fatal("[doSqlSim]new task config error: ", err)
		}
		err = sqlsim.SqlSimTask(taskConfig, nil)
		if err != nil {
			log.Fatal("[doSqlSim]sqlsim task error: ", err)
		}
	case "taskpool":
		taskPoolConfig, err := task.NewTaskPoolConfig(args[3])
		if err != nil {
			log.Fatal("[doSqlSim]new task pool config error: ", err)
		}
		err = sqlsim.SqlSimTaskPool(taskPoolConfig)
		if err != nil {
			log.Fatal("[doSqlSim]sqlsim task pool error: ", err)
		}
	default:
		log.Fatal("[doSqlSim]please use task, taskpool")
	}
}