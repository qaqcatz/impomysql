package testsqls

import (
	"errors"
	"fmt"
	"github.com/qaqcatz/impomysql/connector"
	"io/ioutil"
	"path"
	"runtime"
)

// sudo docker run -itd --name test -p 13306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql
const (
	host = "127.0.0.1"
	port = 13306
	username = "root"
	password = "123456"
	dbname = "TEST"
)

// InitDBTEST:
//   DROP DATABASE IF EXISTS TEST
//   CREATE DATABASE TEST
func InitDBTEST() error {
	conn, err := connector.NewConnector(host, port, username, password, "")
	if err != nil {
		return err
	}
	result := conn.ExecSQL("DROP DATABASE IF EXISTS " + dbname)
	if result.Err != nil {
		return err
	}
	result = conn.ExecSQL("CREATE DATABASE " + dbname)
	if result.Err != nil {
		return err
	}
	return nil
}

// InitDBTEST:
//   CREATE DATABASE IF NOT EXISTS TEST
func EnsureDBTEST() error {
	conn, err := connector.NewConnector(host, port, username, password, "")
	if err != nil {
		return err
	}
	result := conn.ExecSQL("CREATE DATABASE IF NOT EXISTS TEST " + dbname)
	if result.Err != nil {
		return err
	}
	return nil
}

func GetConnector() (*connector.Connector, error) {
	conn, err := connector.NewConnector(host, port, username, password, dbname)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// SQLExec: Execute the sql, print the result into standard output stream.
func SQLExec(sql string) error {
	conn, err := GetConnector()
	if err != nil {
		return err
	}
	fmt.Println("Exec SQL:", sql)
	result := conn.ExecSQL(sql)
	if result.Err != nil {
		return result.Err
	}
	fmt.Println("Exec result:", result.ToString())
	return nil
}

// table benchmark:

// InitTableCOMPANY:
//   DROP TABLE IF EXISTS COMPANY
//   CREATE TABLE COMPANY (ID INT, NAME TEXT, AGE INT, CITY TEXT)
//   INSERT INTO COMPANY VALUES
//   (1, 'A', 18, 'a'), (2, 'B', 19, 'b'), (3, 'C', 20, 'c'),
//   (4, 'A', 19, 'c'), (5, 'A', 19, 'c'), (6, 'B', 18, 'b')
func InitTableCOMPANY() error {
	conn, err := GetConnector()
	if err != nil {
		return err
	}
	result := conn.ExecSQL("DROP TABLE IF EXISTS COMPANY")
	if result.Err != nil {
		return result.Err
	}
	result = conn.ExecSQL("CREATE TABLE COMPANY (ID INT, NAME TEXT, AGE INT, CITY TEXT)")
	if result.Err != nil {
		return result.Err
	}
	result = conn.ExecSQL("INSERT INTO COMPANY VALUES (1, 'A', 18, 'a'), (2, 'B', 19, 'b'), " +
		"(3, 'C', 20, 'c'), (4, 'A', 19, 'c'), (5, 'A', 19, 'c'), (6, 'B', 18, 'b')")
	if result.Err != nil {
		return result.Err
	}
	return nil
}

// sql benchmark:
const (
	SQLAGG = "SELECT S, G, CITY FROM ( " +
		"   SELECT SUM(ID+1) AS S, GROUP_CONCAT(NAME ORDER BY NAME DESC) AS G, CITY " +
		"   FROM COMPANY " +
		"   GROUP BY CITY " +
		"   HAVING COUNT(DISTINCT AGE) >= 1 " +
		") AS T " +
		"WHERE T.S > 0;"
	SQLWindow = "SELECT " +
		"   ID AS id, CITY, AGE, " +
		"   SUM(AGE) OVER w " +
		"   AS sum_age, " +
		"   AVG(AGE) OVER (PARTITION BY CITY ORDER BY ID ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) " +
		"   AS avg_age, " +
		"   ROW_NUMBER() OVER (PARTITION BY CITY ORDER BY ID) " +
		"   AS rn " +
		"   FROM COMPANY " +
		"   WINDOW w AS (PARTITION BY CITY ORDER BY ID ROWS UNBOUNDED PRECEDING)"
	SQLSelectValue = "SELECT 1"
	SQLSelectValue2 = "SELECT 1.0001"
	SQLSelectValue3 = "SELECT 'a'"
	SQLSubQuery = "SELECT * FROM COMPANY WHERE ID = (SELECT ID FROM COMPANY WHERE ID = 1)"
	SQLSubQuery2 = "SELECT * FROM COMPANY WHERE ID = ANY (SELECT ID FROM COMPANY WHERE ID > 1)"
	SQLSubQuery3 = "SELECT * FROM COMPANY WHERE ID NOT IN (SELECT ID FROM COMPANY WHERE ID IN (1, 2))"
	SQLSubQuery4 = "SELECT * FROM COMPANY WHERE ID > ALL (SELECT ID FROM COMPANY WHERE ID <= 1)"
)

// sql file benchmark:

const (
	SQLFileQuote = "quote"
	SQLFileTest = "test"
	SQLFileAgg = "agg"
	SQLFileWindow = "window"
)

// ReadSQLFile: read the sql file under testsqls with the help of runtime.Caller().
//
// The third return value is the absolute filepath,
// you can use it to get the actual location of the file
func ReadSQLFile(sqlFileName string) ([]byte, error, string) {
	sqlFileName += ".sql"
	if _, file, _, ok := runtime.Caller(0); !ok {
		return nil, errors.New("ReadSQLFile: runtime.Caller(0) error "), ""
	} else {
		sqlFilePath := path.Join(file, "../", sqlFileName)
		data, err := ioutil.ReadFile(sqlFilePath)
		if err != nil {
			return nil, errors.New("ReadSQLFile: read " + sqlFilePath + " error: " + err.Error()), ""
		}
		return data, nil, sqlFilePath
	}
}


