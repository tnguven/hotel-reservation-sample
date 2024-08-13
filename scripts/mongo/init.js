print('Start ##############################################################');

db = db.getSiblingDB('hotel-io');
db.createCollection('users');
db.createCollection('hotels');

print('END #################################################################');
