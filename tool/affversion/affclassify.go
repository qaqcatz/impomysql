package affversion

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/qaqcatz/impomysql/task"
	"github.com/qaqcatz/nanoshlib"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type VListPair struct{
	O1V string `json:"o1v"`
	BugList []string `json:"bugList"`
}

// AffClassify: classify bugs according to the versions they affect.
//
// Specifically, for each bug,
// we will calculate the oldest reproducible version `o1v` and use it for classification
// if the bug can not be reproduced on the previous version of `o1v` (and no error)
//
// Make sure you have done `affversion` or `affdbdeployer`, we will query the database `affversion.db`.
// You also need to provide `dbdeployer`, which will tell us the order of each version
//
// So the command is:
//   `./impomysql affclassify dbDeployerPath dbJsonPath taskPoolConfigPath`
//
// We will create `affclassify.json` in taskPoolPath. It is an array of {`o1v`, bug list}.
//
// We will also create a directory `affclassify` in taskPoolPath.
// For each `o1v`, we will save the first detected bug in `affclassify`
func AffClassify(dbDeployerPath string, dbJsonPath string, config *task.TaskPoolConfig) {
	dbDeployerAbsPath, err := filepath.Abs(dbDeployerPath)
	if err != nil {
		panic("[AffClassify]path abs error: " + err.Error())
	}
	dbJsonAbsPath, err := filepath.Abs(dbJsonPath)
	if err != nil {
		panic("[AffClassify]path abs error: " + err.Error())
	}

	// check sqlite database affversion
	taskPoolPath := config.GetTaskPoolPath()
	affVersionDBPath := path.Join(taskPoolPath, "affversion.db")
	exists, err := pathExists(affVersionDBPath)
	if err != nil {
		panic("[AffClassify]check affversion.db error: " + err.Error() + ": " + taskPoolPath)
	}
	if !exists {
		panic("[AffClassify]affversion.db does not exists: " + taskPoolPath)
	}

	// open
	affVersionDB, err := sql.Open("sqlite3", affVersionDBPath)
	if err != nil {
		panic("[AffVersionTask]open database error: " + err.Error())
	}
	defer affVersionDB.Close()

	// get images list(old -> new)
	outStream, errStream, err := nanoshlib.Exec(dbDeployerAbsPath + " -cfg " + dbJsonAbsPath + " ls " + config.DBMS, -1)
	if err != nil {
		panic("[AffClassify]dbdeployer ls "+config.DBMS+" error" + err.Error() + ": " + errStream)
	}
	images := strings.Split(strings.TrimSpace(outStream), "\n")
	images = images[1:]

	// gen order map, the larger the value, the newer the version
	orderMap := make(map[string]int)
	for i, image := range images {
		orderMap[image] = i
	}

	// select * from affversion where status = 1
	rows, err := affVersionDB.Query(`SELECT * FROM affversion WHERE status=1;`)
	if err != nil {
		panic("[AffClassify]select bugs error: " + err.Error())
	}
	defer rows.Close()

	// oldest1Map: status=1, taskId/sqlsim/bugJsonName -> oldest version
	oldest1Map := make(map[string]string)
	for rows.Next() {
		var taskId int
		var bugJsonName string
		var version string
		var status int
		err = rows.Scan(&taskId, &bugJsonName, &version, &status)
		if err != nil {
			panic("[AffClassify]scan row error: " + err.Error())
		}

		bug := "task-"+strconv.Itoa(taskId)+"/sqlsim/"+bugJsonName
		if status == 1 {
			if v, ok := oldest1Map[bug]; ok {
				if vcmp(orderMap, v, version) > 0 {
					oldest1Map[bug] = version
				}
			} else {
				oldest1Map[bug] = version
			}
		} else {
			panic("[AffClassify]status must be 1")
		}
	}
	if rows.Err() != nil {
		panic("[AffClassify]rows err: " + rows.Err().Error())
	}

	// select * from affversion where status = 0
	rows, err = affVersionDB.Query(`SELECT * FROM affversion WHERE status=0;`)
	if err != nil {
		panic("[AffClassify]select bugs error: " + err.Error())
	}
	defer rows.Close()

	// status0Map: status=0, taskId/sqlsim/bugJsonName@version -> true
	status0Map := make(map[string]bool)
	for rows.Next() {
		var taskId int
		var bugJsonName string
		var version string
		var status int
		err = rows.Scan(&taskId, &bugJsonName, &version, &status)
		if err != nil {
			panic("[AffClassify]scan row error: " + err.Error())
		}

		bugVersion := "task-"+strconv.Itoa(taskId)+"/sqlsim/"+bugJsonName+"@"+version
		if status == 0 {
			status0Map[bugVersion] = true
		} else {
			panic("[AffClassify]status must be 0")
		}
	}
	if rows.Err() != nil {
		panic("[AffClassify]rows err: " + rows.Err().Error())
	}

	// `o1v` -> bugList
	vListPairMap := make(map[string][]string)
	for bug, o1v := range oldest1Map {
		pv := preV(images, orderMap, o1v)
		if pv == "" {
			if bugList, ok := vListPairMap[o1v]; ok {
				bugList = append(bugList, bug)
				vListPairMap[o1v] = bugList
			} else {
				vListPairMap[o1v] = []string{bug}
			}
		} else {
			if _, ok := status0Map[bug+"@"+pv]; ok {
				if bugList, ok := vListPairMap[o1v]; ok {
					bugList = append(bugList, bug)
					vListPairMap[o1v] = bugList
				} else {
					vListPairMap[o1v] = []string{bug}
				}
			}
		}
	}

	// gen json
	vListPairs := make([]*VListPair, 0)
	for o1v, bugList := range vListPairMap {
		vListPairs = append(vListPairs, &VListPair{
			O1V: o1v,
			BugList: bugList,
		})
	}

	// write json
	jsonData, err := json.Marshal(vListPairs)
	if err != nil {
		panic("[AffClassify]marshal error: " + err.Error())
	}
	affClassifyPath := path.Join(taskPoolPath, "affClassify.json")
	err = ioutil.WriteFile(affClassifyPath, jsonData, 0777)
	if err != nil {
		panic("[AffClassify]write json error: " + err.Error())
	}

	// We will also create a directory `affclassify` in taskPoolPath.
	// For each `o1v`, we will save the first detected bug in `affclassify`
	affPackPath := path.Join(taskPoolPath, "affclassify")
	_ = os.RemoveAll(affPackPath)
	_ = os.Mkdir(affPackPath, 0777)
	for o1v, bugList := range vListPairMap {
		if v, ok := orderMap[o1v]; ok {
			// find the first reported bug
			var firstBug *task.BugReport = nil
			firstBugJsonPath := ""
			var firstBugTime time.Time
			for _, bug := range bugList {
				bugJsonPath := path.Join(taskPoolPath, bug)
				bugReport, err := task.NewBugReport(bugJsonPath)
				if err != nil {
					panic("[AffClassify]read bug error: " + err.Error() + ": " + bugJsonPath)
				}
				bugTime := parseTimeStr(bugReport.ReportTime)
				if firstBug == nil || firstBugTime.After(bugTime) {
					firstBug = bugReport
					firstBugJsonPath = bugJsonPath
					firstBugTime = bugTime
				}
			}

			sqlsimPath := path.Dir(firstBugJsonPath)
			taskPath := path.Dir(sqlsimPath)
			taskName := path.Base(taskPath)

			if v > 9999 {
				panic("[AffClassify]v > 9999? " + strconv.Itoa(v))
			}
			// mkdir affPackPath/v-o1v-taskName
			o1v := strings.ReplaceAll(o1v, "/", "@")
			vPath := path.Join(affPackPath, fmt.Sprintf("%04d-%s-%s", v, o1v, taskName))
			_ = os.Mkdir(vPath, 0777)

			// cp ddl
			ddlPath := path.Join(taskPath, "output.data.sql")
			_, errStream, err := nanoshlib.Exec("cp "+ddlPath+" "+vPath, -1)
			if err != nil {
				panic("[AffClassify]cp ddl error: " + errStream)
			}
			// cp bug json
			_, errStream, err = nanoshlib.Exec("cp "+firstBugJsonPath+" "+vPath, -1)
			if err != nil {
				panic("[AffClassify]cp bug json error: " + errStream)
			}
			// cp bug log
			firstBugLogPath := firstBugJsonPath[0:len(firstBugJsonPath)-5]+".log"
			_, errStream, err = nanoshlib.Exec("cp "+firstBugLogPath+" "+vPath, -1)
			if err != nil {
				panic("[AffClassify]cp bug log error: " + errStream)
			}
		} else {
			panic("[AffClassify]can not found " + o1v)
		}
	}
}

// v1 == v2: 0
// v1 newer than v2: 1
// v1 older than v2: -1
func vcmp(orderMap map[string]int, v1 string, v2 string) int {
	vv1 := -1
	if v, ok := orderMap[v1]; ok {
		vv1 = v
	} else {
		panic("[vcmp]can not found " + v1)
	}
	vv2 := -1
	if v, ok := orderMap[v2]; ok {
		vv2 = v
	} else {
		panic("[vcmp]can not found " + v2)
	}
	if vv1 == vv2 {
		return 0
	} else if vv1 < vv2 {
		return -1
	} else {
		return 1
	}
}

func preV(images []string, orderMap map[string]int, version string) string {
	vv := -1
	if v, ok := orderMap[version]; ok {
		vv = v
	} else {
		panic("[preV]can not found " + version)
	}

	if vv == 0 {
		return ""
	}
	return images[vv-1]
}

func parseTimeStr(timeStr string) time.Time {
	idx := strings.Index(timeStr, " m=+")
	if idx != -1 {
		timeStr = timeStr[:idx]
	}
	t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", timeStr)
	if err != nil {
		panic("[parseTimeStr]parse time error: " + err.Error())
	}
	return t
}