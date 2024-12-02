# Descriptions
This is a repo for me to play around with.

The initial idea came to mind when I noticed I didn't resume any of my project due to me keep finding interesting things in addition to my laziness an my inability to work on multiple project (╥﹏╥). Plus rewriting the same stuff that I would need (logs, db, etc.) an maintaining them in multiple repos is a pain.

With that in mind this repo will be build with the mind of a monolith in mind (since I'm lazy) but have the ability to be easily separated into a microservice after I figure out on how to to do it. Would probably happen by playing around in the application layer which is located in the app folder. Probably would clode the api folder, change the Application struct to include connections to other services which involves `map[any]any` and add the process of sending it to another microservice. Would probably need to add another interface in ports specifically for a service too. 

## Architecture
The repo uses Ports and Adapters (or Hexagonal Architecture). Would probably explain it here if got time in the future.

### File Structure
```
- cmd/{name} -> app binary
- commmon
  - config -> config files and config method
  - schema -> db migrations
  - scripts -> external scripts that will be used
    - data.sql -> test data
    - install.sh -> install script for dev
    - request.http -> postman but in http file (use any rest client personally use [rest.nvim](https://github.com/rest-nvim/rest.nvim/tree/main) )
- internal
  - adapters -> handling connections
    - primary -> incoming stuff such as incoming http request (handled by http server)
    - secondary -> outgoing stuff such as db connections (handled by drivers)
  - applications -> handle flows for services **(except for core)**
    - core -> core business logic such as how encrytion works, or password validation
  - logger -> logs
  - models -> structs, models, constants for stuff
  - utils -> utility function (it's bad I know but I suck at naming. Good luck future kids of mine)
  - test -> contains end to end testing **(NOT READY SO DON'T BOTHER COPYING)**
- Makefile -> make file
- config.json -> the actual config
  
```

### Stack
- [go-chi](https://github.com/go-chi/chi/): for http routes
- [validator](https://github.com/go-playground/validator): json validation
- [golang-migrate](https://github.com/golang-migrate/migrate): db migration
- [uuid](https://github.com/google/uuid): uuid generation
- [pgx](https://github.com/jackc/pgx/v5): PostgreSQL driver
- [zerolog](https://github.com/rs/zerolog): structured logging
- [viper](https://github.com/spf13/viper): read configs
- [test-container](https://github.com/testcontainers/testcontainers-go): test container for integration testing
- [crypto](https://golang.org/x/crypto): cryptography
- [lumberjack](https://gopkg.in/natefinch/lumberjack.v2): rolling files (automatically write new file for each day)
