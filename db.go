package gorms

import (
	"errors"
	"path/filepath"

	"golang.org/x/exp/constraints"

	"gorm.io/gorm"
)

type Model[T constraints.Integer] struct {
	ID T `gorm:"primaryKey;autoIncrement" json:"id" yaml:"id" form:"id"`
}

// Superset of gorm.DB
type DB struct {
	*gorm.DB
	// file path where this DB is created
	file string
	// error occurred while calling method
	err error
}

func (db *DB) File() string {
	return db.file
}

// PkgPath returns package path where this DB is created.
func (db *DB) PkgPath() string {
	return filepath.Dir(db.file)
}

func (db *DB) Error() error {
	return db.err
}

func (db *DB) Succeed() bool {
	return db.err == nil
}

// Found returns whether record is found.
func (db *DB) Found() bool {
	return !errors.Is(db.err, gorm.ErrRecordNotFound)
}

func (db *DB) AutoMigrate(dst ...any) *DB {
	Ref.Init(dst...)
	db.err = db.DB.AutoMigrate(dst...)
	return db
}

// Close fetch the original sql.DB and attempt to close it.
//
// It returns nil if the DB has been closed successfully.
//
// Otherwise an error will be returned.
func (db *DB) Close() error {
	if db.DB == nil {
		return nil
	}
	// fetch original sql.DB
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	// close DB
	err = sqlDB.Close()
	if err != nil {
		return err
	}
	// remove
	db.DB = nil
	return nil
}

func (db *DB) GetInstance(r *gorm.DB) *DB {
	return &DB{
		DB:   r,
		file: db.file,
		err:  r.Error,
	}
}

func (db *DB) First(dest any, conds ...any) *DB {
	db.err = db.DB.First(dest, conds...).Error
	return db
}

func (db *DB) Find(dest any, conds ...any) *DB {
	db.err = db.DB.Find(dest, conds...).Error
	return db
}

func (db *DB) Select(fields ...string) (tx *DB) {
	tx = db.GetInstance(db.DB.Select(fields))
	if tx.err != nil {
		panic(tx.err)
	}
	return
}

func (db *DB) Debug() *DB {
	return db.GetInstance(db.DB.Debug())
}

func (db *DB) Do(fn func(db *DB)) *DB {
	fn(db)
	return db
}

type caller bool

func (c caller) First(fn func()) caller {
	if bool(c) {
		fn()
	}
	return c
}

func (c caller) Create(fn func()) caller {
	if !bool(c) {
		fn()
	}
	return c
}

func (db *DB) FirstOrCreate(x any, conds ...any) caller {
	if db.First(x, conds...).Found() {
		return true
	}
	db.Create(x)
	return false
}

func (db *DB) PreloadType(in any, m map[string][]any) (tx *DB) {
	tx = db.GetInstance(db.DB.Model(in))
	if tx.Statement.Preloads == nil {
		if m != nil {
			tx.Statement.Preloads = m
		} else {
			tx.Statement.Preloads = make(map[string][]any)
		}
	} else {
		for k, v := range m {
			tx.Statement.Preloads[k] = v
		}
	}
	for _, key := range Ref.Get(in) {
		if _, ok := tx.Statement.Preloads[key]; !ok {
			tx.Statement.Preloads[key] = []any{}
		}
	}
	return
}

func splitConds(conds []any) (map[string][]any, []any) {
	if len(conds) > 0 {
		if m, ok := conds[0].(map[string][]any); ok {
			return m, conds[1:]
		}
	}
	return nil, conds
}

func (db *DB) Preload(t any, conds ...any) *DB {
	m, conds := splitConds(conds)
	return db.PreloadType(t, m).First(t, conds...)
}

func (db *DB) Preloads(t any, conds ...any) *DB {
	m, conds := splitConds(conds)
	return db.PreloadType(t, m).Find(t, conds...)
}

func NewDB(file string, gormDB *gorm.DB) *DB {
	return &DB{
		DB:   gormDB,
		file: file,
	}
}
