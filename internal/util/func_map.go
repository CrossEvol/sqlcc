package util

import "github.com/crossevol/sqlcc/internal/common"

func IsID(columnPair common.ColumnPair) bool {
	return columnPair.IsAutoIncrement || columnPair.IsPrimaryKey
}

func IsNotID(columnPair common.ColumnPair) bool {
	return !columnPair.IsAutoIncrement && !columnPair.IsPrimaryKey
}
