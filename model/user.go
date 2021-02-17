package model

import (
	"github.com/snail007/gmc"
	gcore "github.com/snail007/gmc/core"
	gmap "github.com/snail007/gmc/util/map"
	"strconv"
	"sync"
)

type UserModel struct {
	db    gcore.Database
	table string
	pkey  string
	once  *sync.Once
}

var (
	User = NewUserModel()
)

func NewUserModel() *UserModel {
	u := &UserModel{
		table: "user",
		pkey:  "user_id",
		once:  &sync.Once{},
	}
	return u
}

func (s *UserModel) DB() gcore.Database {
	if s.db == nil {
		s.once.Do(func() {
			s.db = gmc.DB.DB()
		})
	}
	return s.db
}

func (s *UserModel) GetByID(id string) (ret gmap.Mss, error error) {
	db := s.DB()
	rs, err := db.Query(db.AR().From(s.table).Where(gmap.M{
		s.pkey: id,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	ret = rs.Row()
	return
}

func (s *UserModel) GetBy(where gmap.M) (ret gmap.Mss, error error) {
	db := s.DB()
	rs, err := db.Query(db.AR().From(s.table).Where(where).Limit(0, 1))
	if err != nil {
		return
	}
	ret = rs.Row()
	return
}

func (s *UserModel) MGetByIDs(ids []string, orderBy ...interface{}) (ret gmap.Mss, error error) {
	db := s.DB()
	ar := db.AR().From(s.table).Where(gmap.M{
		s.pkey: ids,
	})
	if col, by := s.OrderBy(orderBy...); col != "" {
		ar.OrderBy(col, by)
	}
	rs, err := db.Query(ar)
	if err != nil {
		return
	}
	ret = rs.Row()
	return
}

func (s *UserModel) MGetBy(where gmap.M, orderBy ...interface{}) (ret []gmap.Mss, error error) {
	db := s.DB()
	ar := db.AR().From(s.table).Where(where).Limit(0, 1)
	if col, by := s.OrderBy(orderBy...); col != "" {
		ar.OrderBy(col, by)
	}
	rs, err := db.Query(ar)
	if err != nil {
		return
	}
	ret = rs.Rows()
	return
}

func (s *UserModel) DeleteBy(where gmap.M) (cnt int64, err error) {
	db := s.DB()
	rs, err := db.Exec(db.AR().Delete(s.table, where))
	if err != nil {
		return
	}
	cnt = rs.RowsAffected()
	return
}

func (s *UserModel) DeleteByIDs(ids []string) (cnt int64, err error) {
	db := s.DB()
	rs, err := db.Exec(db.AR().Delete(s.table, gmap.M{
		s.pkey: ids,
	}))
	if err != nil {
		return
	}
	cnt = rs.RowsAffected()
	return
}

func (s *UserModel) Insert(data gmap.M) (cnt int64, err error) {
	db := s.DB()
	rs, err := db.Exec(db.AR().Insert(s.table, data))
	if err != nil {
		return
	}
	cnt = rs.RowsAffected()
	return
}

func (s *UserModel) InsertBatch(data []gmap.M) (cnt int64, err error) {
	db := s.DB()
	rs, err := db.Exec(db.AR().InsertBatch(s.table, data))
	if err != nil {
		return
	}
	cnt = rs.RowsAffected()
	return
}

func (s *UserModel) UpdateByIDs(ids []string, data gmap.M) (cnt int64, err error) {
	db := s.DB()
	rs, err := db.Exec(db.AR().Update(s.table, data, gmap.M{
		s.pkey: ids,
	}))
	if err != nil {
		return
	}
	cnt = rs.RowsAffected()
	return
}

func (s *UserModel) UpdateBy(where, data gmap.M) (cnt int64, err error) {
	db := s.DB()
	rs, err := db.Exec(db.AR().Update(s.table, data, where))
	if err != nil {
		return
	}
	cnt = rs.RowsAffected()
	return
}

func (s *UserModel) Page(where gmap.M, offset, length int, orderBy ...interface{}) (ret []gmap.Mss, total int, err error) {
	db := s.DB()
	ar := db.AR().Select("count(*) as total").From(s.table)
	if len(where) > 0 {
		ar.Where(where)
	}
	rs, err := db.Query(ar)
	if err != nil {
		return
	}
	t, _ := strconv.ParseInt(rs.Value("total"), 10, 64)
	total = int(t)
	ar = db.AR().From(s.table).Where(where).Limit(offset, length)
	if len(where) > 0 {
		ar.Where(where)
	}
	if col, by := s.OrderBy(orderBy...); col != "" {
		ar.OrderBy(col, by)
	}
	rs, err = db.Query(ar)
	if err != nil {
		return
	}
	ret = rs.Rows()
	return
}

func (s *UserModel) OrderBy(orderBy ...interface{}) (col, by string) {
	if len(orderBy) > 0 {
		switch val := orderBy[0].(type) {
		case gmap.M:
			for k, v := range val {
				col, by = k, v.(string)
				break
			}
		case gmap.Mss:
			for k, v := range val {
				col, by = k, v
				break
			}
		}
	}
	return
}
