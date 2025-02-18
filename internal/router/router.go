package router

import (
	"com.mx/crud/api/v0/auth"
	"com.mx/crud/api/v0/reservations"
	"com.mx/crud/config/database"
	"com.mx/crud/internal/middleware"
	"com.mx/crud/internal/repository"
	"com.mx/crud/internal/service"
	"github.com/gofiber/fiber/v2"

	"com.mx/crud/api/v0/condominiums"
	"com.mx/crud/api/v0/maintenances"
	"com.mx/crud/api/v0/payments"
	"com.mx/crud/api/v0/users"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {

	api := app.Group("/api")

	// Crear repositorios y servicios
	userRepo := repository.NewUserRepository(database.DB)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, userService)
	tokenRepo := repository.NewTokenRepository(database.DB)
	tokenService := service.NewTokenService(tokenRepo)

	condominiumRepo := repository.NewCondominiumRepository(database.DB)
	buildingRepo := repository.NewBuildingRepository(database.DB)
	paymentRepo := repository.NewPaymentRepository(database.DB)
	buildingService := service.NewBuildingService(buildingRepo, condominiumRepo)
	paymentService := service.NewPaymentService(paymentRepo)
	reservationRepo := repository.NewReservationRepository(database.DB)
	residentRepo := repository.NewResidentRepository(database.DB)
	maintenanceRepo := repository.NewMaintenanceRepository(database.DB)
	apartmentRepo := repository.NewApartmentRepository(database.DB)

	reservationService := service.NewReservationService(reservationRepo)
	residentService := service.NewResidentService(residentRepo, userRepo)
	maintenanceService := service.NewMaintenanceService(maintenanceRepo)
	apartmentService := service.NewApartmentService(apartmentRepo, buildingRepo)
	condominiumService := service.NewCondominiumService(condominiumRepo, buildingRepo)
	// Rutas de usuario
	userRoute := api.Group("/users")
	userRoute.Get("/:id", users.GetUserHandler(userService))
	userRoute.Post("", users.CreateUserHandler(userService))
	userRoute.Get("", users.GetAllUsersHandler(userService))
	userRoute.Patch("/:id", users.UpdateUserHandler(userService))
	userRoute.Delete("/:id", users.DeleteUserHandler(userService))

	// Rutas de condominios
	condominiumRoute := api.Group("/condominiums", middleware.Protected(tokenService))

	condominiumRoute.Get("", condominiums.GetAllCondominiumsHandler(condominiumService))
	condominiumRoute.Post("", condominiums.CreateCondominiumHandler(condominiumService))
	condominiumRoute.Get("/:id", condominiums.GetCondominiumByIDHandler(condominiumService))
	condominiumRoute.Put("/:id", condominiums.UpdateCondominiumHandler(condominiumService))
	condominiumRoute.Delete("/:id", condominiums.DeleteCondominiumHandler(condominiumService))
	condominiumRoute.Get("/:id/apartments", condominiums.GetAllApartmentsByCondominiumHandler(apartmentService))
	// Rutas de edificios dentro de condominios
	buildingRoute := condominiumRoute.Group("/:id/buildings")
	buildingRoute.Get("", condominiums.GetAllBuildingsHandler(buildingService))
	buildingRoute.Post("", condominiums.CreateBuildingsHandler(buildingService))
	buildingRoute.Get("/:idBuilding", condominiums.GetBuildingsByIDHandler(buildingService))
	buildingRoute.Patch("/:idBuilding", condominiums.UpdateBuildingsHandler(buildingService))
	buildingRoute.Delete("/:idBuilding", condominiums.DeleteBuildingsHandler(buildingService))

	// Rutas de apartamentos dentro de edificios
	apartmentRoute := buildingRoute.Group("/:idBuilding/apartments")
	apartmentRoute.Get("", condominiums.GetAllApartmentsHandler(apartmentService))
	apartmentRoute.Get("/:idApartment", condominiums.GetApartmentHandler(apartmentService))
	apartmentRoute.Post("", condominiums.CreateApartmentHandler(apartmentService, buildingService))
	apartmentRoute.Patch("/:idApartment", condominiums.UpdateApartmentHandler(apartmentService))
	apartmentRoute.Delete("/:idApartment", condominiums.DeleteApartmentHandler(apartmentService))

	// Rutas de residentes dentro de apartamentos
	residentRoute := apartmentRoute.Group("/:idApartment/residents")
	residentRoute.Get("", condominiums.GetAllResidentsHandler(residentService))
	residentRoute.Post("", condominiums.CreateResidentHandler(residentService, userService))
	residentRoute.Get("/:idResident", condominiums.GetResidentByIDHandler(residentService))
	residentRoute.Patch("/:idResident", condominiums.UpdateResidentHandler(residentService))
	residentRoute.Delete("/:idResident", condominiums.DeleteResidentHandler(residentService))

	// Rutas de pagos
	paymentRoute := residentRoute.Group("/:idResident/payments")
	paymentRoute.Get("/", payments.GetAllPaymentsHandler(paymentService))
	paymentRoute.Post("/", payments.CreatePaymentHandler(paymentService, residentService))
	paymentRoute.Get("/:id", payments.GetPaymentByIDHandler(paymentService, residentService))
	paymentRoute.Patch("/:id", payments.UpdatePaymentHandler(paymentService, residentService))
	paymentRoute.Delete("/:id", payments.DeletePaymentHandler(paymentService, residentService))

	// Rutas de reservas
	reservationRoute := api.Group("/reservations")
	reservationRoute.Get("/", reservations.GetAllReservationsHandler(reservationService))
	reservationRoute.Post("/", reservations.CreateReservationHandler(reservationService))
	reservationRoute.Get("/:id", reservations.GetReservationByIDHandler(reservationService))
	reservationRoute.Patch("/:id", reservations.UpdateReservationHandler(reservationService))
	reservationRoute.Delete("/:id", reservations.DeleteReservationHandler(reservationService))

	// Rutas de mantenimientos
	maintenanceRoute := api.Group("/maintenances")
	maintenanceRoute.Get("/", maintenances.GetAllMaintenancesHandler(maintenanceService))
	maintenanceRoute.Post("/", maintenances.CreateMaintenanceHandler(maintenanceService))
	maintenanceRoute.Get("/:id", maintenances.GetMaintenanceByIDHandler(maintenanceService))
	maintenanceRoute.Patch("/:id", maintenances.UpdateMaintenanceHandler(maintenanceService))
	maintenanceRoute.Delete("/:id", maintenances.DeleteMaintenanceHandler(maintenanceService))

	// Auth
	authRoute := api.Group("/auth")
	authRoute.Post("/login", auth.LoginHandler(authService, tokenService))
	authRoute.Post("/logout", auth.LogoutHandler(tokenService))
	authRoute.Post("/register", auth.RegisterHandler(authService))
	authRoute.Post("/refresh-token", auth.RefreshTokenHandler(authService, tokenService))

}
