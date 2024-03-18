package model

type User struct {
    ID       int    `json:"user_id"`
    Username string `json:"username"`
    Password string `json:"password"`
    RoleID   int    `json:"role_id"`
}