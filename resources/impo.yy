{
    -- We have made a lot of effort to analyze the features of mysql:
    -- * analyzed https://dev.mysql.com/doc/refman/8.0/en/
    -- * analyzed all 175 ast.Node of tidb parser(https://github.com/pingcap/tidb/tree/v5.4.2/parser),
    -- of which 57 nodes are related to query, including 31 operators and 274 functions.
    -- Nevertheless, we may still be ill-considered. Please contact us via github issue to help us improve impo.
    --
    -- According to the MySQL features supported by this .yy file, we will (manually) implement the visitor of impo, see mutation/stage2/mutatevisitor.go
    -- This .yy file can also be used to generate random sqls, see https://github.com/pingcap/go-randgen.
    --
    -- If you use this .yy to generate random sqls, note that:
    -- * the number of fields in each SELECT statement is 3.
    -- * the more nested SELECT, the easier it is to produce empty results, so we do not support recursive subqueries.
    -- * only generate 1 column in subqueries with EXISTS | IN | comparison_operator + ANY, SOME, ALL.
    -- * only generate 2 column in ROW subqueries with IN.
    -- * do not support value SELECT.
    -- * in order to prevent the results from being too large, we only generate 1 UNION, 1 JOIN.
    -- * USING and NATURAL JOIN often produce an empty result, we do not support them.
    -- * complex AND, OR structures will affect the effectiveness of impo, so we only support these fixed structures:
    --   x, x AND x, x OR x, (x OR x) AND (x OR x), (x AND x) OR (x AND x).
    --   Similarly, we do not support recursive binary/unary operations.
    -- * only generate 1 non-recursive WITH.
    -- * only generate 1 column in ORDER BY.
    -- * do not support DISTINCT + ORDER BY.
    -- * only support some fixed expressions in LIKE, REGEXP, and do not support ESCAPE, BINARY.
    -- * do not support FOR in index hints.
    --
    -- lua functions:
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
    expr
    | (expr) IS TRUE
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
    literal_or_identifier
    | function_call

literal_or_identifier:
    # boost weight
    literal
    | identifier
    | identifier
    | identifier

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
    | { print(math.random()) }

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

function_call:
    math_functions
    | time_functions
    | string_functions
    | information_functions

math_functions:
    {print("abs")}(numl)
    | {print("acos")}(numl)
    | {print("asin")}(numl)
    | {print("atan")}(numl)
    | {print("atan2")}(numl)
    | {print("ceil")}(numl)
    | {print("ceiling")}(numl)
    | {print("cos")}(numl)
    | {print("cot")}(numl)
    | {print("crc32")}(literal)
    | {print("degrees")}(numl)
    | {print("floor")}(numl)
    | {print("ln")}(numl)
    | {print("log")}(numl)
    | {print("log2")}(numl)
    | {print("log10")}(numl)
    | {print("pi()")}
    | {print("radians")}(numl)
    | {print("round")}(numl)
    | {print("sign")}(numl)
    | {print("sin")}(numl)
    | {print("sqrt")}(numl)
    | {print("tan")}(numl)

time_functions:
    {print("adddate")}(_date, interval_expr)
    | {print("addtime")}(_datetime, _time)
    | {print("date")}(_datetime)
    | {print("date_add")}(_date, interval_expr)
    | {print("date_sub")}(_date, interval_expr)
    | {print("datediff")}(_datetime, _date)
    | {print("day")}(_date)
    | {print("dayname")}(_date)
    | {print("dayofmonth")}(_date)
    | {print("dayofweek")}(_date)
    | {print("dayofyear")}(_date)
    | {print("from_days")}(_int)
    | {print("from_unixtime")}(_int)
    | {print("hour")}(_time)
    | {print("microsecond")}(_datetime)
    | {print("minute")}(_datetime)
    | {print("month")}(_date)
    | {print("monthname")}(_date)
    | {print("quarter")}(_date)
    | {print("sec_to_time")}(_int)
    | {print("second")}(_time)
    | {print("subdate")}(_date, interval_expr)
    | {print("subtime")}(_datetime, _time)
    | {print("time")}(_datetime)
    | {print("time_to_sec")}(_time)
    | {print("timediff")}(_datetime, _datetime)
    | {print("timestamp")}(_date)
    | {print("to_days")}(_date)
    | {print("to_seconds")}(_datetime)
    | {print("unix_timestamp")}(_datetime)
    | {print("week")}(_date)
    | {print("weekday")}(_datetime)
    | {print("weekofyear")}(_date)
    | {print("year")}(_date)
    | {print("yearweek")}(_date)
    | {print("last_day")}(_datetime)

string_functions:
    {print("ASCII")}(literal_or_identifier)
    | {print("BIN")}(literal_or_identifier)
    | {print("CONCAT")}(literal_or_identifier, literal_or_identifier, literal_or_identifier)
    | {print("CONCAT_WS")}(literal_or_identifier, literal_or_identifier, literal_or_identifier)
    | {print("FIELD")}(literal_or_identifier, literal_or_identifier, literal_or_identifier, literal_or_identifier)
    | {print("INSTR")}(literal_or_identifier, literal_or_identifier)
    | {print("LCASE")}(literal_or_identifier)
    | {print("LEFT")}(literal_or_identifier, _digit)
    | {print("LENGTH")}(literal_or_identifier)
    | {print("LOCATE")}(literal_or_identifier, literal_or_identifier)
    | {print("LOWER")}(literal_or_identifier)
    | {print("LTRIM")}(literal_or_identifier)
    | {print("OCT")}(literal_or_identifier)
    | {print("OCTET_LENGTH")}(literal_or_identifier)
    | {print("ORD")}(literal_or_identifier)
    | {print("QUOTE")}(literal_or_identifier)
    | {print("REPEAT")}(literal_or_identifier, _digit)
    | {print("REPLACE")}(literal_or_identifier, literal_or_identifier, literal_or_identifier)
    | {print("REVERSE")}(literal_or_identifier)
    | {print("RIGHT")}(literal_or_identifier, _digit)
    | {print("RTRIM")}(literal_or_identifier)
    | {print("SPACE")}(_digit)
    | {print("STRCMP")}(literal_or_identifier, literal_or_identifier)
    | {print("SUBSTRING")}(literal_or_identifier, _digit)
    | {print("SUBSTR")}(literal_or_identifier, _digit)
    | {print("TO_BASE64")}(literal_or_identifier)
    | {print("TRIM")}(literal_or_identifier)
    | {print("UPPER")}(literal_or_identifier)
    | {print("UCASE")}(literal_or_identifier)
    | {print("HEX")}(literal_or_identifier)
    | {print("UNHEX")}(literal_or_identifier)
    | {print("BIT_LENGTH")}(literal_or_identifier)
    | {print("CHAR_LENGTH")}(literal_or_identifier)
    | {print("SOUNDEX")}(literal_or_identifier)

information_functions:
    {print("CHARSET")}(literal_or_identifier)
    | {print("COERCIBILITY")}(literal_or_identifier)
    | {print("COLLATION")}(literal_or_identifier)
    | {print("format_bytes")}(literal_or_identifier)


