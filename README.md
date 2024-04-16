# Generic Dating
This is the backend of a generic dating app written in Go. It has been mainly an exercice for me to see how I fare with Go. It uses JWT to authenticate users with an expiration of 24h. Users are randomly created with `/user/create` endpoint.
The rest of endpoints expect you to provide a bearer token you can get by logging in with a valid user at `/user/login`.
For more info on the currenrly implemented endpoints, check `API Routes` section.


## How to build and Run
- Duplicate `.env.example` into `.env`. Change variables if desired, not needed.
- Run `docker-compose up --build -d` and that should spawn a mysql database and the go api. The first wants to bind port `3306` and the second `3000` to your local ones. Be sure they are not in use or change them in `docker/docker-compose.yml` accordingly.


## API Routes

* GET `/user/create` -> Generates a random new user
```
{
	"ID": <int>,
	"CreatedAt": <datetime>,
	"UpdatedAt": <datetime>,
	"DeletedAt": null,
	"Email": <string>,
	"Password": <string>,  // This will be the unencrypted password
	"Name": <string>,
	"Gender": <string>,
	"Age": <int>
}
```

* POST `/user/login` -> Expects a JSON with `email` and `password`
```
{
    "token": <string>
}
```
This token will need to be provided on the rest of secured endpoints as a `Bearer` token on the headers as Authentication

* GET `/discover` -> Returns all swippable profiles. It accepts url query params such as `min_age`, `max_age` and `gender` that allow for filtering. These are sorted by `distance` and a hidden `attractiveness score`.
```
{
    "results": [
        {
            "id": <integer>,
            "name": <string>,
            "gender": <string>,
            "age": <integer>,
            "distanceFromMe" <int>
        },
        ...
    ]
}
```

* POST `/swipe/<id>` -> Expects a json such as `{"preference": <bool>}`
```
{
    "matched": <bool>,
    "matchID": <integer> // Only if matched is True
}
```

## Improvements

- **TESTING!** Although I have manually tested the code with postman, good code should ALWAYS come with a set of tests
- `Signup` for users, so we can create users and not rely only on the fake generation of them
- `Logout` and `refresh tokens` to have a better user experience
- `Delete` user endpoint. Currently there is not way to remove a user other than connecting to the DB directly
- Add `swagger` so the API is self documented within the code and can serve its own schema. Something like gin-swagger should be good [gin-swagger](https://github.com/swaggo/gin-swagger)
- We could add a `max_distance` filtering parameter to discovery
- Implement a way to notify the other user in a match that a match has been made from one of its previous swipes
- `Matches` endpoint so a user can see with whom it matched so they can see each other email? Maybe implement a chat system?
- Allow user to select desired units for distance, now it defaults to Kms and cannot be changed
- Integrate [google maps api](https://github.com/googlemaps/google-maps-services-go) to get latitute longitude given an address. I have not done it as that requires an API key to work and that would mean if anyone wanted to run this code would need to generate one for themselves.
