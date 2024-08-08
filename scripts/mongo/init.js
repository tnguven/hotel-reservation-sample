print('Start ##############################################################');

db = db.getSiblingDB('hotel-io');
// db.createUser({
//   user: 'admin',
//   pwd: 'secret',
//   roles: [{ role: 'dbOwner', db: 'hotel-io' }],
// });
db.createCollection('users');
db.createCollection('hotels');

print('END #################################################################');
