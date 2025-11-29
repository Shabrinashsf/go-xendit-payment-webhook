package migrations

import (
	"github.com/Shabrinashsf/go-xendit-payment-webhook/migrations/seeds"
	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) error {
	if err := seeds.ListProductSeeder(db); err != nil {
		return err
	}
	return nil
}
