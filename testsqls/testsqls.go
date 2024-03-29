package testsqls

import (
	"fmt"
	"github.com/qaqcatz/impomysql/connector"
)

// dbms
const (
	host     = "127.0.0.1"
	username = "root"
	password = "123456"
	dbname   = "TEST"

	MySQL     = "mysql"
	MariaDB   = "mariadb"
	TiDB      = "tidb"
	OceanBase = "oceanbase"

	mySQLPort     = 13306
	mariaDBPort   = 23306
	tiDBPort      = 4000
	oceanBasePort = 2881
)

// GetConnector: get connector for the DBMS, default: MySQL
func GetConnector(DBMS string) (*connector.Connector, error) {
	myPort := 0
	switch DBMS {
	case MySQL:
		myPort = mySQLPort
	case MariaDB:
		myPort = mariaDBPort
	case TiDB:
		myPort = tiDBPort
	case OceanBase:
		myPort = oceanBasePort
	default:
		myPort = mySQLPort
	}
	conn, err := connector.NewConnector(host, myPort, username, password, dbname)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// SQLExec: Execute the sql, print the result into standard output stream.
func SQLExec(sql string, DBMS string) error {
	conn, err := GetConnector(DBMS)
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
func InitTableCOMPANY(DBMS string) error {
	conn, err := GetConnector(DBMS)
	if err != nil {
		return err
	}
	result := conn.ExecSQL("DROP TABLE IF EXISTS COMPANY")
	if result.Err != nil {
		return result.Err
	}
	result = conn.ExecSQL("CREATE TABLE COMPANY (ID INT, NAME TEXT, AGE INT, CITY TEXT, KEY(ID), KEY(AGE))")
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
	SQLSelectValue  = "SELECT 1"
	SQLSelectValue2 = "SELECT 1.0001"
	SQLSelectValue3 = "SELECT 'a'"
	SQLSubQuery     = "SELECT * FROM COMPANY WHERE ID = (SELECT ID FROM COMPANY WHERE ID = 1)"
	SQLSubQuery2    = "SELECT * FROM COMPANY WHERE ID = ANY (SELECT ID FROM COMPANY WHERE ID > 1)"
	SQLSubQuery3    = "SELECT * FROM COMPANY WHERE ID NOT IN (SELECT ID FROM COMPANY WHERE ID IN (1, 2))"
	SQLSubQuery4    = "SELECT * FROM COMPANY WHERE ID > ALL (SELECT ID FROM COMPANY WHERE ID <= 1)"
	SQLSubQuery5    = "SELECT * FROM COMPANY WHERE NOT EXISTS (SELECT ID FROM COMPANY WHERE ID <= 1)"
	// In MySQL, JOIN, CROSS JOIN, and INNER JOIN are syntactic equivalents (they can replace each other).
	SQLJOIN   = "SELECT * FROM COMPANY JOIN (SELECT * FROM COMPANY WHERE ID = 1) AS T1 ON COMPANY.ID > T1.ID"
	SQLJOIN2  = "SELECT * FROM COMPANY LEFT OUTER JOIN (SELECT * FROM COMPANY WHERE ID = 2) AS T1 ON COMPANY.ID > T1.ID"
	SQLJOIN3  = "SELECT * FROM COMPANY RIGHT OUTER JOIN (SELECT * FROM COMPANY WHERE ID = 2) AS T1 ON COMPANY.ID > T1.ID"
	SQLJOIN4  = "SELECT * FROM COMPANY STRAIGHT_JOIN (SELECT * FROM COMPANY WHERE ID = 2) AS T1 ON COMPANY.ID > T1.ID"
	SQLJOIN5  = "SELECT * FROM COMPANY NATURAL JOIN (SELECT * FROM COMPANY WHERE ID = 2) AS T1"
	SQLJOIN6  = "SELECT * FROM COMPANY NATURAL LEFT JOIN (SELECT * FROM COMPANY WHERE ID = 2) AS T1"
	SQLLIMIT  = "SELECT * FROM COMPANY LIMIT 2147483647,1"
	SQLLIMIT2 = "SELECT * FROM COMPANY LIMIT 1"
	SQLUNION  = "SELECT * FROM COMPANY UNION ALL SELECT * FROM (SELECT * FROM COMPANY UNION SELECT * FROM COMPANY) AS T1"
	SQLUNION2 = "(SELECT * FROM COMPANY UNION ALL SELECT * FROM (SELECT * FROM COMPANY UNION SELECT * FROM COMPANY) AS T1) " +
		"UNION ALL " +
		"SELECT * FROM (SELECT * FROM COMPANY UNION SELECT * FROM COMPANY) AS T1"
	SQLWITH  = "WITH XX AS (SELECT * FROM COMPANY) SELECT * FROM XX"
	SQLWITH2 = "WITH RECURSIVE fibonacci (n, fib_n, next_fib_n) AS " +
		"(SELECT 1, 0, 1 UNION ALL SELECT n + 1, next_fib_n, fib_n + next_fib_n " +
		"FROM fibonacci WHERE n < 10 ) SELECT * FROM fibonacci"
	SQLWITH3       = "WITH XX AS (SELECT * FROM COMPANY UNION ALL SELECT * FROM COMPANY) SELECT * FROM XX"
	SQLIN          = "SELECT 1 IN (1, 2, 3)"
	SQLIN2         = "SELECT * FROM COMPANY WHERE ID IN (1, 2, 3)"
	SQLWHERE       = "SELECT * FROM COMPANY WHERE TRUE"
	SQLHAVING      = "SELECT * FROM COMPANY HAVING TRUE"
	SQLLIKE        = "SELECT * FROM COMPANY WHERE 'abc' NOT LIKE 'A_%' ESCAPE '_'"
	SQLRegExp      = "SELECT * FROM COMPANY WHERE 'abc' NOT REGEXP '^A[B]*C$'"
	SQLBetween     = "SELECT * FROM COMPANY WHERE ID BETWEEN 1 AND 3"
	SQLBetween2    = "SELECT * FROM COMPANY WHERE ID BETWEEN '1' AND '3'"
	SQLHint        = "SELECT ID, AGE FROM COMPANY USE INDEX(ID, AGE)"
	SQLOrderBy     = "SELECT ID, NAME FROM COMPANY ORDER BY ID"
	SQLDropTable   = "DROP TABLE IF EXISTS CT;"
	SQLCreateTable = "CREATE TABLE IF NOT EXISTS CT (`pk` int primary key," +
		"`c1` bigint," +
		"`c1_u` bigint unsigned," +
		"key (`c1`)," +
		"key (`c1_u`)) character set utf8\npartition by hash(pk)\npartitions 2;"
	SQLINSERT    = "INSERT INTO CT VALUES (0,1,1),(1,2,2),(2,3,3);"
	SQLBinaryOp  = "SELECT * FROM COMPANY WHERE ID = 1 OR ID = 2;"
	SQLBinaryOp2 = "SELECT * FROM COMPANY WHERE (AGE > 18 AND AGE < 20) AND 2 > 1 OR 2 < 1 AND (3 > 1 AND 4 > 1) OR 3 < 1 OR 4 < 1;"
	SQLTimeFunc = "SELECT to_seconds('1998-06-24 00:00:01') FROM COMPANY WHERE to_seconds('1998-05-01 00:00:01') > 0;"
	SQLStrFunc = "SELECT to_base64('123');"
	SQLInfoFunc = "SELECT format_bytes('123');"
	)
