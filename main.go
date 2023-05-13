package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/harold/trip-golang/models"
	"github.com/harold/trip-golang/storage"
	"github.com/joho/godotenv"

	"gorm.io/gorm"
)

type Trip struct {
	PassengerName string `json:"passengername"`
	Destination   string `json:"destination"`
	PickupPoint   string `json:"pickuppoint"`
	PhoneNumber   string `json:"phonenumber"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateTrip(context *fiber.Ctx) error {
	trip := Trip{}

	err := context.BodyParser(&trip)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	if trip.PhoneNumber != "" {
		phoneNumberInt, err := strconv.Atoi(trip.PhoneNumber)
		if err != nil {
			context.Status(http.StatusUnprocessableEntity).JSON(
				&fiber.Map{"message": "invalid phone number"})
			return err
		}
		trip.PhoneNumber = strconv.Itoa(phoneNumberInt)
	}

	err = r.DB.Create(&trip).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create trip"})
		return nil
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "trip has been added"})
	return nil
}

func (r *Repository) DeleteTrip(context *fiber.Ctx) error {
	tripModel :=
		models.Trip{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Delete(tripModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete trip plan",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "trip deleted successfully",
	})
	return nil
}

func (r *Repository) GetTrip(context *fiber.Ctx) error {
	tripModel := &[]models.Trip{}

	err := r.DB.Find(tripModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not get trip data",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "trip successfully fetch", "data": tripModel})
	return nil
}

func (r *Repository) GetTripByID(context *fiber.Ctx) error {
	id := context.Params("id")
	tripModel := &models.Trip{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot found",
		})
		return nil
	}

	fmt.Println("the ID is", id)

	err := r.DB.Where("id = ?", id).First(tripModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not get trip",
		})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "student id fetched successfully",
		"data":    tripModel,
	})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_trip", r.CreateTrip)
	api.Delete("delete_trip/:id", r.DeleteTrip)
	api.Get("/get_trips/:id", r.GetTripByID)
	api.Get("/trips", r.GetTrip)

}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("could not load the database")
	}
	err = models.MigrateTrips(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
