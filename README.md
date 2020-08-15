# bookish-palm-tree

Following a [go oauth tutorial](https://www.loginradius.com/engineering/blog/google-authentication-with-golang-and-goth/).

## Authentication provider links

Create the appropriate oauth keys for your application, and set the values in the environment.

[Google](https://console.developers.google.com/).

The trick appears to be to have `http://lvh.me:3000/auth/google/callback` in the authorised redirect URI's.
