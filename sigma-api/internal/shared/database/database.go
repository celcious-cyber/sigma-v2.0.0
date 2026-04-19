package database

import (
	"log"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sigma-api/internal/core/models"
)

var DB *gorm.DB

// InitDB initializes the high-performance SQLite connection with WAL mode and CGO-free driver.
func InitDB(dsn string) {
	var err error

	// High-Performance DSN parameters:
	// _busy_timeout=5000: Wait up to 5s if the database is locked
	// _journal_mode=WAL: Enable Write-Ahead Logging for concurrent reads
	db, err := gorm.Open(sqlite.Open(dsn+"?_busy_timeout=5000&_journal_mode=WAL"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		PrepareStmt: true, // Cache prepared statements for zero-allocation performance
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get generic database object: %v", err)
	}

	// Connection Tuning for SQLite:
	// Sequential write guarantee to avoid "database is locked" errors.
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Optimization PRAGMAs
	db.Exec("PRAGMA journal_mode=WAL;")
	db.Exec("PRAGMA synchronous=NORMAL;")
	db.Exec("PRAGMA foreign_keys=ON;")
	db.Exec("PRAGMA cache_size=-64000;") // 64MB Cache

	// Auto-Migration
	log.Println("Running Auto-Migrations...")
	err = db.AutoMigrate(
		&models.User{},
		&models.Student{},
		&models.Classroom{},
		&models.Teacher{},
		&models.Alumni{},
		&models.Subject{},
		&models.TeacherSchedule{},
		&models.Attendance{},
		&models.Assessment{},
		&models.StudyPeriod{},
		&models.AcademicCalendar{},
		&models.QuranMemorization{},
		&models.PaymentCategory{},
		&models.Invoice{},
		&models.Payment{},
		&models.IncomingLetter{},
		&models.OutgoingLetter{},
		&models.LetterCode{},
		&models.Visitor{},
		&models.Announcement{},
		&models.Medicine{},
		&models.Disease{},
		&models.FitnessRecord{},
		&models.Book{},
		&models.Borrowing{},
		&models.Setting{},
	)
	if err != nil {
		log.Fatalf("Auto-Migration failed: %v", err)
	}

	DB = db
	log.Printf("High-Performance SQLite initialized in WAL mode at: %s", dsn)
}
