// Package testsqls: sql benchmark for testing.
// Prepare:
//   sudo docker run -itd --name mysqltest -p 13306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql:8.0.30
//   sudo docker run -itd --name mariadbtest -p 23306:3306 -e MYSQL_ROOT_PASSWORD=123456 --privileged=true mariadb:10.11.1-rc
//   sudo docker run -itd --name tidbtest -p 4000:4000 pingcap/tidb:v6.4.0
//   + SET PASSWORD = '123456';
//   sudo docker run -itd --name oceanbasetest -p 2881:2881 oceanbase/oceanbase-ce:4.0.0.0
//   + SET PASSWORD = PASSWORD('123456');
// We will use database TEST for testing.
// Make sure there is no important data in TEST, as we will automatically clear it.
package testsqls
