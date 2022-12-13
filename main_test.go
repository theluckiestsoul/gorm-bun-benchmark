package main

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"database/sql"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	//"github.com/uptrace/bun/driver/sqliteshim"
)

type User struct {
	ID   int
	Name string
}

func setupBun() *bun.DB {
	sqldb, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())
	var sb strings.Builder
	sb.WriteString(`
	DROP TABLE IF EXISTS users;
	CREATE TABLE users(id INTEGER PRIMARY KEY, name TEXT);
	`)
	for i := 0; i < 10; i++ {
		sb.WriteString(fmt.Sprintf(`INSERT INTO users (name) VALUES ("John %d");`, i))
	}
	db.Exec(sb.String())

	return db
}

func setupGorm() *gorm.DB {
	db, err := gorm.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		panic("failed to connect database")
	}
	var sb strings.Builder
	sb.WriteString(`
	DROP TABLE IF EXISTS users;
	CREATE TABLE users(id INTEGER PRIMARY KEY, name TEXT);
	`)
	for i := 0; i < 10; i++ {
		sb.WriteString(fmt.Sprintf(`INSERT INTO users (name) VALUES ("John %d");`, i))
	}
	db.Exec(sb.String())

	return db
}

func BenchmarkGormInsert(b *testing.B) {
	// Set up a connection to the database
	db := setupGorm()
	defer db.Close()

	// Execute the query and measure the time taken
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.Create(&User{
			Name: "John",
		})
	}
}

func BenchmarkUptraceBunInsert(b *testing.B) {
	// Set up a connection to the database
	db := setupBun()
	defer db.Close()

	// Execute the query and measure the time taken
	b.ResetTimer()
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		db.NewInsert().Model(&User{
			Name: "John",
		}).Exec(ctx)
	}
}

func BenchmarkGormQuery(b *testing.B) {
	// Set up a connection to the database
	db := setupGorm()
	defer db.Close()

	// Execute the query and measure the time taken
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var users []User
		db.Find(&users)
	}
}

func BenchmarkUptraceBunQuery(b *testing.B) {
	// Set up a connection to the database
	db := setupBun()
	defer db.Close()

	// Execute the query and measure the time taken
	b.ResetTimer()
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		var users []User
		db.NewSelect().Model(&users).Scan(ctx)
	}
}

func BenchmarkGormUpdate(b *testing.B) {
	// Set up a connection to the database
	db := setupGorm()
	defer db.Close()

	// Execute the query and measure the time taken
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.Model(&User{}).Where("id = ?", i).Update("name", "John")
	}
}

func BenchmarkUptraceBunUpdate(b *testing.B) {
	// Set up a connection to the database
	db := setupBun()
	defer db.Close()

	// Execute the query and measure the time taken
	b.ResetTimer()
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		db.NewUpdate().Model(&User{}).Where("id = ?", i).Set("name = ?", "John").Exec(ctx)
	}
}
