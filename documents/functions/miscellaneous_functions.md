# MySQL Document

https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html

# Overview

```shell
	[1]: only tidb features
	[2]: supported features of impo
	[3]: unsupported features of impo
	
	[3]AnyValue        = "any_value"
	[1]DefaultFunc     = "default_func"
	[2]InetAton        = "inet_aton"
	[2]InetNtoa        = "inet_ntoa"
	[2]Inet6Aton       = "inet6_aton"
	[2]Inet6Ntoa       = "inet6_ntoa"
	[1]IsFreeLock      = "is_free_lock"
	[2]IsIPv4          = "is_ipv4"
	[2]IsIPv4Compat    = "is_ipv4_compat"
	[2]IsIPv4Mapped    = "is_ipv4_mapped"
	[2]IsIPv6          = "is_ipv6"
	[1]IsUsedLock      = "is_used_lock"
	[2]IsUUID          = "is_uuid"
	[3]MasterPosWait   = "master_pos_wait"
	[2]NameConst       = "name_const"
	[1]ReleaseAllLocks = "release_all_locks"
	[3]Sleep           = "sleep"
	[3]UUID            = "uuid"
	[3]UUIDShort       = "uuid_short"
	[3]UUIDToBin       = "uuid_to_bin"
	[3]BinToUUID       = "bin_to_uuid"
	[1]VitessHash      = "vitess_hash"
	[1]GetLock     = "get_lock"
	[1]ReleaseLock = "release_lock"
```

# any_value

This function is useful for `GROUP BY` queries when the [`ONLY_FULL_GROUP_BY`](https://dev.mysql.com/doc/refman/8.0/en/sql-mode.html#sqlmode_only_full_group_by) SQL mode is enabled, for cases when MySQL rejects a query that you know is valid for reasons that MySQL cannot determine. The function return value and type are the same as the return value and type of its argument, but the function result is not checked for the [`ONLY_FULL_GROUP_BY`](https://dev.mysql.com/doc/refman/8.0/en/sql-mode.html#sqlmode_only_full_group_by) SQL mode.

> cannot support GROUP BY

# default_func

> not found

# inet_aton

Given the dotted-quad representation of an IPv4 network address as a string, returns an integer that represents the numeric value of the address in network byte order (big endian). [`INET_ATON()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_inet-aton) returns `NULL` if it does not understand its argument, or if *`expr`* is `NULL`.

**example**

```sql
mysql> select inet_aton('127.0.0.1');
+------------------------+
| inet_aton('127.0.0.1') |
+------------------------+
|             2130706433 |
+------------------------+
1 row in set (0.00 sec)

mysql> select inet_aton('127');
+------------------+
| inet_aton('127') |
+------------------+
|              127 |
+------------------+
1 row in set (0.00 sec)

mysql> select inet_aton(127);
+----------------+
| inet_aton(127) |
+----------------+
|            127 |
+----------------+
1 row in set (0.00 sec)

mysql> select inet_aton('a.a.a.a');
+----------------------+
| inet_aton('a.a.a.a') |
+----------------------+
|                 NULL |
+----------------------+
1 row in set, 1 warning (0.00 sec)

mysql> select inet_aton(BINARY 'hello world');
+---------------------------------+
| inet_aton(BINARY 'hello world') |
+---------------------------------+
|                            NULL |
+---------------------------------+
1 row in set, 2 warnings (0.00 sec)

```

# inet_ntoa

Given a numeric IPv4 network address in network byte order, returns the dotted-quad string representation of the address as a string in the connection character set. [`INET_NTOA()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_inet-ntoa) returns `NULL` if it does not understand its argument.

**example**

```sql
mysql> SELECT INET_NTOA(167773449);
+----------------------+
| INET_NTOA(167773449) |
+----------------------+
| 10.0.5.9             |
+----------------------+
1 row in set (0.00 sec)

mysql> SELECT INET_NTOA(1);
+--------------+
| INET_NTOA(1) |
+--------------+
| 0.0.0.1      |
+--------------+
1 row in set (0.00 sec)

mysql> SELECT INET_NTOA(0);
+--------------+
| INET_NTOA(0) |
+--------------+
| 0.0.0.0      |
+--------------+
1 row in set (0.00 sec)

mysql> SELECT INET_NTOA(BINARY 'hello world');
+---------------------------------+
| INET_NTOA(BINARY 'hello world') |
+---------------------------------+
| 0.0.0.0                         |
+---------------------------------+
1 row in set, 2 warnings (0.00 sec)

mysql> SELECT INET_NTOA('a.1.a.1');
+----------------------+
| INET_NTOA('a.1.a.1') |
+----------------------+
| 0.0.0.0              |
+----------------------+
1 row in set, 1 warning (0.00 sec)

mysql> SELECT INET_NTOA(NULL);
+-----------------+
| INET_NTOA(NULL) |
+-----------------+
| NULL            |
+-----------------+
1 row in set (0.00 sec)
```

# inet6_aton 

Given an IPv6 or IPv4 network address as a string, returns a binary string that represents the numeric value of the address in network byte order (big endian). Because numeric-format IPv6 addresses require more bytes than the largest integer type, the representation returned by this function has the [`VARBINARY`](https://dev.mysql.com/doc/refman/8.0/en/binary-varbinary.html) data type: [`VARBINARY(16)`](https://dev.mysql.com/doc/refman/8.0/en/binary-varbinary.html) for IPv6 addresses and [`VARBINARY(4)`](https://dev.mysql.com/doc/refman/8.0/en/binary-varbinary.html) for IPv4 addresses. If the argument is not a valid address, or if it is `NULL`, [`INET6_ATON()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_inet6-aton) returns `NULL`.

**example**

```mysql
mysql> select INET6_ATON('fdfe::5a55:caff:fefa:9089');
+----------------------------------------------------------------------------------+
| INET6_ATON('fdfe::5a55:caff:fefa:9089')                                          |
+----------------------------------------------------------------------------------+
| 0xFDFE0000000000005A55CAFFFEFA9089                                               |
+----------------------------------------------------------------------------------+
1 row in set (0.00 sec)
```

# inet6_ntoa

Given an IPv6 or IPv4 network address represented in numeric form as a binary string, returns the string representation of the address as a string in the connection character set. If the argument is not a valid address, or if it is `NULL`, [`INET6_NTOA()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_inet6-ntoa) returns `NULL`.

**example**

```sql
mysql> select inet6_ntoa(UNHEX('FDFE0000000000005A55CAFFFEFA9089'));
+-------------------------------------------------------+
| inet6_ntoa(UNHEX('FDFE0000000000005A55CAFFFEFA9089')) |
+-------------------------------------------------------+
| fdfe::5a55:caff:fefa:9089                             |
+-------------------------------------------------------+
1 row in set (0.00 sec)
```

# is_free_lock

> not found

# is_ipv4

Returns 1 if the argument is a valid IPv4 address specified as a string, 0 otherwise. Returns `NULL` if *`expr`* is `NULL`.

**example**

```sql
mysql> SELECT IS_IPV4('10.0.5.9'), IS_IPV4('10.0.5.256');
+---------------------+-----------------------+
| IS_IPV4('10.0.5.9') | IS_IPV4('10.0.5.256') |
+---------------------+-----------------------+
|                   1 |                     0 |
+---------------------+-----------------------+
1 row in set (0.00 sec)

mysql> SELECT IS_IPV4('10.0.5');
+-------------------+
| IS_IPV4('10.0.5') |
+-------------------+
|                 0 |
+-------------------+
1 row in set (0.00 sec)

mysql> SELECT IS_IPV4('1');
+--------------+
| IS_IPV4('1') |
+--------------+
|            0 |
+--------------+
1 row in set (0.00 sec)

mysql> SELECT IS_IPV4(NULL);
+---------------+
| IS_IPV4(NULL) |
+---------------+
|          NULL |
+---------------+
1 row in set (0.00 sec)
```

# is_ipv4_compat

This function takes an IPv6 address represented in numeric form as a binary string, as returned by [`INET6_ATON()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_inet6-aton). It returns 1 if the argument is a valid IPv4-compatible IPv6 address, 0 otherwise (unless *`expr`* is `NULL`, in which case the function returns `NULL`). IPv4-compatible addresses have the form `::*`ipv4_address`*`.	

**example**

```sql
mysql> SELECT IS_IPV4_COMPAT(INET6_ATON('::10.0.5.9'));
+------------------------------------------+
| IS_IPV4_COMPAT(INET6_ATON('::10.0.5.9')) |
+------------------------------------------+
|                                        1 |
+------------------------------------------+
1 row in set (0.00 sec)

mysql> SELECT IS_IPV4_COMPAT(INET6_ATON('::ffff:10.0.5.9'));
+-----------------------------------------------+
| IS_IPV4_COMPAT(INET6_ATON('::ffff:10.0.5.9')) |
+-----------------------------------------------+
|                                             0 |
+-----------------------------------------------+
1 row in set (0.00 sec)

mysql> SELECT IS_IPV4_COMPAT(123);
+---------------------+
| IS_IPV4_COMPAT(123) |
+---------------------+
|                   0 |
+---------------------+
1 row in set (0.00 sec)

mysql> SELECT IS_IPV4_COMPAT(NULL);
+----------------------+
| IS_IPV4_COMPAT(NULL) |
+----------------------+
|                 NULL |
+----------------------+
1 row in set (0.00 sec)
```

# is_ipv4_mapped

This function takes an IPv6 address represented in numeric form as a binary string, as returned by [`INET6_ATON()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_inet6-aton). It returns 1 if the argument is a valid IPv4-mapped IPv6 address, 0 otherwise, unless *`expr`* is `NULL`, in which case the function returns `NULL`. IPv4-mapped addresses have the form `::ffff:*`ipv4_address`*`.

**example**

```sql
mysql> SELECT IS_IPV4_MAPPED(INET6_ATON('::10.0.5.9'));
+------------------------------------------+
| IS_IPV4_MAPPED(INET6_ATON('::10.0.5.9')) |
+------------------------------------------+
|                                        0 |
+------------------------------------------+
1 row in set (0.00 sec)

mysql> SELECT IS_IPV4_MAPPED(INET6_ATON('::ffff:10.0.5.9'));
+-----------------------------------------------+
| IS_IPV4_MAPPED(INET6_ATON('::ffff:10.0.5.9')) |
+-----------------------------------------------+
|                                             1 |
+-----------------------------------------------+
1 row in set (0.00 sec)
```

# is_ipv6

Returns 1 if the argument is a valid IPv6 address specified as a string, 0 otherwise, unless *`expr`* is `NULL`, in which case the function returns `NULL`. This function does not consider IPv4 addresses to be valid IPv6 addresses.

**example**

```sql
SELECT IS_IPV6('10.0.5.9'), IS_IPV6('::1');
+---------------------+----------------+
| IS_IPV6('10.0.5.9') | IS_IPV6('::1') |
+---------------------+----------------+
|                   0 |              1 |
+---------------------+----------------+
1 row in set (0.00 sec)
```

# is_used_lock

> not found

# is_uuid

Returns 1 if the argument is a valid string-format UUID, 0 if the argument is not a valid UUID, and `NULL` if the argument is `NULL`.

“Valid” means that the value is in a format that can be parsed. That is, it has the correct length and contains only the permitted characters (hexadecimal digits in any lettercase and, optionally, dashes and curly braces). This format is most common:

```none
aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee
```

These other formats are also permitted:

```none
 aaaaaaaabbbbccccddddeeeeeeeeeeee
{aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee}
```

**example**

```sql
mysql> SELECT IS_UUID('6ccd780c-baba-1026-9564-5b8c656024db');
+-------------------------------------------------+
| IS_UUID('6ccd780c-baba-1026-9564-5b8c656024db') |
+-------------------------------------------------+
|                                               1 |
+-------------------------------------------------+
1 row in set (0.00 sec)

mysql> SELECT IS_UUID('6ccd780c-baba-1026-9564-5b8c6560');
+---------------------------------------------+
| IS_UUID('6ccd780c-baba-1026-9564-5b8c6560') |
+---------------------------------------------+
|                                           0 |
+---------------------------------------------+
1 row in set (0.01 sec)

mysql> SELECT IS_UUID('1');
+--------------+
| IS_UUID('1') |
+--------------+
|            0 |
+--------------+
1 row in set (0.00 sec)
```

# master_pos_wait

This function is for control of source/replica synchronization. It blocks until the replica has read and applied all updates up to the specified position in the source's binary log. From MySQL 8.0.26, [`MASTER_POS_WAIT()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_master-pos-wait) is deprecated and the alias [`SOURCE_POS_WAIT()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_source-pos-wait) should be used instead. In releases before MySQL 8.0.26, use [`MASTER_POS_WAIT()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_master-pos-wait).

The return value is the number of log events the replica had to wait for to advance to the specified position. The function returns `NULL` if the replication SQL thread is not started, the replica's source information is not initialized, the arguments are incorrect, or an error occurs. It returns `-1` if the timeout has been exceeded. If the replication SQL thread stops while [`MASTER_POS_WAIT()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_master-pos-wait) is waiting, the function returns `NULL`. If the replica is past the specified position, the function returns immediately.

If the binary log file position has been marked as invalid, the function waits until a valid file position is known. The binary log file position can be marked as invalid when the [`CHANGE REPLICATION SOURCE TO`](https://dev.mysql.com/doc/refman/8.0/en/change-replication-source-to.html) option `GTID_ONLY` is set for the replication channel, and the server is restarted or replication is stopped. The file position becomes valid after a transaction is successfully applied past the given file position. If the applier does not reach the stated position, the function waits until the timeout. Use a [`SHOW REPLICA STATUS`](https://dev.mysql.com/doc/refman/8.0/en/show-replica-status.html) statement to check if the binary log file position has been marked as invalid.

On a multithreaded replica, the function waits until expiry of the limit set by the [`replica_checkpoint_group`](https://dev.mysql.com/doc/refman/8.0/en/replication-options-replica.html#sysvar_replica_checkpoint_group), [`slave_checkpoint_group`](https://dev.mysql.com/doc/refman/8.0/en/replication-options-replica.html#sysvar_slave_checkpoint_group), [`replica_checkpoint_period`](https://dev.mysql.com/doc/refman/8.0/en/replication-options-replica.html#sysvar_replica_checkpoint_period) or [`slave_checkpoint_period`](https://dev.mysql.com/doc/refman/8.0/en/replication-options-replica.html#sysvar_slave_checkpoint_period) system variable, when the checkpoint operation is called to update the status of the replica. Depending on the setting for the system variables, the function might therefore return some time after the specified position was reached.

If binary log transaction compression is in use and the transaction payload at the specified position is compressed (as a `Transaction_payload_event`), the function waits until the whole transaction has been read and applied, and the positions have updated.

If a *`timeout`* value is specified, [`MASTER_POS_WAIT()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_master-pos-wait) stops waiting when *`timeout`* seconds have elapsed. *`timeout`* must be greater than 0; a zero or negative *`timeout`* means no timeout.

The optional *`channel`* value enables you to name which replication channel the function applies to. See [Section 17.2.2, “Replication Channels”](https://dev.mysql.com/doc/refman/8.0/en/replication-channels.html) for more information.

This function is unsafe for statement-based replication. A warning is logged if you use this function when [`binlog_format`](https://dev.mysql.com/doc/refman/8.0/en/replication-options-binary-log.html#sysvar_binlog_format) is set to `STATEMENT`.

> cannot support

# name_const

Returns the given value. When used to produce a result set column, [`NAME_CONST()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_name-const) causes the column to have the given name. The arguments should be constants.

<font color="red">Note that the arguments should be constants!</font>

**example**

```sql
mysql> SELECT NAME_CONST('myname', 14);
+--------+
| myname |
+--------+
|     14 |
+--------+
1 row in set (0.02 sec)

mysql> SELECT NAME_CONST('myname', C1) FROM T;
ERROR 1210 (HY000): Incorrect arguments to NAME_CONST
```

# release_all_locks

> not found

# sleep

Sleeps (pauses) for the number of seconds given by the *`duration`* argument, then returns 0. The duration may have a fractional part. If the argument is `NULL` or negative, [`SLEEP()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_sleep) produces a warning, or an error in strict SQL mode.

**example**

```sql
mysql> select sleep(3);
+----------+
| sleep(3) |
+----------+
|        0 |
+----------+
1 row in set (3.00 sec)
```

> may cause long waits, do not support.

# uuid

Returns a Universal Unique Identifier (UUID) generated according to RFC 4122, “A Universally Unique IDentifier (UUID) URN Namespace” (http://www.ietf.org/rfc/rfc4122.txt).

A UUID is designed as a number that is globally unique in space and time. Two calls to [`UUID()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_uuid) are expected to generate two different values, even if these calls are performed on two separate devices not connected to each other.

Warning

Although [`UUID()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_uuid) values are intended to be unique, they are not necessarily unguessable or unpredictable. If unpredictability is required, UUID values should be generated some other way.

[`UUID()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_uuid) returns a value that conforms to UUID version 1 as described in RFC 4122. The value is a 128-bit number represented as a `utf8mb3` string of five hexadecimal numbers in `aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee` format:

- The first three numbers are generated from the low, middle, and high parts of a timestamp. The high part also includes the UUID version number.

- The fourth number preserves temporal uniqueness in case the timestamp value loses monotonicity (for example, due to daylight saving time).

- The fifth number is an IEEE 802 node number that provides spatial uniqueness. A random number is substituted if the latter is not available (for example, because the host device has no Ethernet card, or it is unknown how to find the hardware address of an interface on the host operating system). In this case, spatial uniqueness cannot be guaranteed. Nevertheless, a collision should have *very* low probability.

  The MAC address of an interface is taken into account only on FreeBSD, Linux, and Windows. On other operating systems, MySQL uses a randomly generated 48-bit number.

**example**

```sql
mysql> SELECT UUID();
+--------------------------------------+
| UUID()                               |
+--------------------------------------+
| bbb3898b-5113-11ed-92a3-0242ac110002 |
+--------------------------------------+
1 row in set (0.00 sec)
```

> uncertain, do not support

# uuid_short

Returns a “short” universal identifier as a 64-bit unsigned integer. Values returned by [`UUID_SHORT()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_uuid-short) differ from the string-format 128-bit identifiers returned by the [`UUID()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_uuid) function and have different uniqueness properties. The value of [`UUID_SHORT()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_uuid-short) is guaranteed to be unique if the following conditions hold:

- The [`server_id`](https://dev.mysql.com/doc/refman/8.0/en/replication-options.html#sysvar_server_id) value of the current server is between 0 and 255 and is unique among your set of source and replica servers
- You do not set back the system time for your server host between [**mysqld**](https://dev.mysql.com/doc/refman/8.0/en/mysqld.html) restarts
- You invoke [`UUID_SHORT()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_uuid-short) on average fewer than 16 million times per second between [**mysqld**](https://dev.mysql.com/doc/refman/8.0/en/mysqld.html) restarts

The [`UUID_SHORT()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_uuid-short) return value is constructed this way:

```clike
  (server_id & 255) << 56
+ (server_startup_time_in_seconds << 24)
+ incremented_variable++;
```

**example**

```sql
mysql> SELECT UUID_SHORT();
+--------------------+
| UUID_SHORT()       |
+--------------------+
| 100013969170759680 |
+--------------------+
1 row in set (0.00 sec)
```

>  uncertain, do not support

# uuid_to_bin

Converts a string UUID to a binary UUID and returns the result. (The [`IS_UUID()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_is-uuid) function description lists the permitted string UUID formats.) The return binary UUID is a [`VARBINARY(16)`](https://dev.mysql.com/doc/refman/8.0/en/binary-varbinary.html) value. If the UUID argument is `NULL`, the return value is `NULL`. If any argument is invalid, an error occurs.

[`UUID_TO_BIN()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_uuid-to-bin) takes one or two arguments:

- The one-argument form takes a string UUID value. The binary result is in the same order as the string argument.
- The two-argument form takes a string UUID value and a flag value:
  - If *`swap_flag`* is 0, the two-argument form is equivalent to the one-argument form. The binary result is in the same order as the string argument.
  - If *`swap_flag`* is 1, the format of the return value differs: The time-low and time-high parts (the first and third groups of hexadecimal digits, respectively) are swapped. This moves the more rapidly varying part to the right and can improve indexing efficiency if the result is stored in an indexed column.

Time-part swapping assumes the use of UUID version 1 values, such as are generated by the [`UUID()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_uuid) function. For UUID values produced by other means that do not follow version 1 format, time-part swapping provides no benefit. For details about version 1 format, see the [`UUID()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_uuid) function description.

**example**

```sql
mysql> SELECT UUID_TO_BIN('6ccd780c-baba-1026-9564-5b8c656024db', 1);
+----------------------------------------------------------------------------------------------------------------+
| UUID_TO_BIN('6ccd780c-baba-1026-9564-5b8c656024db', 1)                                                         |
+----------------------------------------------------------------------------------------------------------------+
| 0x1026BABA6CCD780C95645B8C656024DB                                                                             |
+----------------------------------------------------------------------------------------------------------------+
1 row in set (0.00 sec)


mysql> SELECT UUID_TO_BIN('6ccd780c-baba-1026-9564-5b8c656024db');
+----------------------------------------------------------------------------------------------------------+
| UUID_TO_BIN('6ccd780c-baba-1026-9564-5b8c656024db')                                                      |
+----------------------------------------------------------------------------------------------------------+
| 0x6CCD780CBABA102695645B8C656024DB                                                                       |
+----------------------------------------------------------------------------------------------------------+
1 row in set (0.00 sec)

mysql> SELECT UUID_TO_BIN('6ccd780c-baba-1026-9564-');
ERROR 1411 (HY000): Incorrect string value: '6ccd780c-baba-1026-9564-' for function uuid_to_bin
```

> need to provide a uuid, otherwise an error will occur, not supported for convenience

# bin_to_uuid

[`BIN_TO_UUID()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_bin-to-uuid) is the inverse of [`UUID_TO_BIN()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_uuid-to-bin). It converts a binary UUID to a string UUID and returns the result. The binary value should be a UUID as a [`VARBINARY(16)`](https://dev.mysql.com/doc/refman/8.0/en/binary-varbinary.html) value. The return value is a string of five hexadecimal numbers separated by dashes. (For details about this format, see the [`UUID()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_uuid) function description.) If the UUID argument is `NULL`, the return value is `NULL`. If any argument is invalid, an error occurs.

[`BIN_TO_UUID()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_bin-to-uuid) takes one or two arguments:

- The one-argument form takes a binary UUID value. The UUID value is assumed not to have its time-low and time-high parts swapped. The string result is in the same order as the binary argument.
- The two-argument form takes a binary UUID value and a swap-flag value:
  - If *`swap_flag`* is 0, the two-argument form is equivalent to the one-argument form. The string result is in the same order as the binary argument.
  - If *`swap_flag`* is 1, the UUID value is assumed to have its time-low and time-high parts swapped. These parts are swapped back to their original position in the result value.

For usage examples and information about time-part swapping, see the [`UUID_TO_BIN()`](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_uuid-to-bin) function description.

**example**

```sql
mysql> SELECT bin_to_uuid(0x6CCD780CBABA102695645B8C656024DB);
+-------------------------------------------------+
| bin_to_uuid(0x6CCD780CBABA102695645B8C656024DB) |
+-------------------------------------------------+
| 6ccd780c-baba-1026-9564-5b8c656024db            |
+-------------------------------------------------+
1 row in set (0.00 sec)

mysql> SELECT bin_to_uuid(1);
ERROR 1411 (HY000): Incorrect string value: '1' for function bin_to_uuid
```

> need to provide a binary uuid, otherwise an error will occur, not supported for convenience

# vitess_hash

> not found

# get_lock

> not found

# release_lock

> not found