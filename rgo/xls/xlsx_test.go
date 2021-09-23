package xls

import (
	"log"
	"testing"
	"time"
)

type Member struct {
	ID            int        `gorm:"id" json:"id" xlsx:"0"`
	Name          string     `gorm:"name" json:"name" xlsx:"1"`
	Email         string     `gorm:"email" json:"email" xlsx:"2"`
	Password      string     `gorm:"password" json:"password" xlsx:"3"`
	RememberToken string     `gorm:"remember_token" json:"remember_token" xlsx:"4"`
	CreatedAt     *time.Time `xlsx:"-"`
	UpdatedAt     *time.Time `xlsx:"-"`
}

func testXlsxImport(t *testing.T) {
	t.Run("xxxlsx", func(t *testing.T) {
		var members []Member
		xlsx := NewXlsx()
		err := xlsx.SetSrcDir("rgo/public/files/").
			SetFileName("我的excel.xlsx").
			SetExceptFirstLine(true).
			Import(&members)
		if err != nil {
			t.Error(err)
		}
		log.Println(members)
		log.Println(xlsx.Cnt)
	})
}

func testXlsxExport(t *testing.T) {
	t.Run("xxxlsxs", func(t *testing.T) {
		users := make([]Member, 0)
		users = append(users, Member{
			ID:       1,
			Name:     "oldda",
			Email:    "17710818223@163.com",
			Password: "xxx",
		})

		xlsx := NewXlsx()
		errs := xlsx.SetSheetLimit(1).
			SetFileName("我的excel").
			SetSheetName("我的Sheet").
			SetDstDir("rgo/public/files").
			SetHeader(&[]string{"ID", "姓名", "邮箱", "密码", "附加"}).
			SetData(users).
			Export() //导出
		if errs != nil {
			t.Error(errs)
		}
		log.Println(xlsx.Cnt)
	})
}
