# Folloxers

Since Elon is desperately trying to dig himself out of a financial hole, he has decided to restrict the free twitter API heavily. It has been seen time and time again that when you offer a free service, bait folks in, and then put the squeeze on them, they get a little mad.

Who could be surprised that mr tony stark le epic bacon doggo is following in the steps of Ballmer before him.


## How to Set Up

- Do `cp sample.env .env`
- Go to this page https://x.com/username/followers
- Open your Developer Panel to the Network tab
- Scroll down on the page until you see a GET request to /i/api/graphql/{something}/Followers
- Click on that and you'll see many sensitive fields you need to copy out to the `.env` file.

Env Table:
|Env Key|How to Find It|
|---|---|
|`CURSOR_FLOOR`|Copy the first value out of the variables.cursor URL param. There will be 2 numbers separated by a `\|` character.|
|`CURSOR_CIEL`|Copy the second value out of the variables.cursor URL param. There will be 2 numbers separated by a `\|` character.|
|`X_CSRF_TOKEN`|Value of header named `x-csrf-token`|
|`X_BEARER_AUTH_TOKEN`|Value of header named `authorization`, WITHOUT the `Bearer ` prefix.|
|`X_AUTH_TOKEN`|Cookie Header > `auth_token=` value|
|`X_KDT`|Cookie Header > `kdt=` value|
|`X_USER_ID`|URL Param `variables` > `userId`|

Once you have all these, it's basically set-and-forget. I haven't seen the session expire, ever.


## Limitations

The API is rate limited to 50 requests per 5ish minutes. Good job guys! That's actually a pretty good control. Hats off. 

But if you're patient, the script has logic to respect and follow the `429` wait parameters.
