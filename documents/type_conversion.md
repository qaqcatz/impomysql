## Type Conversion in Expression Evaluation

https://dev.mysql.com/doc/refman/8.0/en/type-conversion.html

When an operator is used with operands of different types, type conversion occurs to make the operands compatible. Some conversions occur implicitly. For example, MySQL automatically converts strings to numbers as necessary, and vice versa.

```sql
mysql> SELECT 1+'1';
        -> 2
mysql> SELECT CONCAT(2,' test');
        -> '2 test'
```

It is also possible to convert a number to a string explicitly using the [`CAST()`](https://dev.mysql.com/doc/refman/8.0/en/cast-functions.html#function_cast) function. Conversion occurs implicitly with the [`CONCAT()`](https://dev.mysql.com/doc/refman/8.0/en/string-functions.html#function_concat) function because it expects string arguments.

```sql
mysql> SELECT 38.8, CAST(38.8 AS CHAR);
        -> 38.8, '38.8'
mysql> SELECT 38.8, CONCAT(38.8);
        -> 38.8, '38.8'
```

See later in this section for information about the character set of implicit number-to-string conversions, and for modified rules that apply to `CREATE TABLE ... SELECT` statements.

The following rules describe how conversion occurs for comparison operations:

- If one or both arguments are `NULL`, the result of the comparison is `NULL`, except for the `NULL`-safe [`<=>`](https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_equal-to) equality comparison operator. For `NULL <=> NULL`, the result is true. No conversion is needed.

- If both arguments in a comparison operation are strings, they are compared as strings.

- If both arguments are integers, they are compared as integers.

- Hexadecimal values are treated as binary strings if not compared to a number.

- If one of the arguments is a [`TIMESTAMP`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) or [`DATETIME`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) column and the other argument is a constant, the constant is converted to a timestamp before the comparison is performed. This is done to be more ODBC-friendly. This is not done for the arguments to [`IN()`](https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_in). To be safe, always use complete datetime, date, or time strings when doing comparisons. For example, to achieve best results when using [`BETWEEN`](https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_between) with date or time values, use [`CAST()`](https://dev.mysql.com/doc/refman/8.0/en/cast-functions.html#function_cast) to explicitly convert the values to the desired data type.

  A single-row subquery from a table or tables is not considered a constant. For example, if a subquery returns an integer to be compared to a [`DATETIME`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) value, the comparison is done as two integers. The integer is not converted to a temporal value. To compare the operands as [`DATETIME`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html) values, use [`CAST()`](https://dev.mysql.com/doc/refman/8.0/en/cast-functions.html#function_cast) to explicitly convert the subquery value to [`DATETIME`](https://dev.mysql.com/doc/refman/8.0/en/datetime.html).

- If one of the arguments is a decimal value, comparison depends on the other argument. The arguments are compared as decimal values if the other argument is a decimal or integer value, or as floating-point values if the other argument is a floating-point value.

- In all other cases, the arguments are compared as floating-point (double-precision) numbers. For example, a comparison of string and numeric operands takes place as a comparison of floating-point numbers.

For information about conversion of values from one temporal type to another, see [Section 11.2.7, “Conversion Between Date and Time Types”](https://dev.mysql.com/doc/refman/8.0/en/date-and-time-type-conversion.html).

Comparison of JSON values takes place at two levels. The first level of comparison is based on the JSON types of the compared values. If the types differ, the comparison result is determined solely by which type has higher precedence. If the two values have the same JSON type, a second level of comparison occurs using type-specific rules. For comparison of JSON and non-JSON values, the non-JSON value is converted to JSON and the values compared as JSON values. For details, see [Comparison and Ordering of JSON Values](https://dev.mysql.com/doc/refman/8.0/en/json.html#json-comparison).

The following examples illustrate conversion of strings to numbers for comparison operations:

```sql
mysql> SELECT 1 > '6x';
        -> 0
mysql> SELECT 7 > '6x';
        -> 1
mysql> SELECT 0 > 'x6';
        -> 0
mysql> SELECT 0 = 'x6';
        -> 1
```

For comparisons of a string column with a number, MySQL cannot use an index on the column to look up the value quickly. If *`str_col`* is an indexed string column, the index cannot be used when performing the lookup in the following statement:

```sql
SELECT * FROM tbl_name WHERE str_col=1;
```

The reason for this is that there are many different strings that may convert to the value `1`, such as `'1'`, `' 1'`, or `'1a'`.

Comparisons between floating-point numbers and large values of `INTEGER` type are approximate because the integer is converted to double-precision floating point before comparison, which is not capable of representing all 64-bit integers exactly. For example, the integer value 253 + 1 is not representable as a float, and is rounded to 253 or 253 + 2 before a float comparison, depending on the platform.

To illustrate, only the first of the following comparisons compares equal values, but both comparisons return true (1):

```sql
mysql> SELECT '9223372036854775807' = 9223372036854775807;
        -> 1
mysql> SELECT '9223372036854775807' = 9223372036854775806;
        -> 1
```

When conversions from string to floating-point and from integer to floating-point occur, they do not necessarily occur the same way. The integer may be converted to floating-point by the CPU, whereas the string is converted digit by digit in an operation that involves floating-point multiplications. Also, results can be affected by factors such as computer architecture or the compiler version or optimization level. One way to avoid such problems is to use [`CAST()`](https://dev.mysql.com/doc/refman/8.0/en/cast-functions.html#function_cast) so that a value is not converted implicitly to a float-point number:

```sql
mysql> SELECT CAST('9223372036854775807' AS UNSIGNED) = 9223372036854775806;
        -> 0
```

For more information about floating-point comparisons, see [Section B.3.4.8, “Problems with Floating-Point Values”](https://dev.mysql.com/doc/refman/8.0/en/problems-with-float.html).

The server includes `dtoa`, a conversion library that provides the basis for improved conversion between string or [`DECIMAL`](https://dev.mysql.com/doc/refman/8.0/en/fixed-point-types.html) values and approximate-value ([`FLOAT`](https://dev.mysql.com/doc/refman/8.0/en/floating-point-types.html)/[`DOUBLE`](https://dev.mysql.com/doc/refman/8.0/en/floating-point-types.html)) numbers:

- Consistent conversion results across platforms, which eliminates, for example, Unix versus Windows conversion differences.
- Accurate representation of values in cases where results previously did not provide sufficient precision, such as for values close to IEEE limits.
- Conversion of numbers to string format with the best possible precision. The precision of `dtoa` is always the same or better than that of the standard C library functions.

Because the conversions produced by this library differ in some cases from non-`dtoa` results, the potential exists for incompatibilities in applications that rely on previous results. For example, applications that depend on a specific exact result from previous conversions might need adjustment to accommodate additional precision.

The `dtoa` library provides conversions with the following properties. *`D`* represents a value with a [`DECIMAL`](https://dev.mysql.com/doc/refman/8.0/en/fixed-point-types.html) or string representation, and *`F`* represents a floating-point number in native binary (IEEE) format.

- *`F`* -> *`D`* conversion is done with the best possible precision, returning *`D`* as the shortest string that yields *`F`* when read back in and rounded to the nearest value in native binary format as specified by IEEE.
- *`D`* -> *`F`* conversion is done such that *`F`* is the nearest native binary number to the input decimal string *`D`*.

These properties imply that *`F`* -> *`D`* -> *`F`* conversions are lossless unless *`F`* is `-inf`, `+inf`, or `NaN`. The latter values are not supported because the SQL standard defines them as invalid values for [`FLOAT`](https://dev.mysql.com/doc/refman/8.0/en/floating-point-types.html) or [`DOUBLE`](https://dev.mysql.com/doc/refman/8.0/en/floating-point-types.html).

For *`D`* -> *`F`* -> *`D`* conversions, a sufficient condition for losslessness is that *`D`* uses 15 or fewer digits of precision, is not a denormal value, `-inf`, `+inf`, or `NaN`. In some cases, the conversion is lossless even if *`D`* has more than 15 digits of precision, but this is not always the case.

Implicit conversion of a numeric or temporal value to string produces a value that has a character set and collation determined by the [`character_set_connection`](https://dev.mysql.com/doc/refman/8.0/en/server-system-variables.html#sysvar_character_set_connection) and [`collation_connection`](https://dev.mysql.com/doc/refman/8.0/en/server-system-variables.html#sysvar_collation_connection) system variables. (These variables commonly are set with [`SET NAMES`](https://dev.mysql.com/doc/refman/8.0/en/set-names.html). For information about connection character sets, see [Section 10.4, “Connection Character Sets and Collations”](https://dev.mysql.com/doc/refman/8.0/en/charset-connection.html).)

This means that such a conversion results in a character (nonbinary) string (a [`CHAR`](https://dev.mysql.com/doc/refman/8.0/en/char.html), [`VARCHAR`](https://dev.mysql.com/doc/refman/8.0/en/char.html), or [`LONGTEXT`](https://dev.mysql.com/doc/refman/8.0/en/blob.html) value), except in the case that the connection character set is set to `binary`. In that case, the conversion result is a binary string (a [`BINARY`](https://dev.mysql.com/doc/refman/8.0/en/binary-varbinary.html), [`VARBINARY`](https://dev.mysql.com/doc/refman/8.0/en/binary-varbinary.html), or [`LONGBLOB`](https://dev.mysql.com/doc/refman/8.0/en/blob.html) value).

For integer expressions, the preceding remarks about expression *evaluation* apply somewhat differently for expression *assignment*; for example, in a statement such as this:

```sql
CREATE TABLE t SELECT integer_expr;
```

In this case, the table in the column resulting from the expression has type [`INT`](https://dev.mysql.com/doc/refman/8.0/en/integer-types.html) or [`BIGINT`](https://dev.mysql.com/doc/refman/8.0/en/integer-types.html) depending on the length of the integer expression. If the maximum length of the expression does not fit in an [`INT`](https://dev.mysql.com/doc/refman/8.0/en/integer-types.html), [`BIGINT`](https://dev.mysql.com/doc/refman/8.0/en/integer-types.html) is used instead. The length is taken from the `max_length` value of the [`SELECT`](https://dev.mysql.com/doc/refman/8.0/en/select.html) result set metadata (see [C API Basic Data Structures](https://dev.mysql.com/doc/c-api/8.0/en/c-api-data-structures.html)). This means that you can force a [`BIGINT`](https://dev.mysql.com/doc/refman/8.0/en/integer-types.html) rather than [`INT`](https://dev.mysql.com/doc/refman/8.0/en/integer-types.html) by use of a sufficiently long expression:

```sql
CREATE TABLE t SELECT 000000000000000000000;
```