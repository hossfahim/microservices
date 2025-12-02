db = db.getSiblingDB("admin");
db.auth("admin_root", "password_root");

db = db.getSiblingDB("ridenow_rides");

db.createUser({
  user: "user_app",
  pwd: "strong_app_password",
  roles: [{ role: "readWrite", db: "ridenow_rides" }],
});

db.createCollection("rides");

db.rides.insertMany([
  {
    passenger_id: "passenger-001",
    driver_id: "driver-001",
    from_zone: "Downtown",
    to_zone: "Airport",
    price: 25.50,
    status: "ASSIGNED",
    payment_status: "PENDING",
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    passenger_id: "passenger-002",
    driver_id: "driver-002",
    from_zone: "Suburbs",
    to_zone: "City Center",
    price: 18.75,
    status: "IN_PROGRESS",
    payment_status: "PENDING",
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    passenger_id: "passenger-003",
    driver_id: "driver-003",
    from_zone: "Beach",
    to_zone: "Hotel District",
    price: 32.00,
    status: "COMPLETED",
    payment_status: "CAPTURED",
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    passenger_id: "passenger-004",
    driver_id: "driver-004",
    from_zone: "University",
    to_zone: "Train Station",
    price: 15.25,
    status: "CANCELLED",
    payment_status: "REFUNDED",
    created_at: new Date(),
    updated_at: new Date(),
  },
]);

