# impomysql
Detecting Logic Bugs in MySQL through Implication Oracle.

## 1. What is logical bug?

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

## 2. What is Implication Oracle?

In the above example, we changed `HAVING f1 != 0`  to `HAVING 1`.

In theory, the predicate of sql1 → the predicate of sql2, and the result of sql1 ⊆ the result of sql2. 

If the actual result does not satisfy this relationship, we consider that there is a logical bug.

Although the idea is simple, some features make it difficult to implement, such as aggregate functions, window functions, type conversion, LIMIT, LEFT/RIGHT JOIN, flow control operations, etc.

We will discuss these features in our paper:

```
todo
```

## 3. How to use?

### 3.1 As a tool



### 3.2 As a go package



## 4. How to read the source code?



If you want to test all modules, you should use serial test: `-p 1`