package controllers

import "math"

// 分页结构体
type Pagination struct {
	PageIndex int
	PageSize  int
	Total     int64
	PageCount int
	Pages     []int
	PrePage   int
	NextPage  int
}

// 分页工具函数
func PageTool(pageCount int, pageIndex int) []int {
	if pageCount <= 0 {
		return []int{}
	}

	if pageCount <= 5 {
		pages := make([]int, pageCount)
		for i := range pages {
			pages[i] = i + 1
		}
		return pages
	}

	if pageIndex <= 3 {
		return []int{1, 2, 3, 4, 5}
	}

	if pageIndex >= pageCount-2 {
		return []int{pageCount - 4, pageCount - 3, pageCount - 2, pageCount - 1, pageCount}
	}

	return []int{pageIndex - 2, pageIndex - 1, pageIndex, pageIndex + 1, pageIndex + 2}
}

// 创建分页函数
func CreatePagination(count int64, pageIndex, pageSize int) *Pagination {
	pageCount := int(math.Ceil(float64(count) / float64(pageSize)))

	if pageIndex > pageCount {
		pageIndex = pageCount
	}
	if pageCount == 0 {
		pageIndex = 1
	}

	prePage := pageIndex - 1
	if prePage < 1 {
		prePage = 1
	}

	nextPage := pageIndex + 1
	if nextPage > pageCount {
		nextPage = pageCount
	}

	return &Pagination{
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Total:     count,
		PageCount: pageCount,
		Pages:     PageTool(pageCount, pageIndex),
		PrePage:   prePage,
		NextPage:  nextPage,
	}
}
