package oracle

//func TestDoubleCheck(t *testing.T) {
//	conn, err := testsqls.GetConnector()
//	if err != nil {
//		t.Fatal(err.Error())
//	}
//	sql1 := "(SELECT (DATE_ADD(~`f5`, INTERVAL 1 DAY)) AS `f1`,(COERCIBILITY(`f5`)) AS `f2`,(`f4`) AS `f3` FROM (SELECT `col_decimal(40, 20)_undef_unsigned` AS `f4`,`col_decimal(40, 20)_key_signed` AS `f5`,`col_float_undef_signed` AS `f6` FROM `table_3_utf8_2` USE INDEX (`col_double_key_unsigned`)) AS `t1` WHERE ((((`f5`) BETWEEN `f6` AND FLOOR(0.3360101420486982)) IS FALSE) AND ((`f4`)>=ANY (SELECT `col_char(20)_key_signed` FROM `table_5_utf8_2` IGNORE INDEX (`col_decimal(40, 20)_key_signed`, `col_decimal(40, 20)_key_unsigned`)))) OR (((1) NOT BETWEEN `f5` AND ACOS(5)) AND (NOT (CAST((`f6`) AS CHAR) REGEXP _UTF8MB4'^[0-9]'))) HAVING ((NOT ((`f2`) IN (COLLATION(`f1`),COLLATION(`f1`),_UTF8MB4'2010-07-22'))) OR ((CAST((COERCIBILITY(_UTF8MB4'04:35:34')) AS CHAR) NOT REGEXP _UTF8MB4'[a-z]+[0-9]+') IS TRUE) OR (CAST((LOG(0.9167666653516418)) AS CHAR) NOT REGEXP _UTF8MB4'[0-9]+[a-z]+')) IS TRUE ORDER BY `f6`) UNION ALL (SELECT (ACOS(0.9294061433657331)) AS `f1`,(FORMAT_BYTES(`f9`)) AS `f2`,(-`f7`+LN(3)) AS `f3` FROM (SELECT (`f12`<<CONCAT(`f12`, 1, `f12`)) AS `f10`,(~9) AS `f8`,(COT(8)&DATEDIFF(_UTF8MB4'2003-02-22 02:11:56', _UTF8MB4'2016-06-09')) AS `f11` FROM (SELECT `col_double_undef_signed` AS `f12`,`col_char(20)_key_signed` AS `f13`,`col_bigint_key_signed` AS `f14` FROM `table_3_utf8_2`) AS `t2` WHERE (((CAST((QUARTER(_UTF8MB4'2003-11-01')) AS CHAR) LIKE _UTF8MB4'%1%') IS FALSE) AND (NOT ((`f12`-`f14`)<(_UTF8MB4'i'))) AND (((-`f12`)!=(FORMAT_BYTES(`f14`)*`f12`)) IS TRUE)) IS FALSE ORDER BY `f13`) AS `t3` NATURAL JOIN (SELECT (SUBSTRING(_UTF8MB4'you''re', 6)) AS `f7`,(-NULL) AS `f15`,(FORMAT_BYTES(`f18`)) AS `f9` FROM (SELECT `col_varchar(20)_undef_signed` AS `f16`,`col_bigint_key_signed` AS `f17`,`col_bigint_undef_unsigned` AS `f18` FROM `table_5_utf8_2` USE INDEX (`col_varchar(20)_key_signed`, `col_bigint_key_signed`)) AS `t4` WHERE ((NOT ((`f16`) NOT BETWEEN _UTF8MB4'2013-06-27 15:27:55' AND 0)) OR ((CAST((DEGREES(-1646855295281958936)) AS CHAR) LIKE _UTF8MB4'%0%') IS FALSE) OR (NOT (CAST((ASIN(4)) AS CHAR) REGEXP _UTF8MB4'[0-9]+[a-z]+'))) IS TRUE) AS `t5`)"
//	sql2 := "(SELECT (DATE_ADD(~`f5`, INTERVAL 1 DAY)) AS `f1`,(COERCIBILITY(`f5`)) AS `f2`,(`f4`) AS `f3` FROM (SELECT `col_decimal(40, 20)_undef_unsigned` AS `f4`,`col_decimal(40, 20)_key_signed` AS `f5`,`col_float_undef_signed` AS `f6` FROM `table_3_utf8_2` USE INDEX (`col_double_key_unsigned`)) AS `t1` WHERE ((((`f5`) BETWEEN `f6` AND FLOOR(0.3360101420486982)) IS FALSE) AND ((`f4`)>=ALL (SELECT `col_char(20)_key_signed` FROM `table_5_utf8_2` IGNORE INDEX (`col_decimal(40, 20)_key_signed`, `col_decimal(40, 20)_key_unsigned`)))) OR (((1) NOT BETWEEN `f5` AND ACOS(5)) AND (NOT (CAST((`f6`) AS CHAR) REGEXP _UTF8MB4'^[0-9]'))) HAVING ((NOT ((`f2`) IN (COLLATION(`f1`),COLLATION(`f1`),_UTF8MB4'2010-07-22'))) OR ((CAST((COERCIBILITY(_UTF8MB4'04:35:34')) AS CHAR) NOT REGEXP _UTF8MB4'[a-z]+[0-9]+') IS TRUE) OR (CAST((LOG(0.9167666653516418)) AS CHAR) NOT REGEXP _UTF8MB4'[0-9]+[a-z]+')) IS TRUE ORDER BY `f6`) UNION ALL (SELECT (ACOS(0.9294061433657331)) AS `f1`,(FORMAT_BYTES(`f9`)) AS `f2`,(-`f7`+LN(3)) AS `f3` FROM (SELECT (`f12`<<CONCAT(`f12`, 1, `f12`)) AS `f10`,(~9) AS `f8`,(COT(8)&DATEDIFF(_UTF8MB4'2003-02-22 02:11:56', _UTF8MB4'2016-06-09')) AS `f11` FROM (SELECT `col_double_undef_signed` AS `f12`,`col_char(20)_key_signed` AS `f13`,`col_bigint_key_signed` AS `f14` FROM `table_3_utf8_2`) AS `t2` WHERE (((CAST((QUARTER(_UTF8MB4'2003-11-01')) AS CHAR) LIKE _UTF8MB4'%1%') IS FALSE) AND (NOT ((`f12`-`f14`)<(_UTF8MB4'i'))) AND (((-`f12`)!=(FORMAT_BYTES(`f14`)*`f12`)) IS TRUE)) IS FALSE ORDER BY `f13`) AS `t3` NATURAL JOIN (SELECT (SUBSTRING(_UTF8MB4'you''re', 6)) AS `f7`,(-NULL) AS `f15`,(FORMAT_BYTES(`f18`)) AS `f9` FROM (SELECT `col_varchar(20)_undef_signed` AS `f16`,`col_bigint_key_signed` AS `f17`,`col_bigint_undef_unsigned` AS `f18` FROM `table_5_utf8_2` USE INDEX (`col_varchar(20)_key_signed`, `col_bigint_key_signed`)) AS `t4` WHERE ((NOT ((`f16`) NOT BETWEEN _UTF8MB4'2013-06-27 15:27:55' AND 0)) OR ((CAST((DEGREES(-1646855295281958936)) AS CHAR) LIKE _UTF8MB4'%0%') IS FALSE) OR (NOT (CAST((ASIN(4)) AS CHAR) REGEXP _UTF8MB4'[0-9]+[a-z]+'))) IS TRUE) AS `t5`)"
//	result1 := conn.ExecSQLS(sql1)
//	result2 := conn.ExecSQLS(sql2)
//	if result1.Err != nil {
//		t.Log("err1", result1.Err)
//	} else {
//		t.Log("no err1")
//	}
//	if result2.Err != nil {
//		t.Log("err2", result2.Err)
//	} else {
//		t.Log("no err2")
//	}
//	_, errBuf, err := conn.ExecSQLX(sql1, 60000)
//	if err != nil {
//		t.Fatal(string(errBuf))
//	}
//	_, _, err = conn.ExecSQLX(sql1, 60000)
//	if err != nil {
//		t.Fatal(err.Error())
//	}
//	if !DoubleCheck(conn, sql1, sql2, result1.Err != nil, result2.Err != nil) {
//		t.Fatal("!double check")
//	}
//	t.Log("double check ok")
//}