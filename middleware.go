package main

import (
    "github.com/gofiber/fiber/v3"
)

func JWTProtected() fiber.Handler {
    return func(c fiber.Ctx) error {
        tokenString := c.Get("Authorization")
        if tokenString == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
        }

        tokenString = tokenString[len("Bearer "):] // Remove "Bearer " prefix if present
        token, err := ValidateJWT(tokenString)
        if err != nil || !token.Valid {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
        }

        return c.Next()
    }
}
