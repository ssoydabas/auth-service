# Authentication Service API

A comprehensive authentication service providing user management, authentication, and authorization capabilities.

## Overview

This service provides a RESTful API for managing user accounts, authentication, and authorization. It's built with Go and uses JWT for secure authentication.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

The API uses Bearer token authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your_token>
```

## API Endpoints

### Account Management

#### Create Account
- **POST** `/accounts`
- Creates a new user account
- Required fields:
  - email
  - first_name (2-50 characters)
  - last_name (2-50 characters)
  - password (min 8 characters)
  - phone

#### Authenticate Account
- **POST** `/accounts/authenticate`
- Authenticates user and returns JWT token
- Required fields:
  - password (min 8 characters)
- Optional fields:
  - email
  - phone

#### Get Current Account
- **GET** `/accounts/me`
- Returns current user's account details
- Requires authentication

#### Get Account by ID
- **GET** `/accounts/{id}`
- Returns account details for specified ID

### Password Management

#### Request Password Reset
- **POST** `/accounts/set-reset-password-token`
- Initiates password reset process
- Accepts either email or phone

#### Reset Password
- **POST** `/accounts/reset-password`
- Resets password using token
- Required fields:
  - token
  - password (min 8 characters)
  - confirm_password (min 8 characters)

### Email Verification

#### Verify Email
- **POST** `/accounts/verify-email`
- Verifies email address using token
- Required fields:
  - token

## Response Formats

### Success Response
```json
{
"data": {
// Response data
}
}
```

### Error Response
```json
{
"code": 400,
"message": "Error message"
}
```

### Validation Error Response
```json
{
"code": 400,
"message": "Validation error",
"errors": [
{
"field": "field_name",
"message": "validation message"
}
]
}
```

## HTTP Status Codes

- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `404` - Not Found
- `500` - Internal Server Error

## Documentation

For more detailed information about specific endpoints:
- Account Operations: https://docs.example.com/accounts
- Authentication Process: https://docs.example.com/authentication
- Authorization Guide: https://docs.example.com/authorization

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/ssoydabas/auth-service/LICENSE) file for details.

## Contact

- **Developer**: Sertan Soydabas
- **Email**: ssoydabas41@gmail.com
- **GitHub**: https://github.com/ssoydabas