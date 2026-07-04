# make
a web-based platform for managing safety seminars and special-use room access

## documentation
Any hosted instance has swagger API documentation available at the endpoint `/swagger/index.html`

## environment variables

make is configured entirely through environment variables. These can be supplied using a `.env` file or passed directly to the executable or Docker container.

| Variable            | Description |
|---------------------|-------------|
| `MAKE_HOST`         | Public URL of the application (ex `http://localhost:8080`) |
| `MAKE_OIDC_ID`      | OIDC public ID |
| `MAKE_OIDC_SECRET`  | OIDC secret |
| `MAKE_LDAP_BIND_DN` | LDAP bind DN used for directory queries |
| `MAKE_LDAP_PASS`    | Password for the LDAP bind account |
| `MAKE_DB_HOST`      | Database hostname ex (`postgres.example.com`) |
| `MAKE_DB_NAME`      | Database name |
| `MAKE_DB_USER`      | Database username |
| `MAKE_DB_PASS`      | Database password |
| `MAKE_ENABLE_WRITE` | Enables the worker and will add the configured LDAP roles to users upon training submission |
| `DEV`               | Exables frontend VITE proxy (only for local development) |

Boilerplate .env:
```
MAKE_OIDC_ID=
MAKE_OIDC_SECRET=
MAKE_HOST=
MAKE_LDAP_BIND_DN=
MAKE_LDAP_PASS=
MAKE_DB_HOST=
MAKE_DB_NAME=
MAKE_DB_USER=
MAKE_DB_PASS=
MAKE_DISABLE_WRITE=
DEV=
```

## worker
make features a worker responsible for synchronizing user training status with LDAP group memberships. The worker runs every 10 minutes and whenever a training is submitted.

When the worker is ran, it will perform the following tasks:
1. Fetch all user beloning the the LDAP groups associated with each area
2. Verifies that each user's training for that area is still valid
3. Removes the corresponding LDAP role if the user's training is no longer valid
4. Fetches all users who have completed a training, and determines which areas they have access to
5. Adds any missing LDAP groups idempotently for users with access to an area

A user's training can become valid for several reasons:
- The training was revoked by an admin
- The training has expired

The worker will NOT attempt to modify users who inherit group memberships hierarchically through LDAP (FreeIPA nested memberships). Only users with directly managed group memberships are modified (usually manually entered through FreeIPA or through make). Hierarchically managed group memberships always take precedence over those managed by make.

## building for production

make utilizes Go's embed features to allow for the entire project, including database 
migrations and the entire frontend, to be embedded inside the built executable. 
Once built, nothing other than the binary is needed to run the project.

### docker

The docker image automates building the frontend, backend, and creates a minimal docker 
image (37MB). the internal container port 8080 was chosen to allow for rootless 
operation of the container.

To build the docker image, run the following:
```bash
docker build --tag make .
```

Once built, the image can be run with:
```bash
docker run --env-file .env -p 8080:8080 make
```

Or, you can define the env variables manually:
```bash
docker run -p 8080:8080 \
  -e MAKE_OIDC_ID="$MAKE_OIDC_ID" \
  -e MAKE_OIDC_SECRET="$MAKE_OIDC_SECRET" \
  -e MAKE_HOST="$MAKE_HOST" \
  -e MAKE_LDAP_BIND_DN="$MAKE_LDAP_BIND_DN" \
  -e MAKE_LDAP_PASS="$MAKE_LDAP_PASS" \
  -e MAKE_DB_DRIVER="$MAKE_DB_DRIVER" \
  -e MAKE_DB_HOST="$MAKE_DB_HOST" \
  -e MAKE_DB_NAME="$MAKE_DB_NAME" \
  -e MAKE_DB_USER="$MAKE_DB_USER" \
  -e MAKE_DB_PASS="$MAKE_DB_PASS" \
  -e DEV="$DEV" \
  makedotcsh
```

### binary

To manually build a binary, you must first build the frontend. In the `web/` directory, run:
```bash
npm run build
```

Now, back in the project root, build the backend. go will automatically embed
the frontend into the executable
```bash
go build .
```

Or, if you wish for a smaller binary:
```bash
CGO_ENABLED=0 go build -ldflags="-s -w" .
```

## local development

Please ensure you meet the following requirements:
 - Node v20+
 - go 1.26+
   - air (`go install github.com/air-verse/air@latest`)
 - PostgreSQL database

make uses a pretty cool development method, where the Vite development frontend
is proxied through the gin HTTP server. This allows for all of Vite's development
features to still work while keeping the entire project running through the http
server, which handles authentication and keeps everything on the same port.

### starting development webserver
For development, enter `web/` and complete the following:

- install web dependencies
  ```bash
  npm i
  ```

- run development server
  ```bash
  npm run dev
  ```

The development server will begin at port 5173. Do not use this port to access the project.

### starting backend
In a new terminal, complete the following:

- install dependencies
   ```bash
   go mod download
   ```

- begin development backend server
   ```bash
   air
   ```

## LICENSE

Copyright (C) 2026 Tyler Severino  
Developed for Computer Science House

This project is licensed under the GNU Affero General Public License v3.0 or later (AGPL-3.0-or-later).
See the `LICENSE` file for details.
