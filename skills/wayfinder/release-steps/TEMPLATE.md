# Release:Next

**Ships:** everything on `dev` that is not on `main`. Last release: x.y.z (#PR, YYYY-MM-DD).
**Maps waiting:** #N {FeatureName}, ... (Phase `built`, work on `dev`).

Each line is a human step that production needs for work already merged to `dev`. The `release-steps` skill adds them as work lands and re-checks them before each release. Tick a box only when the step is done on prod. A step done on dev only stays unticked and says so.

## Before merge
Steps prod needs before the new code deploys: infra or secrets the app reads at startup, third-party setup the new code calls, data that a migration expects to exist.

## After deploy
Steps that need the new code running: Scheduler jobs that call new routes, data imports, backfills, maps VM or GeoServer changes, settings to fill in, cutovers.

## Verify on prod
One smoke check per feature, run by a human after deploy.

## Dropped
Steps that later work made unnecessary, each with the commit or file that made it so.
