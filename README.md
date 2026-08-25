# Sport Court Booking

## Project Description

Sport Court Booking is a RESTful API for managing sport courts, customer bookings, and booking payments. Users can create accounts, authenticate with JWT, browse and manage courts, create bookings, and record payments. Administrators can update booking statuses and access resources according to the authorization rules implemented by the API.

The service runs locally on `http://localhost:8080` by default. It can be deployed as a Go serverless function on Vercel.

## Technology Stack

- **Language:** Go `1.25.0`
- **Web framework:** Gin `1.12.0`
- **Database:** PostgreSQL
- **Database driver:** `lib/pq`
- **Database migrations:** `sql-migrate` with embedded SQL files
- **Configuration:** Viper, `.env` file, or system environment variables
- **Authentication:** JWT (`golang-jwt/jwt`) with bcrypt password hashing
- **Validation:** Gin request binding and validator tags

## Folder Structure

```text
.
├── configs/                 # Environment and application configuration
├── databases/
│   ├── connection/          # PostgreSQL connection setup
│   └── migration/           # Embedded database migration runner and SQL files
├── handler/                 # HTTP handlers and request/response mapping
├── helpers/
│   ├── common/              # Shared responses, validation, and password helpers
│   └── constant/            # Shared errors and constants
├── middlewares/             # JWT authentication middleware
├── model/
│   ├── dto/request/          # Incoming API request DTOs
│   ├── dto/response/         # Outgoing API response DTOs
│   └── entity/               # Database/domain entities
├── repository/              # Database access layer
├── routes/                  # HTTP route registration
├── service/                 # Business logic layer
├── go.mod                   # Go module and dependency declarations
└── main.go                  # Application entry point
```

## Flow

1. The application loads configuration from `.env` or the system environment.
2. A PostgreSQL connection is initialized using `DATABASE_URL` and `DB_ENGINE`.
3. Embedded SQL migrations are applied automatically at startup.
4. Repositories, services, handlers, and routes are constructed using dependency injection.
5. Gin starts the HTTP server on port `8080`.
6. Public authentication endpoints issue access tokens after successful authentication.
7. Protected endpoints require an `Authorization: Bearer <JWT_TOKEN>` header.
8. Requests flow through the middleware, handler, service, and repository layers before returning a JSON response.

Required configuration values:

```env
DATABASE_URL=<POSTGRES_CONNECTION_URL>
DB_ENGINE=postgres
jwt_secret_key=<JWT_SECRET>
```

Start the service with:

```bash
go run .
```

## Deploy to Vercel

Create a Vercel project from this repository and use the default settings. Vercel detects `api/index.go` as the serverless entry point and `vercel.json` forwards all routes to it.

Set these environment variables in the Vercel project settings:

```env
DATABASE_URL=<POSTGRES_CONNECTION_URL>
DB_ENGINE=postgres
jwt_secret_key=<JWT_SECRET>
```

The PostgreSQL database must be hosted on a service reachable from Vercel. Database migrations run when the function instance initializes. Because Vercel functions are serverless, do not use `go run .` as the Vercel start command.

## Available APIs

All paths are relative to `<API_BASE_URL>`. Unless marked **Public**, an endpoint requires a valid JWT bearer token.

### Authentication

| Method | Path | Access | Description | Request body |
| --- | --- | --- | --- | --- |
| `POST` | `/login` | Public | Authenticate a user and return login data | `email`, `password` |
| `POST` | `/sign-up` | Public | Create a user account | `email`, `password`, `re_type_password`, `name` |

### Courts

| Method | Path | Access | Description | Request body |
| --- | --- | --- | --- | --- |
| `POST` | `/api/courts` | JWT | Create a court | `name`, `type`, `price_per_hour`, `location` |
| `GET` | `/api/courts` | JWT | List courts | None |
| `GET` | `/api/courts/:id` | JWT | Get a court by ID | None |
| `PUT` | `/api/courts/:id` | JWT | Partially update a court | Any of `name`, `type`, `price`, `location` |
| `DELETE` | `/api/courts/:id` | JWT | Delete a court | None |

Supported court types are `futsal` and `badminton`.

### Bookings

| Method | Path | Access | Description | Request body |
| --- | --- | --- | --- | --- |
| `POST` | `/api/bookings` | JWT | Create a booking for the authenticated user | `court_id`, `booking_date`, `start_time`, `end_time` |
| `GET` | `/api/bookings` | JWT | List bookings for the authenticated user | None |
| `GET` | `/api/bookings/:id` | JWT | Get a booking by ID, subject to ownership or admin access | None |
| `PUT` | `/api/bookings/status/:id` | Admin JWT | Update a booking status | `status` |

Supported booking statuses are `confirmed`, `cancelled`, and `completed`.

### Payments

| Method | Path | Access | Description | Request body |
| --- | --- | --- | --- | --- |
| `POST` | `/api/payments` | JWT | Record a payment for a booking owned by the authenticated user | `booking_id`, `amount`, `method` |
| `GET` | `/api/payments/:booking_id` | JWT | Get payment details for a booking, subject to ownership or admin access | None |
| `GET` | `/api/payments` | JWT | List payments for the authenticated user | None |

Supported payment methods are `cash`, `transfer`, and `qris`.

> Request and response formats may evolve with the DTOs in `model/dto/request` and `model/dto/response`. Add `<REQUEST_EXAMPLE>` and `<RESPONSE_EXAMPLE>` here when an API contract is formally defined.
