package database

import (
	"strconv"
	"strings"
)

func generateInsertPlaceholder(columnCount int, rowCount int, initialParameterNumber int) string {
	rows := make([]string, 0, rowCount)
	parameterCount := initialParameterNumber
	for i := 0; i < rowCount ; i++ {
		rowParameters := make([]string, 0, columnCount)
		for j := 0; j < columnCount ; j++ {
			rowParameters = append(rowParameters, "$" + strconv.Itoa(parameterCount))
			parameterCount++
		}
		rowPlaceholder := "(" + strings.Join(rowParameters, ", ") + ")"
		rows = append(rows, rowPlaceholder)
	}

	return strings.Join(rows, ", ")
}

func generateTupleFilterPlaceholder(columnList []string, tupleCount int) string {
	orFilterList := make([]string, 0, tupleCount)
	parameterCount := 1
	for i := 0; i < tupleCount; i++ {
		andFilterList := make([]string, 0, len(columnList))
		for _, column := range columnList {
			andFilterList = append(andFilterList, column + " = " + "$" + strconv.Itoa(parameterCount))
			parameterCount++
		}
		andFilter := "(" + strings.Join(andFilterList, " AND ") + ")"
		orFilterList = append(orFilterList, andFilter)
	}

	return strings.Join(orFilterList, " OR ")
}
