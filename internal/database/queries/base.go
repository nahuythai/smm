package queries

import (
	"fmt"
	"regexp"
	"smm/pkg/logging"
	"smm/pkg/request"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	logger = logging.GetLogger()
)

type queryOption struct {
	skip       int64
	limit      int64
	onlyFields bson.M
	sort       bson.D
}

type QueryOption interface {
	SetPagination(p *request.Pagination)
	QuerySkip() *int64
	QueryLimit() *int64
	QueryOnlyField() bson.M
	SetOnlyField(fields ...string)
	AddSort(map[string]int)
	QuerySort() bson.D
}

func NewOption() QueryOption {
	return &queryOption{
		onlyFields: make(bson.M),
		sort:       make(bson.D, 0),
	}
}

func (o *queryOption) SetPagination(p *request.Pagination) {
	o.skip = p.Skip
	o.limit = p.Limit
}

func (o *queryOption) QuerySkip() *int64 {
	return &o.skip
}

func (o *queryOption) QueryLimit() *int64 {
	return &o.limit
}

func (o *queryOption) SetOnlyField(fields ...string) {
	for _, field := range fields {
		o.onlyFields[field] = 1
	}
}

func (o *queryOption) QueryOnlyField() bson.M {
	return o.onlyFields
}

func (o *queryOption) AddSort(sortMap map[string]int) {
	for field, sortType := range sortMap {
		if sortType != QuerySortTypeAsc && sortType != QuerySortTypeDesc {
			sortType = QuerySortTypeDesc
		}
		o.sort = append(o.sort, bson.E{Key: field, Value: sortType})
	}
}

func (o *queryOption) QuerySort() bson.D {
	return o.sort
}

const (
	QueryFilterMethodEqual = "$eq"
	QueryFilterMethodRegex = "$regex"
)

const (
	QuerySortTypeAsc  = 1
	QuerySortTypeDesc = -1
)

type FilterOption interface {
	BuildAndQuery() bson.M
	AddFilter(filter ...Filter)
}

type Filter struct {
	Value  interface{}
	Method string
	Field  string
}

type filterOption struct {
	filters []Filter
}

func NewFilterOption() FilterOption {
	return &filterOption{
		filters: make([]Filter, 0),
	}
}

func (f *filterOption) BuildAndQuery() bson.M {
	if len(f.filters) == 0 {
		return bson.M{}
	}
	query := make([]bson.M, 0, len(f.filters))
	for _, filter := range f.filters {
		if filter.Method == QueryFilterMethodRegex {
			filter.Value = primitive.Regex{Pattern: regexp.QuoteMeta(fmt.Sprintf("%v", filter.Value)), Options: "i"}
		}
		query = append(query, bson.M{
			filter.Field: bson.M{filter.Method: filter.Value},
		})
	}
	return bson.M{"$and": query}
}

func (f *filterOption) BuildOrQuery() bson.M {
	if len(f.filters) == 0 {
		return bson.M{}
	}
	query := make([]bson.M, 0, len(f.filters))
	for _, filter := range f.filters {
		if filter.Method == QueryFilterMethodRegex {
			filter.Value = regexp.QuoteMeta(fmt.Sprintf("%v", filter.Value))
		}
		query = append(query, bson.M{
			filter.Field: bson.M{filter.Method: filter.Value},
		})
	}
	return bson.M{"$or": query}
}

func (f *filterOption) AddFilter(filter ...Filter) {
	f.filters = append(f.filters, filter...)
}
