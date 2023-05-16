package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"messengerAPI/controllers"
	"messengerAPI/initializers"
	"messengerAPI/middleware"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables \n", err.Error())
	}
	initializers.ConnectDB(&config)
}

func main() {
	app := fiber.New()

	//app.Use(cors.New(cors.Config{
	//	AllowOrigins:     "http://localhost:8888, http://localhost:3000, https://c7b4-188-255-11-131.ngrok-free.app",
	//	AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	//	AllowMethods:     "GET, POST, PUT, DELETE",
	//	AllowCredentials: true,
	//}))
	app.Use(cors.New())

	app.Post("/upload", middleware.DeserializeUser, controllers.UploadFile)

	messageRoutes := app.Group("/chats", middleware.DeserializeUser)
	messageRoutes.Get("/", controllers.GetConversations)
	messageRoutes.Post("/", controllers.CreateConversation)
	messageRoutes.Get("/:id", controllers.GetMessages)
	messageRoutes.Post("/:id", controllers.CreateMessage)
	messageRoutes.Get("/messages/:id/delete", controllers.DeleteMessage)
	messageRoutes.Post("/messages/:id/attach", controllers.AttachFile)
	messageRoutes.Post("/messages/:id/read", controllers.ReadMessage)
	messageRoutes.Put("/:id", controllers.EditConversation)
	messageRoutes.Delete("/messages/:id/delete", controllers.DeleteMessageForAll)
	messageRoutes.Post("/:id/participants/add", controllers.AddParticipant)
	messageRoutes.Get("/:id/participants", controllers.GetConversationParticipants)
	messageRoutes.Delete("/:id/participants/delete", controllers.DeleteParticipant)
	messageRoutes.Get("/:id/attachments", controllers.GetAttachments)

	//messageRoutes.Post("/dialogs/:id", controllers.CreateDialogMessage)
	//messageRoutes.Get("/dialogs/:id", controllers.GetDialogMessages)
	//messageRoutes.Get("/dialogs/", controllers.GetDialogs)

	userRoutes := app.Group("/users", middleware.DeserializeUser)
	userRoutes.Get("/", controllers.GetUsers)
	userRoutes.Get("/:id", controllers.GetUser)
	userRoutes.Put("/:id", controllers.UpdateUser)
	userRoutes.Put("/:id/change-password", controllers.ChangePassword)
	userRoutes.Put("/:id/change-decrypt-password", controllers.ChangeDecryptPassword)
	userRoutes.Get("/current-user", controllers.GetLoggedUser)
	userRoutes.Post("/search", controllers.SearchUsers)
	userRoutes.Post("/filter", controllers.FilterUsers)

	//userRoutes.Get("/:id/create-dialog", controllers.CreateDialog)

	authRoutes := app.Group("/auth")
	authRoutes.Post("/login", controllers.SignInUser)
	authRoutes.Post("/register", controllers.SignUpUser)
	authRoutes.Post("/logout", middleware.DeserializeUser, controllers.LogoutUser)

	app.Static("/uploads", "/uploads")

	log.Fatal(app.Listen("127.0.0.1:8888"))
}
