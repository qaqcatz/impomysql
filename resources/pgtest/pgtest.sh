../go-randgen gentest -Z ./post_num.zz.lua -Y ./post_num.yy -Q 10 --seed 1 -B

sed -i '/^key (/d' output.data.sql
sed -i 's/ double/ float/g' output.data.sql
sed 's/'\`'/"/g' output.data.sql > tmp.data.sql
mv tmp.data.sql output.data.sql
sed -i '10s/.*/\"col_decimal(40, 20)_key_signed\" decimal(40, 20)  /' output.data.sql
sed -i '22s/.*/\"col_decimal(40, 20)_key_signed\" decimal(40, 20)  /' output.data.sql

sed 's/'\`'/"/g' output.rand.sql > tmp.rand.sql
mv tmp.rand.sql output.rand.sql