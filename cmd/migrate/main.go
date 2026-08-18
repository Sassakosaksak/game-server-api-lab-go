package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URLが設定されていません")
	}

	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		log.Fatalf("Migrationフォルダのパス取得に失敗しました: %v", err)
	}

	migration, err := migrate.New("file://"+filepath.ToSlash(migrationsPath), databaseURL)
	if err != nil {
		log.Fatalf("Migrationの準備に失敗しました: %v", err)
	}
	defer migration.Close()

	if err := migration.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("適用するMigrationはありません")
			return
		}

		log.Fatalf("Migrationの適用に失敗しました: %v", err)
	}

	log.Println("Migrationを適用しました")
}
