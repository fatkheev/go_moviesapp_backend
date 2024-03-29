package model

type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

type SuccessResponse struct {
    Data interface{} `json:"data"`
}