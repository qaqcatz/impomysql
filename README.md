# impomysql

[![Go Reference](https://pkg.go.dev/badge/github.com/qaqcatz/impomysql.svg)](https://pkg.go.dev/github.com/qaqcatz/impomysql)

Detecting Logic Bugs in mysql through Implication Oracle.

Also supports DBMS compatible with mysql syntax, such as mariadb, tidb, oceanbase.

## 1. What is logical bug

see this bug report as an example:

https://bugs.mysql.com/bug.php?id=108937

In theory, the result of sql1 ⊆ the result of sql2:

```sql
SELECT c1-DATE_SUB('2008-05-25', INTERVAL 1 HOUR_MINUTE) AS f1 FROM t HAVING f1 != 0; -- sql1
SELECT c1-DATE_SUB('2008-05-25', INTERVAL 1 HOUR_MINUTE) AS f1 FROM t HAVING 1; -- sql2
```

Because the `HAVING 1` in sql2 is always true, but the `HAVING f1 != 0` in sql1 may be false. 

However, the date value changed after changing `HAVING f1 != 0` to `HAVING 1`, this is a logical bug:

```sql
mysql> SELECT c1-DATE_SUB('2008-05-25', INTERVAL 1 HOUR_MINUTE) AS f1 FROM t HAVING f1 != 0; -- sql1
+------------+
| f1         |
+------------+
| -1928.8181 |
|  -1995.009 |
|      -2007 |
+------------+
3 rows in set (0.00 sec)

mysql> SELECT c1-DATE_SUB('2008-05-25', INTERVAL 1 HOUR_MINUTE) AS f1 FROM t HAVING 1; -- sql2
+---------------------+
| f1                  |
+---------------------+
| -20080524235820.816 |
| -20080524235887.008 |
|     -20080524235899 |
+---------------------+
3 rows in set (0.00 sec)
```

## 2. What is Implication Oracle

In the above example, we changed `HAVING f1 != 0`  to `HAVING 1`.

In theory, the predicate of sql1 → the predicate of sql2, and the result of sql1 ⊆ the result of sql2. 

If the actual result does not satisfy this relationship, we consider that there is a logical bug.

Although the idea is simple, some features make it difficult to implement, such as aggregate functions, window functions, type conversion, LIMIT, LEFT/RIGHT JOIN, flow control operations, etc.

We will discuss these features in our paper:

```
todo  
```

You can also see the source code for more details:

* mutation/doc.go

* mutation/stage1/doc.go
* mutation/stage2/doc.go
* mutation/stage2/mutatevisitor.go
* resources/impo.yy

## 3. How to use

### 3.1 build

It is recommended to use `golang 1.16.2`.

```shell
git clone --depth=1 https://github.com/qaqcatz/impomysql.git
cd impomysql
go build
```

<font color="red">In the following we will refer to the path of `impomysql` as `${IMPOHOME}`</font>

Now you will see an executable file `${IMPOHOME}/impomysql`.

### 3.2 start your DBMS

For example, you can start mysql with docker:

```shell
sudo docker run -itd --name test -p 13306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql:8.0.30
```

You can also compile and install the DBMS yourself.

### 3.3 run task

We treat DBMS testing as `task`.

#### quick start

We assume you have executed `sudo docker run -itd --name mysqltest -p 13306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql:8.0.30`

```shell
cd ${IMPOHOME}
./impomysql task ./resources/taskconfig.json
```

#### input

You need to provide a configuration file. see `${IMPOHOME}/resources/raskconfig.json`:

```json
{
  "outputPath": "./output",
  "dbms": "mysql",
  "taskId": 1,
  "host": "127.0.0.1",
  "port": 13306,
  "username": "root",
  "password": "123456",
  "dbname": "TEST",
  "seed": 123456,
  "ddlPath": "./resources/ddl.sql",
  "dmlPath": "./resources/dml.sql"
}
```

* `outputPath`, `dbms`, `taskId`: we will save the result in `outputPath`/`dbms`/task-`taskId`. `taskId` >= 0.

* `host`, `port`, `username`, `password`, `dbname`: we will create a database connector with dsn `username`:`password`@tcp(`host`:`port`)/`dbname`, and init database `dbname`.

* `seed`: random seed. If seed <= 0, we will use the current time.

* `ddlPath`: sql file responsible for creating data. see `${IMPOHOME}/resources/ddl.sql`:

  ```sql
  create table t (c1 double);
  insert into t values (79.1819),(12.991),(1);
  ```

  We will init database according to `ddlPath`.

* `dmlPath`: sql file responsible for querying data. see `${IMPOHOME}/resources/dml.sql`:

  ```sql
  SELECT c1-DATE_SUB('2008-05-25', INTERVAL 1 HOUR_MINUTE) AS f1 FROM t HAVING f1 != 0;
  SELECT 1;
  SELECT 'abc';
  ```

  For each sql statement in `dmlPath`, we will do some mutations according to Implication Oracle to detect logical bugs.

Note that:

* the paths in `taskconfig.json` are relative to `${IMPOHOME}`(for example, `./output` is actually `${IMPOHOME}/output`). You can also use absolute paths. Actually, we will automatically convert these paths to absolute paths before executing the `task`.
* we only focus on  `SELECT` statements in your `dmlPath`, which means we will ignore some of your sqls such as EXPLAIN, PREPARE.
* It is recommended not to use sql with side effects, such as assign operations, SELECT INTO.

#### output

You will see a new directory `${IMPOHOME}/output/mysql/task-1`. 

>  Actually we will remove the old directory and create a new directory.

There is a directory named `bugs` in `${IMPOHOME}/output/mysql/task-1`, and two files named `bug-0-0-FixMHaving1U.log` and `bug-0-0-FixMHaving1U.json` respectively in `bugs`.

We will save logical bugs in `bugs`. For each bug, we will create two files: bug-`bugId`-`sqlId`-`mutationName`.log and bug-`bugId`-`sqlId`-`mutationName`.json. `bugId` is the bug number(start from 0) during this task, `sqlId` is the original sql number(start from 0) in `dmlPath`, `mutationName` is the name of mutation.

* bug-`bugId`-`sqlId`-`mutationName`.log: save the mutation name, original sql, original result, mutated sql, mutated result, and the relationship between the original result and the mutated result we expect. For example:

  ```sql
  **************************************************
  [MutationName] FixMHaving1U
  **************************************************
  [IsUpper] true
  **************************************************
  [OriginalResult]
  ColumnName(ColumnType)s:  f1(DOUBLE)
  row 0: -1928.8181
  row 1: -1995.009
  row 2: -2007
  2.422742ms
  **************************************************
  [MutatedResult]
  ColumnName(ColumnType)s:  f1(DOUBLE)
  row 0: -20080524235820.816
  row 1: -20080524235887.008
  row 2: -20080524235899
  1.250519ms
  **************************************************
  
  -- OriginalSql
  SELECT `c1`-DATE_SUB(_UTF8MB4'2008-05-25', INTERVAL 1 HOUR_MINUTE) AS `f1` FROM `t` HAVING `f1`!=0;
  -- MutatedSql
  SELECT `c1`-DATE_SUB(_UTF8MB4'2008-05-25', INTERVAL 1 HOUR_MINUTE) AS `f1` FROM `t` HAVING 1;
  
  ```

  `[IsUpper] true` means that the mutated result  should ⊆ the original result. It is clear that the actual execution result violates this relationship.

  >  `[IsUpper] false` means that the original result should ⊆ the mutated result.

* bug-`bugId`-`sqlId`-`mutationName`.json: json format of bug-`bugId`-`sqlId`-`mutationName` .log exclude execution result. For example:

  ```json
  {
    "reportTime": "2022-11-13 23:26:33.51294115 +0800 CST m=+0.200207850",
    "bugId": 0,
    "sqlId": 0,
    "mutationName": "FixMHaving1U",
    "isUpper": true,
    "originalSql": "SELECT `c1`-DATE_SUB(_UTF8MB4'2008-05-25', INTERVAL 1 HOUR_MINUTE) AS `f1` FROM `t` HAVING `f1`!=0",
    "mutatedSql": "SELECT `c1`-DATE_SUB(_UTF8MB4'2008-05-25', INTERVAL 1 HOUR_MINUTE) AS `f1` FROM `t` HAVING 1"
  }
  ```

Additionally, there are two files in `${IMPOHOME}/output/mysql/task-1`:

* `task.log`: task log file, from which you can get task progress, task error during execution, and logic bugs.

* `result.json`: If the task executes successfully, you will get `result.json` like:

   ```json
   {
     "startTime": "2022-11-18 19:15:29.437554842 +0800 CST m=+0.001609271",
     "ddlSqlsNum": 2,
     "dmlSqlsNum": 3,
     "endInitTime": "2022-11-18 19:15:29.476101614 +0800 CST m=+0.040156084",
     "stage1ErrNum": 0,
     "stage1ExecErrNum": 0,
     "stage2ErrNum": 0,
     "stage2UnitNum": 5,
     "stage2UnitErrNum": 0,
     "stage2UnitExecErrNum": 0,
     "impoBugsNum": 1,
     "saveBugErrNum": 0,
     "endTime": "2022-11-18 19:15:29.478823678 +0800 CST m=+0.042878143"
   }
   ```
   
   This file is used for debugging, from which you can get the task's start time(`startTime`), end time(`endTime`), and the number of logical bugs we detected(`impoBugsNum`).

### 3.4 run task with go-randgen

A `task` can automatically generate `ddlPath` and `dmlPath` with the help of [go-randgen](https://github.com/pingcap/go-randgen), you need to build it first.

#### build go-randgen

```shell
git clone https://github.com/pingcap/go-randgen.git
cd go-randgen
go get -u github.com/jteeuwen/go-bindata/...
make all
```

Now you will see an executable file `go-randgen`, copy it to `${IMPOHOME}/resources`.

#### quick start

We assume you have executed `sudo docker run -itd --name mysqltest -p 13306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql:8.0.30`

```shell
cd ${IMPOHOME}
./impomysql task ./resources/taskrdgenconfig.json
```

#### input

Next, modify the configuration of `task`. see  `${IMPOHOME}/resources/taskrdgenconfig.json`:

```json
{
  "outputPath": "./output",
  "dbms": "mysql",
  "taskId": 1,
  "host": "127.0.0.1",
  "port": 13306,
  "username": "root",
  "password": "123456",
  "dbname": "TEST",
  "seed": 123456,
  "rdGenPath": "./resources/go-randgen",
  "zzPath": "./resources/impo.zz.lua",
  "yyPath": "./resources/impo.yy",
  "queriesNum": 100,
  "needDML": true
}
```

We removed `ddlPath` and `dmlPath`, added `randGenPath`, `zzPath`, `yyPath`, `queriesNum`, `needDML`:

* `randGenPath`: the path of your go-randgen executable file.

* `zzPath`, `yyPath`: `go-randgen`  will generate a ddl sql file `output.data.sql` according to `zzPath`, and generate a dml sql file  `output.rand.sql` according to `yyPath`. 

  We have provided a default zz file `impo.zz.lua` and a default yy file `impo.yy` in `${IMPOHOME}/resources`. It is recommended to use these default files.

* `queriesNum`: the number of sqls in `output.rand.sql`.

* `needDML`: if `needDML` is false, we will delete `output.rand.sql` at the end of `task` .  It is recommended to set this value to false, because the size of `output.rand.sql` is usually very large(about 10MB with 10000 sqls).

Note that:

* Similarly, the paths in `taskrdgenconfig.json` are relative to `${IMPOHOME}`. You can also use absolute paths. Actually, we will automatically convert these paths to absolute paths before executing the `task`.

* For go-randgen, we actually execute the following command:

  ```shell
  cd outputPath/dbms/task-taskId && randGenPath gentest -Z zzPath -Y yyPath -Q queriesNum --seed seed -B
  ```

* If you used both (non empty) `rdGenPath` and `ddlPath`, `dmlPath`, we will run `task` with `go-randgen`, and set `ddlPath` to  `outputPath/dbms/task-taskId/output.data.sql`, set `dmlPath` to `outputPath/dbms/task-taskId/output.rand.sql`.

#### output

In addition to `bugs`, `task.log`, `result.json`, you will also see `output.data.sql`, `output.rand.sql`.

Of course, if you set `needDML` to false, we will delete `output.rand.sql`.

### 3.5 run task pool

`taskpool` can continuously run tasks in parallel. Make sure you can run task with [go-randgen](https://github.com/pingcap/go-randgen).

#### quick start

We assume you have executed `sudo docker run -itd --name mysqltest -p 13306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql:8.0.30`

```shell
cd ${IMPOHOME}
./impomysql taskpool ./resources/taskpoolconfig.json
```

#### input

See `${IMPOHOME}/resources/taskpoolconfig.json`:

```json
{
  "outputPath": "./output",
  "dbms": "mysql",
  "host": "127.0.0.1",
  "port": 13306,
  "username": "root",
  "password": "123456",
  "dbPrefix": "TEST",
  "seed": 123456,
  "randGenPath": "./resources/go-randgen",
  "zzPath": "./resources/impo.zz.lua",
  "yyPath": "./resources/impo.yy",
  "queriesNum": 100,
  "threadNum": 4,
  "maxTasks": 16,
  "maxTimeS": 60
}
```

* `outputPath`,`dbms`,`host`,`port`,`username`,`password`,`randGenPath`,`zzPath`,`yyPath`,`queriesNum`: same as `task`
* `threadNum`: the number of threads(coroutines). 
* `maxTasks`:  maximum number of tasks, <= 0 means no limit.
* `maxTimeS`: maximum time(second), <=0 means no limit.
* `dbPrefix`: for each thread we will create a database connector, the dbname of each connector is `dbPrefix`+thread id.
* `seed`: the seed of each task is `seed`+task id.

Note that:

* `taskpool` will continuously run tasks with go-randgen in parallel, and we will set `needDML` to false.
* It is recommended to set `queriesNum` to a large value(>=10000, a task with `queriesNum`=10000 will take about 5~10 minutes), otherwise you will get a lot of task directories.
* We will stop all tasks if the database crashes.

#### output

In `${IMPOHOME}/output/mysql`, you will not only see the task directories, but also:

* task-`taskId`-config.json: the configuration file of task-`taskId`.

* `taskpool.log`:  taskpool log file, from which you can get taskpool progress, task error during execution, and logic bugs.

* `result.json`:  If the taskpool executes successfully, you will get `result.json` like:

  ```json
  {
    "startTime": "2022-11-18 19:33:20.290851516 +0800 CST m=+0.001550968",
    "totalTaskNum": 19,
    "finishedTaskNum": 16,
    "errorTaskNum": 0,
    "errorTaskIds": [],
    "bugsNum": 4,
    "bugTaskIds": [
      0,
      6,
      11
    ],
    "endTime": "2022-11-18 19:33:28.446037916 +0800 CST m=+8.156737393"
  }
  ```
  
  This file is used for debugging, from which you can get the taskpool's start time(`startTime`), end time(`endTime`), the number of logical bugs we detected(`bugsNum`) and their taskId(`bugTaskIds`).

