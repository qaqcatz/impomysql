package mayaffect

// MayAffect:
// Verify whether the bugs detected by tasks can be reproduced on the specified version of DBMS.
//
// dbmsOutputPath: the OutputPath of your tasks + '/' + the DBMS of your tasks, for example, ./output/mysql
//
// version: the specified version of DBMS, needs to be a unique string, it is recommended to use tag or commit id.
//
// portOrDSN: You need to deploy the specified version of DBMS to portOrDSN in advance.
// We will automatically judge whether you provide a port or a dsn.
//   - if you provide a port, we will create a connector with dsn "root:123456@tcp(127.0.0.1:port)/TEST",
//	 - if you provide a dsn, we will create a connector with your dsn.
//
// Before introducing whereVersionEQ, you need to know how MayAffect works:
//
// We will create a sqlite database `mayaffect.db` in dbmsOutputPath with a table:
//   CREATE TABLE `mayaffect` (`bugJsonPath` TEXT, `version` TEXT);
// If `mayaffect.db` does not exist, we will create database `mayaffect.db` and table `mayaffect`,
// then recursively traverse each bug in dbmsOutputPath, and update table `mayaffect`:
//   INSERT INTO `mayaffect` VALUES (bugJsonPath, "");
//
// When you use MayAffect, we will:
//   SELECT `bugJsonPath` FROM `mayaffect` WHERE `version` = whereVersionEQ
// Obviously, If whereVersionEQ="", you will get all bugJsonPaths.
//
// For each bugJsonPath, we will verify whether the bug can be reproduced on the specified version of DBMS.
//
// If it can be reproduced, we will first:
//   SELECT * FROM `mayaffect` WHERE `bugJsonPath`=bugJsonPath AND `version`=version
// If it returns an empty result, then we:
//   INSERT INTO `mayaffect` VALUES (bugJsonPath, version);
// This is done to ensure that each row in is unique.
//
// Now you understand how MayAffect works, you can query the table `mayaffect` to get the information you want.
func MayAffect(dbmsOutputPath string, version string, portOrDSN string, whereVersionEQ string) error {
	return nil
}