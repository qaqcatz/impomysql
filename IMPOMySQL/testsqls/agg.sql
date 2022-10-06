create table table_10_utf8_4 (
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
partitions 4;
insert into table_10_utf8_4 values (0,82.1847,1,39.0425,38.1089,-1,1,94.1106,1.009,12.991,19755,-13064,0,1,79.1429,-2,1,"well",'3
','-0','e'),(1,1,20.0078,-9.183,68.1957,1,2,1,0.0001,12.991,2,71.0510,1,-1,2,12.991,12.991,'3	','1','3	','-0'),(2,-2,1,-21247,1.009,2,1.009,0.0001,36.0002,-2,2,-0,0.0001,-2,0.1598,47.1515,1.009,'3	','w','-1','e'),(3,-9.183,1,1.1384,2,15.1271,12.991,-2,0.0001,36.1270,79.1819,0.0001,0.0001,3.1387,52.0818,-0,0.0001,'1','3	','0','0'),(4,79,12.991,107,2,-0.0001,0,1.009,1.009,34,1,-1,69.0208,1,2,120,12.991,'3	','-1',"if",'b'),(5,-2,1,-9.183,1,12.991,0.0001,53,12.991,1.009,12.991,12.991,0.0001,-0.0001,12.991,0.0001,2,'3
','p','0','3	'),(6,-0.0001,12.991,1.009,1.009,-9.183,2,0,1,-2,1,2,1,2,1.009,2,12.991,'3
','0','k','0'),(7,1.009,0.0001,-1,12.991,2,47,2,0,12.991,12.991,1.009,0,1.009,1.009,-0.0001,6949,'-1','	3','1','m'),(8,-0,1,0,0,0.0001,28.1237,12.991,0,12.991,12.991,-0,12.991,2,2,2,1.009,'0','	3','0','	3'),(9,45.0855,1,38.1166,1,1.009,80.0284,2,122,0.0001,0,-1,11130,0,1,1,0,"know",'-0','
3','3
');
create table table_10_utf8_6 (
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
partitions 6;
insert into table_10_utf8_6 values (0,-0,0.0001,-9.183,35,0.0001,1.009,7,2,103,1,12.991,0.0001,-15067,1.009,15308,3.0556,'-1','a','3	','3	'),(1,31,2,-0.0001,2,99.0884,0,-1,33.1253,-2,1,-1,12.991,-0.0001,0,12.991,1.009,'0','1',"with",'
3'),(2,-83,1.009,-0,29208,12.991,0.0001,1,1.009,-3746,0,-1,77.0545,2,1,0.0001,0,'-1','	3','-0','-1'),(3,-0,98.1262,-7287,2,17796,2,-2,27,-0.0001,53.0610,-109,107,-9.183,0,2,32.1208,'
3','1','	3','1'),(4,1.009,0.0001,-2,12.991,2,0.0001,-9.183,12.991,2,1.009,-1,12.991,-8882,1,0.0001,12.991,"i",'-0','3
','3	'),(5,-0.0001,1.009,1.009,53.1114,-0.0001,12.991,-2,2,1,2,-0,1,12.991,1,-0.0001,2,'
3','1','m','	3'),(6,-0.0001,1,1.009,2,0,0.0001,2,2,12.991,14.1552,-1,0.0001,86.0712,12.991,1,0,'-0','0','1','	3'),(7,-0.0001,0,9.0894,0,57.1401,62.0262,-19151,2,-2,54,0,0,52.0740,12.991,0.0001,12.991,'1','y','
3','3
'),(8,-1,1.009,12.991,19.0121,0,0.0001,1,12.991,1,12.991,1,2,-2,2,0,0,'1','	3','-0',"or"),(9,-0.0001,23,1,12.991,12.991,10.1794,12.991,0,105,4,-9.183,0,0.0001,1.009,-0.0001,0,'3
','
3',"the",'0');
create table table_10_utf8_undef (
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
) character set utf8 ;
insert into table_10_utf8_undef values (0,0,2,-2,0.0001,-0.0001,1,0,12.991,-9.183,0,-9.183,2509,2,0,2,0,'z','3	','	3','
3'),(1,0,2,0.0001,0.0001,12.991,0.0001,1.009,67.1778,80.1179,1,2,1.009,-118,0.0001,91,59,'
3',"can't",'	3','-0'),(2,8923,51.1190,0,1,0.0001,2,-1,2,13566,12.991,68,85,-0.0001,12.991,-0.0001,79.0454,'0',"all",'
3','0'),(3,-1,14119,100.1527,0.0001,-0,2,12.991,4.1635,22.1428,0.0001,-9.183,1.009,-1,22979,-0,0.0001,'-0','1','a','	3'),(4,-28769,0.0001,-2,2,1.009,2,1,12.991,12.991,110,0.0001,77.1356,1.009,92,-17566,0.0001,"that",'1','-1','e'),(5,0.0001,1.009,127,1,-9.183,1.009,1,1,0,0,-1,0.0001,1,1.009,12.991,2,'-0','3
',"was",'-0'),(6,2,1.009,0,12.991,-9.183,1,12.991,1.009,1,47.0676,-9.183,0.0001,2,1,59.1341,2,'1','	3','3	','-0'),(7,-2,22738,0,0.0001,41.1472,1,0,12.991,-2,1,-100,42,-29083,0,-1,8832,'	3','3
','3	','
3'),(8,-2,0.0001,-1,0.0001,-9.183,2,-9.183,1,0,0,-69,0,-0,1.009,1,2,'-0','
3','n','3	'),(9,-2,0,-9.183,0,12.991,2,2,0,-0.0001,1,-1,2,0,2,31652,44.1224,"he",'0','0','
3');
create table table_10_latin1_4 (
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
) character set latin1 
partition by hash(pk)
partitions 4;
insert into table_10_latin1_4 values (0,12.991,65,0,0,-2,1,0.0001,1.009,-30415,0.0001,12.991,47.0317,-9.183,1.009,0,1,'1','0','-1','-1'),(1,-5,10914,2,0.0001,-9.183,2,2,0,1,12.991,0.0001,12.991,5908,0,75.0603,2,'0','0','0','	3'),(2,-0,0.0001,1,1,2,2,12.991,79,-0,77.0086,1,2,2,12.991,2,75,'-0','3
','0','
3'),(3,1,0,-29,93.1784,0.0001,54.0394,1.009,0,0,1,1.009,0.0001,0.0001,1,95.1781,12.991,'1','x','3	',"time"),(4,1.1275,82.0442,-9.183,56.1492,10.0362,2,-9.183,30.0062,12.991,97.0778,3.1648,1.009,1.009,88.1716,-2,0.0001,'1','1','l',"of"),(5,-9.183,1.009,12.991,12.991,-9.183,13871,9969,1,-1,44.0285,-0.0001,9606,0,2,-0,0.0001,'v','3
','3	','1'),(6,1,12.991,1,12.991,2,12.991,-10600,1,-1,15994,2,12.991,98.1734,0,1.009,1,'-0','	3','-1',"that"),(7,0.0001,0.0001,0,84,-0,2,-2,76.0071,12.991,2,-2,1.1939,-120,1,-9.183,11.1973,'	3','0','v','-0'),(8,-0.0001,12.991,0,1,-68,12.991,77.1519,25.0749,-0,2,2,0,-1,1,0,0.0001,"could",'3	','	3','-1'),(9,1,19110,60.0976,1.009,1.009,1,-2,1.009,12.991,70,-0,1,68.0971,0,0.0001,0.0001,'1','0','1','0');
create table table_10_latin1_6 (
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
) character set latin1 
partition by hash(pk)
partitions 6;
insert into table_10_latin1_6 values (0,-9.183,13.1791,-2,1.009,-9.183,2,-2,10,-2,66.0143,-0,10.0015,0.0001,25.0460,1.009,0,'-0','3
','0','	3'),(1,-9.183,51.0614,83,0.0001,12.991,0,-2,26.0490,114,2,-2,0,2,1,-1,42.1290,'1','-1','1','3	'),(2,-2,12.991,-2370,0.0001,2,1.009,1.009,1.009,-23193,46.1058,-67,59.0325,-1,1.009,-1,0,"say",'0','0','
3'),(3,8472,1.009,-111,36.0148,0.0001,1,0.0001,48,-18346,0,0,12.991,19.1518,12.991,-0,25595,'v','3
','-1','k'),(4,12.991,0,2,50,2,1,2,77,2,12.991,1.009,0.0001,12.991,1,0,13,'3	','3	',"if",'3	'),(5,1,12.991,2,0.0001,74.0862,1,12.991,0.0001,1,0,12.991,21.1617,44.1628,2,40.0911,2,'
3','-1','b','3
'),(6,12.991,12.991,1,28.0227,98,0,-22,1.009,74,12.991,1.009,5488,1,0.0001,-0,7.0142,'-0','-0','-1','3
'),(7,-1,6.0507,-0.0001,1,-0,1,-2,12.991,1.009,1.009,-9.183,80,-47,52.0884,-9.183,63.1623,'0',"hey",'z','3	'),(8,-0,1.009,0.0001,0,0,24911,2,1.009,-1,0,-0,0.0001,1,1.0216,2,1,'1','
3','1','
3'),(9,969,1,12.991,0.0001,11.0906,1.009,0,91.0202,-31735,59.1366,-9.183,1.009,-0,7868,0,1,'	3','-0','-1','t');
create table table_10_latin1_undef (
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
) character set latin1 ;
insert into table_10_latin1_undef values (0,0.0001,1,6936,0.0001,-1,10720,12.991,12.991,12.991,2,-0.0001,100.1708,8968,0.0001,-9.183,1,'3
','3	','w','r'),(1,0.0001,0.0001,-0,0.0001,1.009,12.991,-9.183,2,-0,1,74.1581,1.009,12.991,1,2,0.0001,"I'll",'3	','s','0'),(2,76,1.009,-0,1.009,91.0417,1,47.1841,1.009,-0,2,0.0001,2,12.991,26,-1,47.0638,"it",'1','
3','	3'),(3,12.991,1.009,12.991,12.991,-41,21.1905,-0,1,-1,1,-30088,0.0001,-1,1.009,0,12.991,'0','-1','3
','1'),(4,25.1231,15.0327,-0,12.991,-0,29122,-9.183,12.991,-9.183,1.009,-1,1,12.991,1,-9.183,2,'-1','0','0','
3'),(5,42,73,18.1800,0.0001,-0.0001,0.0001,-0,0.0001,34.0434,12.991,-8164,0,2,0,-0,1.009,'1','3
','3
','	3'),(6,-63,1.009,-23260,12.991,-2,1.009,0,0.0001,1,1.009,-1,2,-9.183,0,1,1,'-1','
3','k','3
'),(7,-2,14.1222,1,1.009,-0.0001,126,-2,2,0,0,1,1.009,1.009,2,-1,0.0001,'3	','0','1','3	'),(8,1.009,1.009,-1,19,1,1.009,-1,12.991,0,1,0.0001,99.0359,0,0,2,0,'3	',"can't",'3	','1'),(9,0,1,-1,0.0001,-9.183,0,1,12.991,-0,1.009,12.991,12.991,-2,0.0001,-9.183,1,'-0','
3','3	','d');
create table table_10_binary_4 (
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
) character set binary 
partition by hash(pk)
partitions 4;
insert into table_10_binary_4 values (0,-0.0001,0,-0,0.0001,-6966,0,12.991,2,0.0001,1.009,-1,99.1981,1.009,1.009,12.991,12.991,'1','3
','3
','	3'),(1,0,12.991,-9.183,83.0989,2,1,12.991,0.0001,-21,19.1985,-1,0,0.0001,1.009,12.991,1.009,'-1','3	','3	','	3'),(2,-0,1.009,1,12.991,-0,1.009,-0.0001,12.991,1,2,1,0,-9.183,0.0001,0,1.009,'3	','
3','0','	3'),(3,-1,0,0,1.009,-0,12.991,59.0759,22.0340,-0,18.1469,-0.0001,0.0001,22363,2,-1,1,'0','-1','
3','1'),(4,2,63.1528,0,59.0123,1,1,2,2,-9.183,0,-127,97.1695,-2,0.0001,0.0001,1.009,"a",'3
','0','-0'),(5,-1,1.009,-1,1.009,-0,65,1,0.0001,12.991,46.0283,-1,0,-9.183,0,-0,29028,'
3','3	',"me",'3
'),(6,-0.0001,0,-0.0001,2,12.991,1.009,12.991,19,-0,0,-2,2,21.1526,0.0001,2,1.009,'-1','	3','0','0'),(7,-0.0001,1,1,2,1.009,0.0001,19,0.0001,-9.183,0,1.009,1,0.0001,2,1,12.991,'q','3	','
3','3	'),(8,2,31507,-0,18.1224,-0.0001,12.991,122,0.0001,-0,0,2,1.009,2,2,0,1,'-1','0','
3','m'),(9,1,12.991,90.1325,12.991,-18,2,12.991,1,-9.183,1,1,0,-9.183,68,-2,0,'
3','
3','3	','0');
create table table_10_binary_6 (
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
) character set binary 
partition by hash(pk)
partitions 6;
insert into table_10_binary_6 values (0,-3699,0,70,19642,12.991,2,-0.0001,0.0001,-8275,69,-9.183,122,-0.0001,10,2,1,"yeah",'-1','0','-1'),(1,12.991,20877,-0,1,-0.0001,0.0001,0,116,12.991,0,58,51,2,12.991,-1,1,'3
','1','0',"now"),(2,-9.183,12.991,0.0001,0.0001,53.0072,2,-2,1.009,50.1433,1.009,-9.183,2,42,6066,-0.0001,12.991,'3	',"hey",'
3','3
'),(3,-2,2,1,0,22279,2,-0.0001,12.991,0.0001,0.0001,1.009,1.009,-0,3.1877,-1,0,'1','3
','1',"here"),(4,-58,0,0,12.991,-2,1.009,-0,1.009,0,1.009,1,0.0001,-0.0001,1,-0,12441,'3	','	3','3	',"got"),(5,-32333,1,1,12.991,0.0001,32643,-0,16077,-1,2,2,1,1.009,12.991,1.009,1,"here",'1','k','
3'),(6,12.991,0,-9.183,12.991,-9.183,0.0001,1.009,0,2,2,12.991,0,2,1.009,-1,46,'3	','
3','-1','	3'),(7,-9.183,0.0001,-9.183,89.1873,-2,2,88.0997,0,1,0.0001,-0,1.009,12.991,1,-1,13284,'0','-0','o','3
'),(8,2,1.009,8683,12.991,-2,1.009,1.009,87.0259,-2,1,0,2,1,0.0001,-2,23203,'1','-1','0','3	'),(9,-0,1,-1,0,-0.0001,1.009,-1,12.991,1.009,2,-0,12.991,-0.0001,1,1.009,1,'1','
3','-0','	3');
create table table_10_binary_undef (
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
) character set binary ;
insert into table_10_binary_undef values (0,59.0853,0.0001,45,1.009,-1,0.0001,2,59.0957,0.0001,29.0134,-0.0001,2,-17882,12.991,-0.0001,12.991,'3	','x','1','3
'),(1,0.0001,22.1290,12.991,12.991,-0.0001,0,2,1,-11857,2,12.991,1,0.0001,117,1.009,0.0001,'k','3
','-1','3	'),(2,12.991,22719,2,0.0001,0.0001,0,15102,9.1283,24916,0,25,0.0001,0,0.0001,1,1,'3	','3	','t','-1'),(3,-9.183,110,1.009,2,2,0,-26270,0.0001,-0.0001,2,-2,1.009,71.1226,0.0001,-0.0001,0.1338,'3
',"it",'1','t'),(4,0.0001,23805,-1,39,0,0.0001,0,7727,2,0.0001,22474,0,-1,1,-0,0,'-1',"I'm",'
3','-1'),(5,-880,2,-2,2,1.009,12.991,-9.183,0.0001,0.0001,23.0704,0,12.991,-886,53.0325,-2,12.991,'-0',"I'm",'	3','1'),(6,-0,31.1867,1.009,72.1683,12.991,12.991,-2,12.991,2,121,-0.0001,0.0001,1,0,-2,1.009,'0','-0','-1','
3'),(7,-0,71.0825,-1,12.991,1.009,1,1.009,11393,-9.183,1.009,58.1888,0,1.009,52,-128,0,'q','3
','1','	3'),(8,2,0.0001,-0,1,2,0,-9.183,0,-2,1.009,0.0001,97.1617,-2,73.1320,1,2,'-1','	3',"hey","yes"),(9,-0,1,-37,45.1163,0.0001,1.009,2,12.991,2,12.991,-9.183,2,0,1,-2,2,"one",'-1','3	','3
');
create table table_30_utf8_4 (
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
partitions 4;
insert into table_30_utf8_4 values (0,-1,12.991,15721,0.0001,-0,26050,1.009,0.0001,1.009,12.991,0.0001,0,27163,0.0001,-2,0,'3
','3	','0','3	'),(1,0.0001,51.0312,1,2,1,0,0.0001,1,-0,2,-9.183,12.991,1,12.991,0.0001,0.0001,'0','3
','r','1'),(2,21.0903,0.0001,1.009,1,-9.183,46,-0.0001,0.0001,-5924,0.0001,22607,1,-1,2,17910,2,'0','-0','	3','-0'),(3,-1,19,-32643,2,81.1358,0.0001,-0,57,-0,2,-62,1.009,-30104,6276,0,1.009,'
3','3	','1','1'),(4,46,17,1,0,-1,1.009,-9.183,1.009,1.009,29568,11049,1.009,1,92.1752,0,0.0001,'3
','3
','-1','
3'),(5,2,1.009,-2,0.0001,0.0001,1,-1,0,-0.0001,0,-2,18,-2,0.0001,57,0,'0','-0','
3','-1'),(6,-0,17887,118,2,0.0001,45.1049,-9.183,0.0001,0.0001,34.0291,12.991,1,-1,2,0,1,"be",'3	','
3','-0'),(7,32,108,-41,16234,78,12.991,-1,2,-9.183,121,1.009,6.1363,-31,20.0766,-52,1,"a",'
3','0',"now"),(8,-0.0001,5307,2.0152,99.0746,2,56,1.009,85,-0.0001,0.0001,-2,2,2,1.009,0.0001,0,'g','0','3	','3	'),(9,23613,0,12.991,73,0.0001,12.991,42,12.991,2,0.0001,-13748,0.0001,74.1479,6.0951,2,126,'-0','-1','-0','3	'),(10,-2,0.0001,2,1.009,-0.0001,1.009,11.1646,0.0001,0.0001,1.009,12.991,95,2,14930,-9.183,2,'	3','1','3
','3
'),(11,0.0001,0,1,0.0001,-9.183,2,17691,2,-0.0001,12.991,1,49.1971,-0,2,-0,1.009,'e','	3','1','3
'),(12,12.991,1.009,-0.0001,0.0001,27,1,-0.0001,86.0094,24068,2,-0,25125,1.009,18415,12.991,2,'3	','3	','-0',"him"),(13,-2,0.0001,12.991,12.991,2,2,2,30.0659,1,0,-0.0001,1.009,1,1.009,0.0001,12.991,'0','3	','-1','3	'),(14,-13424,0,48.1438,2,12.991,0.0001,1,12.991,-2,0.0001,12.991,7.0137,-9.183,73.0722,0.0001,77.0761,'-0','-1','	3','
3'),(15,-2,1.009,23631,0.0001,-0.0001,12.991,-2,0,-0.0001,101,-9.183,2,0,0.0001,0.0001,29.1314,'1','-1','
3','0'),(16,-1,2,1.009,23.1514,-25,5,22272,0,-29595,2,-9.183,0,-0,2,12.991,0,'
3','-0',"there",'3
'),(17,-0,1.009,-2,1,-14201,32.0987,1,0.0001,-2,15126,31854,1,-9.183,28046,1,12.991,'-0','-0','-1','3
'),(18,40,24274,1.009,1.009,0,1.009,-2,12.991,2,47.0629,0,44.1680,0.0001,2,-0,1,'3
','-1','	3','3
'),(19,0.0001,2,-2,12.991,2,31139,35.1749,0.0001,-0.0001,12.991,1,1,12.991,0,2,1,'1','1','m','	3'),(20,-2,97.1381,-2,12.991,1.009,12.991,1,0,-9.183,2,0.0001,2,45.1861,28.1611,-15,1.009,'f','-1','0','3
'),(21,-0,0.0001,0.0001,0,-2,1.009,1.009,1,1.009,0.0001,0,15.1735,2,0,-0.0001,1.009,"we",'-1','3	','	3'),(22,42.0597,93.1712,-9.183,0.0001,-0.0001,0.0001,2,1,0,12.991,-2,12.991,-2,0.0001,0,1,'-0','-1','-1','3	'),(23,-0,12.991,0.0001,30517,-1,1,1,1.009,-0.0001,0.0001,0.0001,0,-1,0.0001,-9.183,1.009,'3	','	3','0','
3'),(24,-0,71.1447,29591,2,-0,1.009,11,86,2,81.1723,-0.0001,0.0001,0.0001,0.0001,-74,8006,'-1','-1','-0','	3'),(25,-2,0,19,12.991,1,72.1963,0,101,59.0447,12.991,0.0001,0,-9.183,1,0.0001,12.991,'	3','-1','3	',"would"),(26,-9.183,1,1,0.0001,12.991,0.0001,0,1,1.009,0,1,7351,-0.0001,12.991,2,29291,"when",'3	',"that's","had"),(27,0,0,-0.0001,2,-2374,1,28032,0.0001,0.0001,0.0001,1,2,0.0001,12.991,-9.183,12.991,'0','1','1','z'),(28,0,0.0001,1,12.991,-9.183,1.009,1.009,1,-0,2,-2,12.991,-0,1.1774,0,12.991,'3
','
3',"her",'-1'),(29,-9.183,0,-2,0,-9.183,4363,-0,0,0.0001,2,-0.0001,27.1110,-1,64,-2,2,'-0','
3','-0','0');
create table table_30_utf8_6 (
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
partitions 6;
insert into table_30_utf8_6 values (0,-1,0,0,0,1,1,-18737,1.009,-2,12.991,-2,0,12.991,1.009,1,88.0587,'	3','-0','-1','1'),(1,-0.0001,2,1,1,68.0215,1.009,58.0051,2.1638,2,1.009,2,8.0861,-0,0,-1,2,'	3','3
','
3','v'),(2,-1,1,-9.183,1,2,0.0001,-84,1.009,-1,1.009,0,2,-0,11844,-0,2,'
3',"when",'	3','
3'),(3,26.1566,0,-9.183,0.0001,0,12.991,12.991,0.0001,7330,44.0540,-0,0.0001,-0.0001,2,-0.0001,67.0204,'1','-1','-0','1'),(4,69,9.1872,-29302,22,0.0001,0,2,2,0,1,12.991,2,12.991,0,0,9,'j','3
','	3','3	'),(5,12.991,1,40.0540,79,2,1.009,12.991,50,0,86.1235,-0.0001,0,-1,98,0,5.0135,'j',"not",'-0','a'),(6,93.1153,12.991,-1,78.1536,-2,1,1.009,0.0001,-0,110,47.1286,46.1455,2,1,1,2,'3
','3	','3
','-0'),(7,2,27,0,0,-2,0.0001,-21,41,-1,0.0001,-2,5.1324,-0,76.1094,-0.0001,1,'d','3
','3
','3	'),(8,-20,12.991,1,43.0272,-9.183,0.0001,0,74,-9.183,1,-2,0.0001,-0,0.0001,-94,19196,'	3','y','	3','
3'),(9,0.0001,0.0001,1.009,12.991,46,1.009,2,69.1057,76.1027,12.991,-0.0001,56.1721,-9.183,1,-1,2,'3	','	3','q','1'),(10,26588,0,12.991,2,0.0001,73.1909,-0.0001,38.0826,10989,19,75,1.009,-0,33.1338,1.009,0,'-1','3	','	3','l'),(11,-0.0001,0,-0.0001,12.991,-2,2,-1,74.0647,-85,36,12.991,12.991,66,2,-2,30056,'c','1','s','-1'),(12,-10170,1.009,0.0001,0.0001,1.009,96,1.009,1.009,-1,1.009,-1,0.0001,12.991,0.0001,0,1,'1',"just",'0','0'),(13,43,2,-1,9,-2,2,45.0121,2,-0.0001,47,0.0001,1.009,12.991,0,0,1,"my",'
3','	3','	3'),(14,12.991,0,2,0,1,0.0001,0,0.0001,12.991,1.009,-0,85.0673,-0.0001,21.0028,-30013,0,"didn't","how",'z','0'),(15,0,12.991,12.991,16.0231,-9.183,1.009,-1,1.009,-0.0001,0,75,12.991,0,1,-9.183,12.991,'0','-1','3
','
3'),(16,7.1365,1.009,-2,2,-9.183,1,4478,27.0134,-2,12586,-1,0.0001,12.991,0.0001,12.991,0,'3	','-1','i','0'),(17,-0.0001,1.009,12.991,127,12.991,39.1860,92.1982,0,2,1.009,-1,0,35.1787,0.0001,12.991,5.0999,'0','d','0','	3'),(18,0,45.1373,-0.0001,1.009,1,1,1,12.991,1.009,0.0001,-1,1.009,12.991,12.991,-0,12.991,'0','3	','n','3
'),(19,12.991,98,-2,1,0.0001,0.0001,99,12.991,1.009,0,1.009,1,-0,0.0001,0,0,'
3','3	','1',"got"),(20,-1,1.009,0,1,1.009,1,-2,44.0092,1.009,0.0001,-0,116,1.009,2,-122,1,'3	','1',"then",'-1'),(21,-9.183,30587,-0,1.009,1,0,1.009,86.0293,2,1.009,0,10184,9.0345,1.009,-2,0.0001,"or",'-1','0','0'),(22,0,45.1197,-0.0001,5.0693,90,1.009,12.991,29.0798,-0,28.0921,12.991,0.0001,59.0348,83.1559,64,1.009,'-1','
3','	3','m'),(23,-1,1,12.991,0.0001,40,12.991,1.009,1.009,0.0001,2,2,67.0146,55.1630,16.1737,1,0.0001,'3
','-0','	3','	3'),(24,2,0.0001,12.991,1.009,54,1,-1,0,0,63.1111,2,94.1353,1.009,17,-1,2,'3	','3
','-1','0'),(25,-128,114,1,2,-0.0001,0,-2,12.991,0.0001,1,-2,12.991,-0.0001,0,1.009,1.009,'3
','-1','
3','
3'),(26,0,1.009,43.1691,12.991,0.0001,9.0536,12.991,73.0238,12.991,12.991,51,0,1,0.0001,-0.0001,93.1556,'-0','1','-0','0'),(27,19.1331,2,0.0001,1,-25,1.009,2,0.0001,12.991,0,5.1259,1,1,12.991,0,63.0470,'-1','-0',"I'm",'0'),(28,2,1,-2,1.009,1,60.1452,-2,1.009,12.991,55,59.1223,31.1023,1,1.009,-2,1.009,'
3','1',"like",'	3'),(29,12.991,0.0001,0,0.0001,12.991,0,0.0001,2,-0.0001,1,-2,1,-0,2,0.0001,63.0600,'	3','
3',"will",'-1');
create table table_30_utf8_undef (
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
) character set utf8 ;
insert into table_30_utf8_undef values (0,-2,0.0001,0.0001,1.009,-1,12.991,-23977,1,1.009,1,-0.0001,1,57.1252,10669,-9.183,2,'-0','3
','g','0'),(1,12.991,2,52,1,1.009,2,2,2,21454,2,-1,38.0884,12.991,1,-20705,22032,'0','d','	3','-0'),(2,-14102,0,93.0918,30.0630,2,2,2,11999,2,1,-30,2,-10102,12.991,35.0162,0,'
3','	3','	3','3	'),(3,0.0001,27634,16117,0.0001,50.0854,92.1168,1,2,-0,5.1286,0.0001,2,12.991,0.0001,0.0001,12.991,'r','-0','
3','-0'),(4,-0,0,-21,36.0533,47.0517,12.991,-1,2,-0.0001,1.009,2922,0,-9.183,0,2,0.0001,'3
','-1','-0','-1'),(5,-1,0,-9.183,2,1.009,21625,0.0001,112,31.1396,1.009,-1,2,12.991,73,-22398,37.0379,'3
','	3','-1',"it's"),(6,2,1,0,1.009,-1,31,89.0237,12.991,1,12.991,2,0,1.009,9.0990,61.1783,1.009,'1','-0','	3','1'),(7,36.0287,29.1579,34,1,-0,0.0001,-0.0001,0.0001,-14977,0.0001,61,28934,0,0,-1,2,'	3','-1','	3','
3'),(8,-2,69.1818,-2,1,-2,12.991,-2238,14.0136,33.1547,1,-1,2,12.991,12.991,-2,12.991,'0','
3','0','-0'),(9,0.0001,0.0001,82.1243,88.0069,1,18.0444,0,12.991,-1,23374,-0.0001,33.1409,1.009,4,1.009,7.0090,'-0','	3','3	','u'),(10,-0,12.991,-29560,1,1.009,32.0279,-1,1.009,1,116,-0,41.1811,-20793,12.991,0.0001,12.991,'	3','1','
3',"as"),(11,-2,47.1703,12.991,4864,12.991,32.0815,-2,1,-0.0001,2,1,14.0200,-0,0,1.009,1.009,'	3','-0','3
','1'),(12,-9.183,58.1637,12.991,2,12.991,0,-0,68.0763,-0.0001,12.991,0.0001,22,41.1688,1,-127,2,'1','e',"yes",'-1'),(13,-0,1,1.009,77,2,97.0288,2,12.991,0,0.0001,1.009,12.991,0,0,0,1.009,'1',"had",'0','-1'),(14,0,0.0001,92.0834,36,1.009,12.991,2,2,-9.183,2,0,1.1859,2,1.009,-0.0001,1.009,'	3',"so",'1','n'),(15,-1,1,-0,1.009,2,29.1055,0.0001,0.0001,0,43,10.1847,20686,2,1,1.009,12.991,'1','3	','f','1'),(16,-0,1.009,-1,0.0001,-0,0,-9.183,12.991,1.009,11574,-0.0001,1.009,-0.0001,12.991,7.1392,2,'
3','-0','	3','1'),(17,-0.0001,0,-0.0001,0.0001,-0.0001,0.0001,1,0,-0,1,-2,1,-0.0001,1,1,56.1598,'3	','-0','	3','
3'),(18,0.0001,31.1784,12.991,12.991,-1,0,-0,74,1,0.0001,-2,53,-1,18211,-0,1.009,'1',"back",'-1','-0'),(19,-0.0001,1,77.0905,1,1.009,2,12.991,32675,12.991,1,0,2,-9.183,0,-0,0.0001,'c','0','3
',"is"),(20,-9.183,1,-2,1.009,-9.183,0.0001,1,12.991,-1,2,121,0,-1,0,0.0001,57.1324,"look",'3
','-0','0'),(21,-0.0001,1.009,1.009,1,12.991,0,0.0001,1,-2244,2,15,8,0,1,-0,1,'h','0',"this",'	3'),(22,1,2,67.0938,88,-0,77,-0,1.009,-0,2,0,63.1954,0,0,1.009,12.991,'-0','0','
3','
3'),(23,0.0001,1.009,-15157,12.991,-2,65.1155,-118,1.009,12.991,0,-2,1.009,-0.0001,0.0001,-9.183,0.0001,"are",'3
','-1','3
'),(24,1.009,26703,1,0,-9.183,2,0.0001,89.0639,12.991,1.009,0.0001,0,89.1324,1.009,31.1127,119,'
3','-0','0','-1'),(25,0,1,8383,3,2620,59,-0.0001,12.991,-0,1.009,0.0001,0,-0,12.991,92.0776,1,'3	','	3',"on",'-1'),(26,-0.0001,2,0.0001,5583,12.991,0.0001,-2,0.0001,0,83,0.0001,0,96.0616,0.0001,99.1790,1.009,'	3','-0','f','q'),(27,-1,0.0001,2,0,-27,1.009,0,1,-1,1.009,21,12.991,37,12.991,1,58.0685,"there",'
3','-0','-1'),(28,1.009,0,1.009,0,2,45,0,0.0001,12.991,12.991,-1,0,2,12.991,16.0903,0.0001,'-1','3
','-0',"when"),(29,-2,1,-0,0.0001,2,0.0001,-0.0001,26.1406,-0,0.0001,-0,1,1.009,1.009,12.991,2,"did",'3	','-0','-1');
create table table_30_latin1_4 (
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
) character set latin1 
partition by hash(pk)
partitions 4;
insert into table_30_latin1_4 values (0,-9.183,2217,1.009,0.0001,1.009,92.0223,12.991,45.0989,-0.0001,100.0603,0.0001,2,1.009,0,0.0001,0.0001,'
3','1','
3','1'),(1,2,12.991,1.009,0.0001,-0.0001,2,-2,12.991,0,0.0001,-9.183,76.1019,0.0001,96.1097,-0.0001,15,'3
','
3','w',"his"),(2,12.991,100.1764,-1,1,-2,1.009,113,27846,70.1808,2,0.0001,2,1,71,-20450,1,'0','y','c',"that"),(3,12.991,0,1.009,12.991,0.0001,1.009,87.0312,0,-0.0001,1.009,-0.0001,0,0.0001,1.009,-23192,0.0001,'3	','3
','u','
3'),(4,-1,0,0.0001,0,1,12.991,2926,2,17450,2,13406,12.991,-0.0001,36.0565,-9.183,2,"now",'-0','3
','-0'),(5,-0,25.0219,2491,0.0001,-0,1,12.991,0,-1,72.0778,-9.183,12.991,-9.183,0,0.0001,23.0069,'3	','v','0','1'),(6,2,1.009,105,53.0593,-0.0001,69.1163,2,14,-9.183,4.1427,1,0.0001,-118,0.0001,1.009,0,'
3','c','
3','	3'),(7,-0,2,-2,1,122,1,-9.183,1,1.009,1,-16,2,-0.0001,66,1,1,'-0','-0','0','-0'),(8,2,0,96.0631,1,-0.0001,1.009,12.991,25323,-0.0001,15,-15316,0,80.0029,0,1,88.1070,'1','3	','0','-1'),(9,70.1713,41,12.991,0.0001,12.991,12.991,-0.0001,24062,0.0001,0.0001,-0,0.0001,12.991,0,0.0001,0.0001,'	3','
3','1','1'),(10,1,1.009,36.0343,0,2,1.009,-22791,0,-2120,71.0902,-2,0,-9.183,0,-8674,12.991,'0','-1','-1',"it's"),(11,2,0.0001,-19098,6.1371,0.0001,7510,0,2,0,1,12.991,1,2,0,-15,36,"out",'3
','
3','n'),(12,-0,2,48.0668,1,-9.183,12.991,0,16.1310,1,1,-0,6,-1,1,-1,98,'3	','i','0','3	'),(13,-1,0,-1,0.0001,71.0494,8.0021,-0.0001,12.991,-9.183,0,24.1643,20670,-2,2,94.0957,1.009,'
3','-0','h','
3'),(14,-1,2,-0,0,2,12.991,12.991,12.991,0.0001,12.991,9.0441,2,1.009,0.0001,12.991,80.1216,'3
','-0','1','n'),(15,2,2,-1,1.009,-0,26.1199,1.009,2,14,49.0936,36.0770,0.0001,1.009,113,-1,2,'3
','
3','3	','3
'),(16,-1,2,9769,22313,0.0001,12.991,35.0971,45.1928,2,1,-2,2,-4009,1.009,12,20872,'h','r','h','-0'),(17,-1,0,-0.0001,1,2,98.0522,2,0.0001,-0,0,2,1.009,80.1364,1,2,0,'	3','-0','w','-1'),(18,-2,0.0001,-2,1.009,-1,1.009,-6896,12.991,-2,2,0.0001,1,-1,13.1196,-124,2,'-1','0','3	','1'),(19,0,1.009,0.0001,58.1115,2,0,1.009,5.1879,-1,1487,1.009,42.1614,1,1,0.0001,0,'0',"are",'
3','
3'),(20,-2,15.1800,-25990,2,2,1.009,-0.0001,2,12.991,0,-2,0,-28376,1.009,-1,51.0145,'o','-1','3	','1'),(21,0,33.1756,3831,0,1,0.0001,-0.0001,15813,-20739,79.0581,-7798,3.1792,8,12.991,-0.0001,12.991,'3
','0','3
','
3'),(22,-0.0001,0.0001,-19,1,0,12.991,12620,2,-2,1.009,-0,12.991,-0.0001,12.991,-0,1,'3
','1','
3','-0'),(23,2,12.991,-0,0.0001,12.991,12.991,2,12.991,52.0791,1.009,2,1,56,1.009,0.0001,1.009,'
3','j','3
','
3'),(24,2,0.0001,0,2,-2,1,1.009,1.009,-0.0001,1,0,12.991,-1,1.009,12.991,0.0001,'i','3
','-1',"could"),(25,1.009,1,-19121,1,-2,1.009,57.0407,12.991,12.991,26314,1.009,1,-2,1,1.009,0.0001,'3
','0',"he's",'	3'),(26,1,23.0509,0.0001,52.0317,2,118,0.0001,0,-35,0.0001,-2,1.009,87,86.1017,12.991,2,'0','1','3	',"see"),(27,-29866,10.1998,12.991,12.991,1,2,12.991,0,12.991,13367,0,117,49,1.009,2,1,'3
','-1','-0','-0'),(28,12.991,12.991,6588,87.1089,36.1931,0,1.009,49.0127,-0.0001,53.0615,1,0,0.0001,15653,-9.183,2,'1','1','-0','1'),(29,0,0.0001,-0,1,-1,12.991,-9.183,2,-9.183,0,1,0,0,12.991,-59,16344,'-0','m','0',"me");
create table table_30_latin1_6 (
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
) character set latin1 
partition by hash(pk)
partitions 6;
insert into table_30_latin1_6 values (0,1.009,0.0001,0.0001,39.0002,0,0,2,12.991,8.1012,10463,99.0808,89,-31914,89.0376,0,1,'0','
3','-0','0'),(1,19.1370,1,35.1142,53.1058,-2,0,-9.183,12.991,-9.183,0,12.991,2,0,2,12.991,50.0218,'0','1','j','x'),(2,-2,88.1853,2,45,-0,1,-0,0,29.0645,0,-9.183,20537,-9.183,0.0001,62,18.0931,'-0','-1','
3',"one"),(3,1,12.991,0.0001,1,-0.0001,2,93.1963,0.0001,1.1163,0.0001,-0,12.991,-2,0.0001,-22450,1,'
3','n','0','-0'),(4,-2,12.991,-0,0,26.1557,109,-0,1,1,1,0.0001,5437,2,12.991,0,2,'1',"say",'-0','u'),(5,-9.183,2,16257,67,2,48.1344,-49,25707,-10276,5493,-2,0.0001,2,1,-2,12.991,'0','1','3	','-1'),(6,6032,1,-0.0001,1,0.0001,1.009,-2,0,1,116,-1,7315,87.0926,12.991,1,59.1314,'k','-1','	3','o'),(7,-29,12.991,-0,2,12.991,59.1147,-9.183,0,20250,0,2,1,-2,93.1317,1,12.991,'3
','-0','-1','3	'),(8,0,1.009,-2,69.0983,0,1,-0.0001,14349,-0.0001,0.0001,-0.0001,12.991,0.0001,70.1499,-15834,12.991,'-0','-1','3	','0'),(9,2,47,1,1,87.1282,21569,-0.0001,12.991,12.991,0,0,1,-94,0.0001,2,1.009,'3	',"then","are","right"),(10,1.009,21648,-0,12.991,4180,2,0,1,2,0,-2,2,-9.183,2,0,12.991,'
3','3	','0',"did"),(11,24078,0.0001,-784,1.009,-2,2,-0.0001,1,-9.183,55,0,1,22.0858,2,-1,0.0001,'0','
3','	3','3	'),(12,1,2,1,1.009,-20588,1,1.009,12.991,20.1149,1,2,1.009,-0.0001,58.0616,-2,70,'0','1','-1','f'),(13,-0,12.991,-66,0.0001,11,55.1166,-1,53,1,2,0.0001,12.991,-9.183,0.0001,2,1,'
3','3	','e','0'),(14,28,2,-0.0001,38.0579,1.009,42,-4590,37.0342,0,1,-0,12.991,0.0001,45.1072,-0,24.0361,'3	','3
','u','3	'),(15,1.009,0,-1,0.0001,-2,2,-2,0,0,2,12.991,1.009,-9.183,1,38.1115,0,'n','
3','-1','-0'),(16,-2,19.1728,-9.183,54.1490,-9.183,0.0001,89.1554,0,-2,2,2,1.009,-0.0001,1.009,1,12.991,"why",'3
','1','
3'),(17,26.1004,2,1,1.009,0.0001,12.991,-13,1.009,-2,12.991,-0.0001,0,12.991,1,-0,0.0001,'-1','0','0','3	'),(18,-8359,2.0410,-0.0001,0.0001,-2,12.991,5811,0,1.009,1.009,49,52,-1,57.0137,-1,2,'-1','3
','3	','
3'),(19,-0,1,-0.0001,1.009,122,0,2,2,66,97.1479,23,0.0001,-2,22,-9.183,0.0001,'
3','0','3	','-0'),(20,1.009,0.0001,-18482,56.1560,-13596,0.0001,-0,2,-0,1,-0.0001,1,12.991,0.0001,36,12.991,'0','3
','
3','-1'),(21,0.0001,12.991,1.009,1,-0.0001,0,2,2,-2,25.1320,-1,0.0001,20.1633,1,-0.0001,12.991,'
3','1','0','3
'),(22,-9.183,1.009,-9.183,1,-2,1.009,106,1.009,2,0,-0,12.991,12.991,0.0001,-8827,2,"know",'1','	3','-1'),(23,12.991,1.009,77,1.009,-9.183,3.0610,0.0001,29.0883,-0.0001,0,-0,2,26.0541,12.991,0,0.0001,'	3','3
','
3',"my"),(24,11.1778,12.991,97.1309,0.0001,12.991,1.009,1.1090,1,-1,2,-0.0001,0.0001,0,1,-85,1,'z','-0','-0','-0'),(25,-0,1.009,1.009,0,1.009,0,28.0797,1,7.0416,29.1591,-0.0001,1,56.1185,12.991,-9.183,1.009,'0','	3','-0','0'),(26,-2,541,0.0001,12.991,-0.0001,1,0,19545,2,0,20.1301,85.0497,-9.183,1,87.1609,1,'-0','-0','1','d'),(27,12.991,1,2,2,22002,0.0001,2,52.0756,-1,34.1814,52.0411,15050,-27412,12.991,2,13.1315,'3	','-1','	3','
3'),(28,26916,2,-2,0.0001,-0,2,-2,98.0883,-0.0001,1,-0.0001,1.009,0.0001,86.0790,-0,1.009,'	3','
3','	3','0'),(29,-0,1,-29136,1,-0.0001,0,0,1.009,-1,107,-0.0001,1.009,12.991,83.0858,-0.0001,1,'1','
3','1',"do");
create table table_30_latin1_undef (
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
) character set latin1 ;
insert into table_30_latin1_undef values (0,29.0520,0.0001,-0.0001,0,12.991,0,-1,43.1866,2,0.0001,-9.183,0,1.009,2,0,0,'1','	3','	3','x'),(1,-67,0.0001,0.0001,1.009,-0,0.0001,12.0862,1.009,2.0286,12.991,-9.183,1,1,22795,12.991,1.009,'0','m','0','1'),(2,12.991,66,-0.0001,0.0001,1,0,12.991,2,-9.183,0.0001,0,12.991,0.0001,0.0001,1.009,2,'3	',"or","I'll","if"),(3,-2,2,1,0.0001,0.0001,1,-85,1,1.009,12.991,-2,2,12.991,2,-9.183,1,'-0','d','-1','	3'),(4,-9.183,1.009,0.0001,12.991,8947,1,0,2,55.0949,30186,1,115,-9.183,2,-1,0,"she",'1',"with",'1'),(5,123,0.0001,4931,12.1885,1.009,1,-1,0.0001,-0,1,1,2,45.1963,1.009,30,55,'	3','
3','0','1'),(6,1,9,1.009,2,-2,27.0058,0.0001,78.1255,-2,12.991,77.0634,0.0001,8.1279,12.991,0,1.009,'	3','-0','-0','	3'),(7,69,0.0001,-0,12.991,1,86.1254,-1,46.1977,2,82.0035,-9.183,0.0001,-0.0001,2,0,0,'3
','-0','
3','3
'),(8,85,0.0001,12.991,0.0001,0.0001,0.0001,-9.183,1,112,12.991,104,0.0001,1,0,86,1,'-1','	3','0','
3'),(9,0,0,2,12.991,2,67.0097,-0.0001,12.991,96.1871,1,1,1,0,80.0722,-0,12.991,'	3','f','	3','0'),(10,10849,12.991,-0.0001,1.009,-2,0,0,98,-0,0.0001,-9.183,2,1,0.0001,-9.183,21623,'0',"her",'-0','	3'),(11,0,1,1,1.009,0.0001,1,-122,0.0001,12.991,0.0001,25604,2,56.0089,1,-9.183,12.991,'0',"could",'	3','	3'),(12,-19848,13141,6709,125,-0,2,1,2,19894,1,-0,60.0486,-2,1,91.0587,0.0001,'1',"had",'0','	3'),(13,12.991,0,95.0029,2,-14189,96.0096,99.0229,8584,1,12.991,1.009,55.0258,-21771,0.0001,-2,1.009,'1','
3','i','
3'),(14,2,0.0001,1,12.991,2,2,1,2,-0,1.009,-0.0001,1.009,0,0.0001,-1,88.0905,'0','1','0','
3'),(15,56,12.991,-0.0001,44,68.0872,112,0.0001,2,12.991,1.009,-1,1.009,-2,2,12.991,1.009,'
3','
3','1','3	'),(16,29368,12.991,-0,26,1,0,0,15455,0,26534,12.991,12.991,-0,2,0.0001,1.009,'
3','0','3
','3	'),(17,-1,1.009,-0.0001,99,-2,0,-0.0001,56,-9.183,0.0001,0,2,-0.0001,0,-9.183,0,'3	','3	','3
','	3'),(18,12.991,0,-0.0001,18.0946,0.0001,45.1081,1.009,1.009,9,1.009,-2,90.0945,-0.0001,0.0001,13.1990,0.0001,'3	','-1','
3','1'),(19,-9.183,12.991,0,2,-9.183,14166,0.0001,17.1989,-0,1,0,0.0001,1,0.0001,-0,1,'r','-1','3	','t'),(20,-0.0001,1,1.009,1.009,2,0.0001,15419,2,1,47.1709,-0,12.991,12.991,1,-1,2,'-0','w','
3','-0'),(21,12.991,1,0.0001,0,0.0001,2,1,1.009,-2,419,-2,1,-0,84.0898,-0,67.0059,'-1','3	','3	','0'),(22,0.0001,0.0001,70.0923,114,115,105,2,1.009,-9.183,2487,-1,12.991,-9.183,13.0019,1,67.1039,'0',"with",'1','	3'),(23,-819,0,-9.183,0,121,12.991,79.1487,0,95.0175,15.0998,2,0.0001,-0.0001,1,1,1,'	3','1','1','3	'),(24,2,0,-2,94,2,0,0.0001,2,0,30.0643,-47,1.009,-0.0001,0,82,2,'0',"we","can",'	3'),(25,1.009,0,12.991,1,1,51.0946,-1,1.009,-1,1.009,1,0,1.009,12.991,1,22889,'-1','w','	3','s'),(26,54,1.009,-0.0001,0.0001,0.0001,1,1,0,-0.0001,7425,-125,12.991,1.009,2,1.009,1,'1','1','-0','	3'),(27,0,0.0001,0,0,1,1,-91,1.009,-0.0001,126,1,88,-0.0001,1,-2,1.009,'0',"the",'3
','3
'),(28,-0.0001,12.991,-17,0.0001,-0.0001,0,109,2,-2,1,-0.0001,1.009,0.0001,1.009,53.1786,1,'1','	3',"something",'3
'),(29,2,12.991,9.1305,1.009,-2,0.0001,-30807,1,-0.0001,0.0001,-251,1,-0.0001,12.991,12.991,1,'1','3	','	3','	3');
create table table_30_binary_4 (
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
) character set binary 
partition by hash(pk)
partitions 4;
insert into table_30_binary_4 values (0,75.1477,95.0009,-9.183,0,-0.0001,22,-0.0001,1,-72,2,1.009,78.1318,-16016,15732,2,93,'
3','0','	3','
3'),(1,-9.183,12.991,1.009,0.0001,1.009,28,-0,0.0001,-0,2,0,12.991,-2,1,11855,1.009,'
3','0','-1',"had"),(2,-9.183,2,-2,2,-0.0001,974,-1,12.991,-1,1.009,0.0001,1,-9.183,2,0,1,'-0','0','-1','1'),(3,0,0.0001,-0.0001,1,79.1644,2,-1,0.0001,80.0037,1.009,100.1167,36.0115,0.0001,12.991,-6689,12.991,'-0','
3','w','	3'),(4,0,2,12.991,0.0001,-7491,86,1,12.991,1,10.1381,124,1.009,-68,76.1084,-0.0001,0.0001,'3	',"i",'3
','m'),(5,-0,63.1707,-1,2,1.009,1,0,29040,-1,1.009,-2,12.991,-9.183,2,0,12.991,'-1','1','3
','
3'),(6,1,12.991,1.009,0.0001,-9.183,1.009,-24035,1,-2,1,1,0,-1,0.0001,0.0001,2,'
3','1','3	','1'),(7,-1,12.991,0.0001,1,1.009,12.991,5925,0,-31760,18574,0,1.009,-9.183,0,0,0.0001,'0','-0','1','1'),(8,1.009,6097,0.0001,0,19249,12.991,2,1,12.991,91,-2,2.1416,-9.183,1,0,13.0255,'3	','
3','3	','0'),(9,12.991,12.991,-1,1.009,121,12.991,2,1,-2,0,0,1.009,0.0001,1,-20028,93.1039,'
3','3	','	3',"who"),(10,12.991,0,-29890,0,-6295,12.991,12.991,0.0001,-9.183,0.0001,12.991,1,12.991,2,-9.183,1,'-1','
3','	3','z'),(11,-1,2,1,78.0615,-9.183,1.009,12.991,2,-2,0.0001,0,12.991,-0.0001,0,-9.183,1.009,'3
','0','c','1'),(12,-0.0001,1,-0.0001,1,12.991,0.0001,2,97,24888,2,-9.183,42.0194,-0,1.009,0.0001,12.991,'3	',"who",'1','1'),(13,1,4.0123,2,1,-0.0001,88,2,19485,22.1562,104,-76,12.991,-1,0,66.1403,0.0001,'g','0',"going",'	3'),(14,15.1797,1,-1,12.991,13675,0.0001,1,1,-0.0001,2,-0,0,-0,0,-81,2,"if",'0','0','3	'),(15,0,0,0,0.0001,-9.183,1.009,1.009,1.009,0,1.009,0,12.991,-2,1,-1,26.0553,'3
','	3','
3','3
'),(16,-0,0,0,0.0001,0.0001,0.0001,1.009,0.0001,0,29048,84.0753,0.0001,-30002,12.991,-1,0,'1','-0','3	','-0'),(17,-2,12.991,1,20.1614,0.0001,0.0001,0,0,0,0.0001,2,1.009,99.1331,0,80.1206,0,'	3','-1','-1','-0'),(18,2.1027,2,-9,12.991,0.0001,12.991,2,23,-112,0.0001,25065,0.0582,12.991,0,-2,1,'3
','3
','1','s'),(19,0,1.009,0.0001,0.0001,14331,1,0,2,-0.0001,12.991,8,67.0203,12.991,2,-2,1,"I'll",'-0','-1','3	'),(20,23467,1,2,0,-1,0.0001,0.0001,0,-2,1,2,2,-1,53,2,1,'1','-0','3	','3	'),(21,1,0,2,0.0001,1,1,108,92,1,0.0001,79.1183,2,-2,12.991,-2,61.1318,'3	','0','
3','
3'),(22,-1,14715,0,2,9687,0,-1,1288,12.991,0,1,12.991,0.0001,0.0001,-9.183,2,"not",'-0','0','3
'),(23,10416,0.0001,-2,30,-0.0001,1,61.0106,12.991,2,1,-0,1.009,19,2,0.0001,1.009,'0','	3','1','-1'),(24,0.0001,0.0001,0,2,-9.183,1.009,1.009,1.009,1.009,0.0001,12.991,1.009,-9.183,0,1,44.0135,'p','-1','-1','-0'),(25,-0.0001,1,-0,0.0001,31859,15491,1,1,-2,20.1594,0,2,0.0001,2,-1,0.0001,'1','-1','	3','3
'),(26,-1,6.1576,2,75.0812,-9.183,90.0260,-0,0.0001,1.009,0.0001,1,1,-33,37.0686,1,12.991,'-0',"good",'3
','3	'),(27,37.0393,0.0001,-0.0001,1,-2,1.009,-2,2,12.991,2,-9.183,20343,0,2,0.0001,1,'
3','0',"his","how"),(28,-5,1,-9.183,2,-1,0,-0.0001,1,-1,118,-2,5.1761,-1,2,-9.183,1,'
3','-0','0','3
'),(29,-2,2,85,6.0107,-1,27.1481,-9.183,2,-1,1.009,0.0001,0,-17,2,-9.183,0,'3	','
3','-0','t');
create table table_30_binary_6 (
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
) character set binary 
partition by hash(pk)
partitions 6;
insert into table_30_binary_6 values (0,12.991,0,-40,12.991,0,1,2,0,-0.0001,123,-4,54,19.1077,16304,10623,1.009,'0','v','o',"been"),(1,-0.0001,1.009,1.009,2,2,1.009,-0.0001,12.991,12.991,2.0856,-67,12.991,0.0001,78.1332,-2,946,'0','0','3
','	3'),(2,-14819,52.0972,-0,0.0001,-9.183,12.991,1.009,0,-517,0,-0,12.991,-2,2,0,2,'	3',"can",'1',"yeah"),(3,0,1.009,0,12.991,0,12.991,-0.0001,1.009,-17712,3746,-80,2,-9.183,1.009,0,74.0912,"got",'-1',"going",'-0'),(4,12.991,0,1.009,27682,12.991,0,-1,0.0001,2,19.1804,12.991,96.1668,1.009,75,-123,0,'3
','-0','0','3	'),(5,1,0.0001,1.009,19.0570,1.009,1,0,12.991,-9.183,1.009,0.0001,72.0144,2,0.0001,-0.0001,1.009,'0','s','3	','	3'),(6,-2,2,12.991,2,-0.0001,1,-2,13,-1,12.991,2,0,0.0001,0.0001,0.0001,1.009,'
3','1','0','3	'),(7,-9.183,0,86.0861,0.0001,0.0001,1.009,2,5426,-0.0001,0,-2,0,79.0732,12.991,1,1,'k',"yes",'3	','l'),(8,12.991,1,-0.0001,0,-1,0,1.009,0.0001,9.1696,1,0,26704,-2,76,6947,2,'
3','-1','3
','3	'),(9,1,0.0001,0.0001,0.0001,-2,1,-2,12.991,-2,1,-9.183,12.991,0.0001,2,8035,96.0829,'	3','	3','	3','3
'),(10,67,12.991,-1,12.991,-9.183,12.991,-1,0,1,1.009,-1,1.009,0,2,1.009,1,'	3','
3','-1','0'),(11,-1,69.1308,-9.183,1,1.009,0,-1,43.0526,-1,2,0.0001,1.009,1.009,12.991,0.0001,0,'3	','	3','	3','	3'),(12,0,1,-9.183,2,-0,0.0001,2,65,-0,1.009,-9.183,7634,71,0.0001,-9.183,12.991,'3	','
3','-0','3
'),(13,18.1647,32351,1,41.1082,16703,12.991,1.009,72.1919,-0,1,1.009,0.0001,12.991,1,-2,6.0985,'0','3
',"some",'-0'),(14,-6798,0,0.0001,71.1672,0.0001,1.009,1,1.009,-9.183,118,-0,1,-1,0.0001,0.0001,76,'
3','-1','1','-1'),(15,1.009,1.009,2,0,1,12.991,-9.183,10913,2,2,2,21,-0.0001,2,-9.183,20.1156,'e','3	','
3','3
'),(16,2,12.991,-31791,1.009,-0,12.991,-0.0001,12.991,0.0001,0.0001,1.009,15530,1,1.009,-0.0001,7.0528,'	3',"about",'3
','
3'),(17,2,1,-2,15.0847,85.0195,2,9688,12.991,1,1,-1,12.991,-2,52.0099,12.991,2,'0','3
','3
',"ok"),(18,69.1296,2,-0,2,-0.0001,2,-1,1.009,-9.183,1,-0.0001,0.0001,-0,1.009,-1,2,'	3','3
','	3','1'),(19,1.009,0.0001,-3511,12.991,-9.183,1,45.1496,28537,12.991,0.0001,-2,2,1,1,-9.183,12.991,'0','j',"from",'	3'),(20,1.009,0.0001,1,2,32,47.1090,-9.183,12.991,-2,0,-42,1.009,0,2,-0.0001,17820,'0','
3','
3','0'),(21,-1,12.991,0.0001,2,12.991,1,-2,6.1948,1,1,12.991,1,0.0001,46,2,0.0001,'3
','3
','	3','1'),(22,2,2,1.009,1.009,1,1,-1,1.009,1.009,23.0836,15.1072,2,1,0,-2,1.009,'
3','-1','3
','3	'),(23,-9.183,0,1.009,82,-9.183,12.0022,2,0.0001,-2,1.009,1,90,0.0001,2,0.0001,86.0802,'0','	3','
3','-1'),(24,1,55,-9.183,67.0449,-2,0,-9.183,12.991,-0.0001,1,1633,0.0001,1.009,1,-0,0,'
3',"say",'-1','	3'),(25,12.991,19551,2,5119,1.009,1.009,-53,0.0001,99.1684,0,1,0,0.0001,1,60.1825,0.0001,'-1','3
','d','
3'),(26,17567,0.0001,20203,72.0121,-26311,1,-82,0,2,24.1553,1,74.0913,0.0001,0,22.0231,2,'0',"don't","up",'1'),(27,1.009,0.0001,-1,2,-2,2,12.991,1,28488,18002,3,68.0057,-0.0001,2,2,2,'3	','v','-1','	3'),(28,53,14709,-1,1.009,6260,1,1,1,-86,8.1518,1,34,12.991,0.0001,-1,0,'	3','	3','1','-1'),(29,2,0.0001,0,59.1775,-2,0,12.991,26684,-9.183,12.991,-0,1.009,1.009,0,2,5.1021,'i','3
','3
','	3');
create table table_30_binary_undef (
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
) character set binary ;
insert into table_30_binary_undef values (0,12.991,1,-127,0.0001,1.009,1,-1,1,-0,12.991,9699,12.991,2,0.0001,2,1.009,'-1','	3','3
','-0'),(1,0.0001,106,2,1.009,231,88.1203,-1,70.0888,1.009,1,-2,1,-9.183,1.009,-0,28.0711,'1','	3','0','1'),(2,2,0.0001,-23325,1.009,1.009,0.0001,-0,12.991,-0,100.1200,-9.183,1,-1,2,-2,0.0001,'-0','r','0','-0'),(3,0.0001,1,-53,19927,-2,82,-0.0001,62,-9.183,1,0.0001,1,1.009,106,12.991,1.009,'3
','0',"yeah","him"),(4,-1,0,17797,4,-1,12.991,0,0.0001,0.0130,1.009,1.009,2.0163,12853,0,-9.183,0.0960,'
3','d','-0',"think"),(5,-9.183,2,-2,38,-9.183,0,0.0001,1,-645,2,1,63,-9.183,12.991,-0,116,'3	','3	',"your",'1'),(6,-0.0001,1,-0,78,0,15252,53,2,-0.0001,1.009,1.009,67.0827,0,1.009,10,1.009,'-0',"just","know",'-0'),(7,1,40.0288,2,1.009,-8103,1.009,-1,12.991,0.0001,0.0001,-10247,79.0040,48.0909,1.009,1,10,'-0','	3','3	','1'),(8,0.0001,1,2,12.991,-9.183,12.991,2,2,50.0333,2,12.991,0,-20714,0,-1,0,'1',"my",'-0','0'),(9,12.991,2,2,2,16789,0,1.009,111,12.991,54.0361,-2,1.009,-0,2,2,4525,'-0','-1','3	','-0'),(10,1.009,1.009,1,1.009,-2,1.1137,0.0001,12.991,1,1.009,-0.0001,2,-0,0.0001,2,12.991,'0','1','-0','1'),(11,-0.0001,1,1,8,11951,1.009,-0.0001,52.0244,0.0001,12.991,-2,0.0001,1,79.1369,-0,12.991,'3
','-0','-0','0'),(12,-2,0.0001,0,60,-2,70,-9.183,0,-0,2,0,1,-9.183,1.009,-1,1,'	3','-0','c','
3'),(13,1,1.009,90.0045,10.0513,0,10,-0,1,0.0001,0,1.009,102,1.009,4428,1.009,2,'3
','m','-0','3	'),(14,5.0454,1.009,1.009,1,-9.183,0,-2,2,-1,57.1188,-2,0.0001,-1,0,-2,2,'0','0','0','3
'),(15,2,12.991,-0,2,1.009,55.0385,1,0.0001,1,7038,27.0820,1,0.0001,1.009,-9.183,39,'i','	3','-0','
3'),(16,-0,0,22.0400,1,-1,0.0001,-2,0,75.0523,8,-9.183,12.991,-9.183,0.0001,-0,0.0001,'3
','1','0','1'),(17,-47,0,-0,2,-0,2,81.1906,1.009,-9.183,1.009,1,0.0001,-2,0,81.0486,26,'3
','0',"in",'-0'),(18,-9.183,1.009,2,1.009,12.991,1,-0,0,1,2,20.1087,12.991,-9.183,57,-9.183,12.991,'-1','s','-0','3
'),(19,-9.183,6602,66.0074,0.0001,-0,1,68.1334,0.0001,-1,24,-9.183,0.0001,1,0.0001,12.991,1,'3	','0','-0','0'),(20,0.0001,12.991,-0.0001,0.0001,-0.0001,0.0001,2,99.0120,-2,1,0.0001,1,0,0,-2,41.0731,'3
','-1','k','1'),(21,1.009,1.009,-0.0001,0,-9.183,12.991,-0.0001,26,-1,2,2,2,0.0001,1.009,0,0.0001,"from",'-0','a','-1'),(22,-2,17818,-2,1,2,0.0001,2,12.991,2,28,1,93.1882,2,0,-9.183,0.0001,'-0','
3',"is",'	3'),(23,-9.183,16.1384,22376,84,0.1162,2,-0.0001,0.0001,-1,1.009,2,0,-0.0001,29.0532,0,1.009,'	3','-0','
3','
3'),(24,12590,45.0125,23.1535,110,-2,1,1,0,0.0001,0.0001,-0,1.009,12.991,0.0001,1.009,0.0001,'
3','3
','3
','3	'),(25,-0,1,1,0.0001,-26290,1,-9.183,1.009,-1,53,0,0,-0.0001,1,1,1.009,'	3','3	','3	','-0'),(26,76,15.1446,2,0.0001,16381,1,19.0818,0,36.1772,63.0051,-1,1,-9.183,27.0342,2,42,"because",'0','3	',"didn't"),(27,-1,72.0336,-0.0001,0.0001,1,41.1205,-2,99.1510,-1,12.991,1.009,12.991,20325,12.991,0.0001,1.009,'3
','-0','3
','
3'),(28,1.009,95.0890,-1,1.009,0.0001,2,12.991,2,1.009,58.1366,-2,1,14158,33.0736,2,12.991,'3	','1','3	',"your"),(29,1.009,2,82.1686,1,0,12.991,-0.0001,12.991,0.0001,0,115,97.0507,0.0001,61,1,44.0416,'3	','3	','1','3
');
SELECT MIN( `col_decimal(40, 20)_undef_signed` ) AS myagg FROM table_10_utf8_undef WHERE (  ( ( 8 - ( 5 = 0 ) ) + 8 ) ) IS TRUE GROUP BY `col_double_undef_signed`;
SELECT MIN( `col_float_undef_signed` ) AS myagg FROM table_30_utf8_undef WHERE (  2 ) IS FALSE GROUP BY `col_bigint_undef_signed`;
SELECT SUM( `col_float_undef_unsigned` ) AS myagg FROM table_30_binary_undef WHERE (  TAN ( `col_decimal(40, 20)_key_signed` ) ) IS FALSE GROUP BY `col_float_key_signed`;
SELECT MAX( `col_float_undef_unsigned` ) AS myagg FROM table_30_latin1_4 GROUP BY `col_bigint_key_signed`;
SELECT SUM( `col_bigint_key_signed` ) AS myagg FROM table_10_latin1_6 WHERE ( NOT `col_bigint_key_signed` ) IS TRUE GROUP BY `col_float_key_unsigned`;
SELECT MIN( `col_bigint_undef_unsigned` ) AS myagg FROM table_10_binary_4 GROUP BY `col_float_undef_unsigned`;
SELECT COUNT( `col_char(20)_undef_signed` ) AS myagg FROM table_30_utf8_6 GROUP BY `col_char(20)_key_signed`;
SELECT MIN( `col_float_undef_unsigned` ) AS myagg FROM table_30_binary_undef GROUP BY `col_float_undef_unsigned`;
SELECT MAX( `col_varchar(20)_key_signed` ) AS myagg FROM table_30_binary_4 GROUP BY `col_bigint_key_unsigned`;
SELECT SUM( `col_decimal(40, 20)_key_signed` ) AS myagg FROM table_30_binary_undef WHERE (  `col_decimal(40, 20)_undef_signed` ) IS NULL GROUP BY `col_char(20)_key_signed`;
SELECT MIN( `col_decimal(40, 20)_undef_signed` ) AS myagg FROM table_30_latin1_4 WHERE (  0 ) IS NULL GROUP BY `col_bigint_key_unsigned`;
SELECT COUNT( `col_bigint_key_signed` ) AS myagg FROM table_10_binary_6 GROUP BY `col_char(20)_undef_signed`;
SELECT MIN( `col_double_undef_signed` ) AS myagg FROM table_30_latin1_undef WHERE (  EXPORT_SET( 4, 1, 7, `col_decimal(40, 20)_undef_unsigned` ) ) IS NULL GROUP BY `col_double_undef_unsigned`;
SELECT COUNT( `col_varchar(20)_undef_signed` ) AS myagg FROM table_10_latin1_undef WHERE ( NOT `col_decimal(40, 20)_key_unsigned` ) IS TRUE GROUP BY `col_decimal(40, 20)_undef_signed`;
SELECT SUM( `col_float_key_signed` ) AS myagg FROM table_30_utf8_4 GROUP BY `col_double_undef_unsigned`;
SELECT MAX( `col_decimal(40, 20)_key_unsigned` ) AS myagg FROM table_10_utf8_4 WHERE ( NOT CASE 7 WHEN `col_float_undef_signed` THEN `col_decimal(40, 20)_undef_signed` ELSE `col_float_key_unsigned` END ) IS TRUE GROUP BY `col_bigint_key_unsigned`;
SELECT MIN( `col_float_undef_signed` ) AS myagg FROM table_10_latin1_4 WHERE ( NOT COALESCE( `col_bigint_undef_signed`, `col_char(20)_key_signed`, 2 ) ) IS FALSE GROUP BY `col_float_undef_signed`;
SELECT SUM( `col_bigint_undef_unsigned` ) AS myagg FROM table_30_utf8_6 GROUP BY `col_bigint_undef_signed`;
SELECT MIN( `col_bigint_key_unsigned` ) AS myagg FROM table_10_latin1_4 WHERE (  3 ) IS NULL GROUP BY `col_double_key_unsigned`;
SELECT SUM( `col_char(20)_key_signed` ) AS myagg FROM table_10_utf8_4 WHERE (  7 ) IS FALSE GROUP BY `col_double_undef_signed`;
SELECT COUNT( `col_double_key_unsigned` ) AS myagg FROM table_30_utf8_undef WHERE ( NOT `col_float_key_unsigned` ) IS TRUE GROUP BY `col_bigint_key_signed`;
SELECT COUNT( `col_bigint_key_unsigned` ) AS myagg FROM table_30_utf8_undef GROUP BY `col_char(20)_key_signed`;
SELECT MIN( `col_char(20)_key_signed` ) AS myagg FROM table_10_latin1_4 GROUP BY `col_double_key_unsigned`;
SELECT MAX( `col_char(20)_undef_signed` ) AS myagg FROM table_10_binary_6 WHERE ( NOT COS ( 8 ) ) IS TRUE GROUP BY `col_float_undef_signed`;
SELECT MIN( `col_double_key_signed` ) AS myagg FROM table_30_utf8_6 WHERE ( NOT 8 ) IS NULL GROUP BY `col_varchar(20)_undef_signed`;
SELECT MAX( `col_float_undef_unsigned` ) AS myagg FROM table_10_latin1_4 WHERE ( NOT HEX( 8 ) ) IS TRUE GROUP BY `col_float_key_unsigned`;
SELECT SUM( `col_decimal(40, 20)_key_unsigned` ) AS myagg FROM table_30_utf8_4 WHERE (  `col_double_undef_unsigned` ) IS NULL GROUP BY `col_bigint_key_unsigned`;
SELECT SUM( `col_float_undef_unsigned` ) AS myagg FROM table_30_latin1_undef GROUP BY `col_varchar(20)_key_signed`;
SELECT MAX( `col_varchar(20)_key_signed` ) AS myagg FROM table_30_binary_6 WHERE (  CASE `col_varchar(20)_key_signed` WHEN 9 THEN `col_bigint_undef_signed` ELSE `col_char(20)_key_signed` END ) IS NULL GROUP BY `col_bigint_key_signed`;
SELECT MAX( `col_float_undef_unsigned` ) AS myagg FROM table_10_binary_6 WHERE (  `col_float_key_unsigned` ) IS FALSE GROUP BY `col_varchar(20)_undef_signed`;
SELECT COUNT( `col_float_key_signed` ) AS myagg FROM table_30_latin1_6 WHERE ( NOT 0 ) IS TRUE GROUP BY `col_double_undef_unsigned`;
SELECT MIN( `col_varchar(20)_undef_signed` ) AS myagg FROM table_30_binary_4 GROUP BY `col_decimal(40, 20)_key_signed`;
SELECT MIN( `col_decimal(40, 20)_key_signed` ) AS myagg FROM table_10_latin1_undef WHERE (  `col_decimal(40, 20)_undef_unsigned` ) IS TRUE GROUP BY `col_decimal(40, 20)_undef_unsigned`;
SELECT MIN( `col_float_undef_unsigned` ) AS myagg FROM table_30_utf8_undef WHERE (  ( `col_float_undef_signed` >= ( 1 DIV ( 3 >= 7 ) ) ) ) IS FALSE GROUP BY `col_decimal(40, 20)_undef_signed`;
SELECT COUNT( `col_float_key_unsigned` ) AS myagg FROM table_30_binary_undef GROUP BY `col_bigint_undef_unsigned`;
SELECT MAX( `col_decimal(40, 20)_key_signed` ) AS myagg FROM table_30_utf8_4 GROUP BY `col_decimal(40, 20)_undef_unsigned`;
SELECT SUM( `col_double_undef_signed` ) AS myagg FROM table_10_binary_4 GROUP BY `col_decimal(40, 20)_undef_unsigned`;
SELECT COUNT( `col_char(20)_key_signed` ) AS myagg FROM table_10_binary_4 WHERE ( NOT TAN ( 7 ) ) IS NULL GROUP BY `col_varchar(20)_key_signed`;
SELECT SUM( `col_decimal(40, 20)_undef_unsigned` ) AS myagg FROM table_30_utf8_4 GROUP BY `col_float_undef_signed`;
SELECT COUNT( `col_varchar(20)_key_signed` ) AS myagg FROM table_30_utf8_6 GROUP BY `col_double_undef_unsigned`;
SELECT SUM( `col_bigint_key_signed` ) AS myagg FROM table_10_utf8_undef WHERE ( NOT LOCATE( `col_decimal(40, 20)_undef_unsigned`, 2 ) ) IS NULL GROUP BY `col_decimal(40, 20)_key_unsigned`;
SELECT SUM( `col_decimal(40, 20)_undef_unsigned` ) AS myagg FROM table_10_utf8_6 GROUP BY `col_bigint_undef_signed`;
SELECT MIN( `col_bigint_key_unsigned` ) AS myagg FROM table_30_utf8_undef GROUP BY `col_bigint_key_unsigned`;
SELECT COUNT( `col_decimal(40, 20)_undef_unsigned` ) AS myagg FROM table_10_binary_6 WHERE ( NOT ISNULL( ( `col_double_key_signed` DIV 3 ) ) ) IS TRUE GROUP BY `col_bigint_undef_unsigned`;
SELECT SUM( `col_double_key_signed` ) AS myagg FROM table_10_latin1_6 WHERE (  COALESCE( `col_float_key_unsigned`, ( 4 + `col_double_key_unsigned` ), 6, 8, ( `col_char(20)_key_signed` >= `col_bigint_undef_signed` ) ) ) IS FALSE GROUP BY `col_decimal(40, 20)_key_signed`;
SELECT MIN( `col_decimal(40, 20)_undef_signed` ) AS myagg FROM table_30_utf8_4 GROUP BY `col_float_key_signed`;
SELECT MAX( `col_float_undef_signed` ) AS myagg FROM table_30_latin1_4 GROUP BY `col_bigint_key_signed`;
SELECT COUNT( `col_double_key_unsigned` ) AS myagg FROM table_10_binary_4 WHERE (  QUOTE( `col_bigint_undef_unsigned` ) ) IS NULL GROUP BY `col_bigint_key_signed`;
SELECT MIN( `col_decimal(40, 20)_undef_signed` ) AS myagg FROM table_30_latin1_undef GROUP BY `col_bigint_undef_signed`;
SELECT MIN( `col_decimal(40, 20)_undef_signed` ) AS myagg FROM table_30_binary_6 WHERE (  3 ) IS FALSE GROUP BY `col_bigint_undef_signed`;
SELECT SUM( `col_bigint_key_unsigned` ) AS myagg FROM table_10_utf8_4 GROUP BY `col_bigint_key_unsigned`;
SELECT COUNT( `col_varchar(20)_key_signed` ) AS myagg FROM table_10_binary_undef WHERE (  `col_float_key_signed` ) IS FALSE GROUP BY `col_decimal(40, 20)_undef_signed`;
SELECT COUNT( `col_bigint_undef_unsigned` ) AS myagg FROM table_10_utf8_6 GROUP BY `col_double_undef_unsigned`;
SELECT COUNT( `col_float_key_unsigned` ) AS myagg FROM table_30_binary_6 WHERE ( NOT ( `col_double_undef_unsigned` < 8 ) ) IS FALSE GROUP BY `col_bigint_undef_unsigned`;
SELECT SUM( `col_varchar(20)_undef_signed` ) AS myagg FROM table_30_binary_6 GROUP BY `col_bigint_key_unsigned`;
SELECT SUM( `col_double_key_unsigned` ) AS myagg FROM table_10_binary_undef GROUP BY `col_float_key_unsigned`;
SELECT MAX( `col_double_key_unsigned` ) AS myagg FROM table_30_binary_4 WHERE ( NOT SIN ( 7 ) ) IS FALSE GROUP BY `col_double_undef_unsigned`;
SELECT MIN( `col_bigint_undef_unsigned` ) AS myagg FROM table_10_latin1_undef GROUP BY `col_double_undef_unsigned`;
SELECT MAX( `col_char(20)_undef_signed` ) AS myagg FROM table_30_utf8_4 GROUP BY `col_float_key_signed`;
SELECT MIN( `col_double_key_signed` ) AS myagg FROM table_10_utf8_6 GROUP BY `col_char(20)_key_signed`;
SELECT MAX( `col_char(20)_key_signed` ) AS myagg FROM table_30_binary_undef WHERE (  CAST( 1 AS DATE ) ) IS FALSE GROUP BY `col_bigint_undef_signed`;
SELECT MIN( `col_decimal(40, 20)_key_unsigned` ) AS myagg FROM table_30_utf8_6 WHERE ( NOT SIN ( `col_double_undef_unsigned` ) ) IS FALSE GROUP BY `col_decimal(40, 20)_key_unsigned`;
SELECT MAX( `col_double_key_signed` ) AS myagg FROM table_30_binary_6 WHERE (  ( `col_float_key_unsigned` < ( `col_float_key_unsigned` / `col_varchar(20)_key_signed` ) ) ) IS NULL GROUP BY `col_double_undef_unsigned`;
SELECT MIN( `col_float_key_unsigned` ) AS myagg FROM table_10_utf8_4 WHERE (  COS ( ( ( `col_bigint_undef_unsigned` / `col_decimal(40, 20)_undef_signed` ) - ( 5 <= ( ( 9 * `col_double_undef_signed` ) = ( `col_double_key_signed` >= `col_double_undef_signed` ) ) ) ) ) ) IS FALSE GROUP BY `col_double_undef_signed`;
SELECT COUNT( `col_bigint_undef_signed` ) AS myagg FROM table_30_utf8_6 GROUP BY `col_bigint_undef_signed`;
SELECT SUM( `col_float_key_signed` ) AS myagg FROM table_10_utf8_4 WHERE (  BINARY `col_double_key_signed` ) IS FALSE GROUP BY `col_decimal(40, 20)_undef_signed`;
SELECT SUM( `col_bigint_undef_unsigned` ) AS myagg FROM table_30_latin1_4 GROUP BY `col_bigint_undef_signed`;
SELECT COUNT( `col_varchar(20)_undef_signed` ) AS myagg FROM table_10_latin1_6 WHERE ( NOT CASE 9 WHEN ( `col_float_key_signed` != `col_char(20)_undef_signed` ) THEN `col_bigint_key_unsigned` END ) IS NULL GROUP BY `col_bigint_key_signed`;
SELECT MIN( `col_decimal(40, 20)_key_unsigned` ) AS myagg FROM table_30_utf8_6 WHERE ( NOT LCASE( `col_decimal(40, 20)_key_unsigned` ) ) IS FALSE GROUP BY `col_float_key_signed`;
SELECT MAX( `col_decimal(40, 20)_undef_unsigned` ) AS myagg FROM table_30_utf8_6 WHERE (  7 ) IS NULL GROUP BY `col_bigint_key_unsigned`;
SELECT MIN( `col_float_key_unsigned` ) AS myagg FROM table_30_utf8_6 WHERE ( NOT GREATEST( `col_decimal(40, 20)_undef_signed`, 7, `col_bigint_key_signed`, ( `col_double_undef_signed` % 8 ), ( `col_double_undef_unsigned` >= ( ( `col_bigint_undef_unsigned` < `col_bigint_key_unsigned` ) - ( 2 % `col_float_undef_unsigned` ) ) ) ) ) IS TRUE GROUP BY `col_double_key_signed`;
SELECT MIN( `col_float_undef_signed` ) AS myagg FROM table_30_latin1_6 WHERE ( NOT ( `col_decimal(40, 20)_key_unsigned` % 7 ) ) IS NULL GROUP BY `col_bigint_key_signed`;
SELECT COUNT( `col_float_key_unsigned` ) AS myagg FROM table_30_binary_6 WHERE ( NOT ( `col_float_key_unsigned` <= ( ( ( 3 % `col_float_undef_unsigned` ) + `col_float_key_signed` ) > ( ( `col_bigint_key_signed` != `col_bigint_key_unsigned` ) != ( `col_bigint_undef_signed` < `col_double_undef_signed` ) ) ) ) ) IS NULL GROUP BY `col_char(20)_undef_signed`;
SELECT MIN( `col_decimal(40, 20)_undef_signed` ) AS myagg FROM table_10_binary_undef WHERE ( NOT ( 4 + `col_bigint_key_unsigned` ) ) IS FALSE GROUP BY `col_double_key_signed`;
SELECT MAX( `col_double_key_unsigned` ) AS myagg FROM table_30_latin1_6 WHERE (  2 ) IS TRUE GROUP BY `col_varchar(20)_key_signed`;
SELECT MAX( `col_decimal(40, 20)_undef_signed` ) AS myagg FROM table_30_latin1_4 GROUP BY `col_decimal(40, 20)_key_unsigned`;
SELECT MIN( `col_bigint_key_unsigned` ) AS myagg FROM table_30_latin1_undef WHERE ( NOT MAKE_SET( ( `col_bigint_key_unsigned` DIV `col_double_key_signed` ), `col_varchar(20)_key_signed`, `col_varchar(20)_undef_signed` ) ) IS TRUE GROUP BY `col_decimal(40, 20)_key_signed`;
SELECT MIN( `col_float_key_signed` ) AS myagg FROM table_30_binary_undef WHERE (  OCT( 6 ) ) IS NULL GROUP BY `col_bigint_key_signed`;
SELECT SUM( `col_decimal(40, 20)_key_signed` ) AS myagg FROM table_30_binary_undef GROUP BY `col_char(20)_key_signed`;
SELECT MIN( `col_double_undef_unsigned` ) AS myagg FROM table_10_latin1_undef GROUP BY `col_bigint_undef_signed`;
SELECT COUNT( `col_decimal(40, 20)_key_unsigned` ) AS myagg FROM table_10_binary_undef WHERE (  `col_varchar(20)_undef_signed` ) IS NULL GROUP BY `col_float_undef_unsigned`;
SELECT SUM( `col_varchar(20)_key_signed` ) AS myagg FROM table_30_latin1_4 WHERE ( NOT ( ( ( `col_decimal(40, 20)_undef_signed` DIV ( `col_double_undef_signed` < `col_char(20)_undef_signed` ) ) < ( 5 - ( 5 = 4 ) ) ) DIV `col_float_key_unsigned` ) ) IS NULL GROUP BY `col_char(20)_undef_signed`;
SELECT MAX( `col_double_key_unsigned` ) AS myagg FROM table_30_utf8_4 GROUP BY `col_float_key_signed`;
SELECT SUM( `col_char(20)_undef_signed` ) AS myagg FROM table_10_latin1_4 WHERE (  ( ( 2 >= `col_float_undef_signed` ) / 0 ) ) IS TRUE GROUP BY `col_float_key_signed`;
SELECT MIN( `col_decimal(40, 20)_key_signed` ) AS myagg FROM table_30_latin1_undef GROUP BY `col_bigint_undef_signed`;
SELECT MAX( `col_bigint_undef_signed` ) AS myagg FROM table_30_binary_4 GROUP BY `col_bigint_undef_signed`;
SELECT MAX( `col_float_key_unsigned` ) AS myagg FROM table_30_utf8_6 WHERE ( NOT ( 1 MOD ( `col_bigint_undef_signed` * 2 ) ) ) IS TRUE GROUP BY `col_char(20)_key_signed`;
SELECT SUM( `col_double_undef_unsigned` ) AS myagg FROM table_10_binary_4 WHERE (  SIN ( `col_double_key_signed` ) ) IS NULL GROUP BY `col_double_undef_signed`;
SELECT COUNT( `col_decimal(40, 20)_key_unsigned` ) AS myagg FROM table_10_binary_6 GROUP BY `col_decimal(40, 20)_undef_unsigned`;
SELECT MIN( `col_bigint_undef_unsigned` ) AS myagg FROM table_30_utf8_undef WHERE ( NOT SIN ( `col_double_undef_unsigned` ) ) IS NULL GROUP BY `col_varchar(20)_undef_signed`;
SELECT COUNT( `col_float_undef_signed` ) AS myagg FROM table_30_latin1_undef GROUP BY `col_decimal(40, 20)_undef_unsigned`;
SELECT COUNT( `col_bigint_key_unsigned` ) AS myagg FROM table_30_utf8_4 WHERE (  ( ( ( 2 % 2 ) - ( 7 < `col_float_undef_unsigned` ) ) MOD `col_float_undef_unsigned` ) ) IS TRUE GROUP BY `col_varchar(20)_undef_signed`;
SELECT MAX( `col_float_undef_unsigned` ) AS myagg FROM table_10_latin1_undef WHERE ( NOT CASE `col_varchar(20)_key_signed` WHEN `col_float_key_signed` THEN `col_bigint_undef_unsigned` END ) IS TRUE GROUP BY `col_bigint_key_unsigned`;
SELECT SUM( `col_bigint_undef_unsigned` ) AS myagg FROM table_30_latin1_4 GROUP BY `col_double_undef_unsigned`;
SELECT MAX( `col_bigint_key_unsigned` ) AS myagg FROM table_30_binary_6 WHERE ( NOT 8 ) IS FALSE GROUP BY `col_double_undef_unsigned`;
SELECT MAX( `col_varchar(20)_key_signed` ) AS myagg FROM table_30_latin1_4 GROUP BY `col_float_undef_signed`;
SELECT MAX( `col_char(20)_key_signed` ) AS myagg FROM table_10_binary_4 GROUP BY `col_decimal(40, 20)_undef_unsigned`;
SELECT SUM( `col_varchar(20)_key_signed` ) AS myagg FROM table_30_utf8_4 GROUP BY `col_decimal(40, 20)_key_signed`;
SELECT MIN( `col_bigint_undef_signed` ) AS myagg FROM table_10_latin1_6 WHERE ( NOT `col_float_key_unsigned` ) IS NULL GROUP BY `col_float_key_signed`;
SELECT MAX( `col_double_undef_unsigned` ) AS myagg FROM table_10_binary_4 GROUP BY `col_bigint_key_unsigned`;