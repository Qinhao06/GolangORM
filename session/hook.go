package session

import (
	"GolangORM/log"
	"reflect"
)

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
)

func (s *Session) CallMethod(MethodName string, value interface{}) (err error) {
	model := s.refSchema.Model
	fm := reflect.ValueOf(model).MethodByName(MethodName)
	if value != nil {
		fm = reflect.ValueOf(value).MethodByName(MethodName)
	}
	if !fm.IsValid() {
		return
	}
	param := []reflect.Value{reflect.ValueOf(s)}
	if v := fm.Call(param); len(v) > 0 {
		if err, ok := v[0].Interface().(error); ok {
			log.Errorf("call method %s error: %v", MethodName, err)
			return err
		}
	}
	return nil
}
