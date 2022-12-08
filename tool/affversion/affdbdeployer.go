package affversion

import (
	"github.com/qaqcatz/impomysql/task"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func AffDBDeployer(dbDeployerPath string, config *task.TaskPoolConfig, threadNum int, portInterval int) error {
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

	return nil
}
