package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const connectionString = "host=localhost user=david password=david dbname=cli port=5000 sslmode=disable TimeZone=Europe/Ljubljana"

type Measurement struct {
	Id        int
	Date      string
	Time      string
	StatusApi string
	StatusDb  string
	Mode      int
}

type MeasurementsTable struct {
	database *gorm.DB
}

func NewMeasurements() (*MeasurementsTable, error) {
	db, err := gorm.Open(postgres.Open(connectionString))
	if err != nil {
		return nil, err
	}
	return &MeasurementsTable{database: db}, nil
}

func (m MeasurementsTable) Close() error {
	db, err := m.database.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (m MeasurementsTable) CreateMeasurement(measurement Measurement) (Measurement, error) {
	return measurement, m.database.Create(&measurement).Error
}
