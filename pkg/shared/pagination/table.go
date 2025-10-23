package pagination

import (
	"errors"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type TableRequest struct {
	SortBy        string
	SortDir       string
	Query         string
	QueryField    []string
	Request       PaginationRequestDto
	Data          interface{}
	Filter        interface{}
	AllowedFields []string
	MapFields     map[string]string
}

type Table struct {
	db *gorm.DB
}

func NewTable(db *gorm.DB) *Table {
	return &Table{
		db,
	}
}

func (t *Table) Pagination(query func(interface{}) *gorm.DB, in interface{}) (interface{}, error) {
	req := in.(*TableRequest)
	req.Request.GenerateFilter()
	data := req.Data
	req.AllowedFields = append(req.AllowedFields, "id", "organization_id", "created_at")

	var (
		wg       sync.WaitGroup
		count    = make(chan int64, 1)
		results  = make(chan interface{}, 1)
		errQuery = make(chan error, 2)
		offset   = req.Request.GetPage() * req.Request.GetSize()
	)
	// err := sortValidation(req)
	// if err != nil {
	// 	return nil, err
	// }
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		var cnt int64
		q := query(req.Request)
		q = whereConditions(q, req)
		q, err = whereFilterConditions(q, req)
		if err != nil {
			errQuery <- err
			return
		}

		err := q.Count(&cnt).Error
		if err != nil {
			errQuery <- err
			return
		}
		count <- cnt

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		q := query(req.Request)
		q = whereConditions(q, req)
		q, err = whereFilterConditions(q, req)
		if err != nil {
			errQuery <- err
			return
		}

		if req.Request.GetSize() > 0 {
			q = q.Offset(int(offset)).Limit(int(req.Request.GetSize()))
		}
		if req.Request.GetSortBy() != "" && req.Request.GetSortDir() != "" {
			q = q.Order(getMappingField(req, req.Request.GetSortBy()) + " " + req.Request.GetSortDir())
		}
		err := q.Find(data).Error
		if err != nil {
			errQuery <- err
			return
		}

		results <- data
	}()

	go func() {
		wg.Wait()
		close(count)
		close(results)
		close(errQuery)
	}()

	for err := range errQuery {
		log.Errorf("Pagination err %v", err)
		return nil, errors.New("dbError")
	}

	totalData := <-count
	return &ResultPagination{
		Data:        <-results,
		Page:        req.Request.GetPage(),
		RowsPerPage: req.Request.GetSize(),
		Count:       totalData,
		TotalPages:  calculateTotalPages(int(totalData), req.Request.GetSize()),
	}, nil
}

func getMappingField(req *TableRequest, field string) string {
	value, exists := req.MapFields[field]

	if !exists {
		return field
	}
	return value
}

// Function untuk menghitung total halaman
func calculateTotalPages(count, rowsPerPage int) int {
	if count == 0 || rowsPerPage == 0 {
		return 0
	}
	totalPages := count / rowsPerPage
	if count%rowsPerPage != 0 {
		totalPages++
	}
	return totalPages
}

func sortValidation(req *TableRequest) error {
	if req.SortBy == "" {
		return errors.New("dbErrorSortBy")
	}
	if req.SortDir == "" {
		return errors.New("dbErrorSortDir")
	}
	return nil
}

func whereConditions(db *gorm.DB, req *TableRequest) *gorm.DB {
	if len(req.Request.GetQuery()) == 0 {
		return db
	}
	condStr := []string{}
	values := make([]interface{}, 0)
	for _, v := range req.QueryField {
		condStr = append(condStr, v+" iLIKE ?")
		values = append(values, "%"+req.Request.GetQuery()+"%")
	}
	return db.Where(strings.Join(condStr, " OR "), values...)
}

func whereFilterConditions(db *gorm.DB, req *TableRequest) (*gorm.DB, error) {
	filters := req.Request.GetFilter()
	if filters == nil {
		return db, nil
	}
	for _, v := range filters {
		operator, err := getOperator(v.Op)
		if err != nil {
			return nil, err
		}
		err = validateAllowedFields(v.Field, req)
		if err != nil {
			return nil, err
		}
		db = db.Where(getMappingField(req, v.Field)+" "+operator+" ?", v.Val)
	}
	return db, nil
}

func validateAllowedFields(field string, req *TableRequest) error {
	isExists := false
	for _, v := range req.AllowedFields {
		if v == field {
			isExists = true
			break
		}
	}
	if !isExists {
		log.Errorf("Field %s not allowed", field)
		return errors.New("dbErrorFieldNotAllowed")
	}
	return nil
}

func getOperator(op string) (string, error) {
	switch op {
	case "eq":
		return "=", nil
	case "lte":
		return "<=", nil
	case "lt":
		return "<", nil
	case "gte":
		return ">=", nil
	case "gt":
		return ">", nil
	}
	return "", errors.New("dbErrorOp")
}
