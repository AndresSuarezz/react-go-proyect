package main

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/AndresSuarezz/react-go-proyect/models"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	app := fiber.New()

	// Conectando a la BD
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017/crudgo"))

	// Verificar si hay un error
	if err != nil {
		panic(err)
	}

	// Conexion a la DB
	coll := client.Database("crudgo").Collection("users")

	app.Use(cors.New())

	// Crear un usuario
	app.Post("/users", func(c *fiber.Ctx) error {
		var user models.User

		//Recibir datos del body
		c.BodyParser(&user)
		//Insertando los datos en la base de datos
		result, err := coll.InsertOne(context.TODO(), bson.D{
			{Key: "name", Value: user.Name},
		})

		if err != nil {
			panic(err)
		}

		return c.JSON(fiber.Map{
			"message": "Creando un usuario en la DB",
			"result":  result,
		})
	})

	// Obtener todos los usuarios
	app.Get("/users", func(c *fiber.Ctx) error {
		//Obtenemos los datos
		results, err := coll.Find(context.TODO(), bson.D{})
		if err != nil {
			panic(err)
		}

		var users []models.User

		// Guardamos en una variable cada uno de los usuarios
		for results.Next(context.TODO()) {
			var user models.User
			results.Decode(&user)
			users = append(users, user)
		}

		return c.JSON(&fiber.Map{
			"users": users,
		})
	})

	// Eliminar Usuario
	app.Delete("/users/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		// Convertir el ID a tipo ObjectId si es necesario
		objID, err := primitive.ObjectIDFromHex(id)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "ID de usuario inválido",
			})
		}

		// Eliminar el usuario de la base de datos
		result, err := coll.DeleteOne(context.TODO(), bson.M{"_id": objID})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error al eliminar usuario",
			})
		}

		// Verificar si se eliminó correctamente
		if result.DeletedCount == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Usuario no encontrado",
			})
		}

		// Usuario eliminado exitosamente
		return c.JSON(fiber.Map{
			"message": "Usuario eliminado",
		})
	})

	// Actualizar Usuario
	app.Put("/users/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		// Convertir el ID a tipo ObjectId si es necesario
		objID, err := primitive.ObjectIDFromHex(id)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "ID de usuario inválido",
			})
		}

		var user models.User

		// Recibir datos del body
		c.BodyParser(&user)

		// Actualizar el usuario en la base de datos
		result, err := coll.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: user.Name},
			}},
		})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error al actualizar usuario",
			})
		}

		// Verificar si se actualizó correctamente
		if result.ModifiedCount == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Usuario no encontrado",
			})
		}

		// Usuario actualizado exitosamente
		return c.JSON(fiber.Map{
			"message": "Usuario actualizado",
		})

	})

	app.Listen(":" + port)
}
