# marketplace

Available by the next base URL: http://194.163.167.58:8080

`GET /points/:userId` - returns calculated stats for a user by it's ID; _Bearer token_ of the user is required;

`POST /points` - calculates new stats by taking into account `points` from the request and `points` from the so-called base file; _Bearer token_ of the user is required;

`POST /new-user` - adds new user to the system. Data must be provided in request body, in the next format:

```
{
  "name": "Arya",
  "password": "winter",
  "role": "user"
}
```

All fields are required.
At the moment 2 roles are supported:
- `admin`;
- `user`.

New user with role `user` may be added as is. For new admin - Bearer token of an exsisting admin is required.

`POST /login` - log in via _Basic auth_ using previousy `username` and `password`

`POST /logout` - processes received _Bearer_ token to log out

At the moment general approach is next.
By `POST /points` service receives _points_ data that must be provided in an array of objects with the next structure:

```
[
  {
    "rate": 0.8154823075825571,
    "timestamp": 1616761300265
  }
]
```

Points from the base file are being shifted to the amount of recieved points and new points are being placed at the beginnig of the collection from the file.
