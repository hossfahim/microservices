# Users Service API

This service manages drivers and passengers for the RideNow microservices application.

## Base URL

- Local development: `http://localhost:8081`
- Docker: `http://localhost:3000`

## Driver Endpoints

### Create Driver

Create a new driver. The driver will be set as available by default.

```bash
curl -X POST http://localhost:8081/drivers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe"
  }'
```

### Get All Drivers

Retrieve all drivers, optionally filtered by availability.

```bash
# Get all drivers
curl -X GET http://localhost:8081/drivers

# Get only available drivers
curl -X GET "http://localhost:8081/drivers?available=true"
```

### Update Driver Status

Update a driver's availability status.

```bash
curl -X PATCH http://localhost:8081/drivers/{driver_id}/status \
  -H "Content-Type: application/json" \
  -d '{
    "is_available": false
  }'
```

**Example:**

```bash
curl -X PATCH http://localhost:8081/drivers/507f1f77bcf86cd799439011/status \
  -H "Content-Type: application/json" \
  -d '{
    "is_available": false
  }'
```

## Passenger Endpoints

### Create Passenger

Create a new passenger.

```bash
curl -X POST http://localhost:8081/passengers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Johnson"
  }'
```

### Get All Passengers

Retrieve all passengers.

```bash
curl -X GET http://localhost:8081/passengers
```

### Get Passenger by ID

Retrieve a specific passenger by their ID.

```bash
curl -X GET http://localhost:8081/passengers/{passenger_id}
```

**Example:**

```bash
curl -X GET http://localhost:8081/passengers/507f1f77bcf86cd799439011
```

### Update Passenger

Update a passenger's information.

```bash
curl -X PUT http://localhost:8081/passengers/{passenger_id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Smith"
  }'
```

**Example:**

```bash
curl -X PUT http://localhost:8081/passengers/507f1f77bcf86cd799439011 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Smith"
  }'
```

### Delete Passenger

Delete a passenger by their ID.

```bash
curl -X DELETE http://localhost:8081/passengers/{passenger_id}
```

**Example:**

```bash
curl -X DELETE http://localhost:8081/passengers/507f1f77bcf86cd799439011
```

## Response Examples

### Driver Response

```json
{
  "id": "507f1f77bcf86cd799439011",
  "name": "John Doe",
  "is_available": true
}
```

### Passenger Response

```json
{
  "id": "507f1f77bcf86cd799439011",
  "name": "Alice Johnson",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

## Notes

- All IDs are MongoDB ObjectIDs (24-character hexadecimal strings)
- Timestamps are in ISO 8601 format (UTC)
- When creating a driver, the `is_available` field is automatically set to `true`
- When creating or updating a passenger, timestamps are automatically managed
