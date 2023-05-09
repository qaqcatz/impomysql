package affversion

import (
	"fmt"
	"github.com/qaqcatz/impomysql/task"
	"github.com/qaqcatz/nanoshlib"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// AffDBDeployer:
// In affversion, we need to manually deploy each version of DBMS.
// Now this work can be done automatically, see https://github.com/qaqcatz/dbdeployer
//
// With dbdeployer, you can use the following command to verify all versions from newestImage to oldestImage:
//   `./impomysql affdbdeployer dbDeployerPath dbJsonPath taskPoolConfigPath threadNum port newestImage oldestImage`
func AffDBDeployer(dbDeployerPath string, dbJsonPath string, config *task.TaskPoolConfig, threadNum int, portStr string,
	newestImage string, oldestImage string) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic("[AffDBDeployer]parse port error: " + err.Error() + ": " + portStr)
	}
	dbDeployerAbsPath, err := filepath.Abs(dbDeployerPath)
	if err != nil {
		panic("[AffDBDeployer]path abs error: " + err.Error())
	}
	dbJsonAbsPath, err := filepath.Abs(dbJsonPath)
	if err != nil {
		panic("[AffDBDeployer]path abs error: " + err.Error())
	}

	// create logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	writers := []io.Writer{
		os.Stdout,
	}
	multiWriter := io.MultiWriter(writers...)
	logger.SetOutput(multiWriter)
	logger.SetLevel(logrus.InfoLevel)

	// get images list(old -> new)
	outStream, errStream, err := nanoshlib.Exec(dbDeployerAbsPath + " -cfg " + dbJsonAbsPath + " ls " + config.DBMS, -1)
	if err != nil {
		panic("[AffDBDeployer]dbdeployer ls "+config.DBMS+" error" + err.Error() + ": " + errStream)
	}
	images := strings.Split(strings.TrimSpace(outStream), "\n")
	images = images[1:]

	// cut newestImage-oldestImage
	newestId := -1
	if newestImage == "" {
		newestId = len(images)-1
	} else {
		for i := len(images)-1; i >= 0; i -= 1 {
			if images[i] == newestImage {
				newestId = i
				break
			}
		}
	}
	oldestId := -1
	if oldestImage == "" {
		oldestId = 0
	} else {
		for i := 0; i < len(images); i += 1 {
			if images[i] == oldestImage {
				oldestId = i
				break
			}
		}
	}
	if newestId == -1 || oldestId == -1 || newestId < oldestId {
		panic("[AffDBDeployer]interval error! check your newestImage: "+newestImage+" or oldestImage: "+oldestImage)
	}
	images = images[oldestId:(newestId+1)]

	fmt.Println("We will do verification on "+strconv.Itoa(len(images))+" versions, " +
		"which is an expensive task! Do you want to continue? (Enter yes)")
	var input string = ""
	_, err = fmt.Scanf("%s", &input)
	if err != nil {
		panic("[AffDBDeployer]Input error! You may provide an empty string!")
	}
	if input != "yes" {
		panic("Stop!")
	}

	// new -> old,
	// dbdeployer -cfg dbJsonPath run dbms image portStr
	// affversion taskpool taskpoolConfig threadNum image preImage@1
	deployErrSkip := make([]int, 0)
	AffErrSkip := make([]int, 0)
	logger.Info("Start!")
	logger.Info("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	for i := len(images)-1; i >= 0; i -= 1 {
		image := images[i]
		logger.Info("**************************************************")
		logger.Info("image ", i, ": ", image)
		logger.Info("**************************************************")
		logger.Info("dbdeployer:")
		err := nanoshlib.ExecStd(dbDeployerAbsPath + " -cfg " + dbJsonAbsPath + " run " + config.DBMS + " " + image + " " +portStr, -1)
		if err != nil {
			logger.Warn("[AffDBDeployer]dbdeployer run " + config.DBMS + " " + image + " " + portStr + " error")
			deployErrSkip = append(deployErrSkip, i)
			continue
		}
		logger.Info("**************************************************")
		logger.Info("affversionpool:")
		err = AffVersionTaskPool(config, threadNum, port, image, "")
		if err != nil {
			logger.Warnf("[AffDBDeployer]aff version task pool error: %+v", err)
			AffErrSkip = append(AffErrSkip, i)
			continue
		}
	}
	logger.Info("[deployErrSkip]")
	for _, skip := range deployErrSkip {
		fmt.Print(strconv.Itoa(skip) + " ")
	}
	fmt.Println()
	logger.Info("[AffErrSkip]")
	for _, skip := range AffErrSkip {
		fmt.Print(strconv.Itoa(skip) + " ")
	}
	fmt.Println()

	logger.Info("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	logger.Info("Finished!")
}
