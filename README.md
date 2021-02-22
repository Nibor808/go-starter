# Go Starter

A starter setup for a REST api with React, Go, MongoDB, and Docker.

No styling or layout for client yet!

## Sessions

Tokens and Sessions expire after 10 minutes.

This can be changed in `api/utils/db.go line 55 -> ttl := int32(600)`

Session ttl is refreshed by hitting an endpoint.

### Add .env in api/ with the following keys and add your values

```
DB_CONN=mongodb://host.docker.internal:27017

DB_NAME=<your-database-name>

DEV_URL=http://localhost:5000

SESSION_SECRET=<anythingyoulike>

JWT_KEY=<anythingyoulike>

ADMIN_EMAIL=<your-email>

MAIL_PASS=<your-email-password>

MAIL_HOST=<your-mail-host> eg smtp.gmail.com
```

### Start app

`docker-compose up --build`

In your console you should see:

api/

`Indexes: email_1 creationTime_1 lastActive_1`

`Connected to DATABASE: DB_NAME`

`Listening on 5000`

client/

`Compiled successfully!`

Visit `http://localhost:3000` to create a user or sign in
