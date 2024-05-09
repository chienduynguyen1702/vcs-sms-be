# VCS Server Management System
## Technical
- Golang
- Docker
- ...
## Folder Stucture
### Gateway
```
Path				Explain
.
├───main.go			Main file to start app
├───configs/		Config env, db connection, migration
├───routes/			Route /api to controller
├───middleware/		Middleware between route and controller
├───controllers/	Controller layer 
├───services/		Service layer
├───repositories/	Models's repositories
├───models/			Defind models with GORM
├───init/			Initialize controllers <- services <- repos <- models
├───dtos/			Data to object
├───utilities/		Some utilities functions
├───docker/			Dockerfile build for this app and hosting local db
├───docs/			Swagger Spec
└───scripts/		Some scripts using while develop
└───

```
