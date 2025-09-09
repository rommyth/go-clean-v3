package handler

// Handlers groups all HTTP handlers for easy routing
type Handlers struct {
    UserHandler *UserHandler
    AuthHandler *AuthHandler
    // Add more here as you create them:
    // TodoHandler      *TodoHandler
    // ProductHandler   *ProductHandler
}