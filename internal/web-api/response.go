package webapi

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"google.golang.org/grpc/codes"
)

// OKStatus returns 'ok' status in response
func OKStatus() string {
	return "ok"
}

// ErrorStatus returns 'error' status in response
func ErrorStatus() string {
	return "error"
}

// APIError is a DTO for single error in response
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}

// ResponseBody is a DTO for API response
// swagger:response ResponseBody
type ResponseBody struct {
	Data   interface{} `json:"data,omitempty"`
	Status string      `json:"status"`
	Errors []APIError  `json:"errors,omitempty"`
}

// APIResponse is internal DAO to build API response
type APIResponse struct {
	Header int
	Body   ResponseBody
}

// OKResponse prepares simple OK response with data
func OKResponse(ctx context.Context, w http.ResponseWriter, data interface{}) {
	apiResponse := &APIResponse{
		Header: http.StatusOK,
		Body: ResponseBody{
			Data:   data,
			Status: OKStatus(),
			Errors: nil,
		},
	}
	makeResponse(ctx, w, apiResponse)
}

// BadInputResponse prepares 400 Bad Request response with one error
func BadInputResponse(ctx context.Context, w http.ResponseWriter, msg string) {
	apiResponse := &APIResponse{
		Header: http.StatusBadRequest,
		Body: ResponseBody{
			Status: ErrorStatus(),
			Errors: []APIError{
				{
					Code:    "bad_input",
					Message: msg,
				},
			},
		},
	}
	makeResponse(ctx, w, apiResponse)
}

// UnauthorizedResponse prepares 401 Unauthorized Error response with one error
func UnauthorizedResponse(ctx context.Context, w http.ResponseWriter, msg string) {
	apiResponse := &APIResponse{
		Header: http.StatusUnauthorized,
		Body: ResponseBody{
			Status: ErrorStatus(),
			Errors: []APIError{
				{
					Code:    "unauthorized_error",
					Message: msg,
				},
			},
		},
	}
	makeResponse(ctx, w, apiResponse)
}

// ForbiddenResponse prepares 403 Forbidden Error response with one error
func ForbiddenResponse(ctx context.Context, w http.ResponseWriter, msg string) {
	apiResponse := &APIResponse{
		Header: http.StatusForbidden,
		Body: ResponseBody{
			Status: ErrorStatus(),
			Errors: []APIError{
				{
					Code:    "forbidden_error",
					Message: msg,
				},
			},
		},
	}
	makeResponse(ctx, w, apiResponse)
}

// InternalErrorResponse prepares 500 Internal Server Error response with one error
func InternalErrorResponse(ctx context.Context, w http.ResponseWriter, msg string) {
	apiResponse := &APIResponse{
		Header: http.StatusInternalServerError,
		Body: ResponseBody{
			Status: ErrorStatus(),
			Errors: []APIError{
				{
					Code:    "github.com/ivanfomichev/bank-app/internal_error",
					Message: msg,
				},
			},
		},
	}
	makeResponse(ctx, w, apiResponse)
}

// NotFoundResponse prepares 404 Internal Server Error response with one error
func NotFoundResponse(ctx context.Context, w http.ResponseWriter, msg string) {
	apiResponse := &APIResponse{
		Header: http.StatusNotFound,
		Body: ResponseBody{
			Status: ErrorStatus(),
			Errors: []APIError{
				{
					Code:    "not_found",
					Message: msg,
				},
			},
		},
	}
	makeResponse(ctx, w, apiResponse)
}

func makeResponse(ctx context.Context, w http.ResponseWriter, apiResponse *APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiResponse.Header)
	if err := json.NewEncoder(w).Encode(apiResponse.Body); err != nil {
		log.Printf("write response failed")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// HTTPStatusFromCode maps GRPC status code to HTTP
func HTTPStatusFromCode(code codes.Code) int {
	var httpCode int
	switch code { //nolint:exhaustive
	case codes.OK:
		httpCode = http.StatusOK
	case codes.Canceled:
		httpCode = http.StatusRequestTimeout
	case codes.InvalidArgument, codes.FailedPrecondition, codes.OutOfRange:
		httpCode = http.StatusBadRequest
	case codes.DeadlineExceeded:
		httpCode = http.StatusGatewayTimeout
	case codes.NotFound:
		httpCode = http.StatusNotFound
	case codes.AlreadyExists, codes.Aborted:
		httpCode = http.StatusConflict
	case codes.PermissionDenied:
		httpCode = http.StatusForbidden
	case codes.Unauthenticated:
		httpCode = http.StatusUnauthorized
	case codes.ResourceExhausted:
		httpCode = http.StatusTooManyRequests
	case codes.Unimplemented:
		httpCode = http.StatusNotImplemented
	case codes.Unavailable:
		httpCode = http.StatusServiceUnavailable
	default:
		httpCode = http.StatusInternalServerError
	}
	return httpCode
}
