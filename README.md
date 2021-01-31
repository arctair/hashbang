# hashbang
This starter gifts you a runnable barebones Kotlin Ktor application with tests.
## Eject
`./eject`
Ejecting renames packages, artifacts, services, deployments, and the project name from hashbang to the name of this directory.
## Run the tests
`./gradlew cleanTest test`
### Against a deployment
`./gradlew cleanTest test -Dhasbang_baseUrl=http://hashbang.arctair.com/api`
## Build, deploy, verify
`scripts/cd`
The cd script executes these steps:
1. Build jar file
1. Build and push Docker image
1. Create or update dark (non-production) Kubernetes service and deployment
1. Run blackbox tests against dark deployment baseUrl
1. Swap dark deployment and live deployment
