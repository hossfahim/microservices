db = db.getSiblingDB("admin");
db.auth("admin_root", "password_root");

db = db.getSiblingDB("ridenow_users");

db.createUser({
  user: "user_app",
  pwd: "strong_app_password",
  roles: [{ role: "readWrite", db: "ridenow_users" }],
});

db.createCollection("drivers");

db.drivers.insertMany([
  {
    name: "Rick Sanchez",
    is_available: true,
  },
  {
    name: "Morty Smith",
    is_available: true,
  },
  {
    name: "Summer Smith",
    is_available: false,
  },
  {
    name: "Beth Smith",
    is_available: false,
  },
]);

db.createCollection("passengers");

db.passengers.insertMany([
  {
    name: "Jerry Smith",
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    name: "Unity",
    created_at: new Date(),
    updated_at: new Date(),
  },
]);