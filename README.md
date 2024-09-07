# Google Contacts Birthday Notification

This is a simple service that upon starting sends a notification (as an email)
with upcoming birthdays. These birthdays, filtered for birthdays tomorrow or in
two weeks, are extracted from your Google contacts.

It is intended that the application runs once a day (e.g. via a chronjob) so that
you are daily informed of upcoming birthdays.

## Config

Credentials for the email provider and the Google OAuth2 data must be entered
in the `config.yml`, that is generated if it does not exist yet.
The config looks as follows:

```yml
mail:
  port: 465
  host: mail.com
  user: sender@mail.com
  sender: sender@mail.com
  receiver: receiver@qetz.de
  password: password
  tls: false
  secure: false
peopleApi:
  credentials:
    clientId: clientId.apps.googleusercontent.com
    projectId: projectId
    authUri: https://accounts.google.com/o/oauth2/auth
    tokenUri: https://oauth2.googleapis.com/token
    authProviderCertUrl: https://www.googleapis.com/oauth2/v1/certs
    clientSecret: clientSecret
    redirectUris:
      - http://localhost
  token:
    accessToken: accessToken
    tokenType: Bearer
    refreshToken: refreshToken
    expiry: 2024-09-01T19:01:56.8593584+02:00
```

### Google Credentials

Although this is obviously not the correct way to use OAuth, since this application
runs in a secure environment and Google does not offer a more appropriate method,
this solution was chosen.

To create all necessary Google credentials, you must first create a
project in Google Cloud.
From this, you will receive a `credentials.json`,
which entries you can insert in the `config.yml`.
After that, you need to generate an access and refresh token. One way to do this
is via the [token script](/scripts/token.go). You can implement this script
in a starting function and it will guide you through the authentication process.
The printed tokens must be entered in the `config.yml` as well.
As of today, there is no expiration for refresh tokens, so you should only have
to do this initial setup once.