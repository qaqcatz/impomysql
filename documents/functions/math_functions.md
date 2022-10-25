# MySQL Document

https://dev.mysql.com/doc/refman/8.0/en/numeric-functions.html

# Overview

```shell
[1]: only tidb features
[2]: supported features of impo
[3]: unsupported features of impo
[4]: todo
	
[2]Abs      = "abs"
[2]Acos     = "acos"
[2]Asin     = "asin"
[2]Atan     = "atan"
[2]Atan2    = "atan2"
[2]Ceil     = "ceil"
[2]Ceiling  = "ceiling"
[4]Conv     = "conv"
[2]Cos      = "cos"
[2]Cot      = "cot"
[2]CRC32    = "crc32"
[2]Degrees  = "degrees"
[4]Exp      = "exp"
[2]Floor    = "floor"
[2]Ln       = "ln"
[2]Log      = "log"
[2]Log2     = "log2"
[2]Log10    = "log10"
[2]PI       = "pi"
[2]Pow      = "pow"
[2]Power    = "power"
[2]Radians  = "radians"
[3]Rand     = "rand"
[2]Round    = "round"
[2]Sign     = "sign"
[2]Sin      = "sin"
[2]Sqrt     = "sqrt"
[2]Tan      = "tan"
[4]Truncate = "truncate"
```

# abs

Returns the absolute value of *`X`*, or `NULL` if *`X`* is `NULL`.

The result type is derived from the argument type. An implication of this is that [`ABS(-9223372036854775808)`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_abs) produces an error because the result cannot be stored in a signed `BIGINT` value.

```sql
mysql> SELECT ABS(2);
        -> 2
mysql> SELECT ABS(-32);
        -> 32
```

This function is safe to use with [`BIGINT`](https://dev.mysql.com/doc/refman/8.0/en/integer-types.html) values.

# acos

Returns the arc cosine of *`X`*, that is, the value whose cosine is *`X`*. Returns `NULL` if *`X`* is not in the range `-1` to `1`, or if *`X`* is `NULL`.

```sql
mysql> SELECT ACOS(1);
        -> 0
mysql> SELECT ACOS(1.0001);
        -> NULL
mysql> SELECT ACOS(0);
        -> 1.5707963267949
```

# asin

Returns the arc sine of *`X`*, that is, the value whose sine is *`X`*. Returns `NULL` if *`X`* is not in the range `-1` to `1`, or if *`X`* is `NULL`.

```sql
mysql> SELECT ASIN(0.2);
        -> 0.20135792079033
mysql> SELECT ASIN('foo');

+-------------+
| ASIN('foo') |
+-------------+
|           0 |
+-------------+
1 row in set, 1 warning (0.00 sec)

mysql> SHOW WARNINGS;
+---------+------+-----------------------------------------+
| Level   | Code | Message                                 |
+---------+------+-----------------------------------------+
| Warning | 1292 | Truncated incorrect DOUBLE value: 'foo' |
+---------+------+-----------------------------------------+
```

# atan

Returns the arc tangent of *`X`*, that is, the value whose tangent is *`X`*. Returns *`NULL`* if *`X`* is `NULL`

```sql
mysql> SELECT ATAN(2);
        -> 1.1071487177941
mysql> SELECT ATAN(-2);
        -> -1.1071487177941
```

# atan2

Returns the arc tangent of the two variables *`X`* and *`Y`*. It is similar to calculating the arc tangent of `*`Y`* / *`X`*`, except that the signs of both arguments are used to determine the quadrant of the result. Returns `NULL` if *`X`* or *`Y`* is `NULL`.

```sql
mysql> SELECT ATAN(-2,2);
        -> -0.78539816339745
mysql> SELECT ATAN2(PI(),0);
        -> 1.5707963267949
```

# ceil

[`CEIL()`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_ceil) is a synonym for [`CEILING()`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_ceiling).

# ceiling

Returns the smallest integer value not less than *`X`*. Returns `NULL` if *`X`* is `NULL`.

```sql
mysql> SELECT CEILING(1.23);
        -> 2
mysql> SELECT CEILING(-1.23);
        -> -1
```

For exact-value numeric arguments, the return value has an exact-value numeric type. For string or floating-point arguments, the return value has a floating-point type.

# conv

[`CONV(*`N`*,*`from_base`*,*`to_base`*)`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_conv)

Converts numbers between different number bases. Returns a string representation of the number *`N`*, converted from base *`from_base`* to base *`to_base`*. Returns `NULL` if any argument is `NULL`. The argument *`N`* is interpreted as an integer, but may be specified as an integer or a string. The minimum base is `2` and the maximum base is `36`. If *`from_base`* is a negative number, *`N`* is regarded as a signed number. Otherwise, *`N`* is treated as unsigned. [`CONV()`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_conv) works with 64-bit precision.

`CONV()` returns `NULL` if any of its arguments are `NULL`.

```sql
mysql> SELECT CONV('a',16,2);
        -> '1010'
mysql> SELECT CONV('6E',18,8);
        -> '172'
mysql> SELECT CONV(-17,10,-18);
        -> '-H'
mysql> SELECT CONV(10+'10'+'10'+X'0a',10,10);
        -> '40'
```

> todo

# cos

Returns the cosine of *`X`*, where *`X`* is given in radians. Returns `NULL` if *`X`* is `NULL`.

```sql
mysql> SELECT COS(PI());
        -> -1
```

# cot

Returns the cotangent of *`X`*. Returns `NULL` if *`X`* is `NULL`.

```sql
mysql> SELECT COT(12);
        -> -1.5726734063977
mysql> SELECT COT(0);
        -> out-of-range error
```

# crc32

Computes a cyclic redundancy check value and returns a 32-bit unsigned value. The result is `NULL` if the argument is `NULL`. The argument is expected to be a string and (if possible) is treated as one if it is not.

```sql
mysql> SELECT CRC32('MySQL');
        -> 3259397556
mysql> SELECT CRC32('mysql');
        -> 2501908538
```

# degrees

Returns the argument *`X`*, converted from radians to degrees. Returns `NULL` if *`X`* is `NULL`.

```sql
mysql> SELECT DEGREES(PI());
        -> 180
mysql> SELECT DEGREES(PI() / 2);
        -> 90
```

# exp

Returns the value of *e* (the base of natural logarithms) raised to the power of *`X`*. The inverse of this function is [`LOG()`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_log) (using a single argument only) or [`LN()`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_ln).

If *`X`* is `NULL`, this function returns `NULL`.

```sql
mysql> SELECT EXP(2);
        -> 7.3890560989307
mysql> SELECT EXP(-2);
        -> 0.13533528323661
mysql> SELECT EXP(0);
        -> 1
```

> todo (DOUBLE value is out of range)

# floor

Returns the largest integer value not greater than *`X`*. Returns `NULL` if *`X`* is `NULL`.

```sql
mysql> SELECT FLOOR(1.23), FLOOR(-1.23);
        -> 1, -2
```

For exact-value numeric arguments, the return value has an exact-value numeric type. For string or floating-point arguments, the return value has a floating-point type.

# ln

Returns the natural logarithm of *`X`*; that is, the base-*e* logarithm of *`X`*. If *`X`* is less than or equal to 0.0E0, the function returns `NULL` and a warning “Invalid argument for logarithm” is reported. Returns `NULL` if *`X`* is `NULL`.

```sql
mysql> SELECT LN(2);
        -> 0.69314718055995
mysql> SELECT LN(-2);
        -> NULL
```

This function is synonymous with [`LOG(*`X`*)`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_log). The inverse of this function is the [`EXP()`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_exp) function.

# log

If called with one parameter, this function returns the natural logarithm of *`X`*. If *`X`* is less than or equal to 0.0E0, the function returns `NULL` and a warning “Invalid argument for logarithm” is reported. Returns `NULL` if *`X`* or *`B`* is `NULL`.

The inverse of this function (when called with a single argument) is the [`EXP()`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_exp) function.

```sql
mysql> SELECT LOG(2);
        -> 0.69314718055995
mysql> SELECT LOG(-2);
        -> NULL
```

If called with two parameters, this function returns the logarithm of *`X`* to the base *`B`*. If *`X`* is less than or equal to 0, or if *`B`* is less than or equal to 1, then `NULL` is returned.

```sql
mysql> SELECT LOG(2,65536);
        -> 16
mysql> SELECT LOG(10,100);
        -> 2
mysql> SELECT LOG(1,100);
        -> NULL
```

[`LOG(*`B`*,*`X`*)`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_log) is equivalent to [`LOG(*`X`*) / LOG(*`B`*)`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_log).

# log2

Returns the base-2 logarithm of `*`X`*`. If *`X`* is less than or equal to 0.0E0, the function returns `NULL` and a warning “Invalid argument for logarithm” is reported. Returns `NULL` if *`X`* is `NULL`.

```sql
mysql> SELECT LOG2(65536);
        -> 16
mysql> SELECT LOG2(-100);
        -> NULL
```

[`LOG2()`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_log2) is useful for finding out how many bits a number requires for storage. This function is equivalent to the expression [`LOG(*`X`*) / LOG(2)`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_log).

# log10

Returns the base-10 logarithm of *`X`*. If *`X`* is less than or equal to 0.0E0, the function returns `NULL` and a warning “Invalid argument for logarithm” is reported. Returns `NULL` if *`X`* is `NULL`.

```sql
mysql> SELECT LOG10(2);
        -> 0.30102999566398
mysql> SELECT LOG10(100);
        -> 2
mysql> SELECT LOG10(-100);
        -> NULL
```

[`LOG10(*`X`*)`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_log10) is equivalent to [`LOG(10,*`X`*)`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_log).

# pi

Returns the value of π (pi). The default number of decimal places displayed is seven, but MySQL uses the full double-precision value internally.

```sql
mysql> SELECT PI();
        -> 3.141593
mysql> SELECT PI()+0.000000000000000000;
        -> 3.141592653589793116
```

# pow

Returns the value of *`X`* raised to the power of *`Y`*. Returns `NULL` if *`X`* or *`Y`* is `NULL`.

```sql
mysql> SELECT POW(2,2);
        -> 4
mysql> SELECT POW(2,-2);
        -> 0.25
```

> todo (DOUBLE value is out of range)

# power

This is a synonym for [`POW()`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_pow).

> todo (DOUBLE value is out of range)

# radians

Returns the argument *`X`*, converted from degrees to radians. (Note that π radians equals 180 degrees.) Returns `NULL` if *`X`* is `NULL`.

```sql
mysql> SELECT RADIANS(90);
        -> 1.5707963267949
```

# rand

Returns a random floating-point value *`v`* in the range `0` <= *`v`* < `1.0`. To obtain a random integer *`R`* in the range *`i`* <= *`R`* < *`j`*, use the expression [`FLOOR(*`i`* + RAND() * (*`j`*`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_floor) − `*`i`*))`. For example, to obtain a random integer in the range the range `7` <= *`R`* < `12`, use the following statement:

```sql
SELECT FLOOR(7 + (RAND() * 5));
```

> uncertain

# round

[`ROUND(*`X`*)`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_round), [`ROUND(*`X`*,*`D`*)`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_round)

Rounds the argument *`X`* to *`D`* decimal places. The rounding algorithm depends on the data type of *`X`*. *`D`* defaults to 0 if not specified. *`D`* can be negative to cause *`D`* digits left of the decimal point of the value *`X`* to become zero. The maximum absolute value for *`D`* is 30; any digits in excess of 30 (or -30) are truncated. If *`X`* or *`D`* is `NULL`, the function returns `NULL`.

```sql
mysql> SELECT ROUND(-1.23);
        -> -1
mysql> SELECT ROUND(-1.58);
        -> -2
mysql> SELECT ROUND(1.58);
        -> 2
mysql> SELECT ROUND(1.298, 1);
        -> 1.3
mysql> SELECT ROUND(1.298, 0);
        -> 1
mysql> SELECT ROUND(23.298, -1);
        -> 20
mysql> SELECT ROUND(.12345678901234567890123456789012345, 35);
        -> 0.123456789012345678901234567890
```

The return value has the same type as the first argument (assuming that it is integer, double, or decimal). This means that for an integer argument, the result is an integer (no decimal places):

```sql
mysql> SELECT ROUND(150.000,2), ROUND(150,2);
+------------------+--------------+
| ROUND(150.000,2) | ROUND(150,2) |
+------------------+--------------+
|           150.00 |          150 |
+------------------+--------------+
```

[`ROUND()`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_round) uses the following rules depending on the type of the first argument:

- For exact-value numbers, [`ROUND()`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_round) uses the “round half away from zero” or “round toward nearest” rule: A value with a fractional part of .5 or greater is rounded up to the next integer if positive or down to the next integer if negative. (In other words, it is rounded away from zero.) A value with a fractional part less than .5 is rounded down to the next integer if positive or up to the next integer if negative.
- For approximate-value numbers, the result depends on the C library. On many systems, this means that [`ROUND()`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_round) uses the “round to nearest even” rule: A value with a fractional part exactly halfway between two integers is rounded to the nearest even integer.

The following example shows how rounding differs for exact and approximate values:

```sql
mysql> SELECT ROUND(2.5), ROUND(25E-1);
+------------+--------------+
| ROUND(2.5) | ROUND(25E-1) |
+------------+--------------+
| 3          |            2 |
+------------+--------------+
```

For more information, see [Section 12.25, “Precision Math”](https://dev.mysql.com/doc/refman/8.0/en/precision-math.html).

In MySQL 8.0.21 and later, the data type returned by `ROUND()` (and [`TRUNCATE()`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_truncate)) is determined according to the rules listed here:

- When the first argument is of any integer type, the return type is always [`BIGINT`](https://dev.mysql.com/doc/refman/8.0/en/integer-types.html).

- When the first argument is of any floating-point type or of any non-numeric type, the return type is always [`DOUBLE`](https://dev.mysql.com/doc/refman/8.0/en/floating-point-types.html).

- When the first argument is a [`DECIMAL`](https://dev.mysql.com/doc/refman/8.0/en/fixed-point-types.html) value, the return type is also `DECIMAL`.

- The type attributes for the return value are also copied from the first argument, except in the case of `DECIMAL`, when the second argument is a constant value.

  When the desired number of decimal places is less than the scale of the argument, the scale and the precision of the result are adjusted accordingly.

  In addition, for `ROUND()` (but not for the [`TRUNCATE()`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_truncate) function), the precision is extended by one place to accommodate rounding that increases the number of significant digits. If the second argument is negative, the return type is adjusted such that its scale is 0, with a corresponding precision. For example, `ROUND(99.999, 2)` returns `100.00`—the first argument is `DECIMAL(5, 3)`, and the return type is `DECIMAL(5, 2)`.

  If the second argument is negative, the return type has scale 0 and a corresponding precision; `ROUND(99.999, -1)` returns `100`, which is `DECIMAL(3, 0)`.

# sign

Returns the sign of the argument as `-1`, `0`, or `1`, depending on whether *`X`* is negative, zero, or positive. Returns `NULL` if *`X`* is `NULL`.

```sql
mysql> SELECT SIGN(-32);
        -> -1
mysql> SELECT SIGN(0);
        -> 0
mysql> SELECT SIGN(234);
        -> 1
```

# sin

Returns the sine of *`X`*, where *`X`* is given in radians. Returns `NULL` if *`X`* is `NULL`.

```sql
mysql> SELECT SIN(PI());
        -> 1.2246063538224e-16
mysql> SELECT ROUND(SIN(PI()));
        -> 0
```

# sqrt

Returns the square root of a nonnegative number *`X`*. If *`X`* is `NULL`, the function returns `NULL`.

```sql
mysql> SELECT SQRT(4);
        -> 2
mysql> SELECT SQRT(20);
        -> 4.4721359549996
mysql> SELECT SQRT(-16);
        -> NULL
```

# tan

Returns the tangent of *`X`*, where *`X`* is given in radians. Returns `NULL` if *`X`* is `NULL`.

```sql
mysql> SELECT TAN(PI());
        -> -1.2246063538224e-16
mysql> SELECT TAN(PI()+1);
        -> 1.5574077246549
```

# truncate

[`TRUNCATE(*`X`*,*`D`*)`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_truncate)

Returns the number *`X`*, truncated to *`D`* decimal places. If *`D`* is `0`, the result has no decimal point or fractional part. *`D`* can be negative to cause *`D`* digits left of the decimal point of the value *`X`* to become zero. If *`X`* or *`D`* is `NULL`, the function returns `NULL`.

```sql
mysql> SELECT TRUNCATE(1.223,1);
        -> 1.2
mysql> SELECT TRUNCATE(1.999,1);
        -> 1.9
mysql> SELECT TRUNCATE(1.999,0);
        -> 1
mysql> SELECT TRUNCATE(-1.999,1);
        -> -1.9
mysql> SELECT TRUNCATE(122,-2);
       -> 100
mysql> SELECT TRUNCATE(10.28*100,0);
       -> 1028
```

All numbers are rounded toward zero.

In MySQL 8.0.21 and later, the data type returned by `TRUNCATE()` follows the same rules that determine the return type of the `ROUND()` function; for details, see the description for [`ROUND()`](https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_round).

> todo