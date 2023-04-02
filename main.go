package main

import (
	"github.com/gofiber/fiber/v2"
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
	//micro := fiber.New()
	//app.Mount("/api", micro)
	//app.Use(logger.New())
	//app.Use(cors.New(cors.Config{
	//	AllowOrigins:     "http://localhost:3000",
	//	AllowHeaders:     "Origin, Content-Type, Accept",
	//	AllowMethods:     "GET, POST",
	//	AllowCredentials: true,
	//}))

	messageRoutes := app.Group("/messages", middleware.DeserializeUser)
	messageRoutes.Get("/", controllers.GetConversations)
	messageRoutes.Post("/add-conversation", controllers.CreateConversation)
	messageRoutes.Get("/groups/:id", controllers.GetMessages)
	messageRoutes.Post("/groups/:id", controllers.CreateMessage)
	messageRoutes.Get("/:id/delete", controllers.DeleteMessage)
	messageRoutes.Post("/:id/attach", controllers.AttachFile)
	messageRoutes.Put("/:id", controllers.EditConversation)
	messageRoutes.Get("/:id/delete-for-all", controllers.DeleteMessageForAll)
	messageRoutes.Post("/:id/participants/add", controllers.AddParticipant)
	messageRoutes.Get("/:id/participants", controllers.GetConversationParticipants)
	messageRoutes.Delete("/:id/participants/delete", controllers.DeleteParticipant)
	messageRoutes.Post("/dialogs/:id", controllers.CreateDialogMessage)
	messageRoutes.Get("/dialogs/:id", controllers.GetDialogMessages)
	messageRoutes.Get("/dialogs/", controllers.GetDialogs)

	userRoutes := app.Group("/users", middleware.DeserializeUser)
	userRoutes.Get("/", controllers.GetUsers)
	userRoutes.Get("/:id", controllers.GetUser)
	userRoutes.Put("/:id", controllers.UpdateUser)
	userRoutes.Put("/:id/change-password", controllers.ChangePassword)
	userRoutes.Put("/:id/change-decrypt-password", controllers.ChangeDecryptPassword)
	userRoutes.Get("/current-user", controllers.GetLoggedUser)
	userRoutes.Get("/:id/create-dialog", controllers.CreateDialog)
	userRoutes.Post("/search", controllers.SearchUsers)
	userRoutes.Post("/filter", controllers.FilterUsers)
	//routes.Delete("/:id", h.DeleteUser)

	authRoutes := app.Group("/auth")
	authRoutes.Post("/login", controllers.SignInUser)
	authRoutes.Post("/register", controllers.SignUpUser)
	authRoutes.Post("/logout", middleware.DeserializeUser, controllers.LogoutUser)

	app.Static("/", "htdocs")

	//controllers.RegisterRoutes(app, initializers.DB)

	log.Fatal(app.Listen(":3000"))
}
