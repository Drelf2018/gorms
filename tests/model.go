package tests

import (
	"os"
	"strings"
	"time"

	"github.com/Drelf2018/gorms"
	"gorm.io/gorm"
)

var db *gorms.DB

func init() {
	os.Remove("./test.db")
	db = gorms.SetSQLite("./test.db").AutoMigrate(&Attachment{}, &User{})
	db.Create([]User{{
		Face:    Attachment{Path: "face.png"},
		Pendant: Attachment{Path: "pendant.png"},
	}, {
		Face:    Attachment{Path: "face.png"},
		Pendant: Attachment{Path: "pendant.jpg"},
	}, {
		Face:    Attachment{Path: "face.jpg"},
		Pendant: Attachment{Path: "pendant.png"},
	}, {
		Face:    Attachment{Path: "face.jpg"},
		Pendant: Attachment{Path: "pendant.jpg"},
	}})
}

type Attachment struct {
	Path string `gorm:"primaryKey" json:"path"`
	MIME string
}

func (a *Attachment) AfterCreate(*gorm.DB) error {
	go func() {
		// simulate network request
		time.Sleep(time.Second)
		_, a.MIME, _ = strings.Cut(a.Path, ".")
		db.Updates(a)
	}()
	return nil
}

type User struct {
	gorms.Model[int64]

	FacePath string     `json:"-"`
	Face     Attachment `json:"face"`

	PendantPath string     `json:"-"`
	Pendant     Attachment `json:"pendant"`
}
