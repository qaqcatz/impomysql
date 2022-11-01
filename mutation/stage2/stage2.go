package stage2

import (
	"errors"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/qaqcatz/impomysql/connector"
	"strconv"
)

// CalCandidates: see MutateVisitor
func CalCandidates(sql string) (*MutateVisitor, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, errors.New("CalCandidates: " + err.Error())
	}
	if stmtNodes == nil || len(stmtNodes) == 0 {
		return nil, errors.New("CalCandidates: stmtNodes == nil || len(stmtNodes) == 0")
	}
	rootNode := &stmtNodes[0]
	v := &MutateVisitor{
		Root: *rootNode,
		Candidates: make(map[string][]*Candidate)}
	v.visit(*rootNode, 1)
	return v, nil
}

// ImpoMutate: you can choose any candidate to mutate, each mutation has no side effects.
func ImpoMutate(rootNode ast.Node, candidate *Candidate, seed int64) (string, error) {
	var sql []byte = nil
	var err error = nil
	switch candidate.MutationName {
	case FixMDistinctU:
		sql, err = doFixMDistinctU(rootNode, candidate.Node)
	case FixMDistinctL:
		sql, err = doFixMDistinctL(rootNode, candidate.Node)
	case FixMCmpOpU:
		sql, err = doFixMCmpOpU(rootNode, candidate.Node)
	case FixMCmpOpL:
		sql, err = doFixMCmpOpL(rootNode, candidate.Node)
	case FixMUnionAllU:
		sql, err = doFixMUnionAllU(rootNode, candidate.Node)
	case FixMUnionAllL:
		sql, err = doFixMUnionAllL(rootNode, candidate.Node)
	case FixMInNullU:
		sql, err = doFixMInNullU(rootNode, candidate.Node)
	case FixMWhere1U:
		sql, err = doFixMWhere1U(rootNode, candidate.Node)
	case FixMWhere0L:
		sql, err = doFixMWhere0L(rootNode, candidate.Node)
	case FixMHaving1U:
		sql, err = doFixMHaving1U(rootNode, candidate.Node)
	case FixMHaving0L:
		sql, err = doFixMHaving0L(rootNode, candidate.Node)
	case FixMOn1U:
		sql, err = doFixMOn1U(rootNode, candidate.Node)
	case FixMOn0L:
		sql, err = doFixMOn0L(rootNode, candidate.Node)
	case FixMRmUnionAllL:
		sql, err = doFixMRmUnionAllL(rootNode, candidate.Node)
	case RdMLikeU:
		sql, err = doRdMLikeU(rootNode, candidate.Node, seed)
	case RdMLikeL:
		sql, err = doRdMLikeL(rootNode, candidate.Node, seed)
	case RdMRegExpU:
		sql, err = doRdMRegExpU(rootNode, candidate.Node, seed)
	case RdMRegExpL:
		sql, err = doRdMRegExpL(rootNode, candidate.Node, seed)
	}
	if err != nil {
		return "", errors.New("ImpoMutate: " +  err.Error())
	}
	return string(sql), nil
}

// ImpoMutateAndExec: ImpoMutate + exec.
func ImpoMutateAndExec(rootNode ast.Node, candidate *Candidate, seed int64,
	conn connector.Connector) (string, *connector.Result, error) {
	sql, err := ImpoMutate(rootNode, candidate, seed)
	if err != nil {
		return "", nil, errors.New("ImpoMutateAndExec: " + err.Error())
	}
	result := conn.ExecSQLS(sql)
	return sql, result, nil
}

// MutateResult: slice of (mutation name, mutated sql, isUpper, error)
//
// IsUppers: Does the theoretical execution result of
// the current mutated statement increase?
type MutateResult struct {
	MutNames []string
	MutSqls  []string
	IsUppers []bool // (Candidate.U^Candidate.Flag)^1) == 1
	MutErrs  []error
	Err      error

	ExecResults []*connector.Result // exec MutSqls, nil if MutErrs[i] != nil
}

func (mutateResult *MutateResult) ToString() string {
	res := ""
	for i, mutName := range mutateResult.MutNames {
		if i != 0 {
			res += "\n"
		}
		res += "["+strconv.Itoa(i)+"]==========\n"
		res += "[MutName] " + mutName + "\n"
		if mutateResult.MutErrs[i] != nil {
			res += "[MutErr] " + mutateResult.MutErrs[i].Error()
		} else {
			res += "[MutSql] " + mutateResult.MutSqls[i] + "\n"
			res += "[IsUpper] " + strconv.FormatBool(mutateResult.IsUppers[i])
		}
	}
	return res
}

// MutateAll: For the input sql, try all of its mutation points.
// We will save the mutated sqls into *MutateResult.
func MutateAll(sql string, seed int64) *MutateResult {
	mutateResult := &MutateResult {
		MutNames: make([]string, 0),
		MutSqls:  make([]string, 0),
		IsUppers: make([]bool, 0),
		MutErrs:  make([]error, 0),
		Err:      nil,
	}

	v, err := CalCandidates(sql)
	if err != nil {
		mutateResult.Err = err
		return mutateResult
	}

	root := v.Root
	for mutationName, candidateList := range v.Candidates {
		for _, candidate := range candidateList {
			mutateResult.MutNames = append(mutateResult.MutNames, mutationName)
			mutateResult.IsUppers = append(mutateResult.IsUppers, ((candidate.U^candidate.Flag)^1) == 1)
			newSql, err := ImpoMutate(root, candidate, seed)
			mutateResult.MutErrs = append(mutateResult.MutErrs, err)
			mutateResult.MutSqls = append(mutateResult.MutSqls, newSql)
		}
	}

	return mutateResult
}

// MutateAllAndExec: MutateAll and exec.
func MutateAllAndExec(sql string, seed int64, conn *connector.Connector) *MutateResult {
	mutateResult := MutateAll(sql, seed)
	if mutateResult.Err != nil {
		return mutateResult
	}
	mutateResult.ExecResults = make([]*connector.Result, 0)
	for i, sqlm := range mutateResult.MutSqls {
		if mutateResult.MutErrs[i] != nil {
			mutateResult.ExecResults = append(mutateResult.ExecResults, nil)
		} else {
			result := conn.ExecSQLS(sqlm)
			mutateResult.ExecResults = append(mutateResult.ExecResults, result)
		}
	}
	return mutateResult
}
