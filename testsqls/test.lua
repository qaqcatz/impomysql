-- @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
-- .yy env

-- 随机生成一个 0-9 的数字
function _digit()
    return 0
end
-- 随机生成一个 'a' 到 'z' 之间的字母
function _letter()
    return 'a'
end
-- 随机生成一个英文单词
function _english()
    return 'hello'
end
-- 随机生成一个整型
function _int()
    return 10
end
-- 生成yyyy-MM-dd格式的随机日期
function _date()
    return "2021-12-31"
end
-- 随机生成一个年份
function _year()
    return "2021"
end
-- 随机生成一个hh:mm:ss的随机时间
function _time()
    return "23:59:59"
end
-- 随机生成一个yyyy-MM-dd hh:mm:ss的随机时间
function _datetime()
    return "2022-01-01 00:00:00"
end

function _field_list()
    return "`col_bigint_undef_signed`,`col_bigint_undef_unsigned`,`col_bigint_key_signed`,`col_bigint_key_unsigned`,`col_float_undef_signed`,`col_float_undef_unsigned`,`col_float_key_signed`,`col_float_key_unsigned`,`col_double_undef_signed`,`col_double_undef_unsigned`,`col_double_key_signed`,`col_double_key_unsigned`,`col_decimal(40, 20)_undef_signed`,`col_decimal(40, 20)_undef_unsigned`,`col_decimal(40, 20)_key_signed`,`col_decimal(40, 20)_key_unsigned`,`col_char(20)_undef_signed`,`col_char(20)_key_signed`,`col_varchar(20)_undef_signed`,`col_varchar(20)_key_signed`"
end

-- @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
-- .yy {}

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
curAStack = {["size"]=1, [1]=gen3FieldAliases()}

-- child aliases, set by the current node.
-- We will use it to set the aliases of each child node (Except for some special child nodes);
-- It will also be used as the column names of the current query
childAStack = {["size"]=0}

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

-- @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
-- test

print(seed)

print(math.random())
print(math.random())
print(math.random())

print(stackTop(curAStack)[1])
assert(stackTop(curAStack)[1] == "f1")
pushStack(childAStack,{"h","z","y"})
pushStack(childAStack,{"h2","z2","y2"})
print(stackTop(childAStack)[1])
assert(stackTop(childAStack)[1] == "h2")
popStack(childAStack)
print(stackTop(childAStack)[1])
assert(stackTop(childAStack)[1] == "h")

res = mySplit(_field_list())
for i,v in ipairs(res) do
    print(v, i)
end
print("==========")
res = keyFields()
for i,v in ipairs(res) do
    print(v, i)
end
print("==========")
print(optional(keyFields()))

