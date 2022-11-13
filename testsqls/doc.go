// Package testsqls: sql benchmark for testing.
// Prepare:
//   sudo docker run -itd --name test -p 13306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql
// We will use database TEST for testing.
// Make sure there is no important data in TEST, as we will automatically clear it.
package testsqls
