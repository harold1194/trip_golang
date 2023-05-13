package models

import "gorm.io/gorm"

type Trip struct {
	ID             uint    `gorm:"primary key; autoIncrement" json:"id"`
	PassengerName  *string `json:"passengername"`
	Destination    *string `json:"destination"`
	PickupPoint    *string `json:"pickuppoint"`
	PhoneNumber    *int    `json:"phonenumber"`
	PhoneNumberInt *int    `json:"-"`
}

func MigrateTrips(db *gorm.DB) error {
	err := db.AutoMigrate(&Trip{})
	return err
}
