create table table_3_utf8_2 (
`pk` int primary key,
`col_bigint_undef_signed` bigint  ,
`col_bigint_undef_unsigned` bigint unsigned ,
`col_bigint_key_signed` bigint  ,
`col_bigint_key_unsigned` bigint unsigned ,
`col_float_undef_signed` float  ,
`col_float_undef_unsigned` float unsigned ,
`col_float_key_signed` float  ,
`col_float_key_unsigned` float unsigned ,
`col_double_undef_signed` double  ,
`col_double_undef_unsigned` double unsigned ,
`col_double_key_signed` double  ,
`col_double_key_unsigned` double unsigned ,
`col_decimal(40, 20)_undef_signed` decimal(40, 20)  ,
`col_decimal(40, 20)_undef_unsigned` decimal(40, 20) unsigned ,
`col_decimal(40, 20)_key_signed` decimal(40, 20)  ,
`col_decimal(40, 20)_key_unsigned` decimal(40, 20) unsigned ,
`col_char(20)_undef_signed` char(20)  ,
`col_char(20)_key_signed` char(20)  ,
`col_varchar(20)_undef_signed` varchar(20)  ,
`col_varchar(20)_key_signed` varchar(20)  ,
key (`col_bigint_key_signed`),
key (`col_bigint_key_unsigned`),
key (`col_float_key_signed`),
key (`col_float_key_unsigned`),
key (`col_double_key_signed`),
key (`col_double_key_unsigned`),
key (`col_decimal(40, 20)_key_signed`),
key (`col_decimal(40, 20)_key_unsigned`),
key (`col_char(20)_key_signed`),
key (`col_varchar(20)_key_signed`)
) character set utf8 
partition by hash(pk)
partitions 2;
insert into table_3_utf8_2 values (0,82.1847,1,39.0425,38.1089,-1,1,94.1106,1.009,12.991,19755,-13064,0,1,79.1429,-2,1,"well",'3
','-0','e'),(1,1,20.0078,-9.183,68.1957,1,2,1,0.0001,12.991,2,71.0510,1,-1,2,12.991,12.991,'3	','1','3	','-0'),(2,-2,1,-21247,1.009,2,1.009,0.0001,36.0002,-2,2,-0,0.0001,-2,0.1598,47.1515,1.009,'3	','w','-1','e');
create table table_7_utf8_2 (
`pk` int primary key,
`col_bigint_undef_signed` bigint  ,
`col_bigint_undef_unsigned` bigint unsigned ,
`col_bigint_key_signed` bigint  ,
`col_bigint_key_unsigned` bigint unsigned ,
`col_float_undef_signed` float  ,
`col_float_undef_unsigned` float unsigned ,
`col_float_key_signed` float  ,
`col_float_key_unsigned` float unsigned ,
`col_double_undef_signed` double  ,
`col_double_undef_unsigned` double unsigned ,
`col_double_key_signed` double  ,
`col_double_key_unsigned` double unsigned ,
`col_decimal(40, 20)_undef_signed` decimal(40, 20)  ,
`col_decimal(40, 20)_undef_unsigned` decimal(40, 20) unsigned ,
`col_decimal(40, 20)_key_signed` decimal(40, 20)  ,
`col_decimal(40, 20)_key_unsigned` decimal(40, 20) unsigned ,
`col_char(20)_undef_signed` char(20)  ,
`col_char(20)_key_signed` char(20)  ,
`col_varchar(20)_undef_signed` varchar(20)  ,
`col_varchar(20)_key_signed` varchar(20)  ,
key (`col_bigint_key_signed`),
key (`col_bigint_key_unsigned`),
key (`col_float_key_signed`),
key (`col_float_key_unsigned`),
key (`col_double_key_signed`),
key (`col_double_key_unsigned`),
key (`col_decimal(40, 20)_key_signed`),
key (`col_decimal(40, 20)_key_unsigned`),
key (`col_char(20)_key_signed`),
key (`col_varchar(20)_key_signed`)
) character set utf8 
partition by hash(pk)
partitions 2;
insert into table_7_utf8_2 values (0,-9.183,1,1.1384,2,15.1271,12.991,-2,0.0001,36.1270,79.1819,0.0001,0.0001,3.1387,52.0818,-0,0.0001,'1','3	','0','0'),(1,79,12.991,107,2,-0.0001,0,1.009,1.009,34,1,-1,69.0208,1,2,120,12.991,'3	','-1',"if",'b'),(2,-2,1,-9.183,1,12.991,0.0001,53,12.991,1.009,12.991,12.991,0.0001,-0.0001,12.991,0.0001,2,'3
','p','0','3	'),(3,-0.0001,12.991,1.009,1.009,-9.183,2,0,1,-2,1,2,1,2,1.009,2,12.991,'3
','0','k','0'),(4,1.009,0.0001,-1,12.991,2,47,2,0,12.991,12.991,1.009,0,1.009,1.009,-0.0001,6949,'-1','	3','1','m'),(5,-0,1,0,0,0.0001,28.1237,12.991,0,12.991,12.991,-0,12.991,2,2,2,1.009,'0','	3','0','	3'),(6,45.0855,1,38.1166,1,1.009,80.0284,2,122,0.0001,0,-1,11130,0,1,1,0,"know",'-0','
3','3
');