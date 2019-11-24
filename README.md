# beams-pushnotifications-go

Introducing roles and permissions to your chat app built with Pusher Chatkit... 

#### Prerequisites

- Go ( `>=1.12` )
- A [Pusher Beams](https://dash.pusher.com) application.
- `ngrok`

#### Getting started

You have to clone this repository before moving on `git clone git@github.com:adelowo/beams-push-notifications-go.git`.

To run the backend server, you will need to run

```
$ cd server
$ go build
$ ./server
```

> You will need to edit the file located at `server/.env`

To run the Android application. You need to open it in Android Studio and hit the "Play" button

> You  will need to updae line 24 and 25 of [MainActivity.kt](https://github.com/adelowo/beams-pushnotifications-go/blob/master/PusherBeamsSlackWebhook/app/src/main/java/com/example/pusherbeamsslackwebhook/MainActivity.kt)

## Built With

- [Pusher Beams](https:dash.pusher.com) - APIs to build Push notifications easily,
- Golang
- Kotlin
