package task

var IgnoreErrors = map[int]string {
	1064 : "You have an error in your SQL syntax",
	1267 : "Illegal mix of collations xxx and xxx for operation xxx",
	1271 : "Illegal mix of collations for operation xxx",
	1690 : "value out of range",
}