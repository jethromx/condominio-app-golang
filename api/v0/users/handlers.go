package users

import (
	"com.mx/crud/internal/service"
	"com.mx/crud/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// CreateUserHandler maneja la creación de nuevos usuarios
func CreateUserHandler(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		input := new(UserInput)

		err := utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		user := MapUserInputToModel(input)

		hashedPassword, err := utils.GeneratePassword(input.Password)
		if err != nil {
			log.Debug("Error hashing password:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Error hashing password ", nil)
		}
		user.Password = string(hashedPassword)

		// Obtener el ID del usuario autenticado del contexto
		//auditUserID := c.Locals("auditUserID").(uint)
		auditUserID := 1 // Eliminar esta línea cuando se implemente la autenticación
		user, err = userService.CreateUser(user, uint(auditUserID))

		if err != nil {
			log.Debug("Error creating user: ", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Could not create user: "+err.Error(), nil)
		}

		output := MapUserModelToOutput(user)

		return utils.HandleResponse(c, fiber.StatusOK, "User created successfully", output)
	}
}

func GetAllUsersHandler(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		page, pageSize, err := utils.GetPaginationParams(c)
		if err != nil {
			log.Debug("Error getting pagination parameters:", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid pagination parameters "+err.Error(), nil)
		}

		users, totalRecords, err := userService.GetAllUsers(page, pageSize)
		if err != nil {
			log.Debug("Error getting users:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Could not get users: "+err.Error(), nil)
		}

		outputs := MapUserModelsToOutputs(users)

		response := fiber.Map{
			"data":         outputs,
			"page":         page,
			"pageSize":     pageSize,
			"totalPages":   (totalRecords + int64(pageSize) - 1) / int64(pageSize),
			"totalRecords": totalRecords,
		}

		return utils.HandleResponse(c, fiber.StatusOK, "Users retrieved successfully", response)
	}
}

// GetUserHandler maneja la obtención de un usuario por ID
func GetUserHandler(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			log.Debug("Error parsing ID:", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid user ID "+err.Error(), nil)
		}

		user, err := userService.GetUserByID(uint(id))
		if err != nil {
			log.Debug("Error getting user:", err)
			return utils.HandleResponse(c, fiber.StatusNotFound, "User not found "+err.Error(), nil)
		}

		output := MapUserModelToOutput(user)

		return utils.HandleResponse(c, fiber.StatusOK, "User retrieved successfully ", output)
	}
}

// UpdateUserHandler maneja la actualización de un usuario
func UpdateUserHandler(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid user ID:", err.Error())
		}

		input := new(UserInput)

		err = utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		var user = MapUserInputToModel(input)
		user.ID = uint(id) // Asegúrate de establecer el ID correcto

		if err := userService.UpdateUser(user); err != nil {
			log.Debug("Error updating user:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Could not update user: "+err.Error(), nil)
		}

		output := MapUserModelToOutput(user)

		return utils.HandleResponse(c, fiber.StatusOK, "User updated successfully ", output)
	}
}

// DeleteUserHandler maneja la eliminación de un usuario
func DeleteUserHandler(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid user ID: "+err.Error(), nil)
		}

		if err := userService.DeleteUser(uint(id)); err != nil {
			log.Debug("Error deleting user:", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Could not delete user: "+err.Error(), nil)
		}

		return utils.HandleResponse(c, fiber.StatusOK, "User deleted successfully", nil)
	}
}
