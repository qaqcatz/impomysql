# MySQL Document

https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html

# Overview

```shell
[1]: only tidb features
[2]: supported features of impo
[3]: unsupported features of impo

[3]Coalesce = "coalesce"
[3]Greatest = "greatest"
[3]Least    = "least"
[3]Interval = "interval"
```

# coalesce

Returns the first non-`NULL` value in the list, or `NULL` if there are no non-`NULL` values.

The return type of [`COALESCE()`](https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#function_coalesce) is the aggregated type of the argument types.

```sql
mysql> SELECT coalesce(null, null, 1, null);
+-------------------------------+
| coalesce(null, null, 1, null) |
+-------------------------------+
|                             1 |
+-------------------------------+
1 row in set (0.01 sec)

mysql> SELECT coalesce(null, null);
+----------------------+
| coalesce(null, null) |
+----------------------+
|                 NULL |
+----------------------+
1 row in set (0.00 sec)
```

> does not have the property of implication

# greatest

With two or more arguments, returns the largest (maximum-valued) argument. The arguments are compared using the same rules as for [`LEAST()`](https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#function_least).

```sql
mysql> SELECT GREATEST(2,0);
        -> 2
mysql> SELECT GREATEST(34.0,3.0,5.0,767.0);
        -> 767.0
mysql> SELECT GREATEST('B','A','C');
        -> 'C'
mysql> SELECT GREATEST(34.0,3.0,5.0,767.0, 'a', 'b') > 1;
+--------------------------------------------+
| GREATEST(34.0,3.0,5.0,767.0, 'a', 'b') > 1 |
+--------------------------------------------+
|                                          0 |
+--------------------------------------------+
1 row in set (0.00 sec)

mysql> SELECT GREATEST(34.0,3.0,5.0,767.0) > 1;
+----------------------------------+
| GREATEST(34.0,3.0,5.0,767.0) > 1 |
+----------------------------------+
|                                1 |
+----------------------------------+
1 row in set (0.00 sec)
```

> It is difficult to handle type conversion.

# least

- With two or more arguments, returns the smallest (minimum-valued) argument. The arguments are compared using the following rules:

  - If any argument is `NULL`, the result is `NULL`. No comparison is needed.
  - If all arguments are integer-valued, they are compared as integers.
  - If at least one argument is double precision, they are compared as double-precision values. Otherwise, if at least one argument is a [`DECIMAL`](https://dev.mysql.com/doc/refman/8.0/en/fixed-point-types.html) value, they are compared as [`DECIMAL`](https://dev.mysql.com/doc/refman/8.0/en/fixed-point-types.html) values.
  - If the arguments comprise a mix of numbers and strings, they are compared as strings.
  - If any argument is a nonbinary (character) string, the arguments are compared as nonbinary strings.
  - In all other cases, the arguments are compared as binary strings.

  The return type of [`LEAST()`](https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#function_least) is the aggregated type of the comparison argument types.

  ```sql
  mysql> SELECT LEAST(2,0);
          -> 0
  mysql> SELECT LEAST(34.0,3.0,5.0,767.0);
          -> 3.0
  mysql> SELECT LEAST('B','A','C');
          -> 'A'
  ```

> It is difficult to handle type conversion.

# interval

Returns `0` if *`N`* < *`N1`*, `1` if *`N`* < *`N2`* and so on or `-1` if *`N`* is `NULL`. All arguments are treated as integers. It is required that *`N1`* < *`N2`* < *`N3`* < `...` < *`Nn`* for this function to work correctly. This is because a binary search is used (very fast).

```sql
 mysql> SELECT INTERVAL(23, 1, 15, 17, 30, 44, 200);
        -> 3
mysql> SELECT INTERVAL(10, 1, 10, 100, 1000);
        -> 2
mysql> SELECT INTERVAL(22, 23, 30, 44, 200);
        -> 0
```

> does not have the property of implication

