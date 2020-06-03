package homework3

func flatten_matrix_clockwise(matrix [][]int) []int {
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
	return flat
}
