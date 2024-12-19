# TacTicsToe

TacTicsToe is a TicTacToe but with a twist: it's 4x4. To differentiate itself it has a additional rule. When you make the winning move (i.e. place the third tic in a row), a opponent's tic has to be directly adjacent. And diagonals don't count.

## Let's talk tech

TacTicsToe is made using the standard library from GoLang, which I use to handle the queue 'system' and to match players up with another player. When a game is start, a websocket connection will be made using the [Gorilla/Websockets](https://github.com/gorilla/websocket) library. 

The frontend is written in plain HTML and JS. The reason I started this project was because I wanted to familiarize myself with GoLang, and so good UI wasn't really my primary goal.

If you look at the source code, you'll find code written in Rust that can be used to make a bot. I had it implemented, but I thought that it spoiled the game, because there is a gamebreaking tactic. 

## Getting started

1. Install [GoLang](https://go.dev/)

2. Clone the repo
```bash
git clone https://github.com/svrem/tacticstoe-go.git
cd tacticstoe-go
```

3. Start the server:
```bash
go run .
```

3. Fill the .env file (for Google OAuth2 Creds visit [this](https://developers.google.com/identity/protocols/oauth2))
```
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URI=http://localhost:8080/auth/callback/google

JWT_SECRET=
```

4. Go to `http://localhost:8080`