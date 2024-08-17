print("Start ##############################################################");

db = db.getSiblingDB("hotel_io");
db.createCollection("users");
db.createCollection("hotels");
db.createCollection("bookings");
db.createCollection("rooms");

print("END #################################################################");
