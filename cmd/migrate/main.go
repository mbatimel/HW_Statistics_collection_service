package main

import "github.com/mbatimel/HW_Statistics_collection_service/internal/migrate"

func main() {
	err := migrate.RunMigrations()
	if err != nil {
		panic(err)
	}
}
