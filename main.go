package main

import (
	"github.com/qaqcatz/impomysql/task"
	"github.com/qaqcatz/impomysql/tasktool/affversion"
	"github.com/qaqcatz/impomysql/tasktool/ckstable"
	"github.com/qaqcatz/impomysql/tasktool/sqlsim"
	"log"
	"os"
	"strconv"
)

// todo: use urfave/cli
// use task taskConfigPath
// or  taskpool taskPoolConfigPath
// or  ckstable task taskConfigPath execNum
// or  ckstable taskpool taskPoolConfigPath execNum
// or  sqlsim task taskConfigPath
// or  sqlsim taskpool taskPoolConfigPath
// or  affversion task taskConfigPath port version [whereVersionEQ]
// or  affversion taskpool taskPoolConfigPath port version [whereVersionEQ]
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
	case "ckstable":
		doCKStable(args)
	case "sqlsim":
		doSqlSim(args)
	case "affversion":
		doAffVersion(args)
	default:
		log.Fatal("[main]please use task, taskpool, ckstable, sqlsim, affversion")
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

func doCKStable(args []string) {
	if len(args) <= 4 {
		log.Fatal("[doCKStable]len(args) <= 4")
	}
	execNum, err := strconv.Atoi(args[4])
	if err != nil {
		log.Fatal("[doCKStable]parse execNum error: ", err)
	}
	switch args[2] {
	case "task":
		taskConfig, err := task.NewTaskConfig(args[3])
		if err != nil {
			log.Fatal("[doCKStable]new task config error: ", err)
		}
		err = ckstable.CheckStableTask(taskConfig, nil, execNum)
		if err != nil {
			log.Fatal("[doCKStable]ckstable task error: ", err)
		}
	case "taskpool":
		taskPoolConfig, err := task.NewTaskPoolConfig(args[3])
		if err != nil {
			log.Fatal("[doCKStable]new task pool config error: ", err)
		}
		err = ckstable.CheckStableTaskPool(taskPoolConfig, execNum)
		if err != nil {
			log.Fatal("[doCKStable]ckstable task pool error: ", err)
		}
	default:
		log.Fatal("[doCKStable]please use task, taskpool")
	}
}

func doAffVersion(args []string) {
	if len(args) <= 5 {
		log.Fatal("[doAffVersion]len(args) <= 5")
	}
	port, err := strconv.Atoi(args[4])
	if err != nil {
		log.Fatal("[doAffVersion]parse port error: ", err)
	}
	if port <= 0 {
		log.Fatal("[doAffVersion]port <= 0")
	}
	version := args[5]
	whereVersionEQ := ""
	if len(args) > 6 {
		whereVersionEQ = args[6]
	}
	switch args[2] {
	case "task":
		taskConfig, err := task.NewTaskConfig(args[3])
		if err != nil {
			log.Fatal("[doAffVersion]new task config error: ", err)
		}
		err = affversion.AffVersionTask(taskConfig, nil, port, version, whereVersionEQ)
		if err != nil {
			log.Fatal("[doAffVersion]affversion task error: ", err)
		}
	case "taskpool":
		taskPoolConfig, err := task.NewTaskPoolConfig(args[3])
		if err != nil {
			log.Fatal("[doSqlSim]new task pool config error: ", err)
		}
		err = affversion.AffVersionTaskPool(taskPoolConfig, port, version, whereVersionEQ)
		if err != nil {
			log.Fatal("[doAffVersion]affversion task pool error: ", err)
		}
	default:
		log.Fatal("[doAffVersion]please use task, taskpool")
	}
}