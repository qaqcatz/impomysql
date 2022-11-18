package stage1

import (
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
)

// rmUncertain:
// todo: remove uncertain functions
//
//  rand
//
//  curdate
//  current_date
//  current_time
//  current_timestamp
//  curtime
//  localtime
//  localtimestamp
//  now
//  sysdate
//  utc_date
//  utc_time
//  utc_timestamp
//
//  benchmark
//  CONNECTION_ID
//  CURRENT_USER
//  CURRENT_ROLE
//  DATABASE
//  FOUND_ROWS
//  LAST_INSERT_ID
//  ROW_COUNT
//  SCHEMA
//  SESSION_USER
//  SYSTEM_USER
//  USER
//
//  any_value
//  master_pos_wait
//  sleep
//  uuid
//  uuid_short
//  uuid_to_bin
//  bin_to_uuid
//
//  random_bytes
func rmUncertain(in ast.Node) bool {
	return false
}