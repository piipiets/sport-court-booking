# Sport Court Booking

Sport Court Booking is a RESTful API for managing sport courts, customer bookings, and booking payments. Users can create accounts, authenticate with JWT, browse and manage courts, create bookings, and record payments. Administrators can update booking statuses and access resources according to the authorization rules implemented by the API.

## Live URL

- Production: https://sport-court-booking-six.vercel.app
- Local: http://localhost:8080

## Tech Stack

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
├── api/                     # Vercel serverless entry point
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
JWT_SECRET_KEY=<JWT_SECRET>
```

Start the service with:

```bash
go run .
```

## Database Structure

The project uses PostgreSQL with migration files under `databases/migration/sql_migration/`. The schema includes the core booking flow: users, courts, bookings, payments, and review records.

### `users`

Stores system users and their roles.

| Column | Type | Description |
| --- | --- | --- |
| `id` | `BIGSERIAL` | Auto-increment primary key for each user. |
| `name` | `VARCHAR(150)` | Full name of the user. Required. |
| `email` | `VARCHAR(255)` | Unique email address used for login and identity. Required. |
| `password` | `VARCHAR(255)` | Hashed password value. Required. |
| `role` | `VARCHAR(20)` | User role. Allowed values: `user`, `admin`. Default is `user`. |
| `created_at` | `TIMESTAMPTZ` | Account creation timestamp. Automatically set with `now()`. |

Constraints:
- `email` is unique.
- `role` is restricted to `user` or `admin`.

### `courts`

Stores sport court information.

| Column | Type | Description |
| --- | --- | --- |
| `id` | `BIGSERIAL` | Auto-increment primary key for each court. |
| `name` | `VARCHAR(150)` | Court name. Required. |
| `type` | `VARCHAR(50)` | Court category or type, such as `futsal`, `badminton`, or other configured values. |
| `price_per_hour` | `NUMERIC(12, 2)` | Hourly rental price. Must be greater than or equal to `0`. |
| `location` | `VARCHAR(255)` | Court location/address. Required. |
| `created_at` | `TIMESTAMPTZ` | Creation timestamp. Automatically set with `now()`. |

Constraints:
- `price_per_hour` cannot be negative.

### `bookings`

Records a reservation made by a user for a specific court.

| Column | Type | Description |
| --- | --- | --- |
| `id` | `BIGSERIAL` | Auto-increment primary key for each booking. |
| `user_id` | `BIGINT` | Foreign key to `users.id`. Indicates the booking owner. |
| `court_id` | `BIGINT` | Foreign key to `courts.id`. Indicates the reserved court. |
| `booking_date` | `DATE` | Booking date. Required. |
| `start_time` | `TIME` | Start time of the booking. Required. |
| `end_time` | `TIME` | End time of the booking. Required. Must be later than `start_time`. |
| `status` | `VARCHAR(20)` | Current booking status. Allowed values: `pending`, `confirmed`, `cancelled`, `completed`. Default: `pending`. |
| `total_price` | `NUMERIC(12, 2)` | Total cost of the reservation. Must be greater than or equal to `0`. |
| `payment_deadline` | `TIMESTAMPTZ` | Payment due date/time for the booking. |
| `created_at` | `TIMESTAMPTZ` | Booking creation timestamp. Automatically set with `now()`. |

Constraints:
- `user_id` references `users(id)` with cascade delete.
- `court_id` references `courts(id)` with restrict delete.
- `end_time > start_time` must hold.

### `payments`

Stores payment records linked to a booking.

| Column | Type | Description |
| --- | --- | --- |
| `id` | `BIGSERIAL` | Auto-increment primary key for each payment. |
| `booking_id` | `BIGINT` | Unique foreign key to `bookings.id`. A booking can only have one payment record. |
| `amount` | `NUMERIC(12, 2)` | Paid amount. Must be greater than or equal to `0`. |
| `method` | `VARCHAR(20)` | Payment method. Allowed values: `cash`, `transfer`, `qris`. |
| `status` | `VARCHAR(20)` | Payment state. Allowed values: `unpaid`, `paid`, `refunded`. Default: `unpaid`. |
| `paid_at` | `TIMESTAMPTZ` | Timestamp when the payment was completed. Nullable until paid. |

Constraints:
- `booking_id` is unique.
- `booking_id` references `bookings(id)` with cascade delete.

### `reviews`

Stores a user review for a completed booking.

| Column | Type | Description |
| --- | --- | --- |
| `id` | `BIGSERIAL` | Auto-increment primary key for each review. |
| `booking_id` | `BIGINT` | Foreign key to `bookings.id`. One booking can only be reviewed once. |
| `user_id` | `BIGINT` | Foreign key to `users.id`. Indicates the reviewer. |
| `rating` | `SMALLINT` | Rating score from `1` to `5`. |
| `comment` | `TEXT` | Optional review text/content. |
| `created_at` | `TIMESTAMPTZ` | Review creation timestamp. Automatically set with `now()`. |

Constraints:
- `rating` must be between `1` and `5`.
- Each `booking_id` appears only once in the table.

### Relationship Summary

```text
users ──< bookings >── courts
  │                 │
  │                 └──< payments
  │
  └──< reviews
```

- One user can create many bookings.
- One court can have many bookings.
- One booking can have one payment record.
- One booking can have one review.

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
