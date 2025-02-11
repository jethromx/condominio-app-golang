package router

import (
	"com.mx/crud/api/v0/auth"
	"com.mx/crud/api/v0/building"
	"com.mx/crud/api/v0/reservations"
	"com.mx/crud/config/database"
	"com.mx/crud/internal/middleware"
	"com.mx/crud/internal/repository"
	"com.mx/crud/internal/service"
	"github.com/gofiber/fiber/v2"

	"com.mx/crud/api/v0/apartments"
	"com.mx/crud/api/v0/condominiums"
	"com.mx/crud/api/v0/maintenances"
	"com.mx/crud/api/v0/payments"
	"com.mx/crud/api/v0/residents"
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
	residentService := service.NewResidentService(residentRepo)
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
	condominiumRoute.Patch("/:id", condominiums.UpdateCondominiumHandler(condominiumService))
	condominiumRoute.Delete("/:id", condominiums.DeleteCondominiumHandler(condominiumService))

	condominiumRoute.Get("/:id/buildings", condominiums.GetAllBuildingsHandler(buildingService))
	condominiumRoute.Post("/:id/buildings", condominiums.CreateBuildingsHandler(buildingService))
	condominiumRoute.Get("/:id/buildings/:idBuilding", condominiums.GetBuildingsByIDHandler(buildingService))
	condominiumRoute.Patch("/:id/buildings/:idBuilding", condominiums.UpdateBuildingsHandler(buildingService))
	condominiumRoute.Delete("/:id/buildings/:idBuilding", condominiums.DeleteBuildingsHandler(buildingService))

	condominiumRoute.Get("/:id/buildings/:idBuilding/apartments", condominiums.GetAllApartmentsHandler(apartmentService))
	condominiumRoute.Get("/:id/buildings/:idBuilding/apartments/:idApartment", condominiums.GetApartmentHandler(apartmentService))
	condominiumRoute.Post("/:id/buildings/:idBuilding/apartments", condominiums.CreateApartmentHandler(apartmentService, buildingService))
	condominiumRoute.Patch("/:id/buildings/:idBuilding/apartments/:idApartment", condominiums.UpdateApartmentHandler(apartmentService))
	condominiumRoute.Delete("/:id/buildings/:idBuilding/apartments/:idApartment", condominiums.DeleteApartmentHandler(apartmentService))

	// Rutas de edificios
	buildingRoute := api.Group("/buildings", middleware.Protected(tokenService))
	buildingRoute.Get("/", building.GetAllBuildingsHandler(buildingService))
	buildingRoute.Post("/", building.CreateBuildingHandler(buildingService))
	buildingRoute.Get("/:id", building.GetBuildingByIDHandler(buildingService))
	buildingRoute.Patch("/:id", building.UpdateBuildingHandler(buildingService))
	buildingRoute.Delete("/:id", building.DeleteBuildingHandler(buildingService))
	buildingRoute.Get("/:id/apartments", building.GetAllApartmentsHandler(apartmentService))
	buildingRoute.Get("/:id/apartments/:idApartment", building.GetApartmentHandler(apartmentService))
	buildingRoute.Post("/:id/apartments", building.CreateApartmentHandler(apartmentService))
	buildingRoute.Patch("/:id/apartments/:idApartment", building.UpdateApartmentHandler(apartmentService))
	buildingRoute.Delete("/:id/apartments/:idApartment", building.DeleteApartmentHandler(apartmentService))

	apartmentRoute := api.Group("/apartments", middleware.Protected(tokenService))
	apartmentRoute.Get("", apartments.GetAllApartmentsHandler(apartmentService))
	apartmentRoute.Get("/:id", apartments.GetApartmentHandler(apartmentService))
	apartmentRoute.Post("", apartments.CreateApartmentHandler(apartmentService))
	apartmentRoute.Patch("/:id", apartments.UpdateApartmentHandler(apartmentService))
	apartmentRoute.Delete("/:id", apartments.DeleteApartmentHandler(apartmentService))

	// Rutas de residentes
	residentRoute := api.Group("/residents", middleware.Protected(tokenService))
	residentRoute.Get("/", residents.GetAllResidentsHandler(residentService))
	residentRoute.Post("/", residents.CreateResidentHandler(residentService))
	residentRoute.Get("/:id", residents.GetResidentByIDHandler(residentService))
	residentRoute.Patch("/:id", residents.UpdateResidentHandler(residentService))
	residentRoute.Delete("/:id", residents.DeleteResidentHandler(residentService))

	// Rutas de pagos
	paymentRoute := api.Group("/payments")
	paymentRoute.Get("/", payments.GetAllPaymentsHandler(paymentService))
	paymentRoute.Post("/", payments.CreatePaymentHandler(paymentService))
	paymentRoute.Get("/:id", payments.GetPaymentByIDHandler(paymentService))
	paymentRoute.Patch("/:id", payments.UpdatePaymentHandler(paymentService))
	paymentRoute.Delete("/:id", payments.DeletePaymentHandler(paymentService))

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
