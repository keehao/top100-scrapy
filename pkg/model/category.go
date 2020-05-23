package model

import (
	"fmt"
	"strings"
	"top100-scrapy/pkg/preference"
)

type CategoryRow struct {
	ID       int
	Name     string
	URL      string
	Path     string
	ParentID int
}

func FetchCategoryRow(id int, opts *preference.Options) (*CategoryRow, error) {
	row := new(CategoryRow)
	stmt := fmt.Sprintf("select id, name, url, path, parent_id from categories where id = %d", id)
	err := opts.DB.QueryRow(stmt).Scan(&row.ID, &row.Name, &row.URL, &row.Path, &row.ParentID)
	return row, err
}

func BulkilyInsertCategories(set []*CategoryRow, opts *preference.Options) error {
	valueStrings := make([]string, 0, len(set))
	valueArgs := make([]interface{}, 0, len(set)*4)
	for i, item := range set {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		valueArgs = append(valueArgs, item.Name)
		valueArgs = append(valueArgs, item.Path)
		valueArgs = append(valueArgs, item.URL)
		valueArgs = append(valueArgs, item.ParentID)
	}
	stmt := fmt.Sprintf("INSERT INTO categories (name, path, url, parent_id) VALUES %s", strings.Join(valueStrings, ","))
	_, err := opts.DB.Exec(stmt, valueArgs...)
	return err
}

func BuildPath(n int, parent *CategoryRow) (path string) {
	if n < 10 {
		path = fmt.Sprintf("%s.0%d", parent.Path, n)
	} else {
		path = fmt.Sprintf("%s.%d", parent.Path, n)
	}
	return path
}