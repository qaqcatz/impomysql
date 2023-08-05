{
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
    INNER JOIN
    { pushStack( curAStack, { genFieldAlias(), stackTop(childAStack)[2], genFieldAlias() } ) }
    tablesource2
    { popStack(curAStack) }
    ON explicitp
    | { pushStack( curAStack, { genFieldAlias(), stackTop(childAStack)[2], genFieldAlias() } ) }
      tablesource2
      { popStack(curAStack) }
      CROSS JOIN
      { pushStack( curAStack, { stackTop(childAStack)[1], genFieldAlias(), stackTop(childAStack)[3] } ) }
      tablesource2
      { popStack(curAStack) }
    | { pushStack( curAStack, { genFieldAlias(), stackTop(childAStack)[2], genFieldAlias() } ) }
      tablesource2
      { popStack(curAStack) }
      NATURAL JOIN
      { pushStack( curAStack, { stackTop(childAStack)[1], genFieldAlias(), stackTop(childAStack)[3] } ) }
      tablesource2
      { popStack(curAStack) }

tablesource2:
    (select) AS { print( genTableAlias() ) }
    | (tablename) AS { print( genTableAlias() ) }

select:
    { pushStack( childAStack, gen3FieldAliases() )
      pushStack( colNamesStack, stackTop(childAStack) ) }
    SELECT fieldlist FROM tablerefsclause whereclause orderbyclause
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
    FROM _table

whereclause:
    | WHERE explicitp

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
    ((bit_expr) {print("::bigint) != 0")}
    | (bit_expr) comparison_operator (bit_expr)
    | (bit_expr) comparison_operator {print(optional({"ANY", "SOME", "ALL"}))} (subqueryexpr)
    | EXISTS (subqueryexpr)
    | (bit_expr) {print(optional({"NOT", ""}))} IN (subqueryexpr)
    | ROW(bit_expr, bit_expr) {print(optional({"NOT", ""}))} IN (subqueryexpr2)
    | (bit_expr) {print(optional({"NOT", ""}))} IN (simple_expr, simple_expr, simple_expr)
    | (simple_expr) {print(optional({"NOT", ""}))} BETWEEN simple_expr AND simple_expr

comparison_operator:
    {print("=")}
    | {print(">=")}
    | {print(">")}
    | {print("<=")}
    | {print("<")}
    | {print("!=")} # <>

subqueryexpr:
    SELECT _field FROM _table

subqueryexpr2:
    SELECT _field, _field FROM _table

bit_expr:
    bit_expr1
    | bit_expr2

bit_expr1:
    # boost weight
    simple_expr
    | simple_expr
    | simple_expr
    | {print("-")} simple_expr
    | {print("~")} (simple_expr){print("::bigint")}

bit_expr2:
    (bit_expr1) bit_operator2_safe (bit_expr1)
    | (bit_expr1) bit_operator2_safe (bit_expr1) bit_operator2_safe (bit_expr1)
    | (bit_expr1){print("::bigint")} bit_operator2_unsafe (bit_expr1){print("::bigint")}
    | (bit_expr1){print("::bigint")} bit_operator2_unsafe (bit_expr1){print("::bigint")} bit_operator2_unsafe (bit_expr1){print("::bigint")}
    | (bit_expr1){print("::bigint")} {print("<<")} _digit
    | (bit_expr1){print("::bigint")} {print(">>")} _digit

bit_operator2_safe:
    {print("+")}
    | {print("-")}
    | {print("*")}
    | {print("/")}
    | {print("^")}

bit_operator2_unsafe:
    {print("|")}
    | {print("&")}
    | {print("%")}

simple_expr:
    literal_or_identifier
    | function_call

literal_or_identifier:
    # boost weight
    literal
    | identifier
    | identifier
    | identifier

literal:
    numl

numl:
    _digit
    | float0_1
    # | _int

float0_1:
    { print(math.random()) }

identifier:
    columnname

function_call:
    math_functions

math_functions:
    {print("abs")}(numl)
    | {print("acos")}(float0_1)
    | {print("asin")}(float0_1)
    | {print("atan")}(numl)
    | {print("ceil")}(numl)
    | {print("ceiling")}(numl)
    | {print("cos")}(numl)
    | {print("cot")}(numl)
    | {print("degrees")}(numl)
    | {print("floor")}(numl)
    | {print("ln(abs(")} numl {print("))")}
    | {print("log(abs(")} numl {print("))")}
    | {print("log10(abs(")} numl {print("))")}
    | {print("pi()")}
    | {print("round")}(numl)
    | {print("sign")}(numl)
    | {print("sin")}(numl)
    | {print("sqrt(abs(")} numl {print("))")}
    | {print("tan")}(numl)