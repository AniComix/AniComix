package main

import (
	"github.com/AniComix/server/models"
	"github.com/AniComix/server/storage"
	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery, // generate mode
	})
	storage.InitStorage()
	g.UseDB(storage.DB()) // reuse your gorm db

	// Generate basic type-safe DAO API for struct `model.User` following conventions
	g.ApplyBasic(models.Series{})
	g.ApplyBasic(models.Episode{})
	g.ApplyBasic(models.User{})

	// Generate the code
	g.Execute()
}
