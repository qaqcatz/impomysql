{
    -- According to the MySQL features supported by this .yy file, we will (manual) implement the visitor of impo.
    -- Note that:
    -- (1) Although impo does not support the above features, you can still use them.
    --     impo will preprocess the statement and only work on the sub-ast of legal features.
    -- (2) We have made a lot of effort to analyze the features of MYSQL:
    --     * analyze https://dev.mysql.com/doc/refman/8.0/en/
    --     * analyze all 175 ast.Node of tidb parser(https://github.com/pingcap/tidb/tree/v5.4.2/parser),
    --       of which 57 nodes are related to query, including 31 operators and 317 functions.
    --     Nevertheless, we may still be ill-considered. Please contact us via github issue to help us improve impo.
    -- impo works on SELECT statement, focus on WHERE and HAVING(ON).
    -- The unsupported/supported features are as follows: (we will use operations to summarize operators / functions / statements / clauses uniformly)
    -- (1) impo cannot support these features:
    --     * Numerical operations and their descendants, including |, &, ~, <<, >>, +, -, *, /(DIV), %(MOD), ^.
    --     * For logical operations, impo will continue to visit the descendants of OR(||), AND(&&), NOT(!), IS TRUE|FALSE.
    --       When meeting =, >=, >, <=, <, !=, <>, IN, BETWEEN AND, LIKE, REGEXP, impo will mutate the predicate and
    --       stop visiting the the descendants of these operations.
    --       Cannot support XOR, IS NULL, <=> and their descendants.
    --       Cannot support SOUNDEX(), SOUNDS LIKE, see https://dev.mysql.com/doc/refman/8.0/en/string-functions.html#function_soundex
    --     * aggregate functions, window functions, GROUP BY
    --     * logical operations IN the fields of SELECT.
    --     * LEFT|RIGHT JOIN
    --     * LIMIT
    --     * flow control operations, see https://dev.mysql.com/doc/refman/8.0/en/flow-control-statements.html
    --       CASE statement, CASE operator;
    --       IF statement, IF(), IFNULL(), NULLIF() functions, see https://dev.mysql.com/doc/refman/8.0/en/flow-control-functions.html
    --       ITERATE, LEAVE, REPEAT, RETURN, WHILE.
    --     * uncertain functions, such as random function, current time function
    --     * subqueries with value SELECT.
    --
    --     * SELECT INTO(INSERT INTO SELECT)
    --     * variable(SET)
    --     * MATCH
    --     * PREPARE, EXECUTE, {DEALLOCATE | DROP} PREPARE
    --     * EXPLAIN
    --     * Optimizer Hints
    -- (2) impo supports these features:
    --     * WITH [RECURSIVE]
    --     * UNION
    --     * part of JOIN: JOIN(CROSS JOIN, INNER JOIN), STRAIGHT_JOIN, NATURAL JOIN
    --     * WHERE, HAVING, ORDER BY
    --     * part of logical operations: OR(||), AND(&&), NOT(!), IS TRUE|FALSE, =, >=, >, <=, <, !=, <>, IN, BETWEEN AND, LIKE, REGEXP
    --     * subqueries with EXISTS | IN | comparison_operator + ANY, SOME, ALL.
    --       ROW subqueries with IN.
    --     * LIKE, REGEXP
    --     * Index Hints
    --     * others: BINARY, INTERVAL, CHARACTER SET, COLLATE, CONVERT, CAST
    -- This .yy file can also be used to generate random sqls.
    -- Guided by practical concerns, the number of some operations is limited:
    --   * the number of fields in each SELECT statement is 3,
    --     except for subquery, whose column number is 1(normal) or 2(ROW).
    --   * only support one non-recursive WITH.
    --   * in order to prevent the results from being too large, we only support one UNION, one JOIN.
    --   * USING and NATURAL JOIN often produce an empty result, we do not consider them.
    --   * the more nested SELECT, the easier it is to produce empty results,
    --     so we do not support recursive SELECT.
    --   * only support one column in ORDER BY.
    --   * each WHERE must has an explicit IS TRUE/IS FALSE
    --   * complex AND, OR structures will affect the effectiveness of impo, so we only support some fixed structures:
    --     x | x AND x | x OR x | (x OR x) AND (x OR x) | (x AND x) OR (x AND x)
    --     Similarly, we do not support recursive binary/unary operations.
    --   * only support one column in subqueries with EXISTS | IN | comparison_operator + ANY, SOME, ALL.
    --   * only support two column in ROW subqueries with IN.
    --   * LIKE and REGEXP only support some fixed expressions, and do not support ESCAPE, BINARY.
    --   * we use the zz file for charset/collate testing, so we skim charset/collate statements in yy file,
    --     such as CHARACTER SET, COLLATE, CONVERT(we use CAST instead of CONVERT).
    --   * index hints cannot generate FOR.
    --   * only support INTERVAL 1 unit.
    --   Take it easy, their number is not limited in impo.

    -- How this .yy works?
    -- We treat a query statement as a tree, each node is a subquery, a leaf node is also a special query.
    -- For example:
    -- ==================================================
    -- SELECT (1), f1, f2, f3 FROM
    --   (SELECT (2), _field AS f1 FROM _table) AS t1
    --   JOIN
    --   (SELECT (3), f4 AS f2, f5 AS f3 FROM
    --     (SELECT (4), _field AS f4, _field AS f5, _field AS f6 FROM _table) AS t2
    --   WHERE f6 > 0
    --   ) AS t3
    -- ON f1 >= f2
    -- WHERE f3 = 'A'
    -- UNION
    -- SELECT (5), f1, f2, f3 FROM
    --   (SELECT (6), _field AS f1, _field AS f2, _field AS f3 FROM _table) AS t4
    -- WHERE f3 IN (SELECT (7) AS f3 FROM _table)
    -- ==================================================
    --                     query
    --                   /       \
    --          SELECT(1)         SELECT(5)
    --         /         \       /         \
    -- SELECT(2)  SELECT(3)     SELECT(6)  SELECT(7)
    --                |
    --            SELECT(4)
    -- ==================================================
    -- In order to make sure these SELECT fields are valid:
    -- (1) Before visiting each SELECT node, we will generate its aliases first.
    -- (2) We will ensure that each SELECT node contains all aliases specified by its father node.

    -- @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

    -- lua random seed, use the seed of go-randgen
    seed = tonumber(_int())
    math.randomseed(seed)

    -- example: optional({"DISTINCT", ""}), random return "DISTINCT" or ""
    function optional(args)
        cnt = 0
        for i,v in ipairs(args) do
            cnt = cnt + 1
        end
        rd = math.random(cnt)
        return args[rd]
    end

    function mySplit(str)
        str = str .. ","
        res = {}
        i = 1
        l = 1
        flag = 0
        for r = 1, string.len(str) do
            c = string.sub(str,r,r)
            if c == '`' then
                if flag == 0 then
                    flag = 1
                else
                    flag = 0
                end
            end
            if c == ',' and flag == 0 then
                res[i] = string.sub(str, l, r-1)
                i = i+1
                l = r+1
            end
        end
        return res
    end

    function keyFields()
        fields = mySplit(_field_list())
        res = {}
        rid = 1
        for i,v in ipairs(fields) do
            if string.find(fields[i], "key") then
                res[rid] = v
                rid = rid + 1
            end
        end
        return res
    end

    -- give each table a unique alias, only for derived table
    tid = 0
    function genTableAlias()
        tid = tid + 1
        return "t" .. tid
    end

    -- give each field a unique alias
    fid = 0
    function genFieldAlias()
        fid = fid + 1
        return "f" .. fid
    end

    -- return 3*genFieldAlias()
    function gen3FieldAliases()
        return {genFieldAlias(), genFieldAlias(), genFieldAlias()}
    end

    -- current aliases, set by the father node.
    curAStack = {["size"]=0}

    -- child aliases, set by the current node.
    -- We will use it to set the aliases of each child node.
    -- It will also be used as the column names of the current query.
    childAStack = {["size"]=0}

    -- Tell the expr which column names / aliases are available,
    -- so as to ensure the correct semantics of WHERE and HAVING.
    -- Default: stackTop(childAStack)
    colNamesStack = {["size"]=0}

    function pushStack(stack, element)
        stack.size = stack.size + 1
        stack[stack.size] = {}
        for i,v in ipairs(element) do
            stack[stack.size][i] = element[i]
        end
    end

    function popStack(stack, element)
        if stack.size < 1 then
            error("pop: empty stack!")
        end
        stack[stack.size] = nil
        stack.size = stack.size - 1
    end

    function stackTop(stack)
        if stack.size < 1 then
            error("top: empty stack!")
        end
        return stack[stack.size]
    end

    -- for HAVING
    function setStackTop(stack, element)
        if stack.size < 1 then
            error("pop: empty stack!")
        end
        stack[stack.size] = element
    end

    -- @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
}

query:
    {
        -- init each query
        tid = 0
        fid = 0
        curAStack = {["size"]=1, [1]=gen3FieldAliases()}
        childAStack = {["size"]=0}
        colNamesStack = {["size"]=0}
    }
    myquery

myquery:
    withclause

withclause:
    WITH MYWITH AS (setoprselectlist) SELECT * FROM MYWITH
    | setoprselectlist

# UNION, MYSQL cannot support EXCEPT, INTERSECT
setoprselectlist:
    (select2) setoprtype (select2)
    | select2

# union type
setoprtype:
    UNION
    | UNION ALL

# Special SELECT for JOIN
select2:
    { pushStack( childAStack, gen3FieldAliases() )
      pushStack( colNamesStack, stackTop(childAStack) ) }
    SELECT fieldlist FROM tablerefsclause2
    { popStack(childAStack)
      popStack(colNamesStack) }
    | select

tablerefsclause2:
    join2

# Aliases assignment:
# JOIN ON: (1, x, 3), (x, 2, x) onclause
# NATURAL JOIN: (x, 2, x), (1, x, 3)
join2:
    { pushStack( curAStack, { stackTop(childAStack)[1], genFieldAlias(), stackTop(childAStack)[3] } ) }
    tablesource2
    { popStack(curAStack) }
    jointype
    { pushStack( curAStack, { genFieldAlias(), stackTop(childAStack)[2], genFieldAlias() } ) }
    tablesource2
    { popStack(curAStack) }
    onclause
    | { pushStack( curAStack, { genFieldAlias(), stackTop(childAStack)[2], genFieldAlias() } ) }
      tablesource2
      { popStack(curAStack) }
      NATURAL JOIN
      { pushStack( curAStack, { stackTop(childAStack)[1], genFieldAlias(), stackTop(childAStack)[3] } ) }
      tablesource2
      { popStack(curAStack) }

# (1) In MySQL, JOIN, CROSS JOIN, and INNER JOIN are syntactic equivalents (they can replace each other).
# (2) STRAIGHT_JOIN is similar to JOIN, except that the left table is always read before the right table.
# This can be used for those (few) cases for which the join optimizer processes the tables in a suboptimal order.
jointype:
    JOIN
    | CROSS JOIN
    | INNER JOIN
    | STRAIGHT_JOIN

onclause:
    | ON explicitp

tablesource2:
    (select) AS { print( genTableAlias() ) }
    | (tablename) AS { print( genTableAlias() ) }

select:
    { pushStack( childAStack, gen3FieldAliases() )
      pushStack( colNamesStack, stackTop(childAStack) ) }
    SELECT fieldlist FROM tablerefsclause whereclause havingclause orderbyclause
    { popStack(childAStack)
      popStack(colNamesStack) }

# |selectfields| = 3
fieldlist:
    (bit_expr) AS { print( stackTop(curAStack)[1] ) },
    (bit_expr) AS { print( stackTop(curAStack)[2] ) },
    (bit_expr) AS { print( stackTop(curAStack)[3] ) }

tablerefsclause:
    join

# only for aliases stack
join:
    { pushStack( curAStack, stackTop(childAStack) ) }
    tablesource
    { popStack(curAStack) }

tablesource:
    (tablename) AS { print( genTableAlias() ) }

# Special select to make sure the SELECT fields are valid
tablename:
    SELECT
    _field AS { print( stackTop(curAStack)[1] ) },
    _field AS { print( stackTop(curAStack)[2] ) },
    _field AS { print( stackTop(curAStack)[3] ) }
    FROM _table indexhints

# INDEX = KEY
indexhints:
    | USE INDEX ({ print( optional( keyFields() ) ) })
    | USE INDEX ({ print( optional( keyFields() ) ) },{ print( optional( keyFields() ) ) })
    | { print( optional( {"IGNORE", "FORCE"} ) ) } INDEX ({ print( optional( keyFields() ) ) })
    | { print( optional( {"IGNORE", "FORCE"} ) ) } INDEX ({ print( optional( keyFields() ) ) },{ print( optional( keyFields() ) ) })

whereclause:
    | WHERE explicitp

# HAVING can only use column names / aliases which appear in fieldlist
havingclause:
    | HAVING { setStackTop( colNamesStack, stackTop(curAStack) ) } explicitp { setStackTop( colNamesStack, stackTop(childAStack) ) }

explicitp:
    (expr) IS TRUE
    | (expr) IS FALSE

orderbyclause:
    | ORDER BY columnname

columnname:
    { print( optional( stackTop(colNamesStack) ) ) }

expr:
    expr01
    | (expr01) AND (expr01) # &&
    | (expr01) OR (expr01) # ||
    | (expr01) AND (expr01) AND (expr01)
    | (expr01) OR (expr01) OR (expr01)
    | (expr01) AND (expr01) OR (expr01)
    | (expr01) OR (expr01) AND (expr01)
    | ((expr01) OR (expr01)) AND ((expr01) OR (expr01))
    | ((expr01) AND (expr01)) OR ((expr01) AND (expr01))

expr01:
    boolean_primary
    | NOT (boolean_primary) # !
    | (boolean_primary) IS TRUE
    | (boolean_primary) IS FALSE

boolean_primary:
    bit_expr
    | (bit_expr) comparison_operator (bit_expr)
    | (bit_expr) comparison_operator {print(optional({"ANY", "SOME", "ALL"}))} (subqueryexpr)
    | EXISTS (subqueryexpr)
    | (bit_expr) {print(optional({"NOT", ""}))} IN (subqueryexpr)
    | ROW(bit_expr, bit_expr) {print(optional({"NOT", ""}))} IN (subqueryexpr2)
    | (bit_expr) {print(optional({"NOT", ""}))} IN (simple_expr, simple_expr, simple_expr)
    | (simple_expr) {print(optional({"NOT", ""}))} BETWEEN simple_expr AND simple_expr
    | CAST((simple_expr) AS CHAR) {print(optional({"NOT", ""}))} LIKE like_expr
    | CAST((simple_expr) AS CHAR) {print(optional({"NOT", ""}))} REGEXP reg_expr

comparison_operator:
    {print("=")}
    | {print(">=")}
    | {print(">")}
    | {print("<=")}
    | {print("<")}
    | {print("!=")} # <>

subqueryexpr:
    SELECT _field FROM _table indexhints

subqueryexpr2:
    SELECT _field, _field FROM _table indexhints

like_expr:
    {print("'%0%'")}
    | {print("'%1%'")}

reg_expr:
    {print("'^[a-z]'")}
    | {print("'^[0-9]'")}
    | {print("'[a-z]$'")}
    | {print("'[0-9]$'")}
    | {print("'[a-z]+'")}
    | {print("'[0-9]+'")}
    | {print("'[a-z]+[a-z]+'")}
    | {print("'[0-9]+[0-9]+'")}
    | {print("'[a-z]+[0-9]+'")}
    | {print("'[0-9]+[a-z]+'")}

bit_expr:
    bit_expr1
    | bit_expr2

bit_expr2:
    bit_expr1 bit_operator2 bit_expr1
    | bit_expr1 bit_operator2 bit_expr1 bit_operator2 bit_expr1
    | bit_expr1 {print("+")} interval_expr
    | bit_expr1 {print("-")} interval_expr

bit_operator2:
    {print("|")}
    | {print("&")}
    | {print("<<")}
    | {print(">>")}
    | {print("+")}
    | {print("-")}
    | {print("*")}
    | DIV # /
    | MOD # %
    | {print("^")}

bit_expr1:
    # boost weight
    simple_expr
    | simple_expr
    | simple_expr
    | bit_operator1 simple_expr

bit_operator1:
    {print("-")}
    | {print("~")}
    | {print("!")}
    | BINARY # trick

# https://dev.mysql.com/doc/refman/8.0/en/expressions.html#temporal-intervals
interval_expr:
    INTERVAL 1 interval_unit

interval_unit:
    MICROSECOND
    | SECOND
    | MINUTE
    | HOUR
    | DAY
    | WEEK
    | MONTH
    | QUARTER
    | YEAR
    | SECOND_MICROSECOND
    | MINUTE_MICROSECOND
    | MINUTE_SECOND
    | HOUR_MICROSECOND
    | HOUR_SECOND
    | HOUR_MINUTE
    | DAY_MICROSECOND
    | DAY_SECOND
    | DAY_MINUTE
    | DAY_HOUR
    | YEAR_MONTH

simple_expr:
    # boost weight
    literal
    | identifier
    | identifier
    | identifier
    # | function_call

# https://dev.mysql.com/doc/refman/8.0/en/literals.html
literal:
    # boost weight
    strl
    | strl
    | strl
    | numl
    | numl
    | numl
    | datetimel
    | datetimel
    | booll
    | booll
    | nullv

strl:
    _letter
    | _english

numl:
    _digit
    | _int

datetimel:
    _year
    | _date
    | _time
    | _datetime

# hexl:

# bitl:

booll:
    0 # FALSE
    | 1 # TRUE

nullv:
    NULL

identifier:
    columnname