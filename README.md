# compose-env-manager
[![Build Status](https://travis-ci.org/sgeisbacher/compose-env-manager.svg?branch=master)](https://travis-ci.org/sgeisbacher/compose-env-manager)  
compose-env-manager is a wrapper for **docker-compose**.

It preruns a *docker-compose.yml*-file-generator, which renders based on *scenarios* defined in *compose-env-manager.yml*.  
This means you define multiple *scenarios* (scenario *all* is required, which is also the base for other scenarios).  
Other scenarios must be subsets of the *all*-scenario.  
The renderer will detect missing services by diffing the selected with the all-*scenario* and adds them as */etc/hosts*-entries (via docker's extra-hosts-option) in each of the services. So requests to missing services will be routed to your host-machine.  
Now you are able to run parts of your multi-tier-application directly on your host-machine (parts you're currently working on) and all others in docker with docker-compose.  

## Installation
todo
## Project-Configuration
- create a **compose-env-manager.yml** in the root-folder of your project
- add at least the *all*-scenario, e.g.:

```yaml
templateFolderPath: ./templates/
scenarios:
  all:
    - frontend
    - app
    - db
```
- *compose-env-manager* will lookup the service-templates in **./templates/\<service-name\>.yml**. So create your service-templates (like in the example above *frontend.yml*, *app.yml* and *db.yml*) as you know from [docker-compose-file-reference](https://docs.docker.com/compose/compose-file/) in that folder. e.g.:

```yaml
image: nginx
ports:
  - "80:80"
```
## Additional Configuration
### volume-init
If you run *docker-compose up* your database and all other docker-volumes will be empty. So you may want to provide initial-data-zips for your docker-volumes.  
*compose-env-manager init* will download the specified zip and extracts it to your specified target-dir. If target-dir is not empty, *compose-env-manager* will skip this volume-init.  
You can define a *volume-init*-section in your *compose-env-manager.yml* to achieve that, e.g.:

```yaml
volume-init:
    - name: app-data-dir
      source: http://example.org/app-data.zip
      target: ./data/app
    - name: mysql-data-dir
      source: http://example.org/db-data.zip
      target: ./data/mysql
```

> **Note!**  
> - Don't forget to add this data-directory to .gitignore-file ;-)  
> - Currently only HTTP-Zip-Sources are supported!!!

## Sample Configuration
### compose-env-manager.yml

```yaml
scenarios:
  all:
    - frontend
    - app
    - db
  backend:
    - frontend
    - db
  frontend:
    - app
    - db
volume-init:
    - name: app-data-dir
      source: http://example.org/app-data.zip
      target: ./data/app
    - name: mysql-data-dir
      source: http://example.org/db-data.zip
      target: ./data/mysql
```

### ./templates/db.yml

```yaml
image: mysql
ports:
    - "3306:3306"
volumes:
    - ./data/mysql:/var/lib/mysql
```

## RUN
please check:
```bash
compose-env-manager help
```
## TODOS
- check system prerequisites in bash(bat)-wrapper (python3, docker, ...)
- ~~detect required extra-hosts by diffing all-scenario with current-scenario~~
- ~~introduce pytest~~
- ~~stability-adjustments (no args, ...)~~
- add template-loading from remote-server
- remove hardcoded-dns-server-ip
- create .bat-wrapper
- test on windows
- ~~add global extra-hosts to compose-env-manager.yml~~
- ~~check scenario "all" is present (in configuration)~~
- ~~check all template-files of scenario "all" are present~~
- selected scenario:
    - detect missing services (against "all"-scenario)
    - add extra-hosts for missing scenarios
    - add extra-hosts for global external-services
    - add all services to compose-data-map
- gen yaml and write it to docker-compose.yml 
