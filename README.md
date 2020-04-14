## Go Starter
a starter setup for a REST api in golang and mongodb

### Set up a Sendgrid account for sending email 
[Sendgrid](https://sendgrid.com/)

### Add .env at project root with the following keys and add your values
```
DB_CONN=mongodb://localhost:27017

DB_NAME=<your-database-name>

SENDGRID_API_KEY=<from your sendgrid account>

ADMIN_EMAIL=<your email>
```

### Start app (from inside cwd)
```go run main.go```

you should see ```Connected to DATABASE: DB_NAME``` in your console

visit ```http://localhost:8080```