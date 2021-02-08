## Go Starter

A starter setup for a REST api in Go, MongoDB, and Docker

### Set up a Sendgrid account for sending email

[Sendgrid](https://sendgrid.com/)

### Add .env at project root with the following keys and add your values

```
DB_CONN=mongodb://host.docker.internal:27017

DB_NAME=<your-database-name>

DEV_URL=http://localhost:5000

SENDGRID_API_KEY=<from-your-sendgrid-account>

ADMIN_EMAIL=<your-email>

SESSION_SECRET=<anythingyoulike>

JWT_KEY=<anythingyoulike>
```

### Start app

`docker-compose up --build`

In your console you should see:

`Indexes: email_1 creationTime_1 lastActive_1`

`Connected to DATABASE: DB_NAME`

`Listening on 5000`

Visit `http://localhost:5000` for instructions
