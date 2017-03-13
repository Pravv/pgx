package utils

import (
    "reflect"
    "github.com/jmoiron/sqlx/reflectx"
    "errors"
    "strings"
    "sync"
    "database/sql"
)
var _scannerInterface = reflect.TypeOf((*sql.Scanner)(nil)).Elem()

func IsScannable(t reflect.Type) bool {
   // if reflect.PtrTo(t).Implements(_scannerInterface) {
  // //     return true
   // }
    if t.Kind() != reflect.Struct {
        return true
    }

    // it's not important that we use the right mapper for this particular object,
    // we're only concerned on how many exported fields this struct has
    m := Mapper()
    if len(m.TypeMap(t).Index) == 0 {
        return true
    }
    return false
}

func FieldsByTraversal(v reflect.Value, traversals [][]int, values []interface{}, ptrs bool) error {
    v = reflect.Indirect(v)
    if v.Kind() != reflect.Struct {
        return errors.New("argument not a struct")
    }

    for i, traversal := range traversals {
        if len(traversal) == 0 {
            values[i] = new(interface{})
            continue
        }
        f := reflectx.FieldByIndexes(v, traversal)
        if ptrs {
            values[i] = f.Addr().Interface()
        } else {
            values[i] = f.Interface()
        }
    }
    return nil
}

var NameMapper = strings.ToLower
var origMapper = reflect.ValueOf(NameMapper)
var mpr *reflectx.Mapper
var mprMu sync.Mutex
func Mapper() *reflectx.Mapper {
    mprMu.Lock()
    defer mprMu.Unlock()

    if mpr == nil {
        mpr = reflectx.NewMapperFunc("db", NameMapper)
    } else if origMapper != reflect.ValueOf(NameMapper) {
        // if NameMapper has changed, create a new mapper
        mpr = reflectx.NewMapperFunc("db", NameMapper)
        origMapper = reflect.ValueOf(NameMapper)
    }
    return mpr
}