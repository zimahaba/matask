package request

import (
	"matask/internal/model"
	"strconv"
	"time"

	"github.com/lib/pq"
)

func ToTaskFilter(query map[string][]string) model.TaskFilter {
	var filter model.TaskFilter
	if len(query["name"]) > 0 {
		filter.Name = query["name"][0]
	}
	if len(query["type"]) > 0 {
		filter.Type = query["type"][0]
	}
	if len(query["started1"]) > 0 {
		started, err := time.Parse(time.DateOnly, query["started1"][0])
		if err != nil {
			panic(err)
		} else {
			filter.Started1 = pq.NullTime{Time: started, Valid: true}
		}
	}
	if len(query["started2"]) > 0 {
		started, err := time.Parse(time.DateOnly, query["started2"][0])
		if err != nil {
			panic(err)
		} else {
			filter.Started2 = pq.NullTime{Time: started, Valid: true}
		}
	}
	if len(query["ended1"]) > 0 {
		ended, err := time.Parse(time.DateOnly, query["ended1"][0])
		if err != nil {
			panic(err)
		} else {
			filter.Ended1 = pq.NullTime{Time: ended, Valid: true}
		}
	}
	if len(query["ended2"]) > 0 {
		ended, err := time.Parse(time.DateOnly, query["ended2"][0])
		if err != nil {
			panic(err)
		} else {
			filter.Ended2 = pq.NullTime{Time: ended, Valid: true}
		}
	}
	if len(query["page"]) > 0 {
		page, err := strconv.Atoi(query["page"][0])
		if err != nil {
			filter.Page = 0
		} else {
			filter.Page = page
		}
	} else {
		filter.Page = 0
	}
	if len(query["size"]) > 0 {
		size, err := strconv.Atoi(query["size"][0])
		if err != nil {
			filter.Size = 10
		} else {
			filter.Size = size
		}
	} else {
		filter.Size = 10
	}
	if len(query["sortField"]) > 0 {
		filter.SortField = query["sortField"][0]
	} else {
		filter.SortField = "id"
	}
	if len(query["sortDirection"]) > 0 {
		filter.SortDirection = query["sortDirection"][0]
	} else {
		filter.SortDirection = "ASC"
	}
	return filter
}
