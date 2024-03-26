# Chirpy Web Server

This is a basic project that builds a web server in Go to mock the backend 
for a Twitter-like website, and is mainly intended to help me learn how to deal
with the various types of http requests and to get used to making basic CRUD
operations with a mock database (which for this project, is just a JSON file).

This project is not intended to be used publicly, but, should you want to play
with it, you can try it out buy cloning the repository and running the usual
`go build` command followed by `./chirpy`.

The full API for the project requires a JWT and an API key stored in an .env file
that, as a part of standard best-practices, have not been included here. 

What follows is a list of the allowable operations within the mock server:

**/healthz**

* GET: returns if the server is running

**/reset** 

* GET: resets the number of site visits

**/api/chirps**

* GET: returns all of the chirps in the database
    - All GET requests can be modified with the `author_id` and `sort`
      queries to return chirps from a particular author and sort them

* POST: Creates a chirp if a user is authenticated (they must have a JWT from
  from logging in at /api/login)


    **/api/chirps/{id}**
    
    * GET: retrieves a chirp with a particular id

    * DELETE: deletes a chirp with a particular id. The user making the
      request must be authenticated with a current access token

**/api/users**

* POST: creates a new user
* PUT: updates a user. The user making the request must be authenticated with
  a current access token

**/api/login**

* POST: logs a user in and issues them an access (1d) and refresh (60d) token

**/api/refresh**

* POST: re-issues a user an access token (requires a current, valid refresh token)

**/api/revoke**

* POST: revokes the refresh token for the user

**/api/polka/webhooks**

* POST: handles requests from the polka (mock payment processor) webhook. If 
  the event is an upgrade event, then the user is upgraded to Chirpy Red.