package affversion

import (
	"database/sql"
	"encoding/json"
	"github.com/qaqcatz/impomysql/task"
	"github.com/qaqcatz/nanoshlib"
	"io/ioutil"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type VListPair struct{
	O1V string `json:"o1v"`
	BugList []string `json:"bugList"`
}

// AffClassify: classify bugs according to the versions they affect.
//
// Specifically, for each bug,
// we will calculate the oldest reproducible version `o1v` and use it for classification.
//
// Make sure you have done `affversion` or `affdbdeployer`, we will query the database `affversion.db`.
// You also need to provide `dbdeployer`, which will tell us the order of each version
//
// So the command is:
//   `./impomysql affclassify dbDeployerPath dbJsonPath taskPoolConfigPath`
//
// We will create `affclassify.json` in taskPoolPath. It is an array of {`o1v`, bug list}.
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

	// `o1v` -> bugList
	vListPairMap := make(map[string][]string)
	for bug, o1v := range oldest1Map {
		if bugList, ok := vListPairMap[o1v]; ok {
			bugList = append(bugList, bug)
			vListPairMap[o1v] = bugList
		} else {
			vListPairMap[o1v] = []string{bug}
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