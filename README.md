## Go Starter

A starter setup for a REST api in Go, MongoDB, and Docker

### Add .env at project root with the following keys and add your values

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

`Indexes: email_1 creationTime_1 lastActive_1`

`Connected to DATABASE: DB_NAME`

`Listening on 5000`

Visit `http://localhost:5000` for instructions
