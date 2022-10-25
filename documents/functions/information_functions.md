#### Type Specification

+ 0 fix value

+ 1 exist branch structure

+ 2 exist random value
+ 3 depends on machine performance 
+ 4 depends on the most recently executed [`INSERT`](https://dev.mysql.com/doc/refman/8.0/en/insert.html) statement

+ *_none

#### Information Function

|      func name       |   type    |                         description                          |    parameters    |
| :------------------: | :-------: | :----------------------------------------------------------: | :--------------: |
|      BENCHMARK       |     3     | The [`BENCHMARK()`](https://dev.mysql.com/doc/refman/8.0/en/information-functions.html#function_benchmark) function executes the expression *`expr`* repeatedly *`count`* times. | (*`count,expr`*) |
|       CHARSET        |     0     | Returns the character set of the string argument, or `NULL` if the argument is `NULL`. |    (*`str`*)     |
|     COERCIBILITY     |     1     | Returns the collation coercibility value of the string argument. |    (*`str`*)     |
|      COLLATION       |     0     |        Returns the collation of the string argument.         |    (*`str`*)     |
|    CONNECTION_ID     | uncertain | Returns the connection ID (thread ID) for the connection. Every connection has an ID that is unique among the set of currently connected clients. |        ()        |
|     CURRENT_USER     | uncertain | Returns the user name and host name combination for the MySQL account that the server used to authenticate the current client. |        ()        |
|     CURRENT_ROLE     | uncertain | Returns a string containing the current active roles for the current session, separated by commas, or if there are none. |        ()        |
|       DATABASE       | uncertain | Returns the default (current) database name as a string in the character set. |        ()        |
|      FOUND_ROWS      | uncertain | The query modifier and accompanying [`FOUND_ROWS()`](https://dev.mysql.com/doc/refman/8.0/en/information-functions.html#function_found-rows) function |        ()        |
|    LAST_INSERT_ID    | uncertain | With no argument, [`LAST_INSERT_ID()`](https://dev.mysql.com/doc/refman/8.0/en/information-functions.html#function_last-insert-id) returns a (64-bit) value representing the first automatically generated value successfully inserted for an column as a result of the most recently executed [`INSERT`](https://dev.mysql.com/doc/refman/8.0/en/insert.html) statement. |  (), (*`expr`*)  |
|      ROW_COUNT       | uncertain | The [`ROW_COUNT()`](https://dev.mysql.com/doc/refman/8.0/en/information-functions.html#function_row-count) value is similar to the value from the [`mysql_affected_rows()`](https://dev.mysql.com/doc/c-api/8.0/en/mysql-affected-rows.html) C API function and the row count that the [**mysql**](https://dev.mysql.com/doc/refman/8.0/en/mysql.html) client displays following statement execution. |        ()        |
|        SCHEMA        | uncertain | This function is a synonym for [`DATABASE()`](https://dev.mysql.com/doc/refman/8.0/en/information-functions.html#function_database). |        ()        |
|     SESSION_USER     | uncertain | [`SESSION_USER()`](https://dev.mysql.com/doc/refman/8.0/en/information-functions.html#function_session-user) is a synonym for [`USER()`](https://dev.mysql.com/doc/refman/8.0/en/information-functions.html#function_user). |        ()        |
|     SYSTEM_USER      | uncertain | [`SYSTEM_USER()`](https://dev.mysql.com/doc/refman/8.0/en/information-functions.html#function_system-user) is a synonym for [`USER()`](https://dev.mysql.com/doc/refman/8.0/en/information-functions.html#function_user). |        ()        |
|         USER         | uncertain | Returns the current MySQL user name and host name as a string in the character set. `utf8mb3` |        ()        |
|       VERSION        | uncertain |  Returns a string that indicates the MySQL server version.   |        ()        |
|     TiDBVersion      |     *     |                              *                               |        *         |
|    TiDBIsDDLOwner    |     *     |                              *                               |        *         |
|    TiDBDecodePlan    |     *     |                              *                               |        *         |
| TiDBDecodeSQLDigests |     *     |                              *                               |        *         |
|     format_bytes     |     0     | Given a byte count, converts it to human-readable format and returns a string consisting of a value and a units indicator. |       (N)        |
|    FormatNanoTime    |     *     |                              *                               |        *         |