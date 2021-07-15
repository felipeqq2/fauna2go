# Fauna2go

A docker image that automatically initializes a [Fauna](https://docs.fauna.com/fauna/current/integrations/dev.html) database and key.

## Getting started

With docker installed, run:

```
docker run -p 1000:1000 -p 8443:8443 -p 8084:8084 felipeqq2/fauna2go
```

A server will run on http://localhost:1000, serving the key for the most recently created database. A `POST` request will initialize a clean database, and serve its key.
