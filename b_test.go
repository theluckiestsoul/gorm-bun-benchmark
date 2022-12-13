package main

import (
	"context"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	//"github.com/uptrace/bun/driver/sqliteshim"
)

func BenchmarkGorm(b *testing.B) {
	// Set up a connection to the database
	db := setupGorm()
	defer db.Close()
	
	b.Run("Insert", func(b *testing.B) {
		// Execute the query and measure the time taken
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			db.Create(&User{
				Name: "John",
			})
		}
	})

	b.Run("Query", func(b *testing.B) {
		// Execute the query and measure the time taken
		b.ResetTimer()
		var users []User
		for i := 0; i < b.N; i++ {
			db.Find(&users)
		}
	})

	b.Run("Update", func(b *testing.B) {
		// Execute the query and measure the time taken
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			db.Model(&User{}).Where("id = ?", i).Update("name", "John")
		}
	})
}

func BenchmarkUptraceBun(b *testing.B) {
	// Set up a connection to the database
	db := setupBun()
	defer db.Close()

	b.Run("Insert", func(b *testing.B) {
		// Execute the query and measure the time taken
		b.ResetTimer()
		ctx := context.Background()
		for i := 0; i < b.N; i++ {
			db.NewInsert().Model(&User{
				Name: "John",
			}).Exec(ctx)
		}
	})

	b.Run("Query", func(b *testing.B) {
		// Execute the query and measure the time taken
		b.ResetTimer()
		ctx := context.Background()
		var users []User
		for i := 0; i < b.N; i++ {
			db.NewSelect().Model(&users).Scan(ctx)
		}
	})

	b.Run("Update", func(b *testing.B) {
		// Execute the query and measure the time taken
		b.ResetTimer()
		ctx := context.Background()
		for i := 0; i < b.N; i++ {
			db.NewUpdate().Model(&User{}).Where("id = ?", i).Set("name = ?", "John").Exec(ctx)
		}
	})
}
