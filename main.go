package main

import (
	"github.com/qaqcatz/impomysql/task"
	"github.com/qaqcatz/impomysql/tool/affversion"
	"github.com/qaqcatz/impomysql/tool/ckstable"
	"github.com/qaqcatz/impomysql/tool/sqlsim"
	"log"
	"os"
	"strconv"
)

// todo: use urfave/cli
// use task taskConfigPath
// or  taskpool taskPoolConfigPath
// or  ckstable task taskConfigPath execNum
// or  ckstable taskpool taskPoolConfigPath threadNum execNum
// or  sqlsim task taskConfigPath
// or  sqlsim taskpool taskPoolConfigPath threadNum
// or  affversion task taskConfigPath port version [whereVersionStatus]
// or  affversion taskpool taskPoolConfigPath threadNum port version [whereVersionStatus]
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

func doCKStable(args []string) {
	if len(args) <= 2 {
		log.Fatal("[doCKStable]len(args) <= 2")
	}
	switch args[2] {
	case "task":
		// ckstable task taskConfigPath execNum
		if len(args) <= 4 {
			log.Fatal("[doCKStable]len(args) <= 4")
		}
		execNum, err := strconv.Atoi(args[4])
		if err != nil {
			log.Fatal("[doCKStable]parse execNum error: ", err)
		}

		taskConfig, err := task.NewTaskConfig(args[3])
		if err != nil {
			log.Fatal("[doCKStable]new task config error: ", err)
		}
		err = ckstable.CheckStableTask(taskConfig, nil, execNum)
		if err != nil {
			log.Fatal("[doCKStable]ckstable task error: ", err)
		}
	case "taskpool":
		// ckstable taskpool taskPoolConfigPath threadNum execNum
		if len(args) <= 5 {
			log.Fatal("[doCKStable]len(args) <= 5")
		}
		threadNum, err := strconv.Atoi(args[4])
		if err != nil || threadNum <= 0 {
			log.Fatal("parse threadNum error")
		}
		execNum, err := strconv.Atoi(args[5])
		if err != nil {
			log.Fatal("[doCKStable]parse execNum error: ", err)
		}

		taskPoolConfig, err := task.NewTaskPoolConfig(args[3])
		if err != nil {
			log.Fatal("[doCKStable]new task pool config error: ", err)
		}
		err = ckstable.CheckStableTaskPool(taskPoolConfig, threadNum, execNum)
		if err != nil {
			log.Fatal("[doCKStable]ckstable task pool error: ", err)
		}
	default:
		log.Fatal("[doCKStable]please use task, taskpool")
	}
}

func doSqlSim(args []string) {
	if len(args) <= 2 {
		log.Fatal("[doSqlSim]len(args) <= 2")
	}
	switch args[2] {
	case "task":
		// sqlsim task taskConfigPath
		if len(args) <= 3 {
			log.Fatal("[doSqlSim]len(args) <= 3")
		}
		taskConfig, err := task.NewTaskConfig(args[3])
		if err != nil {
			log.Fatal("[doSqlSim]new task config error: ", err)
		}
		err = sqlsim.SqlSimTask(taskConfig, nil)
		if err != nil {
			log.Fatal("[doSqlSim]sqlsim task error: ", err)
		}
	case "taskpool":
		// sqlsim taskpool taskPoolConfigPath threadNum
		if len(args) <= 4 {
			log.Fatal("[doSqlSim]len(args) <= 4")
		}
		threadNum, err := strconv.Atoi(args[4])
		if err != nil || threadNum <= 0 {
			log.Fatal("parse threadNum error")
		}
		taskPoolConfig, err := task.NewTaskPoolConfig(args[3])
		if err != nil {
			log.Fatal("[doSqlSim]new task pool config error: ", err)
		}
		err = sqlsim.SqlSimTaskPool(taskPoolConfig, threadNum)
		if err != nil {
			log.Fatal("[doSqlSim]sqlsim task pool error: ", err)
		}
	default:
		log.Fatal("[doSqlSim]please use task, taskpool")
	}
}

func doAffVersion(args []string) {
	if len(args) <= 2 {
		log.Fatal("[doAffVersion]len(args) <= 2")
	}
	switch args[2] {
	case "task":
		// affversion task taskConfigPath port version [whereVersionStatus]
		if len(args) <= 5 {
			log.Fatal("[doAffVersion]len(args) <= 5")
		}
		port, err := strconv.Atoi(args[4])
		if err != nil || port <= 0 {
			log.Fatal("[doAffVersion]parse port error")
		}
		version := args[5]
		whereVersionStatus := ""
		if len(args) > 6 {
			whereVersionStatus = args[6]
		}

		taskConfig, err := task.NewTaskConfig(args[3])
		if err != nil {
			log.Fatal("[doAffVersion]new task config error: ", err)
		}
		err = affversion.AffVersionTask(taskConfig, nil, port, version, whereVersionStatus)
		if err != nil {
			log.Fatal("[doAffVersion]affversion task error: ", err)
		}
	case "taskpool":
		// affversion taskpool taskPoolConfigPath threadNum port version [whereVersionStatus]
		if len(args) <= 6 {
			log.Fatal("[doAffVersion]len(args) <= 6")
		}
		threadNum, err := strconv.Atoi(args[4])
		if err != nil || threadNum <= 0 {
			log.Fatal("parse threadNum error")
		}
		port, err := strconv.Atoi(args[5])
		if err != nil || port <= 0 {
			log.Fatal("[doAffVersion]parse port error")
		}
		version := args[6]
		whereVersionStatus := ""
		if len(args) > 7 {
			whereVersionStatus = args[7]
		}

		taskPoolConfig, err := task.NewTaskPoolConfig(args[3])
		if err != nil {
			log.Fatal("[doSqlSim]new task pool config error: ", err)
		}
		err = affversion.AffVersionTaskPool(taskPoolConfig, threadNum, port, version, whereVersionStatus)
		if err != nil {
			log.Fatal("[doAffVersion]affversion task pool error: ", err)
		}
	default:
		log.Fatal("[doAffVersion]please use task, taskpool")
	}
}