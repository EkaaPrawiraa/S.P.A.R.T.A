Run migrations using:

golang-migrate migrate \
  -path migrations \
  -database "$DATABASE_URL" up
