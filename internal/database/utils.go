package database

import "fmt"

var orderQuery = " ORDER BY %s %s "

var sortDirectionMap = map[string]string{
	"ASC":  "ASC",
	"DESC": "DESC",
}

func getOrderQuery(filterSortField string, filterSortDirection string, sortFieldMap map[string]string) string {
	sortField := sortFieldMap[filterSortField]
	if sortField == "" {
		sortField = "id"
	}
	sortDirection := sortDirectionMap[filterSortDirection]
	if sortDirection == "" {
		sortDirection = "ASC"
	}
	return fmt.Sprintf(orderQuery, sortField, sortDirection)
}

func getOffset(page int, size int) int {
	var offset int
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * size
	}
	return offset
}

func calculateTotalPages(count int, size int) int {
	totalPages := count / size
	remainder := count % size
	if totalPages == 0 || (totalPages > 0 && remainder > 0) {
		totalPages++
	}
	return totalPages
}
