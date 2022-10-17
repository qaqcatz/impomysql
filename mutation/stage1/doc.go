// Package stage1: The implication oracle cannot handle these features, remove them.
//
// 1. Aggregate functions: For a node *ast.AggregateFuncExpr, we:
//   - set node.F = "";
//   - clear node.Args, add *test_driver.ValueExpr(value 1) to node.Args;
//   - set node.Distinct to false;
//   - set node.Order to nil.
//   In particular, we need to set *ast.SelectStmt.GroupBy to nil
//   to avoid the semantic error caused by removing the aggregate functions.
//   example:
//   ----------input----------
//   SELECT * FROM (
//     SELECT SUM(ID+1) AS S, GROUP_CONCAT(NAME ORDER BY NAME DESC), CITY
//     FROM COMPANY
//     GROUP BY CITY
//     HAVING COUNT(DISTINCT AGE) >= 1
//   ) AS T
//   WHERE T.S > 0;
//   ----------output----------
//   SELECT * FROM (
//     SELECT (1) AS S, (1), CITY
//     FROM COMPANY
//     HAVING (1) >= 1
//   ) AS T
//   WHERE T.S > 0;
// 2. Window functions: *ast.WindowFuncExpr can only appear in *ast.SelectField.
// In particular *ast.WindowSpec can also appear in *ast.SelectStmt. Therefore, we:
//   - Iterate each *ast.FieldList, replace each *ast.WindowFuncExpr with *test_driver.ValueExpr(value 1);
//   - set *ast.SelectStmt.WindowSpecs to nil(not empty, nil!).
//   example:
//   ----------input----------
//   SELECT ID AS id, CITY, AGE,
//   SUM(AGE) OVER w AS sum_age,
//   AVG(AGE) OVER (PARTITION BY CITY ORDER BY ID ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) AS avg_age,
//   ROW_NUMBER() OVER (PARTITION BY CITY ORDER BY ID) AS rn
//   FROM COMPANY
//   WINDOW w AS (PARTITION BY CITY ORDER BY ID ROWS UNBOUNDED PRECEDING)
//   ----------output----------
//   SELECT ID AS id, CITY, AGE,
//   1 AS sum_age,
//   1 AS avg_age,
//   1 AS rn
//   FROM COMPANY
// 3. {LEFT|RIGHT} [OUTER] JOIN -> JOIN
// 4. Remove Limit
// 5. The transformed sql may fail to execute. It is recommended to execute
// the transformed sql to do some verification.
//
// 6. Only Support SELECT statement.
package stage1
