### Conversion Between Date and Time Types

https://dev.mysql.com/doc/refman/8.0/en/date-and-time-type-conversion.html

To some extent, you can convert a value from one temporal type to another. However, there may be some alteration of the value or loss of information. In all cases, conversion between temporal types is subject to the range of valid values for the resulting type. For example, although [`DATE`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html), [`DATETIME`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html), and [`TIMESTAMP`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) values all can be specified using the same set of formats, the types do not all have the same range of values. [`TIMESTAMP`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) values cannot be earlier than `1970` UTC or later than `'2038-01-19 03:14:07'` UTC. This means that a date such as `'1968-01-01'`, while valid as a [`DATE`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) or [`DATETIME`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) value, is not valid as a [`TIMESTAMP`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) value and is converted to `0`.

Conversion of [`DATE`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) values:

- Conversion to a [`DATETIME`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) or [`TIMESTAMP`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) value adds a time part of `'00:00:00'` because the [`DATE`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) value contains no time information.
- Conversion to a [`TIME`](https://dev.mysql.com/doc/refman/8.0/en/time.html) value is not useful; the result is `'00:00:00'`.

Conversion of [`DATETIME`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) and [`TIMESTAMP`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) values:

- Conversion to a [`DATE`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) value takes fractional seconds into account and rounds the time part. For example, `'1999-12-31 23:59:59.499'` becomes `'1999-12-31'`, whereas `'1999-12-31 23:59:59.500'` becomes `'2000-01-01'`.
- Conversion to a [`TIME`](https://dev.mysql.com/doc/refman/8.0/en/time.html) value discards the date part because the [`TIME`](https://dev.mysql.com/doc/refman/8.0/en/time.html) type contains no date information.

For conversion of [`TIME`](https://dev.mysql.com/doc/refman/8.0/en/time.html) values to other temporal types, the value of [`CURRENT_DATE()`](https://dev.mysql.com/doc/refman/8.0/en/date-and-time-functions.html#function_current-date) is used for the date part. The [`TIME`](https://dev.mysql.com/doc/refman/8.0/en/time.html) is interpreted as elapsed time (not time of day) and added to the date. This means that the date part of the result differs from the current date if the time value is outside the range from `'00:00:00'` to `'23:59:59'`.

Suppose that the current date is `'2012-01-01'`. [`TIME`](https://dev.mysql.com/doc/refman/8.0/en/time.html) values of `'12:00:00'`, `'24:00:00'`, and `'-12:00:00'`, when converted to [`DATETIME`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) or [`TIMESTAMP`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) values, result in `'2012-01-01 12:00:00'`, `'2012-01-02 00:00:00'`, and `'2011-12-31 12:00:00'`, respectively.

Conversion of [`TIME`](https://dev.mysql.com/doc/refman/8.0/en/time.html) to [`DATE`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) is similar but discards the time part from the result: `'2012-01-01'`, `'2012-01-02'`, and `'2011-12-31'`, respectively.

Explicit conversion can be used to override implicit conversion. For example, in comparison of [`DATE`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) and [`DATETIME`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) values, the [`DATE`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) value is coerced to the [`DATETIME`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) type by adding a time part of `'00:00:00'`. To perform the comparison by ignoring the time part of the [`DATETIME`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) value instead, use the [`CAST()`](https://dev.mysql.com/doc/refman/8.0/en/cast-functions.html#function_cast) function in the following way:

```sql
date_col = CAST(datetime_col AS DATE)
```

Conversion of [`TIME`](https://dev.mysql.com/doc/refman/8.0/en/time.html) and [`DATETIME`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) values to numeric form (for example, by adding `+0`) depends on whether the value contains a fractional seconds part. [`TIME(*`N`*)`](https://dev.mysql.com/doc/refman/8.0/en/time.html) or [`DATETIME(*`N`*)`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) is converted to integer when *`N`* is 0 (or omitted) and to a `DECIMAL` value with *`N`* decimal digits when *`N`* is greater than 0:

```sql
 mysql> SELECT CURTIME(), CURTIME()+0, CURTIME(3)+0;
+-----------+-------------+--------------+
| CURTIME() | CURTIME()+0 | CURTIME(3)+0 |
+-----------+-------------+--------------+
| 09:28:00  |       92800 |    92800.887 |
+-----------+-------------+--------------+
mysql> SELECT NOW(), NOW()+0, NOW(3)+0;
+---------------------+----------------+--------------------+
| NOW()               | NOW()+0        | NOW(3)+0           |
+---------------------+----------------+--------------------+
| 2012-08-15 09:28:00 | 20120815092800 | 20120815092800.889 |
+---------------------+----------------+--------------------+
```