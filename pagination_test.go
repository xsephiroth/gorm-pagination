package pagination

import (
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var testDB *gorm.DB

func Test_hasNext(t *testing.T) {
	testCases := []struct {
		total  int
		offset int
		limit  int
		want   bool
	}{
		{0, 0, 10, false},
		{1, 0, 10, false},
		{10, 0, 10, false},
		{10, 10, 10, false},
		{20, 10, 10, false},
		{11, 0, 10, true},
		{11, 10, 10, false},
	}

	for _, tc := range testCases {
		got := hasNext(tc.total, tc.offset, tc.limit)
		if got != tc.want {
			t.Errorf("total: %d, offset: %d, limit: %d, want: %v, got: %v",
				tc.total, tc.offset, tc.limit, tc.want, got)
		}
	}
}

func Test_hasPrev(t *testing.T) {
	testCases := []struct {
		totalPage int
		page      int
		want      bool
	}{
		{0, 1, false},
		{1, 1, false},
		{2, 1, false},
		{2, 2, true},
	}
	for _, tc := range testCases {
		got := hasPrev(tc.totalPage, tc.page)
		if got != tc.want {
			t.Errorf("total_page: %d, page: %d, want: %v, got: %v",
				tc.totalPage, tc.page, tc.want, got)
		}
	}
}

func Test_totalPage(t *testing.T) {
	testCases := []struct {
		total int
		limit int
		want  int
	}{
		{0, 10, 0},
		{1, 10, 1},
		{10, 10, 1},
		{11, 10, 2},
	}

	for _, tc := range testCases {
		got := totalPage(tc.total, tc.limit)
		if got != tc.want {
			t.Errorf("total: %d, limit: %d, want: %d, got: %d",
				tc.total, tc.limit, tc.want, got)
		}
	}
}

func TestPagePagination(t *testing.T) {
	for i := 0; i < 11; i++ {
		a := &A{
			Test: i,
		}
		err := testDB.Create(a).Error
		if err != nil {
			t.Error(err)
		}
	}

	db1 := testDB.Model(A{})
	db2 := testDB.Model(A{}).Where("test > ?", 4) // should be [5...10]

	testCases := []struct {
		ReqDB      *gorm.DB
		ReqPage    int
		ReqLimit   int
		TotalPage  int
		Total      int
		Page       int
		Prev       bool
		Next       bool
		ResultsLen int
	}{
		{db1, 0, 5, 3, 11, 1, false, true, 5}, // req page is 0, should correct page to 1
		{db1, 1, 5, 3, 11, 1, false, true, 5}, // req page is 1
		{db1, 2, 5, 3, 11, 2, true, true, 5},  // req page is 2
		{db1, 3, 5, 3, 11, 3, true, false, 1}, // req page is 3, should got 1 result
		{db1, 4, 5, 3, 11, 3, true, false, 1}, // req page is out of 3, should correct page to 3

		{db2, 0, 5, 2, 6, 1, false, true, 5}, // req page is 0, should correct page to 1
		{db2, 1, 5, 2, 6, 1, false, true, 5}, // req page is 1
		{db2, 2, 5, 2, 6, 2, true, false, 1}, // req page is 2
		{db2, 3, 5, 2, 6, 2, true, false, 1}, // req page is 3, should correct page to 1
	}

	for _, tc := range testCases {
		out := make([]A, 0)
		pp := PagePagination(tc.ReqDB, tc.ReqPage, tc.ReqLimit, &out)

		if pp.TotalPage != tc.TotalPage {
			t.Errorf("p: %#v, TotalPage want: %d, got: %d", pp, tc.TotalPage, pp.TotalPage)
		}
		if pp.Total != tc.Total {
			t.Errorf("p: %#v, Total want: %d, got: %d", pp, tc.Total, pp.Total)
		}
		if pp.Page != tc.Page {
			t.Errorf("p: %#v, Page want: %d, got: %d", pp, tc.Page, pp.Page)
		}
		if pp.Next != tc.Next {
			t.Errorf("p: %#v, want next: %v, got: %v", pp, tc.Next, pp.Next)
		}
		if pp.Prev != tc.Prev {
			t.Errorf("p: %#v, want prev: %v, got: %v", pp, tc.Prev, pp.Prev)
		}

		results := pp.Results.(*[]A)
		if len(*results) != tc.ResultsLen {
			t.Errorf("p: %#v, results should be %d, got: %d", pp, tc.ResultsLen, len(*results))
		}
	}

}

type A struct {
	gorm.Model
	Test int
}

func createTestSqlite() (db *gorm.DB) {
	var err error
	db, err = gorm.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	return
}
func TestMain(m *testing.M) {
	db := createTestSqlite()
	defer db.Close()

	testDB = db
	testDB.CreateTable(&A{})

	r := m.Run()
	os.Exit(r)
}
