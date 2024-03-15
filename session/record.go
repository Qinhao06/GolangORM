package session

import (
	"GolangORM/clause"
	"GolangORM/log"
	"errors"
	"reflect"
)

/*
   处理数据库的记录部分，也就是增删改查
*/

func (s *Session) Insert(values ...interface{}) (int64, error) {

	for _, value := range values {
		err := s.CallMethod(BeforeInsert, value)
		if err != nil {
			log.Errorf("call BeforeInsert error: %v", err)
			return 0, err
		}
	}

	recordValues := make([]interface{}, 0)
	for _, value := range values {
		table := s.Model(value).RefSchema()
		s.clause.Set(clause.Insert, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordFields(value))
	}

	s.clause.Set(clause.Values, recordValues...)
	sql, vars := s.clause.Build(clause.Insert, clause.Values)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		log.Errorf("Insert error: %v", err)
		return 0, err
	}

	for _, value := range values {
		err = s.CallMethod(AfterInsert, value)
		if err != nil {
			log.Errorf("call AfterInsert error: %v", err)
			return 0, err
		}
	}

	if err != nil {
		log.Errorf("call AfterInsert error: %v", err)
		return 0, err
	}

	return result.RowsAffected()
}

func (s *Session) Find(values interface{}) error {

	err := s.CallMethod(BeforeQuery, nil)
	if err != nil {
		log.Errorf("call BeforeFind error: %v", err)
		return err
	}

	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefSchema()

	s.clause.Set(clause.Select, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.Select, clause.Where, clause.OrderBy, clause.Limit)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		log.Errorf("Find error: %v", err)
		return err
	}

	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, field := range table.Fields {
			values = append(values, dest.FieldByName(field.Name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			log.Errorf("scan error: %v", err)
			return err
		}

		err := s.CallMethod(AfterQuery, dest.Addr().Interface())
		if err != nil {
			log.Errorf("call AfterQuery error: %v", err)
			return err
		}

		destSlice.Set(reflect.Append(destSlice, dest))

	}
	return rows.Close()
}

func (s *Session) Delete(values ...interface{}) (int64, error) {
	err := s.CallMethod(BeforeDelete, nil)
	if err != nil {
		log.Errorf("call BeforeDelete error: %v", err)
		return 0, err
	}
	s.clause.Set(clause.Delete, s.RefSchema().Name)
	sql, vars := s.clause.Build(clause.Delete, clause.Where)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		log.Errorf("Delete error: %v", err)
		return 0, err
	}

	err = s.CallMethod(AfterDelete, nil)
	if err != nil {
		log.Errorf("call AfterDelete error: %v", err)
		return 0, err
	}

	return result.RowsAffected()
}

// Update support map[string]interface{}
// also support kv list: "Name", "Tom", "Age", 18, ....
func (s *Session) Update(kv ...interface{}) (int64, error) {

	m, ok := kv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}

	for _, v := range m {
		err := s.CallMethod(BeforeUpdate, v)
		if err != nil {
			log.Errorf("call BeforeUpdate err %v", err)
			return 0, err
		}
	}

	s.clause.Set(clause.Update, s.RefSchema().Name, m)
	sql, vars := s.clause.Build(clause.Update, clause.Where)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		log.Errorf("Update error: %v", err)
		return 0, err
	}

	err = s.CallMethod(AfterUpdate, nil)

	return result.RowsAffected()
}

func (s *Session) Count(values ...interface{}) (int64, error) {
	s.clause.Set(clause.Count, s.RefSchema().Name)
	sql, vars := s.clause.Build(clause.Count, clause.Where)
	row := s.Raw(sql, vars...).QueryRow()
	var cnt int64
	if err := row.Scan(&cnt); err != nil {
		log.Errorf("Count error: %v", err)
		return 0, err
	}
	return cnt, nil
}

func (s *Session) Limit(number int) *Session {
	s.clause.Set(clause.Limit, number)
	return s
}

func (s *Session) OrderBy(filedName string) *Session {
	s.clause.Set(clause.OrderBy, filedName)
	return s
}

func (s *Session) Where(desc string, args ...interface{}) *Session {
	s.clause.Set(clause.Where, append(append([]interface{}{}, desc), args...)...)
	return s
}

func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}
