package testsqls

import (
	"errors"
	"fmt"
	"github.com/qaqcatz/impomysql/connector"
	"io/ioutil"
	"log"
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
	SQLSubQuery5 = "SELECT * FROM COMPANY WHERE NOT EXISTS (SELECT ID FROM COMPANY WHERE ID <= 1)"
	// In MySQL, JOIN, CROSS JOIN, and INNER JOIN are syntactic equivalents (they can replace each other).
	SQLJOIN = "SELECT * FROM COMPANY JOIN (SELECT * FROM COMPANY WHERE ID = 1) AS T1 ON COMPANY.ID > T1.ID"
	SQLJOIN2 = "SELECT * FROM COMPANY LEFT OUTER JOIN (SELECT * FROM COMPANY WHERE ID = 2) AS T1 ON COMPANY.ID > T1.ID"
	SQLJOIN3 = "SELECT * FROM COMPANY RIGHT OUTER JOIN (SELECT * FROM COMPANY WHERE ID = 2) AS T1 ON COMPANY.ID > T1.ID"
	SQLJOIN4 = "SELECT * FROM COMPANY STRAIGHT_JOIN (SELECT * FROM COMPANY WHERE ID = 2) AS T1 ON COMPANY.ID > T1.ID"
	SQLJOIN5 = "SELECT * FROM COMPANY NATURAL JOIN (SELECT * FROM COMPANY WHERE ID = 2) AS T1"
	SQLJOIN6 = "SELECT * FROM COMPANY NATURAL LEFT JOIN (SELECT * FROM COMPANY WHERE ID = 2) AS T1"
	SQLLIMIT = "SELECT * FROM COMPANY LIMIT 2147483647,1"
	SQLLIMIT2 = "SELECT * FROM COMPANY LIMIT 1"
	SQLUNION = "SELECT * FROM COMPANY UNION ALL SELECT * FROM (SELECT * FROM COMPANY UNION SELECT * FROM COMPANY) AS T1"
	SQLUNION2 = "(SELECT * FROM COMPANY UNION ALL SELECT * FROM (SELECT * FROM COMPANY UNION SELECT * FROM COMPANY) AS T1) " +
		"UNION ALL " +
		"SELECT * FROM (SELECT * FROM COMPANY UNION SELECT * FROM COMPANY) AS T1"
	SQLWITH = "WITH XX AS (SELECT * FROM COMPANY) SELECT * FROM XX"
	SQLWITH2 = "WITH RECURSIVE fibonacci (n, fib_n, next_fib_n) AS " +
		"(SELECT 1, 0, 1 UNION ALL SELECT n + 1, next_fib_n, fib_n + next_fib_n " +
		"FROM fibonacci WHERE n < 10 ) SELECT * FROM fibonacci"
	SQLIN = "SELECT 1 IN (1, 2, 3)"
	SQLIN2 = "SELECT * FROM COMPANY WHERE ID IN (1, 2, 3)"
	SQLWHERE = "SELECT * FROM COMPANY WHERE TRUE"
	SQLHAVING = "SELECT * FROM COMPANY HAVING TRUE"
	SQLLIKE = "SELECT * FROM COMPANY WHERE 'abc' NOT LIKE 'A_%' ESCAPE '_'"
	SQLRegExp = "SELECT * FROM COMPANY WHERE 'abc' NOT REGEXP '^A[B]*C$'"
	SQLBetween = "SELECT * FROM COMPANY WHERE ID BETWEEN 1 AND 3"
	SQLBetween2 = "SELECT * FROM COMPANY WHERE ID BETWEEN '1' AND '3'"
	SQLBetween3 = "SELECT * FROM COMPANY WHERE NAME BETWEEN 0 AND 'A'"
	SQLEX = "WITH MYWITH AS " +
		"((SELECT (yearweek('2003-11-02') >> f4 - f4) AS f1, ('2006') AS f2, (~ f6) AS f3 " +
		"FROM  (SELECT (BINARY abs(4)) AS f4, (quarter('2001-01-28')) AS f7, (f9 MOD f9 | f9) AS f6 " +
		"FROM  (SELECT `col_double_undef_signed` AS f8, `col_bigint_key_unsigned` AS f9, `col_varchar(20)_undef_signed` AS f10 " +
		"FROM  table_3_utf8_2 FORCE INDEX (`col_float_key_unsigned`)) AS t1  " +
		"WHERE (((-4958736797163969007 - INTERVAL 1 QUARTER) >= (BINARY f10)) " +
		"OR (NOT (CAST((\"in\") AS CHAR) NOT LIKE '%1%'))) IS FALSE " +
		"HAVING  ((NOT ((from_days(1228031832738593057)) >= (- -9014801606300802676 + INTERVAL 1 DAY_HOUR))) " +
		"OR (NOT (log2(0.6050498217840262) - INTERVAL 1 SECOND_MICROSECOND))) " +
		"AND ((((sign(293170577176557332) MOD f6) NOT IN (atan(0.5116825416612745), f6, last_day('2013-02-23 12:51:07'))) IS TRUE) " +
		"OR (NOT (f4)))  ORDER BY f9 ) AS t2  " +
		"INNER JOIN  (SELECT `col_varchar(20)_key_signed` AS f11, `col_char(20)_undef_signed` AS f5, `col_float_key_unsigned` AS f12 " +
		"FROM table_3_utf8_2 USE INDEX (`col_varchar(20)_key_signed`)) AS t3   ) " +
		"UNION ALL (SELECT (NULL - sign(2)) AS f1, (6 + INTERVAL 1 DAY_HOUR) AS f2, (f14 - INTERVAL 1 HOUR_SECOND) AS f3 " +
		"FROM  (SELECT `col_float_undef_unsigned` AS f13, `col_double_undef_signed` AS f14, `col_float_undef_signed` AS f15 FROM table_5_utf8_2 USE INDEX (`col_float_key_signed`,`col_decimal(40, 20)_key_signed`)) AS t4    ORDER BY f15 )) " +
		"SELECT * FROM MYWITH"
	SQLEX2 = "WITH `MYWITH` AS (SELECT (`f4`) AS `f1`,(COLLATION(`f4`)&`f5`) AS `f2`,(BIN(`f4`)) AS `f3` " +
		"FROM (SELECT `col_char(20)_undef_signed` AS `f7`,`col_bigint_undef_signed` AS `f5`,`col_double_undef_signed` AS `f8` " +
		"FROM `table_3_utf8_2` FORCE INDEX (`col_varchar(20)_key_signed`, `col_varchar(20)_key_signed`)) AS `t1` " +
		"NATURAL JOIN (SELECT (ABS(-1529566578300132310)) AS `f4`,(`f11`) AS `f9`,(`f12`*CEILING(5636069043819042262)) AS `f6` " +
		"FROM (SELECT `col_bigint_undef_unsigned` AS `f10`,`col_char(20)_undef_signed` AS `f11`,`col_double_key_signed` AS `f12` " +
		"FROM `table_5_utf8_2`) AS `t2` WHERE ((((FORMAT_BYTES(`f12`)) NOT BETWEEN COLLATION(`f11`) " +
		"AND DAYOFWEEK(_UTF8MB4'2003-07-05')) IS FALSE) " +
		"OR (((RTRIM(`f12`)) NOT BETWEEN FORMAT_BYTES(_UTF8MB4'2004-03-12') " +
		"AND `f10`) IS TRUE) " +
		"OR ((DATE_ADD(`f11`, INTERVAL 1 DAY_HOUR))=(_UTF8MB4'r'))) IS TRUE " +
		"HAVING (((~`f6`%~`f4`) IN (`f9`,`f9`,`f9`)) IS FALSE) IS TRUE) AS `t3`) " +
		"SELECT * FROM `MYWITH`\n"
)

// sql file benchmark:

// getPackagePath: get the package actual path, then you can read files under the path.
func getPackagePath() (string, error) {
	if _, file, _, ok := runtime.Caller(0); !ok {
		return "", errors.New("PackagePath: runtime.Caller(0) error ")
	} else {
		return path.Join(file, "../"), nil
	}
}

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
	packagePath, err := getPackagePath()
	if err != nil {
		return nil, errors.New("ReadSQLFile: getPackagePath() error "), ""
	}
	sqlFilePath := path.Join(packagePath, sqlFileName)
	data, err := ioutil.ReadFile(sqlFilePath)
	if err != nil {
		return nil, errors.New("ReadSQLFile: read " + sqlFilePath + " error: " + err.Error()), ""
	}
	return data, nil, sqlFilePath
}

// test .zz, .yy

const (
	zzTest = "test.zz.lua"
	yyTest = "test.yy"
)

func GetTestZZPath() string {
	packagePath, err := getPackagePath()
	if err != nil {
		log.Fatal(err)
	}
	return path.Join(packagePath, zzTest)
}

func GetTestYYPath() string {
	packagePath, err := getPackagePath()
	if err != nil {
		log.Fatal(err)
	}
	return path.Join(packagePath, yyTest)
}
