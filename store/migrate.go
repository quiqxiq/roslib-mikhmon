package store

import (
	"github.com/quiqxiq/roslib-mikhmon/store/model"
	"gorm.io/gorm"
)

// Migrate menjalankan AutoMigrate untuk semua tabel.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.User{},
		&model.RefreshToken{},
		&model.MikrotikDevice{},
		&model.Transaction{},
		&model.HotspotProfileConfig{},
	); err != nil {
		return err
	}
	// Drop kolom slug yang sudah tidak dipakai (backward-compatible:
	// idempotent — tidak error kalau kolom sudah tidak ada).
	_ = db.Migrator().DropColumn(&model.MikrotikDevice{}, "slug")
	return nil
}
