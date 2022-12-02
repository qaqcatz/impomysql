tables = {
    rows = {3, 7},
    -- SHOW CHARACTER SET;
    charsets = {'utf8'},
    -- partition number
    -- partitions = {2},
}

--tables = {
--    rows = {10, 30},
--    -- SHOW CHARACTER SET;
--    charsets = {'utf8', 'latin1', 'binary'},
--    -- partition number
--    partitions = {4, 6, 'undef'},
--}

fields = {
    types = {'bigint', 'float', 'double', 'decimal(40, 20)',
             'char(20)', 'varchar(20)'},
    sign = {'signed', 'unsigned'}
}

data = {
    numbers = {'tinyint', 'smallint',
               '12.991', '1.009', '-9.183', '2', '-2', '1', '-1', '0', '-0', '0.0001', '-0.0001',
               'decimal',
    },
    strings = {'letter', 'english', '\'0\'', '\'-0\'', '\'1\'', '\'-1\'', '\'3\n\'', '\'3\t\'', '\'\n3\'', '\'\t3\''},
}