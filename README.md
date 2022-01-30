# Kruonis

A timelines app. 

Everything in `cmd` is the backend API, modelling timeline data models.

### Running locally
Set up an ephemeral postgres instance by running `./up.sh`. This will start a Postgres instance on port 5433 for local dev.
Then, you can use a local config.env to run the go app in `cmd/timelinesv2`.

Once done, you can shutdown this postgres instance by running `./down.sh`.

### Deployment 
TBD.
