package homework3

import "errors"

func FlattenMatrixClockwise(matrix [][]int) ([]int, error) {
	if !isValidMatrix(matrix) {
		return nil, errors.New("invalid matrix")
	}
	rowBottom := len(matrix)
	colRight := len(matrix[0])
	flat := make([]int, 0, rowBottom*colRight)
	rowTop, colLeft := 0, 0
	for rowTop < rowBottom && colLeft < colRight {
		for i := colLeft; i < colRight; i++ {
			flat = append(flat, matrix[rowTop][i])
		}
		rowTop += 1
		for i := rowTop; i < rowBottom; i++ {
			flat = append(flat, matrix[i][colRight-1])
		}
		colRight -= 1
		if rowTop < rowBottom {
			for i := colRight - 1; i > colLeft-1; i-- {
				flat = append(flat, matrix[rowBottom-1][i])
			}
			rowBottom -= 1
		}
		if colLeft < colRight {
			for i := rowBottom - 1; i > rowTop-1; i-- {
				flat = append(flat, matrix[i][colLeft])
			}
			colLeft += 1
		}
	}
	return flat, nil
}

func isValidMatrix(matrix [][]int) bool {
	if len(matrix) == 0 {
		return false
	}
	firstRowLen := len(matrix[0])
	for i := 1; i < len(matrix); i++ {
		if len(matrix[i]) != firstRowLen {
			return false
		}
	}
	return true
}
