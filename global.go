package gorms

import (
	"path/filepath"
	"runtime"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var GlobalDB = make(Manager)
var here string

func init() {
	_, here, _, _ = runtime.Caller(0)
}

// Caller returns the file path where this function is called.
//
// The `skip` parameter has a default value 1.
func Caller(skip int) (file string) {
	for file = here; file == here; skip++ {
		_, file, _, _ = runtime.Caller(skip)
	}
	return
}

// Get returns the DB represented by the caller path.
func Get() *DB {
	return GlobalDB.Get(Caller(2))
}

// Relative returns the DB represented by the path relative to the caller path.
func Relative(elem string) *DB {
	return GlobalDB.Get(filepath.Join(Caller(2), elem))
}

// Link links the DB represented by the relative path to the caller path.
func Link(relative string) {
	origin := Caller(2)
	GlobalDB[origin] = GlobalDB.Get(filepath.Join(origin, relative))
}

func SetGormDB(gormDB *gorm.DB) (db *DB) {
	db = &DB{
		DB:   gormDB,
		file: Caller(2),
	}
	GlobalDB.Set(db, db.PkgPath())
	return
}

func SetDialector(dialector gorm.Dialector, opts ...gorm.Option) *DB {
	gormDB, err := gorm.Open(dialector, opts...)
	if err != nil {
		panic(err)
	}
	return SetGormDB(gormDB)
}

func SetSQLite(file string, opts ...gorm.Option) *DB {
	return SetDialector(sqlite.Open(file), opts...)
}

func get() *DB {
	return GlobalDB.Get(Caller(3))
}

func AutoMigrate(dst ...any) *DB {
	return get().AutoMigrate(dst...)
}

func Close() error {
	return get().Close()
}

func Select(fields ...string) *DB {
	return get().Select(fields...)
}

func Exists[T any](conds ...any) bool {
	return get().First(new(T), conds...).Found()
}

func FirstOrCreate[T any](x *T, conds ...any) caller {
	return get().FirstOrCreate(x, conds...)
}

func returning[V any](db *DB) (v V) {
	switch any(v).(type) {
	case *DB:
		return any(db).(V)
	case bool:
		return any(db.Found()).(V)
	case nil:
		switch any(&v).(type) {
		case *error:
			if err := db.Error(); err != nil {
				return err.(V)
			}
			return
		}
	}
	panic("error return value type")
}

func First[T any, V any](conds ...any) (*T, V) {
	dest := new(T)
	return dest, returning[V](get().First(dest, conds...))
}

func MustFirst[T any](conds ...any) (dest *T) {
	dest = new(T)
	get().First(dest, conds...)
	return
}

func Find[T any, V any](conds ...any) (dest []T, v V) {
	v = returning[V](get().Find(&dest, conds...))
	return
}

func MustFind[T any](conds ...any) (dest []T) {
	get().Find(&dest, conds...)
	return
}

func Preload[T any, V any](conds ...any) (*T, V) {
	dest := new(T)
	return dest, returning[V](get().Preload(dest, conds...))
}

func MustPreload[T any](conds ...any) (dest *T) {
	dest = new(T)
	get().Preload(dest, conds...)
	return
}

func Preloads[T any, V any](conds ...any) (dest []T, v V) {
	v = returning[V](get().Preloads(&dest, conds...))
	return
}

func MustPreloads[T any](conds ...any) (dest []T) {
	get().Preloads(&dest, conds...)
	return
}
