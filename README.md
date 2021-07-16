# Fauna2go

A docker image that automatically initializes a [Fauna](https://docs.fauna.com/fauna/current/integrations/dev.html) database and key.

## Getting started

With docker installed, run:

```
docker run -p 1000:1000 -p 8443:8443 -p 8084:8084 -p 8085:8085 felipeqq2/fauna2go
```

A server will run on http://localhost:1000, serving the key for the most recently created database. A `POST` request will initialize a clean database, and serve its key.

http://localhost:8443 and http://localhost:8084 will serve FaunaDB normally.

A proxy will open on http://localhost:8085. All requests to this endpoint will be resolved to http://localhost:8084 (GraphQL endpoint), bypassing authentication (no need to set up "Authentication" header).
