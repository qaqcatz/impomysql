| **category**          | **mutation name** | **upper**                                              | **lower**                                                    |
| --------------------- | ----------------- | ------------------------------------------------------ | ------------------------------------------------------------ |
| Relation              | distinct          | remove distinct                                        | add distinct                                                 |
| Relation              | union_unionall    | union→union all                                        | union all→union                                              |
| Relation              | union             | r→r union r                                            | r1 union r2→r1<br>r1 union r2→r2                             |
| Relation              | union_intersect   | r1 intersect r2→r1 union r2                            | r1 union r2→r1 intersect r2                                  |
| Relation              | minus             | r1 minus r2→r1                                         | r1→r1 minus r1                                               |
| Predicate             | false             | -                                                      | p→false                                                      |
| Predicate             | is_false          | -                                                      | p is {true,false}→false                                      |
| Predicate             | isnot_false       | -                                                      | p is not {true,false}→false                                  |
| Predicate             | true              | p→true                                                 | -                                                            |
| Predicate             | is_true           | p is {true,false}→true                                 | -                                                            |
| Predicate             | isnot_true        | p is not {true,false}→true                             | -                                                            |
| Comparison expression | le_e              | e1=e2→e1<=e2                                           | e1<=e2→e1=e2                                                 |
| Comparison expression | ge_e              | e1=e2→e1>=e2                                           | e1>=e2→e1=e2                                                 |
| Comparison expression | le_l              | e1\<e2→e1<=e2                                          | e1<=e2→e1<e2                                                 |
| Comparison expression | ge_g              | e1>e2→e1>=e2                                           | e1>=e2→e1>e2                                                 |
| Comparison expression | ne_l              | e1\<e2→e1!=e2                                          | e1!=e2→e1<e2                                                 |
| Comparison expression | ne_g              | e1>e2→e1!=e2                                           | e1!=e2→e1>e2                                                 |
| Comparison expression | any_all           | e op all r→e op any r                                  | e op any r→e op all r                                        |
| Comparison expression | like_nw           | char→_<br>char→%                                       | _→random char<br/>%→random char                              |
| Comparison expression | like_ww           | _→%                                                    | %→_                                                          |
| Comparison expression | regexp_prefix     | remove ^                                               | add ^                                                        |
| Comparison expression | regexp_suffix     | remove $                                               | add $                                                        |
| Comparison expression | regexp_nw         | char→.                                                 | .→random char                                                |
| Comparison expression | regexp_ww         | +→\*<br>?→\*                                           | \*→+<br>\*→?                                                 |
| Comparison expression | in_null           | in(e1,e2,...)→in(e1,e2,...,null)                       | in(e1,e2,...,null)→in(e1,e2,...)                             |
| Comparison expression | between           | e between e1 and e2→e1<=e<br>e between e1 and e2→e<=e2 | e between e1 and e2→e1<e and e <= e2<br>e between e1 and e2→e1<=e and e < e2 |