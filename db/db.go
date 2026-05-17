package db
 
import (
    "log"
 
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gotask/config"
    "gotask/models"
)
 
var DB *gorm.DB
 
func Connect(cfg config.Config) {
    var err error
    DB, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
 
    // AutoMigrate creates/updates the tasks table to match the struct
    if err := DB.AutoMigrate(&models.User{}, &models.Task{}); err != nil {
        log.Fatal("Migration failed:", err)
    }
 
    log.Println("Database connected and migrated.")
}
