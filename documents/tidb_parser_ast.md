# tidb_parser_ast

The structure of each node of tidb_paser's ast(**tidb v5.4.2 **), as well as some sql examples.

We only consider the nodes related to query (irrelevant nodes are marked as `skip`).

**how we got this structure?**

>  right click ast.Node in GoLand
>
> ->Go To
>
> ->Implementation(s)
>
> Implementations of Node  (175 usages found)

```golang
// Node is the basic element of the AST.
// Interfaces embed Node should have 'Node' name suffix.
type Node interface {
	// Restore returns the sql text from ast tree
	Restore(ctx *format.RestoreCtx) error
	// Accept accepts Visitor to visit itself.
	// The returned node should replace original node.
	// ok returns false to stop visiting.
	//
	// Implementation of this method should first call visitor.Enter,
	// assign the returned node to its method receiver, if skipChildren returns true,
	// children should be skipped. Otherwise, call its children in particular order that
	// later elements depends on former elements. Finally, return visitor.Leave.
	Accept(v Visitor) (node Node, ok bool)
	// Text returns the original text of the element.
	Text() string
	// SetText sets original text to the Node.
	SetText(text string)
	// SetOriginTextPosition set the start offset of this node in the origin text.
	SetOriginTextPosition(offset int)
	// OriginTextPosition get the start offset of this node in the origin text.
	OriginTextPosition() int
}
```

# executor  (1)

## change.go  (1)

### **skip**(1)

ChangeExec

```golang
// ChangeExec represents a change executor.
type ChangeExec struct {
	baseExecutor
	*ast.ChangeStmt
}
```

# parser/ast  (169)

## advisor.go  (1)

### **skip(1)**

IndexAdviseStmt

```golang
// IndexAdviseStmt is used to advise indexes
type IndexAdviseStmt struct {
	stmtNode

	IsLocal     bool
	Path        string
	MaxMinutes  uint64
	MaxIndexNum *MaxIndexNumClause
	LinesInfo   *LinesClause
}
```

## ast.go  (7)

### ExprNode

```golang
// ExprNode is a node that can be evaluated.
// Name of implementations should have 'Expr' suffix.
type ExprNode interface {
	// Node is embedded in ExprNode.
	Node
	// SetType sets evaluation type to the expression.
	SetType(tp *types.FieldType)
	// GetType gets the evaluation type of the expression.
	GetType() *types.FieldType
	// SetFlag sets flag to the expression.
	// Flag indicates whether the expression contains
	// parameter marker, reference, aggregate function...
	SetFlag(flag uint64)
	// GetFlag returns the flag of the expression.
	GetFlag() uint64

	// Format formats the AST into a writer.
	Format(w io.Writer)
}
```

> Go To Implementation(s):
>
> 38
>
> type FuncNode interface {
> type ValueExpr interface {
> type BetweenExpr struct {
> type BinaryOperationExpr struct {
> type CaseExpr struct {
> type SubqueryExpr struct {
> type CompareSubqueryExpr struct {
> type TableNameExpr struct {
> type ColumnNameExpr struct {
> type DefaultExpr struct {
> type ExistsSubqueryExpr struct {
> type PatternInExpr struct {
> type IsNullExpr struct {
> type IsTruthExpr struct {
> type PatternLikeExpr struct {
> type ParamMarkerExpr interface {
> type ParenthesesExpr struct {
> type PositionExpr struct {
> type PatternRegexpExpr struct {
> type RowExpr struct {
> type UnaryOperationExpr struct {
> type ValuesExpr struct {
> type VariableExpr struct {
> type MaxValueExpr struct {
> type MatchAgainst struct {
> type SetCollationExpr struct {
> type checkExpr struct {
> type FuncCallExpr struct {
> type FuncCastExpr struct {
> type TrimDirectionExpr struct {
> type AggregateFuncExpr struct {
> type WindowFuncExpr struct {
> type TimeUnitExpr struct {
> type GetFormatSelectorExpr struct {
> type ValueExpr struct {
> type ParamMarkerExpr struct {
> type ValueExpr struct {
> type ParamMarkerExpr struct {

### FuncNode

```golang
// FuncNode represents function call expression node.
type FuncNode interface {
	ExprNode
	functionExpression()
}
```

> Go To Implementation(s):
>
> 4
> type FuncCallExpr struct {
> type FuncCastExpr struct {
> type AggregateFuncExpr struct {
> type WindowFuncExpr struct {

### StmtNode

```golang
// StmtNode represents statement node.
// Name of implementations should have 'Stmt' suffix.
type StmtNode interface {
	Node
	statement()
}
```

> Go To Implementation(s):
>
> 86
> type IndexAdviseStmt struct {
> type DDLNode interface {
> type DMLNode interface {
> type SensitiveStmtNode interface {
> type CreateDatabaseStmt struct {
> type AlterDatabaseStmt struct {
> type DropDatabaseStmt struct {
> type CreateTableStmt struct {
> type DropTableStmt struct {
> type DropPlacementPolicyStmt struct {
> type DropSequenceStmt struct {
> type RenameTableStmt struct {
> type CreateViewStmt struct {
> type CreatePlacementPolicyStmt struct {
> type CreateSequenceStmt struct {
> type CreateIndexStmt struct {
> type DropIndexStmt struct {
> type LockTablesStmt struct {
> type UnlockTablesStmt struct {
> type CleanupTableLockStmt struct {
> type RepairTableStmt struct {
> type AlterTableStmt struct {
> type TruncateTableStmt struct {
> type RecoverTableStmt struct {
> type FlashBackTableStmt struct {
> type AlterPlacementPolicyStmt struct {
> type AlterSequenceStmt struct {
> type SelectStmt struct {
> type SetOprStmt struct {
> type LoadDataStmt struct {
> type CallStmt struct {
> type InsertStmt struct {
> type DeleteStmt struct {
> type UpdateStmt struct {
> type ShowStmt struct {
> type SplitRegionStmt struct {
> type TraceStmt struct {
> type ExplainForStmt struct {
> type ExplainStmt struct {
> type PlanReplayerStmt struct {
> type PrepareStmt struct {
> type DeallocateStmt struct {
> type ExecuteStmt struct {
> type BeginStmt struct {
> type BinlogStmt struct {
> type CommitStmt struct {
> type RollbackStmt struct {
> type UseStmt struct {
> type FlushStmt struct {
> type KillStmt struct {
> type SetStmt struct {
> type SetConfigStmt struct {
> type SetPwdStmt struct {
> type ChangeStmt struct {
> type SetRoleStmt struct {
> type SetDefaultRoleStmt struct {
> type CreateUserStmt struct {
> type AlterUserStmt struct {
> type AlterInstanceStmt struct {
> type DropUserStmt struct {
> type CreateBindingStmt struct {
> type DropBindingStmt struct {
> type CreateStatisticsStmt struct {
> type DropStatisticsStmt struct {
> type DoStmt struct {
> type AdminStmt struct {
> type RevokeStmt struct {
> type RevokeRoleStmt struct {
> type GrantStmt struct {
> type GrantProxyStmt struct {
> type GrantRoleStmt struct {
> type ShutdownStmt struct {
> type RestartStmt struct {
> type HelpStmt struct {
> type RenameUserStmt struct {
> type BRIEStmt struct {
> type PurgeImportStmt struct {
> type CreateImportStmt struct {
> type StopImportStmt struct {
> type ResumeImportStmt struct {
> type AlterImportStmt struct {
> type DropImportStmt struct {
> type ShowImportStmt struct {
> type AnalyzeTableStmt struct {
> type DropStatsStmt struct {
> type LoadStatsStmt struct 

### DDLNode

```golang
// DDLNode represents DDL statement node.
type DDLNode interface {
	StmtNode
	ddlStatement()
}
```

>Go To Implementation(s):
>
>23
>
>type CreateDatabaseStmt struct {
>type AlterDatabaseStmt struct {
>type DropDatabaseStmt struct {
>type CreateTableStmt struct {
>type DropTableStmt struct {
>type DropPlacementPolicyStmt struct {
>type DropSequenceStmt struct {
>type RenameTableStmt struct {
>type CreateViewStmt struct {
>type CreatePlacementPolicyStmt struct {
>type CreateSequenceStmt struct {
>type CreateIndexStmt struct {
>type DropIndexStmt struct {
>type LockTablesStmt struct {
>type UnlockTablesStmt struct {
>type CleanupTableLockStmt struct {
>type RepairTableStmt struct {
>type AlterTableStmt struct {
>type TruncateTableStmt struct {
>type RecoverTableStmt struct {
>type FlashBackTableStmt struct {
>type AlterPlacementPolicyStmt struct {
>type AlterSequenceStmt struct {

### DMLNode

```golang
// DMLNode represents DML statement node.
type DMLNode interface {
	StmtNode
	dmlStatement()
}
```

> Go To Implementation(s):
>
> 9
>
> type SelectStmt struct {
> type SetOprStmt struct {
> type LoadDataStmt struct {
> type CallStmt struct {
> type InsertStmt struct {
> type DeleteStmt struct {
> type UpdateStmt struct {
> type ShowStmt struct {
> type SplitRegionStmt struct {

### ResultSetNode

```golang
// ResultSetNode interface has a ResultFields property, represents a Node that returns result set.
// Implementations include SelectStmt, SubqueryExpr, TableSource, TableName, Join and SetOprStmt.
type ResultSetNode interface {
	Node

	resultSet()
}
```

> Go To Implementation(s):
>
> 6
> type Join struct {
> type TableName struct {
> type TableSource struct {
> type SelectStmt struct {
> type SetOprStmt struct {
> type SubqueryExpr struct {

### **SensitiveStmtNode**

```golang
// SensitiveStmtNode overloads StmtNode and provides a SecureText method.
type SensitiveStmtNode interface {
	StmtNode
	// SecureText is different from Text that it hide password information.
	SecureText() string
}
```

>Go To Implementation(s):
>
>8
>
>type SetPwdStmt struct {
>type ChangeStmt struct {
>type CreateUserStmt struct {
>type AlterUserStmt struct {
>type GrantStmt struct {
>type GrantRoleStmt struct {
>type BRIEStmt struct {
>type CreateImportStmt struct {

## ddl.go  (38)

### skip(38)

CreateDatabaseStmt
AlterDatabaseStmt
DropDatabaseStmt
IndexPartSpecification
ReferenceDef
OnDeleteOpt
OnUpdateOpt
ColumnOption
IndexOption
Constraint
ColumnDef
CreateTableStmt
DropTableStmt
DropPlacementPolicyStmt
DropSequenceStmt
RenameTableStmt
TableToTable 
CreateViewStmt
CreatePlacementPolicyStmt
CreateSequenceStmt
IndexLockAndAlgorithm
CreateIndexStmt
DropIndexStmt
LockTablesStmt
UnlockTablesStmt
CleanupTableLockStmt
RepairTableStmt
ColumnPosition
AlterTableSpec
AlterTableStmt
TruncateTableStmt
PartitionOptions
RecoverTableStmt
FlashBackTableStmt
AttributesSpec
StatsOptionsSpec
AlterPlacementPolicyStmt
AlterSequenceStmt

## dml.go  (34)

### Join

```golang
// Join represents table join.
type Join struct {
	node

	// Left table can be TableSource or JoinNode.
	Left ResultSetNode
	// Right table can be TableSource or JoinNode or nil.
	Right ResultSetNode
	// Tp represents join type.
	Tp JoinType
	// On represents join on condition.
	On *OnCondition
	// Using represents join using clause.
	Using []*ColumnName
	// NaturalJoin represents join is natural join.
	NaturalJoin bool
	// StraightJoin represents a straight join.
	StraightJoin   bool
	ExplicitParens bool
}

// join type:
const (
	// CrossJoin is cross join type.
	CrossJoin JoinType = iota + 1
	// LeftJoin is left Join type.
	LeftJoin
	// RightJoin is right Join type.
	RightJoin
)

// 0: none join

// for example:
// SELECT * FROM t1 
// LEFT JOIN 
// (t2 
//  CROSS JOIN 
//  t3 
//  CROSS JOIN 
//  t4) 
//  ON (t2.a = t1.a AND t3.b = t1.b AND t4.c = t1.c)
```

### TableName

```golang
// TableName represents a table name.
type TableName struct {
   node

   Schema model.CIStr
   Name   model.CIStr

   DBInfo    *model.DBInfo
   TableInfo *model.TableInfo

   IndexHints     []*IndexHint
   PartitionNames []model.CIStr
   TableSample    *TableSample
   // AS OF is used to see the data as it was at a specific point in time.
   AsOf *AsOfClause
}
```

### skip(1)

DeleteTableList

```golang
// DeleteTableList is the tablelist used in delete statement multi-table mode.
type DeleteTableList struct {
	node
	Tables []*TableName
}
```

### OnCondition

```golang
// OnCondition represents JOIN on condition.
type OnCondition struct {
	node

	Expr ExprNode
}

// for example: ON (t2.a = t1.a AND t3.b = t1.b AND t4.c = t1.c)
```

### TableSource

```golang
// TableSource represents table source with a name.
type TableSource struct {
	node

	// Source is the source of the data, can be a TableName,
	// a SelectStmt, a SetOprStmt, or a JoinNode.
	Source ResultSetNode

	// AsName is the alias name of the table source.
	AsName model.CIStr
}
```

### WildCardField

```golang
// WildCardField is a special type of select field content.
type WildCardField struct {
   node

   Table  model.CIStr
   Schema model.CIStr
}

// for example: select *(<-WildCardField) from t;
```

###  SelectField

```golang
// SelectField represents fields in select statement.
// There are two type of select field: wildcard
// and expression with optional alias name.
type SelectField struct {
	node

	// Offset is used to get original text.
	Offset int
	// WildCard is not nil, Expr will be nil.
	WildCard *WildCardField
	// Expr is not nil, WildCard will be nil.
	Expr ExprNode
	// AsName is alias name for Expr.
	AsName model.CIStr
	// Auxiliary stands for if this field is auxiliary.
	// When we add a Field into SelectField list which is used for having/orderby clause but the field is not in select clause,
	// we should set its Auxiliary to true. Then the TrimExec will trim the field.
	Auxiliary bool
}

// for example: select a(<-SelectField), b(<-SelectField) from t;
```

### FieldList

```golang
// FieldList represents field list in select statement.
type FieldList struct {
   node

   Fields []*SelectField
}

// for example: select (a, b)(<-FieldList) from t;
```

### TableRefsClause

```golang
// TableRefsClause represents table references clause in dml statement.
type TableRefsClause struct {
	node

	TableRefs *Join
}

// for example: SELECT * FROM (SELECT A FROM T(<-TableRefsClause))(<-TableRefsClause)
// normally: TableRefsClause -> Join -> TableSource
```

### ByItem

```golang
// ByItem represents an item in order by or group by.
type ByItem struct {
	node

	Expr      ExprNode
	Desc      bool
	NullOrder bool
}

// for example: SELECT A FROM T GROUP BY A(<-ByItem) ORDER BY A(<-ByItem);
// normally: GroupByClause->ByItem, OrderByClause->ByItem
```

### GroupByClause

```golang
// GroupByClause represents group by clause.
type GroupByClause struct {
	node
	Items []*ByItem
}

// for example: SELECT A FROM T (GROUP BY A)(<-GroupByClause)
```

### HavingClause

```golang
// HavingClause represents having clause.
type HavingClause struct {
   node
   Expr ExprNode
}

// for example: SELECT A FROM T GROUP BY A (HAVING SUM(A) > 0)(<-HavingClause)
```

### OrderByClause

```golang
// OrderByClause represents order by clause.
type OrderByClause struct {
   node
   Items    []*ByItem
   ForUnion bool
}

// for example: SELECT * FROM t1 WHERE username LIKE 'l%' UNION SELECT * FROM t1 WHERE username LIKE '%m%' ORDER BY score ASC
```

### TableSample

```golang
type TableSample struct {
	node
	SampleMethod     SampleMethodType
	Expr             ExprNode
	SampleClauseUnit SampleClauseUnitType
	RepeatableSeed   ExprNode
}
```

<font color="red">unknown</font>

### WithClause

```golang
type WithClause struct {
   node

   IsRecursive bool
   CTEs        []*CommonTableExpression
}

// for example:
// WITH RECURSIVE cte AS
// (
//   SELECT 1 AS n, CAST('abc' AS CHAR(20)) AS str
//   UNION ALL
//   SELECT n + 1, CONCAT(str, str) FROM cte WHERE n // < 3
// )
```

### SelectStmt

```golang
// SelectStmt represents the select query node.
// See https://dev.mysql.com/doc/refman/5.7/en/select.html
type SelectStmt struct {
   dmlNode

   // SelectStmtOpts wraps around select hints and switches.
   *SelectStmtOpts
   // Distinct represents whether the select has distinct option.
   Distinct bool
   // From is the from clause of the query.
   From *TableRefsClause
   // Where is the where clause in select statement.
   Where ExprNode
   // Fields is the select expression list.
   Fields *FieldList
   // GroupBy is the group by expression list.
   GroupBy *GroupByClause
   // Having is the having condition.
   Having *HavingClause
   // WindowSpecs is the window specification list.
   WindowSpecs []WindowSpec
   // OrderBy is the ordering expression list.
   OrderBy *OrderByClause
   // Limit is the limit clause.
   Limit *Limit
   // LockInfo is the lock type
   LockInfo *SelectLockInfo
   // TableHints represents the table level Optimizer Hint for join type
   TableHints []*TableOptimizerHint
   // IsInBraces indicates whether it's a stmt in brace.
   IsInBraces bool
   // WithBeforeBraces indicates whether stmt's with clause is before the brace.
   // It's used to distinguish (with xxx select xxx) and with xxx (select xxx)
   WithBeforeBraces bool
   // QueryBlockOffset indicates the order of this SelectStmt if counted from left to right in the sql text.
   QueryBlockOffset int
   // SelectIntoOpt is the select-into option.
   SelectIntoOpt *SelectIntoOption
   // AfterSetOperator indicates the SelectStmt after which type of set operator
   AfterSetOperator *SetOprType
   // Kind refer to three kind of statement: SelectStmt, TableStmt and ValuesStmt
   Kind SelectStmtKind
   // Lists is filled only when Kind == SelectStmtKindValues
   Lists []*RowExpr
   With  *WithClause
   // AsViewSchema indicates if this stmt provides the schema for the view. It is only used when creating the view
   AsViewSchema bool
}

// for example:
// (1) SELECT @@optimizer_switch LIKE '%subquery_to_derived=off%';
// (2) SELECT /*+ NO_RANGE_OPTIMIZATION(t3 PRIMARY, f2_idx) */ f1 FROM t3 WHERE f1 > 30 AND f1 < 33;
// (3) SELECT
//          time, subject, val,
//          SUM(val) OVER (PARTITION BY subject ORDER BY time
//                         ROWS UNBOUNDED PRECEDING)
//            AS running_total,
//          AVG(val) OVER (PARTITION BY subject ORDER BY time
//                         ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING)
//            AS running_average
//        FROM observations;
// (4) SELECT
//      val,
//      ROW_NUMBER() OVER w AS 'row_number',
//      RANK()       OVER w AS 'rank',
//      DENSE_RANK() OVER w AS 'dense_rank'
//    FROM numbers
//    WINDOW w AS (ORDER BY val);
// (5) SELECT * INTO @myvar FROM t1;
```

Refer to specific node for details.

### SetOprSelectList

```golang
// SetOprSelectList represents the SelectStmt/TableStmt/ValuesStmt list in a union statement.
type SetOprSelectList struct {
   node

   With             *WithClause
   AfterSetOperator *SetOprType
   Selects          []Node
}

// for example: 
// (
// select x from t1 union 
// (
// select x from t2
// )(<-SetOprSelectList)
// )(<-SetOprSelectList)
```

### SetOprStmt

```golang
// SetOprStmt represents "union/except/intersect statement"
// See https://dev.mysql.com/doc/refman/5.7/en/union.html
// See https://mariadb.com/kb/en/intersect/
// See https://mariadb.com/kb/en/except/
type SetOprStmt struct {
   dmlNode

   IsInBraces bool
   SelectList *SetOprSelectList
   OrderBy    *OrderByClause
   Limit      *Limit
   With       *WithClause
}

// for example: 
// (
// select x from t1 union 
// (select x from t2)
// )(<-SetOprStmt)

// SetOprStmt->SetOprSelectList
```

### Assignment

```golang
// Assignment is the expression for assignment, like a = 1.
type Assignment struct {
   node
   // Column is the column name to be assigned.
   Column *ColumnName
   // Expr is the expression assigning to ColName.
   Expr ExprNode
}

// for example: UPDATE t1 (SET c1 = 2)(<-Assignment) WHERE c1 = @var1:= 1
// note that:
// 'SELECT @var1 := 1, @var2;' is VariableExpr
// 'SET @name = 43;' is VariableAssignment
```

### ColumnNameOrUserVar

```golang
type ColumnNameOrUserVar struct {
   node
   ColumnName *ColumnName
   UserVar    *VariableExpr
}
```

<font color="red">unknown</font>

### skip(5)

LoadDataStmt

```golang
// LoadDataStmt is a statement to load data from a specified file, then insert this rows into an existing table.
// See https://dev.mysql.com/doc/refman/5.7/en/load-data.html
type LoadDataStmt struct {
   dmlNode

   IsLocal           bool
   Path              string
   OnDuplicate       OnDuplicateKeyHandlingType
   Table             *TableName
   Columns           []*ColumnName
   FieldsInfo        *FieldsClause
   LinesInfo         *LinesClause
   IgnoreLines       uint64
   ColumnAssignments []*Assignment

   ColumnsAndUserVars []*ColumnNameOrUserVar
}
```

CallStmt

```golang
// CallStmt represents a call procedure query node.
// See https://dev.mysql.com/doc/refman/5.7/en/call.html
type CallStmt struct {
   dmlNode

   Procedure *FuncCallExpr
}
```

InsertStmt

```golang
// InsertStmt is a statement to insert new rows into an existing table.
// See https://dev.mysql.com/doc/refman/5.7/en/insert.html
type InsertStmt struct {
	dmlNode

	IsReplace   bool
	IgnoreErr   bool
	Table       *TableRefsClause
	Columns     []*ColumnName
	Lists       [][]ExprNode
	Setlist     []*Assignment
	Priority    mysql.PriorityEnum
	OnDuplicate []*Assignment
	Select      ResultSetNode
	// TableHints represents the table level Optimizer Hint for join type.
	TableHints     []*TableOptimizerHint
	PartitionNames []model.CIStr
}
```

DeleteStmt

```golang
// DeleteStmt is a statement to delete rows from table.
// See https://dev.mysql.com/doc/refman/5.7/en/delete.html
type DeleteStmt struct {
   dmlNode

   // TableRefs is used in both single table and multiple table delete statement.
   TableRefs *TableRefsClause
   // Tables is only used in multiple table delete statement.
   Tables       *DeleteTableList
   Where        ExprNode
   Order        *OrderByClause
   Limit        *Limit
   Priority     mysql.PriorityEnum
   IgnoreErr    bool
   Quick        bool
   IsMultiTable bool
   BeforeFrom   bool
   // TableHints represents the table level Optimizer Hint for join type.
   TableHints []*TableOptimizerHint
   With       *WithClause
}
```

UpdateStmt

```golang
// UpdateStmt is a statement to update columns of existing rows in tables with new values.
// See https://dev.mysql.com/doc/refman/5.7/en/update.html
type UpdateStmt struct {
   dmlNode

   TableRefs     *TableRefsClause
   List          []*Assignment
   Where         ExprNode
   Order         *OrderByClause
   Limit         *Limit
   Priority      mysql.PriorityEnum
   IgnoreErr     bool
   MultipleTable bool
   TableHints    []*TableOptimizerHint
   With          *WithClause
}
```

### Limit

```golang
// Limit is the limit clause.
type Limit struct {
   node

   Count  ExprNode
   Offset ExprNode
}

// for example: SELECT NAME FROM T LIMIT 1,2;
```

### skip(1)

ShowStmt

```golang
// ShowStmt is a statement to provide information about databases, tables, columns and so on.
// See https://dev.mysql.com/doc/refman/5.7/en/show.html
type ShowStmt struct {
   dmlNode

   Tp          ShowStmtType // Databases/Tables/Columns/....
   DBName      string
   Table       *TableName  // Used for showing columns.
   Partition   model.CIStr // Used for showing partition.
   Column      *ColumnName // Used for `desc table column`.
   IndexName   model.CIStr
   Flag        int // Some flag parsed from sql, such as FULL.
   Full        bool
   User        *auth.UserIdentity   // Used for show grants/create user.
   Roles       []*auth.RoleIdentity // Used for show grants .. using
   IfNotExists bool                 // Used for `show create database if not exists`
   Extended    bool                 // Used for `show extended columns from ...`

   // GlobalScope is used by `show variables` and `show bindings`
   GlobalScope bool
   Pattern     *PatternLikeExpr
   Where       ExprNode

   ShowProfileTypes []int  // Used for `SHOW PROFILE` syntax
   ShowProfileArgs  *int64 // Used for `SHOW PROFILE` syntax
   ShowProfileLimit *Limit // Used for `SHOW PROFILE` syntax
}
```

### WindowSpec

```golang
// WindowSpec is the specification of a window.
type WindowSpec struct {
	node

	Name model.CIStr
	// Ref is the reference window of this specification. For example, in `w2 as (w1 order by a)`,
	// the definition of `w2` references `w1`.
	Ref model.CIStr

	PartitionBy *PartitionByClause
	OrderBy     *OrderByClause
	Frame       *FrameClause

	// OnlyAlias will set to true of the first following case.
	// To make compatible with MySQL, we need to distinguish `select func over w` from `select func over (w)`.
	OnlyAlias bool
}

// for example: 
// SELECT
//   val,
//   ROW_NUMBER() OVER (w AS 'row_number')(<-WindowSpec),
//   RANK()       OVER (w AS 'rank')(<-WindowSpec),
//   DENSE_RANK() OVER (w AS 'dense_rank')(<-WindowSpec)
// FROM numbers
// WINDOW (w AS (ORDER BY val))(<-WindowSpec);
```

### SelectIntoOption

```golang
type SelectIntoOption struct {
   node

   Tp         SelectIntoType
   FileName   string
   FieldsInfo *FieldsClause
   LinesInfo  *LinesClause
}

// for example: SELECT * INTO @myvar FROM t1;
// seems that tidb v5.4.2 doest not support it well.
```

PartitionByClause

```golang
// PartitionByClause represents partition by clause.
type PartitionByClause struct {
   node

   Items []*ByItem
}

// for example: SELECT SUM(X) OVER(PARTITION BY country)(<-PartitionByClause) AS country_profit FROM T;
```

### FrameClause

```golang
// FrameClause represents frame clause.
type FrameClause struct {
   node

   Type   FrameType
   Extent FrameExtent
}

// FrameType is the type of window function frame.
type FrameType int

// Window function frame types.
// MySQL only supports `ROWS` and `RANGES`.
const (
	Rows = iota
	Ranges
	Groups
)

// FrameExtent represents frame extent.
type FrameExtent struct {
	Start FrameBound
	End   FrameBound
}

// Restore implements Node interface.
func (n *FrameClause) Restore(ctx *format.RestoreCtx) error {
	switch n.Type {
	case Rows:
		ctx.WriteKeyWord("ROWS")
	case Ranges:
		ctx.WriteKeyWord("RANGE")
	default:
		return errors.New("Unsupported window function frame type")
	}
	ctx.WriteKeyWord(" BETWEEN ")
	if err := n.Extent.Start.Restore(ctx); err != nil {
		return errors.Annotate(err, "An error occurred while restore FrameClause.Extent.Start")
	}
	ctx.WriteKeyWord(" AND ")
	if err := n.Extent.End.Restore(ctx); err != nil {
		return errors.Annotate(err, "An error occurred while restore FrameClause.Extent.End")
	}

	return nil
}


// for example: OVER (PARTITION BY subject ORDER BY time (ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING)(<-FrameClause))
```

### FrameBound

```golang
// FrameBound represents frame bound.
type FrameBound struct {
   node

   Type      BoundType
   UnBounded bool
   Expr      ExprNode
   // `Unit` is used to indicate the units in which the `Expr` should be interpreted.
   // For example: '2:30' MINUTE_SECOND.
   Unit TimeUnitType
}

// for example: BETWEEN (1 PRECEDING)(<-FrameBound) AND (1 FOLLOWING)(<-FrameBound)
```

### skip(2)

SplitRegionStmt

```golang
type SplitRegionStmt struct {
   dmlNode

   Table          *TableName
   IndexName      model.CIStr
   PartitionNames []model.CIStr

   SplitSyntaxOpt *SplitSyntaxOption

   SplitOpt *SplitOption
}
```

<font color="red">tidb feature</font>

AsOfClause

```golang
type AsOfClause struct {
	node
	TsExpr ExprNode
}
```

<font color="red">tidb feature</font>

## expressions.go  (27)

### ValueExpr

```golang
// ValueExpr define a interface for ValueExpr.
type ValueExpr interface {
	ExprNode
	SetValue(val interface{})
	GetValue() interface{}
	GetDatumString() string
	GetString() string
	GetProjectionOffset() int
	SetProjectionOffset(offset int)
}
```

> Go To Implementation(s):
>
> 6 
>
> parser/ast  (2 usages found)
> expressions.go  (1 usage found)
>   type ParamMarkerExpr interface {
> expressions_test.go  (1 usage found)
>   type checkExpr struct {
> parser/test_driver  (2 usages found)
>   test_driver.go  (2 usages found)
>     type ValueExpr struct {
>     type ParamMarkerExpr struct {
> types/parser_driver  (2 usages found)
>   value_expr.go  (2 usages found)
>     type ValueExpr struct {
>     type ParamMarkerExpr struct {

### BetweenExpr

```golang
// BetweenExpr is for "between and" or "not between and" expression.
type BetweenExpr struct {
   exprNode
   // Expr is the expression to be checked.
   Expr ExprNode
   // Left is the expression for minimal value in the range.
   Left ExprNode
   // Right is the expression for maximum value in the range.
   Right ExprNode
   // Not is true, the expression is "not between and".
   Not bool
}

// for example: SELECT * FROM T WHERE ID BETWEEN 0 AND 2
```

### BinaryOperationExpr

```golang
// BinaryOperationExpr is for binary operation like `1 + 1`, `1 - 1`, etc.
type BinaryOperationExpr struct {
   exprNode
   // Op is the operator code for BinaryOperation.
   Op opcode.Op
   // L is the left expression in BinaryOperation.
   L ExprNode
   // R is the right expression in BinaryOperation.
   R ExprNode
}
```

### extra--Op

```golang
// Op is opcode type.
type Op int

// 31
// List operators.
const (
	LogicAnd Op = iota + 1
	LeftShift
	RightShift
	LogicOr
	GE
	LE
	EQ
	NE
	LT
	GT
	Plus
	Minus
	And
	Or
	Mod
	Xor
	Div
	Mul
	Not
	Not2
	BitNeg
	IntDiv
	LogicXor
	NullEQ
	In
	Like
	Case
	Regexp
	IsNull
	IsTruth
	IsFalsity
)

var ops = [...]struct {
	name      string
	literal   string
	isKeyword bool
}{
	LogicAnd: {
		name:      "and",
		literal:   "AND",
		isKeyword: true,
	},
	LogicOr: {
		name:      "or",
		literal:   "OR",
		isKeyword: true,
	},
	LogicXor: {
		name:      "xor",
		literal:   "XOR",
		isKeyword: true,
	},
	LeftShift: {
		name:      "leftshift",
		literal:   "<<",
		isKeyword: false,
	},
	RightShift: {
		name:      "rightshift",
		literal:   ">>",
		isKeyword: false,
	},
	GE: {
		name:      "ge",
		literal:   ">=",
		isKeyword: false,
	},
	LE: {
		name:      "le",
		literal:   "<=",
		isKeyword: false,
	},
	EQ: {
		name:      "eq",
		literal:   "=",
		isKeyword: false,
	},
	NE: {
		name:      "ne",
		literal:   "!=", // perhaps should use `<>` here
		isKeyword: false,
	},
	LT: {
		name:      "lt",
		literal:   "<",
		isKeyword: false,
	},
	GT: {
		name:      "gt",
		literal:   ">",
		isKeyword: false,
	},
	Plus: {
		name:      "plus",
		literal:   "+",
		isKeyword: false,
	},
	Minus: {
		name:      "minus",
		literal:   "-",
		isKeyword: false,
	},
	And: {
		name:      "bitand",
		literal:   "&",
		isKeyword: false,
	},
	Or: {
		name:      "bitor",
		literal:   "|",
		isKeyword: false,
	},
	Mod: {
		name:      "mod",
		literal:   "%",
		isKeyword: false,
	},
	Xor: {
		name:      "bitxor",
		literal:   "^",
		isKeyword: false,
	},
	Div: {
		name:      "div",
		literal:   "/",
		isKeyword: false,
	},
	Mul: {
		name:      "mul",
		literal:   "*",
		isKeyword: false,
	},
	Not: {
		name:      "not",
		literal:   "not ",
		isKeyword: true,
	},
	Not2: {
		name:      "!",
		literal:   "!",
		isKeyword: false,
	},
	BitNeg: {
		name:      "bitneg",
		literal:   "~",
		isKeyword: false,
	},
	IntDiv: {
		name:      "intdiv",
		literal:   "DIV",
		isKeyword: true,
	},
	NullEQ: {
		name:      "nulleq",
		literal:   "<=>",
		isKeyword: false,
	},
	In: {
		name:      "in",
		literal:   "IN",
		isKeyword: true,
	},
	Like: {
		name:      "like",
		literal:   "LIKE",
		isKeyword: true,
	},
	Case: {
		name:      "case",
		literal:   "CASE",
		isKeyword: true,
	},
	Regexp: {
		name:      "regexp",
		literal:   "REGEXP",
		isKeyword: true,
	},
	IsNull: {
		name:      "isnull",
		literal:   "IS NULL",
		isKeyword: true,
	},
	IsTruth: {
		name:      "istrue",
		literal:   "IS TRUE",
		isKeyword: true,
	},
	IsFalsity: {
		name:      "isfalse",
		literal:   "IS FALSE",
		isKeyword: true,
	},
}
```

### WhenClause

```golang
// WhenClause is the when clause in Case expression for "when condition then result".
type WhenClause struct {
   node
   // Expr is the condition expression in WhenClause.
   Expr ExprNode
   // Result is the result expression in WhenClause.
   Result ExprNode
}

// for example: 
// SELECT ID, NAME FROM T WHERE (CASE WHEN ID=2 THEN ID>2 ELSE ID>0 END);
// SELECT ID, NAME FROM T ORDER BY (CASE WHEN NAME IS NULL THEN ID ELSE NAME END);
```

### CaseExpr

```golang
// CaseExpr is the case expression.
type CaseExpr struct {
   exprNode
   // Value is the compare value expression.
   Value ExprNode
   // WhenClauses is the condition check expression.
   WhenClauses []*WhenClause
   // ElseClause is the else result expression.
   ElseClause ExprNode
}

// for example: 
// SELECT ID, NAME FROM T WHERE (CASE WHEN ID=2 THEN ID>2 ELSE ID>0 END);
// SELECT ID, NAME FROM T ORDER BY (CASE WHEN NAME IS NULL THEN ID ELSE NAME END);

// CaseExpr -> WhenClause
```

### SubqueryExpr

```golang
// SubqueryExpr represents a subquery.
type SubqueryExpr struct {
	exprNode
	// Query is the query SelectNode.
	Query      ResultSetNode
	Evaluated  bool
	Correlated bool
	MultiRows  bool
	Exists     bool
}

// for example: SELECT * FROM T WHERE ID IN (SELECT ID FROM T WHERE ID > 0);
// SELECT * FROM t_department WHERE EXISTS (SELECT * FROM t_employee WHERE t_employee.dept_id = t_department.did);
// normally: SubqueryExpr->SelectStmt
```

### CompareSubqueryExpr

```golang
// CompareSubqueryExpr is the expression for "expr cmp (select ...)".
// See https://dev.mysql.com/doc/refman/5.7/en/comparisons-using-subqueries.html
// See https://dev.mysql.com/doc/refman/5.7/en/any-in-some-subqueries.html
// See https://dev.mysql.com/doc/refman/5.7/en/all-subqueries.html
type CompareSubqueryExpr struct {
   exprNode
   // L is the left expression
   L ExprNode
   // Op is the comparison opcode.
   Op opcode.Op
   // R is the subquery for right expression, may be rewritten to other type of expression.
   R ExprNode
   // All is true, we should compare all records in subquery.
   All bool
}

// for example: SELECT * FROM T WHERE ID > ANY (SELECT ID FROM T);
```

### TableNameExpr

```golang
// TableNameExpr represents a table-level object name expression, such as sequence/table/view etc.
type TableNameExpr struct {
   exprNode

   // Name is the referenced object name expression.
   Name *TableName
}
// normally: TableSource -> TableName
```

### ColumnName

```golang
// ColumnName represents column name.
type ColumnName struct {
   node
   Schema model.CIStr
   Table  model.CIStr
   Name   model.CIStr
}

// normally: ColumnNameExpr -> ColumnName
```

### ColumnNameExpr

```golang
// ColumnNameExpr represents a column name expression.
type ColumnNameExpr struct {
	exprNode

	// Name is the referenced column name.
	Name *ColumnName

	// Refer is the result field the column name refers to.
	// The value of Refer.Expr is used as the value of the expression.
	Refer *ResultField
}

// normally: ColumnNameExpr -> ColumnName
```

### skip(1)

DefaultExpr

```golang
// DefaultExpr is the default expression using default value for a column.
type DefaultExpr struct {
   exprNode
   // Name is the column name.
   Name *ColumnName
}
```

### ExistsSubqueryExpr

```golang
// ExistsSubqueryExpr is the expression for "exists (select ...)".
// See https://dev.mysql.com/doc/refman/5.7/en/exists-and-not-exists-subqueries.html
type ExistsSubqueryExpr struct {
   exprNode
   // Sel is the subquery, may be rewritten to other type of expression.
   Sel ExprNode
   // Not is true, the expression is "not exists".
   Not bool
}

// example: SELECT * FROM T WHERE EXISTS (SELECT NULL);
```

### PatternInExpr

```golang
// PatternInExpr is the expression for in operator, like "expr in (1, 2, 3)" or "expr in (select c from t)".
type PatternInExpr struct {
   exprNode
   // Expr is the value expression to be compared.
   Expr ExprNode
   // List is the list expression in compare list.
   List []ExprNode
   // Not is true, the expression is "not in".
   Not bool
   // Sel is the subquery, may be rewritten to other type of expression.
   Sel ExprNode
}

// for example: SELECT * FROM T WHERE ID IN (SELECT ID FROM T);
```

### IsNullExpr

```golang
// IsNullExpr is the expression for null check.
type IsNullExpr struct {
	exprNode
	// Expr is the expression to be checked.
	Expr ExprNode
	// Not is true, the expression is "is not null".
	Not bool
}

// for example: SELECT 1 IS NOT NULL;
```

### IsTruthExpr

```golang
// IsTruthExpr is the expression for true/false check.
type IsTruthExpr struct {
	exprNode
	// Expr is the expression to be checked.
	Expr ExprNode
	// Not is true, the expression is "is not true/false".
	Not bool
	// True indicates checking true or false.
	True int64
}

// for example: SELECT 1 IS NOT TRUE;
```

### PatternLikeExpr

```golang
// PatternLikeExpr is the expression for like operator, e.g, expr like "%123%"
type PatternLikeExpr struct {
   exprNode
   // Expr is the expression to be checked.
   Expr ExprNode
   // Pattern is the like expression.
   Pattern ExprNode
   // Not is true, the expression is "not like".
   Not bool

   Escape byte

   PatChars []byte
   PatTypes []byte
}

// for example: SELECT 'AB' LIKE 'A_'
```

### skip(1)

ParamMarkerExpr

https://dev.mysql.com/doc/refman/8.0/en/sql-prepared-statements.html

```golang
// ParamMarkerExpr expression holds a place for another expression.
// Used in parsing prepare statement.
type ParamMarkerExpr interface {
   ValueExpr
   SetOrder(int)
}


// for example: PREPARE stmt1 FROM 'SELECT * FROM T'
```

### ParenthesesExpr

```golang
// ParenthesesExpr is the parentheses expression.
type ParenthesesExpr struct {
   exprNode
   // Expr is the expression in parentheses.
   Expr ExprNode
}

// ()
```

### skip(1)

PositionExpr

```golang
// PositionExpr is the expression for order by and group by position.
// MySQL use position expression started from 1, it looks a little confused inner.
// maybe later we will use 0 at first.
type PositionExpr struct {
   exprNode
   // N is the position, started from 1 now.
   N int
   // P is the parameterized position.
   P ExprNode
   // Refer is the result field the position refers to.
   Refer *ResultField
}
```

<font color="red">unknown</font>

### PatternRegexpExpr

```golang
// PatternRegexpExpr is the pattern expression for pattern match.
type PatternRegexpExpr struct {
   exprNode
   // Expr is the expression to be checked.
   Expr ExprNode
   // Pattern is the expression for pattern.
   Pattern ExprNode
   // Not is true, the expression is "not rlike",
   Not bool

   // Re is the compiled regexp.
   Re *regexp.Regexp
   // Sexpr is the string for Expr expression.
   Sexpr *string
}

// for example: 
// SELECT 'Michael!' REGEXP '.*';
// SELECT 'Michael!' RLIKE '.*';
// Note that: the following statements is FuncCallExpr: 
// SELECT REGEXP_INSTR('dog cat dog', 'dog');
// SELECT REGEXP_LIKE('CamelCase', 'CAMELCASE');
// SELECT REGEXP_REPLACE('a b c', 'b', 'X');
// SELECT REGEXP_SUBSTR('abc def ghi', '[a-z]+');
```

### RowExpr

```golang
// RowExpr is the expression for row constructor.
// See https://dev.mysql.com/doc/refman/5.7/en/row-subqueries.html
type RowExpr struct {
   exprNode

   Values []ExprNode
}

// SELECT * FROM T WHERE (ID, NAME) = ANY (SELECT * FROM T WHERE ID > 1);
// SELECT * FROM T WHERE ROW(ID, NAME) = ANY (SELECT * FROM T WHERE ID > 1);
```

### UnaryOperationExpr

```golang
// UnaryOperationExpr is the expression for unary operator.
type UnaryOperationExpr struct {
   exprNode
   // Op is the operator opcode.
   Op opcode.Op
   // V is the unary expression.
   V ExprNode
}
```

### skip(1)

ValuesExpr

```golang
// ValuesExpr is the expression used in INSERT VALUES.
type ValuesExpr struct {
   exprNode
   // Column is column name.
   Column *ColumnNameExpr
}
```

### VariableExpr

```golang
// VariableExpr is the expression for variable.
type VariableExpr struct {
	exprNode
	// Name is the variable name.
	Name string
	// IsGlobal indicates whether this variable is global.
	IsGlobal bool
	// IsSystem indicates whether this variable is a system variable in current session.
	IsSystem bool
	// ExplicitScope indicates whether this variable scope is set explicitly.
	ExplicitScope bool
	// Value is the variable value.
	Value ExprNode
}

// for example: SELECT @@global.autocommit;
```

### skip(2)

MaxValueExpr

```golang
// MaxValueExpr is the expression for "maxvalue" used in partition.
type MaxValueExpr struct {
	exprNode
}
```

<font color="red">unknown</font>

MatchAgainst

```golang
// MatchAgainst is the expression for matching against fulltext index.
type MatchAgainst struct {
	exprNode
	// ColumnNames are the columns to match.
	ColumnNames []*ColumnName
	// Against
	Against ExprNode
	// Modifier
	Modifier FulltextSearchModifier
}

// for example: SELECT * FROM articles WHERE MATCH (title,body) AGAINST ('MySQL');
```

<font color="red">unknown</font>

### SetCollationExpr

```golang
// SetCollationExpr is the expression for the `COLLATE collation_name` clause.
type SetCollationExpr struct {
	exprNode
	// Expr is the expression to be set.
	Expr ExprNode
	// Collate is the name of collation to set.
	Collate string
}

// SELECT DISTINCT field1 COLLATE utf8mb4_general_ci FROM table1;
```

## expressions_test.go  (1)

### skip(1)

checkExpr

## functions.go  (7)

### FuncCallExpr

```golang
// FuncCallExpr is for function expression.
type FuncCallExpr struct {
	funcNode
	Tp     FuncCallExprType
	Schema model.CIStr
	// FnName is the function name.
	FnName model.CIStr
	// Args is the function args.
	Args []ExprNode
}
```

### extra--FunctionType

```golang
// 287
// List scalar function names.
const (
	LogicAnd           = "and"
	Cast               = "cast"
	LeftShift          = "leftshift"
	RightShift         = "rightshift"
	LogicOr            = "or"
	GE                 = "ge"
	LE                 = "le"
	EQ                 = "eq"
	NE                 = "ne"
	LT                 = "lt"
	GT                 = "gt"
	Plus               = "plus"
	Minus              = "minus"
	And                = "bitand"
	Or                 = "bitor"
	Mod                = "mod"
	Xor                = "bitxor"
	Div                = "div"
	Mul                = "mul"
	UnaryNot           = "not" // Avoid name conflict with Not in github/pingcap/check.
	BitNeg             = "bitneg"
	IntDiv             = "intdiv"
	LogicXor           = "xor"
	NullEQ             = "nulleq"
	UnaryPlus          = "unaryplus"
	UnaryMinus         = "unaryminus"
	In                 = "in"
	Like               = "like"
	Case               = "case"
	Regexp             = "regexp"
	IsNull             = "isnull"
	IsTruthWithoutNull = "istrue" // Avoid name conflict with IsTrue in github/pingcap/check.
	IsTruthWithNull    = "istrue_with_null"
	IsFalsity          = "isfalse" // Avoid name conflict with IsFalse in github/pingcap/check.
	RowFunc            = "row"
	SetVar             = "setvar"
	GetVar             = "getvar"
	Values             = "values"
	BitCount           = "bit_count"
	GetParam           = "getparam"

	// common functions
	Coalesce = "coalesce"
	Greatest = "greatest"
	Least    = "least"
	Interval = "interval"

	// math functions
	Abs      = "abs"
	Acos     = "acos"
	Asin     = "asin"
	Atan     = "atan"
	Atan2    = "atan2"
	Ceil     = "ceil"
	Ceiling  = "ceiling"
	Conv     = "conv"
	Cos      = "cos"
	Cot      = "cot"
	CRC32    = "crc32"
	Degrees  = "degrees"
	Exp      = "exp"
	Floor    = "floor"
	Ln       = "ln"
	Log      = "log"
	Log2     = "log2"
	Log10    = "log10"
	PI       = "pi"
	Pow      = "pow"
	Power    = "power"
	Radians  = "radians"
	Rand     = "rand"
	Round    = "round"
	Sign     = "sign"
	Sin      = "sin"
	Sqrt     = "sqrt"
	Tan      = "tan"
	Truncate = "truncate"

	// time functions
	AddDate          = "adddate"
	AddTime          = "addtime"
	ConvertTz        = "convert_tz"
	Curdate          = "curdate"
	CurrentDate      = "current_date"
	CurrentTime      = "current_time"
	CurrentTimestamp = "current_timestamp"
	Curtime          = "curtime"
	Date             = "date"
	DateLiteral      = "'tidb`.(dateliteral"
	DateAdd          = "date_add"
	DateFormat       = "date_format"
	DateSub          = "date_sub"
	DateDiff         = "datediff"
	Day              = "day"
	DayName          = "dayname"
	DayOfMonth       = "dayofmonth"
	DayOfWeek        = "dayofweek"
	DayOfYear        = "dayofyear"
	Extract          = "extract"
	FromDays         = "from_days"
	FromUnixTime     = "from_unixtime"
	GetFormat        = "get_format"
	Hour             = "hour"
	LocalTime        = "localtime"
	LocalTimestamp   = "localtimestamp"
	MakeDate         = "makedate"
	MakeTime         = "maketime"
	MicroSecond      = "microsecond"
	Minute           = "minute"
	Month            = "month"
	MonthName        = "monthname"
	Now              = "now"
	PeriodAdd        = "period_add"
	PeriodDiff       = "period_diff"
	Quarter          = "quarter"
	SecToTime        = "sec_to_time"
	Second           = "second"
	StrToDate        = "str_to_date"
	SubDate          = "subdate"
	SubTime          = "subtime"
	Sysdate          = "sysdate"
	Time             = "time"
	TimeLiteral      = "'tidb`.(timeliteral"
	TimeFormat       = "time_format"
	TimeToSec        = "time_to_sec"
	TimeDiff         = "timediff"
	Timestamp        = "timestamp"
	TimestampLiteral = "'tidb`.(timestampliteral"
	TimestampAdd     = "timestampadd"
	TimestampDiff    = "timestampdiff"
	ToDays           = "to_days"
	ToSeconds        = "to_seconds"
	UnixTimestamp    = "unix_timestamp"
	UTCDate          = "utc_date"
	UTCTime          = "utc_time"
	UTCTimestamp     = "utc_timestamp"
	Week             = "week"
	Weekday          = "weekday"
	WeekOfYear       = "weekofyear"
	Year             = "year"
	YearWeek         = "yearweek"
	LastDay          = "last_day"
	// TSO functions
	// TiDBBoundedStaleness is used to determine the TS for a read only request with the given bounded staleness.
	// It will be used in the Stale Read feature.
	// For more info, please see AsOfClause.
	TiDBBoundedStaleness = "tidb_bounded_staleness"
	TiDBParseTso         = "tidb_parse_tso"

	// string functions
	ASCII           = "ascii"
	Bin             = "bin"
	Concat          = "concat"
	ConcatWS        = "concat_ws"
	Convert         = "convert"
	Elt             = "elt"
	ExportSet       = "export_set"
	Field           = "field"
	Format          = "format"
	FromBase64      = "from_base64"
	InsertFunc      = "insert_func"
	Instr           = "instr"
	Lcase           = "lcase"
	Left            = "left"
	Length          = "length"
	LoadFile        = "load_file"
	Locate          = "locate"
	Lower           = "lower"
	Lpad            = "lpad"
	LTrim           = "ltrim"
	MakeSet         = "make_set"
	Mid             = "mid"
	Oct             = "oct"
	OctetLength     = "octet_length"
	Ord             = "ord"
	Position        = "position"
	Quote           = "quote"
	Repeat          = "repeat"
	Replace         = "replace"
	Reverse         = "reverse"
	Right           = "right"
	RTrim           = "rtrim"
	Space           = "space"
	Strcmp          = "strcmp"
	Substring       = "substring"
	Substr          = "substr"
	SubstringIndex  = "substring_index"
	ToBase64        = "to_base64"
	Trim            = "trim"
	Translate       = "translate"
	Upper           = "upper"
	Ucase           = "ucase"
	Hex             = "hex"
	Unhex           = "unhex"
	Rpad            = "rpad"
	BitLength       = "bit_length"
	CharFunc        = "char_func"
	CharLength      = "char_length"
	CharacterLength = "character_length"
	FindInSet       = "find_in_set"
	WeightString    = "weight_string"
	Soundex         = "soundex"

	// information functions
	Benchmark            = "benchmark"
	Charset              = "charset"
	Coercibility         = "coercibility"
	Collation            = "collation"
	ConnectionID         = "connection_id"
	CurrentUser          = "current_user"
	CurrentRole          = "current_role"
	Database             = "database"
	FoundRows            = "found_rows"
	LastInsertId         = "last_insert_id"
	RowCount             = "row_count"
	Schema               = "schema"
	SessionUser          = "session_user"
	SystemUser           = "system_user"
	User                 = "user"
	Version              = "version"
	TiDBVersion          = "tidb_version"
	TiDBIsDDLOwner       = "tidb_is_ddl_owner"
	TiDBDecodePlan       = "tidb_decode_plan"
	TiDBDecodeSQLDigests = "tidb_decode_sql_digests"
	FormatBytes          = "format_bytes"
	FormatNanoTime       = "format_nano_time"

	// control functions
	If     = "if"
	Ifnull = "ifnull"
	Nullif = "nullif"

	// miscellaneous functions
	AnyValue        = "any_value"
	DefaultFunc     = "default_func"
	InetAton        = "inet_aton"
	InetNtoa        = "inet_ntoa"
	Inet6Aton       = "inet6_aton"
	Inet6Ntoa       = "inet6_ntoa"
	IsFreeLock      = "is_free_lock"
	IsIPv4          = "is_ipv4"
	IsIPv4Compat    = "is_ipv4_compat"
	IsIPv4Mapped    = "is_ipv4_mapped"
	IsIPv6          = "is_ipv6"
	IsUsedLock      = "is_used_lock"
	IsUUID          = "is_uuid"
	MasterPosWait   = "master_pos_wait"
	NameConst       = "name_const"
	ReleaseAllLocks = "release_all_locks"
	Sleep           = "sleep"
	UUID            = "uuid"
	UUIDShort       = "uuid_short"
	UUIDToBin       = "uuid_to_bin"
	BinToUUID       = "bin_to_uuid"
	VitessHash      = "vitess_hash"
	// get_lock() and release_lock() is parsed but do nothing.
	// It is used for preventing error in Ruby's activerecord migrations.
	GetLock     = "get_lock"
	ReleaseLock = "release_lock"

	// encryption and compression functions
	AesDecrypt               = "aes_decrypt"
	AesEncrypt               = "aes_encrypt"
	Compress                 = "compress"
	Decode                   = "decode"
	DesDecrypt               = "des_decrypt"
	DesEncrypt               = "des_encrypt"
	Encode                   = "encode"
	Encrypt                  = "encrypt"
	MD5                      = "md5"
	OldPassword              = "old_password"
	PasswordFunc             = "password_func"
	RandomBytes              = "random_bytes"
	SHA1                     = "sha1"
	SHA                      = "sha"
	SHA2                     = "sha2"
	Uncompress               = "uncompress"
	UncompressedLength       = "uncompressed_length"
	ValidatePasswordStrength = "validate_password_strength"

	// json functions
	JSONType          = "json_type"
	JSONExtract       = "json_extract"
	JSONUnquote       = "json_unquote"
	JSONArray         = "json_array"
	JSONObject        = "json_object"
	JSONMerge         = "json_merge"
	JSONSet           = "json_set"
	JSONInsert        = "json_insert"
	JSONReplace       = "json_replace"
	JSONRemove        = "json_remove"
	JSONContains      = "json_contains"
	JSONContainsPath  = "json_contains_path"
	JSONValid         = "json_valid"
	JSONArrayAppend   = "json_array_append"
	JSONArrayInsert   = "json_array_insert"
	JSONMergePatch    = "json_merge_patch"
	JSONMergePreserve = "json_merge_preserve"
	JSONPretty        = "json_pretty"
	JSONQuote         = "json_quote"
	JSONSearch        = "json_search"
	JSONStorageSize   = "json_storage_size"
	JSONDepth         = "json_depth"
	JSONKeys          = "json_keys"
	JSONLength        = "json_length"

	// TiDB internal function.
	TiDBDecodeKey       = "tidb_decode_key"
	TiDBDecodeBase64Key = "tidb_decode_base64_key"

	// MVCC information fetching function.
	GetMvccInfo = "get_mvcc_info"

	// Sequence function.
	NextVal = "nextval"
	LastVal = "lastval"
	SetVal  = "setval"
)

type FuncCallExprType int8
```

### FuncCastExpr

```golang
// FuncCastExpr is the cast function converting value to another type, e.g, cast(expr AS signed).
// See https://dev.mysql.com/doc/refman/5.7/en/cast-functions.html
type FuncCastExpr struct {
	funcNode
	// Expr is the expression to be converted.
	Expr ExprNode
	// Tp is the conversion type.
	Tp *types.FieldType
	// FunctionType is either Cast, Convert or Binary.
	FunctionType CastFunctionType
	// ExplicitCharSet is true when charset is explicit indicated.
	ExplicitCharSet bool
}

// for example:
// SELECT CONVERT('test', CHAR CHARACTER SET utf8 COLLATE utf8_bin);
// SELECT CAST('test' AS CHAR CHARACTER SET utf8) COLLATE utf8_bin;
```

### TrimDirectionExpr

```golang
// TrimDirectionExpr is an expression representing the trim direction used in the TRIM() function.
type TrimDirectionExpr struct {
	exprNode
	// Direction is the trim direction
	Direction TrimDirectionType
}

// for example: SELECT TRIM(LEADING 'x' FROM 'xxxbarxxx');
```

### AggregateFuncExpr

```golang
// AggregateFuncExpr represents aggregate function expression.
type AggregateFuncExpr struct {
	funcNode
	// F is the function name.
	F string
	// Args is the function args.
	Args []ExprNode
	// Distinct is true, function hence only aggregate distinct values.
	// For example, column c1 values are "1", "2", "2",  "sum(c1)" is "5",
	// but "sum(distinct c1)" is "3".
	Distinct bool
	// Order is only used in GROUP_CONCAT
	Order *OrderByClause
}
```

### extra--AggregateFuncName

```golang
const (
    // 18
	// AggFuncCount is the name of Count function.
	AggFuncCount = "count"
	// AggFuncSum is the name of Sum function.
	AggFuncSum = "sum"
	// AggFuncAvg is the name of Avg function.
	AggFuncAvg = "avg"
	// AggFuncFirstRow is the name of FirstRowColumn function.
	AggFuncFirstRow = "firstrow"
	// AggFuncMax is the name of max function.
	AggFuncMax = "max"
	// AggFuncMin is the name of min function.
	AggFuncMin = "min"
	// AggFuncGroupConcat is the name of group_concat function.
	AggFuncGroupConcat = "group_concat"
	// AggFuncBitOr is the name of bit_or function.
	AggFuncBitOr = "bit_or"
	// AggFuncBitXor is the name of bit_xor function.
	AggFuncBitXor = "bit_xor"
	// AggFuncBitAnd is the name of bit_and function.
	AggFuncBitAnd = "bit_and"
	// AggFuncVarPop is the name of var_pop function
	AggFuncVarPop = "var_pop"
	// AggFuncVarSamp is the name of var_samp function
	AggFuncVarSamp = "var_samp"
	// AggFuncStddevPop is the name of stddev_pop/std/stddev function
	AggFuncStddevPop = "stddev_pop"
	// AggFuncStddevSamp is the name of stddev_samp function
	AggFuncStddevSamp = "stddev_samp"
	// AggFuncJsonArrayagg is the name of json_arrayagg function
	AggFuncJsonArrayagg = "json_arrayagg"
	// AggFuncJsonObjectAgg is the name of json_objectagg function
	AggFuncJsonObjectAgg = "json_objectagg"
	// AggFuncApproxCountDistinct is the name of approx_count_distinct function.
	AggFuncApproxCountDistinct = "approx_count_distinct"
	// AggFuncApproxPercentile is the name of approx_percentile function.
	AggFuncApproxPercentile = "approx_percentile"
)
```

### WindowFuncExpr

```golang
// WindowFuncExpr represents window function expression.
type WindowFuncExpr struct {
	funcNode

	// F is the function name.
	F string
	// Args is the function args.
	Args []ExprNode
	// Distinct cannot be true for most window functions, except `max` and `min`.
	// We need to raise error if it is not allowed to be true.
	Distinct bool
	// IgnoreNull indicates how to handle null value.
	// MySQL only supports `RESPECT NULLS`, so we need to raise error if it is true.
	IgnoreNull bool
	// FromLast indicates the calculation direction of this window function.
	// MySQL only supports calculation from first, so we need to raise error if it is true.
	FromLast bool
	// Spec is the specification of this window.
	Spec WindowSpec
}
```

### extra--WindowFuncName

```golang
const (
    // 11
	// WindowFuncRowNumber is the name of row_number function.
	WindowFuncRowNumber = "row_number"
	// WindowFuncRank is the name of rank function.
	WindowFuncRank = "rank"
	// WindowFuncDenseRank is the name of dense_rank function.
	WindowFuncDenseRank = "dense_rank"
	// WindowFuncCumeDist is the name of cume_dist function.
	WindowFuncCumeDist = "cume_dist"
	// WindowFuncPercentRank is the name of percent_rank function.
	WindowFuncPercentRank = "percent_rank"
	// WindowFuncNtile is the name of ntile function.
	WindowFuncNtile = "ntile"
	// WindowFuncLead is the name of lead function.
	WindowFuncLead = "lead"
	// WindowFuncLag is the name of lag function.
	WindowFuncLag = "lag"
	// WindowFuncFirstValue is the name of first_value function.
	WindowFuncFirstValue = "first_value"
	// WindowFuncLastValue is the name of last_value function.
	WindowFuncLastValue = "last_value"
	// WindowFuncNthValue is the name of nth_value function.
	WindowFuncNthValue = "nth_value"
)
```

### skip(1)

TimeUnitExpr

```golang
// TimeUnitExpr is an expression representing a time or timestamp unit.
type TimeUnitExpr struct {
	exprNode
	// Unit is the time or timestamp unit.
	Unit TimeUnitType
}
```

<font color="red">unknown</font>

GetFormatSelectorExpr

```golang
// GetFormatSelectorExpr is an expression used as the first argument of GET_FORMAT() function.
type GetFormatSelectorExpr struct {
	exprNode
	// Selector is the GET_FORMAT() selector.
	Selector GetFormatSelectorType
}
```

## misc.go  (51)

### skip(51)

TraceStmt
ExplainForStmt
ExplainStmt
PlanReplayerStmt
PrepareStmt
DeallocateStmt
ExecuteStmt
BeginStmt
BinlogStmt
CommitStmt
RollbackStmt
UseStmt
VariableAssignment
FlushStmt
KillStmt
SetStmt
SetConfigStmt
SetPwdStmt
ChangeStmt
SetRoleStmt
SetDefaultRoleStmt
CreateUserStmt
AlterUserStmt
AlterInstanceStmt
DropUserStmt
CreateBindingStmt
DropBindingStmt
CreateStatisticsStmt
DropStatisticsStmt
DoStmt
AdminStmt
PrivElem
RevokeStmt
RevokeRoleStmt
GrantStmt
GrantProxyStmt
GrantRoleStmt
ShutdownStmt
RestartStmt
HelpStmt
RenameUserStmt
UserToUser
BRIEStmt
PurgeImportStmt
CreateImportStmt
StopImportStmt
ResumeImportStmt
AlterImportStmt
DropImportStmt
ShowImportStmt
TableOptimizerHint

## stats.go  (3)

### skip(3)

AnalyzeTableStmt
DropStatsStmt
LoadStatsStmt

# parser/test_driver  (2)

## test_driver.go  (2)

### ValueExpr

```golang
// ValueExpr is the simple value expression.
type ValueExpr struct {
	ast.TexprNode
	Datum
	projectionOffset int
}

// Datum is a data box holds different kind of data.
// It has better performance and is easier to use than `interface{}`.
type Datum struct {
	k byte        // datum kind.
	i int64       // i can hold int64 uint64 float64 values.
	b []byte      // b can hold string or []byte values.
	x interface{} // x hold all other types.
}

// Kind constants.
const (
	KindNull          byte = 0
	KindInt64         byte = 1
	KindUint64        byte = 2
	KindFloat32       byte = 3
	KindFloat64       byte = 4
	KindString        byte = 5
	KindBytes         byte = 6
	KindBinaryLiteral byte = 7 // Used for BIT / HEX literals.
	KindMysqlDecimal  byte = 8
	KindMysqlDuration byte = 9
	KindMysqlEnum     byte = 10
	KindMysqlBit      byte = 11 // Used for BIT table column values.
	KindMysqlSet      byte = 12
	KindMysqlTime     byte = 13
	KindInterface     byte = 14
	KindMinNotNull    byte = 15
	KindMaxValue      byte = 16
	KindRaw           byte = 17
	KindMysqlJSON     byte = 18
)
```

### skip(1)

ParamMarkerExpr

# planner/core  (1)

## common_plans.go  (1)

### skip(1)

Change

# types/parser_driver  (2)

## value_expr.go  (2)

### ValueExpr

```golang
// ValueExpr is the simple value expression.
type ValueExpr struct {
   ast.TexprNode
   types.Datum
   projectionOffset int
}

// Kind constants.
const (
	KindNull          byte = 0
	KindInt64         byte = 1
	KindUint64        byte = 2
	KindFloat32       byte = 3
	KindFloat64       byte = 4
	KindString        byte = 5
	KindBytes         byte = 6
	KindBinaryLiteral byte = 7 // Used for BIT / HEX literals.
	KindMysqlDecimal  byte = 8
	KindMysqlDuration byte = 9
	KindMysqlEnum     byte = 10
	KindMysqlBit      byte = 11 // Used for BIT table column values.
	KindMysqlSet      byte = 12
	KindMysqlTime     byte = 13
	KindInterface     byte = 14
	KindMinNotNull    byte = 15
	KindMaxValue      byte = 16
	KindRaw           byte = 17
	KindMysqlJSON     byte = 18
)

// Datum is a data box holds different kind of data.
// It has better performance and is easier to use than `interface{}`.
type Datum struct {
	k         byte        // datum kind.
	decimal   uint16      // decimal can hold uint16 values.
	length    uint32      // length can hold uint32 values.
	i         int64       // i can hold int64 uint64 float64 values.
	collation string      // collation hold the collation information for string value.
	b         []byte      // b can hold string or []byte values.
	x         interface{} // x hold all other types.
}
```

### skip(1)

ParamMarkerExpr