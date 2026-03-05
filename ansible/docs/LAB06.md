# 1. Blocks & Tags

## Block Usage in Each Role
Common role (roles/common/tasks/main.yml):

    Two main blocks: Package installation and User creation.

    Each block has its own tags (packages, users) and rescue/always sections.

    The rescue block fixes apt cache issues and retries installation.

    The always block logs completion.

Docker role (roles/docker/tasks/main.yml):

    Blocks for docker_install and docker_config.

    rescue in the installation block waits 10 seconds and retries (handles transient network errors).

    always ensures Docker service is enabled.

Web_app role (roles/web_app/tasks/main.yml):

    The deployment logic is wrapped in a block with rescue (logs failure) and always (shows status).

    Wipe tasks are included separately with their own tag.

## Tag Strategy
Role‑level tags: common, docker, web_app.

Fine‑grained tags:

    packages, users (common)

    docker_install, docker_config (docker)

    web_app_wipe (wipe only)

    compose, deploy (web_app deployment)

## Output showing selective execution with --tags
```sh
$ ansible-playbook playbooks/provision.yml --tags "docker"

PLAY [Provision web servers with common tools and Docker] *************************************************************************

TASK [Gathering Facts] ************************************************************************************************************
ok: [damir-VirtualBox]

TASK [common : Update apt cache] **************************************************************************************************
ok: [damir-VirtualBox]

TASK [common : Install essential packages] ****************************************************************************************
ok: [damir-VirtualBox]

TASK [common : Log package installation completion] *******************************************************************************
changed: [damir-VirtualBox]

TASK [common : Ensure common user exists] *****************************************************************************************
changed: [damir-VirtualBox]

TASK [common : Log common role completion] ****************************************************************************************
changed: [damir-VirtualBox]

TASK [docker : Include Docker installation tasks] *********************************************************************************
included: /home/damir/Desktop/DevOps/DevOps-Core-Course/ansible/roles/docker/tasks/install.yml for damir-VirtualBox

TASK [docker : Remove old Docker packages (if any)] *******************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install required system packages] **********************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add Docker GPG key] ************************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add Docker repository] *********************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install Docker packages] *******************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Include Docker configuration tasks] ********************************************************************************
included: /home/damir/Desktop/DevOps/DevOps-Core-Course/ansible/roles/docker/tasks/config.yml for damir-VirtualBox

TASK [docker : Ensure Docker service is running and enabled] **********************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add user to docker group] ******************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install Python Docker module for Ansible] **************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Verify Docker installation] ****************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Display Docker version] ********************************************************************************************
ok: [damir-VirtualBox] => {
    "msg": "Docker installed: Docker version 29.2.1, build a5c7197"
}

TASK [docker : Ensure Docker service is enabled and started] **********************************************************************
ok: [damir-VirtualBox]

PLAY RECAP ************************************************************************************************************************
damir-VirtualBox           : ok=19   changed=3    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0   

$ ansible-playbook playbooks/provision.yml --skip-tags "common"

PLAY [Provision web servers with common tools and Docker] *************************************************************************

TASK [Gathering Facts] ************************************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Include Docker installation tasks] *********************************************************************************
included: /home/damir/Desktop/DevOps/DevOps-Core-Course/ansible/roles/docker/tasks/install.yml for damir-VirtualBox

TASK [docker : Remove old Docker packages (if any)] *******************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install required system packages] **********************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add Docker GPG key] ************************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add Docker repository] *********************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install Docker packages] *******************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Include Docker configuration tasks] ********************************************************************************
included: /home/damir/Desktop/DevOps/DevOps-Core-Course/ansible/roles/docker/tasks/config.yml for damir-VirtualBox

TASK [docker : Ensure Docker service is running and enabled] **********************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add user to docker group] ******************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install Python Docker module for Ansible] **************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Verify Docker installation] ****************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Display Docker version] ********************************************************************************************
ok: [damir-VirtualBox] => {
    "msg": "Docker installed: Docker version 29.2.1, build a5c7197"
}

TASK [docker : Ensure Docker service is enabled and started] **********************************************************************
ok: [damir-VirtualBox]

PLAY RECAP ************************************************************************************************************************
damir-VirtualBox           : ok=14   changed=0    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0   

$ ansible-playbook playbooks/provision.yml --tags "packages"

PLAY [Provision web servers with common tools and Docker] *************************************************************************

TASK [Gathering Facts] ************************************************************************************************************
ok: [damir-VirtualBox]

TASK [common : Update apt cache] **************************************************************************************************
ok: [damir-VirtualBox]

TASK [common : Install essential packages] ****************************************************************************************
ok: [damir-VirtualBox]

TASK [common : Log package installation completion] *******************************************************************************
changed: [damir-VirtualBox]

TASK [common : Ensure common user exists] *****************************************************************************************
ok: [damir-VirtualBox]

TASK [common : Log common role completion] ****************************************************************************************
changed: [damir-VirtualBox]

PLAY RECAP ************************************************************************************************************************
damir-VirtualBox           : ok=6    changed=2    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0   

$ ansible-playbook playbooks/provision.yml --tags "docker" --check

PLAY [Provision web servers with common tools and Docker] *************************************************************************

TASK [Gathering Facts] ************************************************************************************************************
ok: [damir-VirtualBox]

TASK [common : Update apt cache] **************************************************************************************************
ok: [damir-VirtualBox]

TASK [common : Install essential packages] ****************************************************************************************
ok: [damir-VirtualBox]

TASK [common : Log package installation completion] *******************************************************************************
changed: [damir-VirtualBox]

TASK [common : Ensure common user exists] *****************************************************************************************
ok: [damir-VirtualBox]

TASK [common : Log common role completion] ****************************************************************************************
changed: [damir-VirtualBox]

TASK [docker : Include Docker installation tasks] *********************************************************************************
included: /home/damir/Desktop/DevOps/DevOps-Core-Course/ansible/roles/docker/tasks/install.yml for damir-VirtualBox

TASK [docker : Remove old Docker packages (if any)] *******************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install required system packages] **********************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add Docker GPG key] ************************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add Docker repository] *********************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install Docker packages] *******************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Include Docker configuration tasks] ********************************************************************************
included: /home/damir/Desktop/DevOps/DevOps-Core-Course/ansible/roles/docker/tasks/config.yml for damir-VirtualBox

TASK [docker : Ensure Docker service is running and enabled] **********************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add user to docker group] ******************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install Python Docker module for Ansible] **************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Verify Docker installation] ****************************************************************************************
skipping: [damir-VirtualBox]

TASK [docker : Display Docker version] ********************************************************************************************
skipping: [damir-VirtualBox]

TASK [docker : Ensure Docker service is enabled and started] **********************************************************************
ok: [damir-VirtualBox]

PLAY RECAP ************************************************************************************************************************
damir-VirtualBox           : ok=17   changed=2    unreachable=0    failed=0    skipped=2    rescued=0    ignored=0   
```

## Output showing error handling with rescue block triggered
Here when retriving GPG key rescue block activates:
```
$ ansible-playbook playbooks/provision.yml --tags "docker_install"

PLAY [Provision web servers with common tools and Docker] *************************************************************************

TASK [Gathering Facts] ************************************************************************************************************
ok: [damir-VirtualBox]

TASK [common : Update apt cache] **************************************************************************************************
ok: [damir-VirtualBox]

TASK [common : Install essential packages] ****************************************************************************************
ok: [damir-VirtualBox]

TASK [common : Log package installation completion] *******************************************************************************
changed: [damir-VirtualBox]

TASK [common : Ensure common user exists] *****************************************************************************************
ok: [damir-VirtualBox]

TASK [common : Log common role completion] ****************************************************************************************
changed: [damir-VirtualBox]

TASK [docker : Include Docker installation tasks] *********************************************************************************
included: /home/damir/Desktop/DevOps/DevOps-Core-Course/ansible/roles/docker/tasks/install.yml for damir-VirtualBox

TASK [docker : Remove old Docker packages (if any)] *******************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install required system packages] **********************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add Docker GPG key] ************************************************************************************************
fatal: [damir-VirtualBox]: FAILED! => {"changed": false, "msg": "Failed to download key at https://download.docker.com/linux/ubuntu/gpg: Connection failure: Remote end closed connection without response"}

TASK [docker : Wait 10 seconds before retry] **************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Retry Docker installation] *****************************************************************************************
included: /home/damir/Desktop/DevOps/DevOps-Core-Course/ansible/roles/docker/tasks/install.yml for damir-VirtualBox

TASK [docker : Remove old Docker packages (if any)] *******************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install required system packages] **********************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add Docker GPG key] ************************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add Docker repository] *********************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install Docker packages] *******************************************************************************************
ok: [damir-VirtualBox]

PLAY RECAP ************************************************************************************************************************
damir-VirtualBox           : ok=16   changed=2    unreachable=0    failed=0    skipped=0    rescued=1    ignored=0   
```

## List of all available tags (--list-tags output)
```sh
$ ansible-playbook playbooks/provision.yml --list-tags

playbook: playbooks/provision.yml

  play #1 (webservers): Provision web servers with common tools and Docker      TAGS: []
      TASK TAGS: [always, common, docker, docker_config, docker_install, packages, users]
```





# 2. Docker Compose Migration
## Template Structure
```
---

services:
  {{ app_name }}:
    image: "{{ dockerhub_username }}/{{ app_name }}:latest"
    container_name: "{{ app_name }}-compose"
    ports:
      - "{{ app_port }}:{{ app_internal_port | default(app_port) }}"
    environment:
      APP_NAME: "{{ app_name }}"
      APP_PORT: "{{ app_internal_port | default(app_port) }}"
      ENVIRONMENT: "{{ app_environment.ENVIRONMENT | default('production') }}"

    restart: "{{ restart_policy | default('unless-stopped') }}"
    networks:
      - app_network
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

networks:
  app_network:
    driver: bridge
```

## Role Dependencies
Defined in roles/web_app/meta/main.yml
This ensures Docker is installed before any container deployment – Ansible automatically executes the docker role first.

Before/After Comparison
Aspect	docker run (previous)	Docker Compose (new)
Definition	Imperative commands	Declarative YAML
Multi‑container	Manual linking	Built‑in networks, dependencies
Idempotency	Custom checks	docker_compose module handles it
Updates	Stop/remove/run cycle	docker compose up with changed config
Environment	Passed via env parameter	Managed in the compose file
Networks	Default bridge	Custom network app_network

## Output showing Docker Compose deployment success
```sh
$ ansible-playbook playbooks/deploy.yml --ask-vault-pass
Vault password: 

PLAY [Deploy application to web servers] ******************************************************************************************

TASK [Gathering Facts] ************************************************************************************************************
[WARNING]: Platform linux on host damir-VirtualBox is using the discovered Python interpreter at /usr/bin/python3.10, but future
installation of another Python interpreter could change the meaning of that path. See https://docs.ansible.com/ansible-
core/2.17/reference_appendices/interpreter_discovery.html for more information.
ok: [damir-VirtualBox]

TASK [Display deployment information] *********************************************************************************************
ok: [damir-VirtualBox] => {
    "msg": [
        "Starting deployment of devops-info-service",
        "Image: damirsadykov/devops-info-service:latest",
        "Target host: damir-VirtualBox"
    ]
}

TASK [docker : Include Docker installation tasks] *********************************************************************************
included: /home/damir/Desktop/DevOps/DevOps-Core-Course/ansible/roles/docker/tasks/install.yml for damir-VirtualBox

TASK [docker : Remove old Docker packages (if any)] *******************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install required system packages] **********************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add Docker GPG key] ************************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add Docker repository] *********************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install Docker packages] *******************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Include Docker configuration tasks] ********************************************************************************
included: /home/damir/Desktop/DevOps/DevOps-Core-Course/ansible/roles/docker/tasks/config.yml for damir-VirtualBox

TASK [docker : Ensure Docker service is running and enabled] **********************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add user to docker group] ******************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install Python Docker module for Ansible] **************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Verify Docker installation] ****************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Display Docker version] ********************************************************************************************
ok: [damir-VirtualBox] => {
    "msg": "Docker installed: Docker version 29.3.0, build 5927d80"
}

TASK [docker : Ensure Docker service is enabled and started] **********************************************************************
ok: [damir-VirtualBox]

TASK [web_app : Ensure docker-compose Python library is installed (for older modules)] ********************************************
skipping: [damir-VirtualBox]

TASK [web_app : Create application directory] *************************************************************************************
ok: [damir-VirtualBox]

TASK [web_app : Template docker-compose.yml] **************************************************************************************
changed: [damir-VirtualBox]

TASK [web_app : Log in to Docker Hub (if credentials provided)] *******************************************************************
changed: [damir-VirtualBox]

TASK [web_app : Deploy with Docker Compose] ***************************************************************************************
changed: [damir-VirtualBox]

TASK [web_app : Wait for application to be ready] *********************************************************************************
ok: [damir-VirtualBox]

TASK [web_app : Check health endpoint] ********************************************************************************************
ok: [damir-VirtualBox]

TASK [web_app : Show deployment status] *******************************************************************************************
ok: [damir-VirtualBox] => {
    "msg": [
        "Application deployed with Docker Compose",
        "Project directory: /opt/devops-info-service",
        "Container status: run 'docker ps' on target"
    ]
}

TASK [Verify deployment] **********************************************************************************************************
ok: [damir-VirtualBox]

TASK [Show deployment status] *****************************************************************************************************
ok: [damir-VirtualBox] => {
    "msg": [
        "Deployment completed successfully!",
        "Container status: exited",
        "Container started: 2026-03-05T19:37:33.021512188Z"
    ]
}

PLAY RECAP ************************************************************************************************************************
damir-VirtualBox           : ok=24   changed=3    unreachable=0    failed=0    skipped=1    rescued=0    ignored=0   
```

## Idempotency proof (second run shows "ok" not "changed")
```sh
$ ansible-playbook playbooks/deploy.yml --ask-vault-pass
Vault password: 

PLAY [Deploy application to web servers] ******************************************************************************************

TASK [Gathering Facts] ************************************************************************************************************
[WARNING]: Platform linux on host damir-VirtualBox is using the discovered Python interpreter at /usr/bin/python3.10, but future
installation of another Python interpreter could change the meaning of that path. See https://docs.ansible.com/ansible-
core/2.17/reference_appendices/interpreter_discovery.html for more information.
ok: [damir-VirtualBox]

TASK [Display deployment information] *********************************************************************************************
ok: [damir-VirtualBox] => {
    "msg": [
        "Starting deployment of devops-info-service",
        "Image: damirsadykov/devops-info-service:latest",
        "Target host: damir-VirtualBox"
    ]
}

TASK [docker : Include Docker installation tasks] *********************************************************************************
included: /home/damir/Desktop/DevOps/DevOps-Core-Course/ansible/roles/docker/tasks/install.yml for damir-VirtualBox

TASK [docker : Remove old Docker packages (if any)] *******************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install required system packages] **********************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add Docker GPG key] ************************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add Docker repository] *********************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install Docker packages] *******************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Include Docker configuration tasks] ********************************************************************************
included: /home/damir/Desktop/DevOps/DevOps-Core-Course/ansible/roles/docker/tasks/config.yml for damir-VirtualBox

TASK [docker : Ensure Docker service is running and enabled] **********************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add user to docker group] ******************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install Python Docker module for Ansible] **************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Verify Docker installation] ****************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Display Docker version] ********************************************************************************************
ok: [damir-VirtualBox] => {
    "msg": "Docker installed: Docker version 29.3.0, build 5927d80"
}

TASK [docker : Ensure Docker service is enabled and started] **********************************************************************
ok: [damir-VirtualBox]

TASK [web_app : Ensure docker-compose Python library is installed (for older modules)] ********************************************
skipping: [damir-VirtualBox]

TASK [web_app : Create application directory] *************************************************************************************
ok: [damir-VirtualBox]

TASK [web_app : Template docker-compose.yml] **************************************************************************************
ok: [damir-VirtualBox]

TASK [web_app : Log in to Docker Hub (if credentials provided)] *******************************************************************
changed: [damir-VirtualBox]

TASK [web_app : Deploy with Docker Compose] ***************************************************************************************
ok: [damir-VirtualBox]

TASK [web_app : Wait for application to be ready] *********************************************************************************
skipping: [damir-VirtualBox]

TASK [web_app : Check health endpoint] ********************************************************************************************
skipping: [damir-VirtualBox]

TASK [web_app : Show deployment status] *******************************************************************************************
ok: [damir-VirtualBox] => {
    "msg": [
        "Application deployed with Docker Compose",
        "Project directory: /opt/devops-info-service",
        "Container status: run 'docker ps' on target"
    ]
}

TASK [Verify deployment] **********************************************************************************************************
ok: [damir-VirtualBox]

TASK [Show deployment status] *****************************************************************************************************
ok: [damir-VirtualBox] => {
    "msg": [
        "Deployment completed successfully!",
        "Container status: exited",
        "Container started: 2026-03-05T19:37:33.021512188Z"
    ]
}

PLAY RECAP ************************************************************************************************************************
damir-VirtualBox           : ok=22   changed=1    unreachable=0    failed=0    skipped=3    rescued=0    ignored=0 
```

## Application running and accessible
```sh
$ curl http://192.168.1.8:5000/
{"endpoints":[{"description":"Service information","method":"GET","path":"/"},{"description":"Health check","method":"GET","path":"/health"}],"request":{"client_ip":"192.168.1.9","method":"GET","path":"/","user_agent":"curl/7.81.0"},"runtime":{"current_time":"2026-03-05T20:26:24.910321Z","timezone":"UTC","uptime_human":"0 hours, 2 minutes","uptime_seconds":159},"service":{"description":"DevOps course info service","framework":"Flask","name":"devops-info-service","version":"1.0.0"},"system":{"architecture":"x86_64","cpu_count":1,"hostname":"3e63ccceb7b4","platform":"Linux","platform_version":"#100~22.04.1-Ubuntu SMP PREEMPT_DYNAMIC Mon Jan 19 17:10:19 UTC ","python_version":"3.12.12"}}
```

## Contents of templated docker-compose.yml
``` sh
$ ssh damir@192.168.1.8 cat /opt/devops-info-service/docker-compose.yml
---

services:
  devops-info-service:
    image: "damirsadykov/devops-info-service:latest"
    container_name: "devops-info-service-compose"
    ports:
      - "5000:5000"
    environment:
      APP_NAME: "devops-info-service"
      APP_PORT: "5000"
      ENVIRONMENT: "production"

    restart: "unless-stopped"
    networks:
      - app_network
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

networks:
  app_network:
    driver: bridge
```


# 3. Wipe Logic

## Implementation Details
Controlled by variable web_app_wipe (default false in defaults/main.yml).

Tasks are placed in a separate file wipe.yml and included at the beginning of main.yml.

Tag web_app_wipe allows selective execution.

## Variable + Tag Double Safety
Variable prevents accidental wipe even if the tag is supplied.

Tag allows limiting execution to wipe tasks only (without deploying).

This ensures that wipe can only happen when explicitly requested.

## Output of Scenario 1 showing normal deployment (wipe skipped)
```sh
$ ansible-playbook playbooks/deploy.yml --ask-vault-pass
Vault password: 

PLAY [Deploy application to web servers] ******************************************************************************************

TASK [Gathering Facts] ************************************************************************************************************
[WARNING]: Platform linux on host damir-VirtualBox is using the discovered Python interpreter at /usr/bin/python3.10, but future
installation of another Python interpreter could change the meaning of that path. See https://docs.ansible.com/ansible-
core/2.17/reference_appendices/interpreter_discovery.html for more information.
ok: [damir-VirtualBox]

TASK [Display deployment information] *********************************************************************************************
ok: [damir-VirtualBox] => {
    "msg": [
        "Starting deployment of devops-info-service",
        "Image: damirsadykov/devops-info-service:latest",
        "Target host: damir-VirtualBox"
    ]
}

TASK [docker : Include Docker installation tasks] *********************************************************************************
included: /home/damir/Desktop/DevOps/DevOps-Core-Course/ansible/roles/docker/tasks/install.yml for damir-VirtualBox

TASK [docker : Remove old Docker packages (if any)] *******************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install required system packages] **********************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add Docker GPG key] ************************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add Docker repository] *********************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install Docker packages] *******************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Include Docker configuration tasks] ********************************************************************************
included: /home/damir/Desktop/DevOps/DevOps-Core-Course/ansible/roles/docker/tasks/config.yml for damir-VirtualBox

TASK [docker : Ensure Docker service is running and enabled] **********************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add user to docker group] ******************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install Python Docker module for Ansible] **************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Verify Docker installation] ****************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Display Docker version] ********************************************************************************************
ok: [damir-VirtualBox] => {
    "msg": "Docker installed: Docker version 29.3.0, build 5927d80"
}

TASK [docker : Ensure Docker service is enabled and started] **********************************************************************
ok: [damir-VirtualBox]

TASK [web_app : Include wipe tasks] ***********************************************************************************************
included: /home/damir/Desktop/DevOps/DevOps-Core-Course/ansible/roles/web_app/tasks/wipe.yml for damir-VirtualBox

TASK [web_app : Stop and remove containers with Docker Compose] *******************************************************************
skipping: [damir-VirtualBox]

TASK [web_app : Remove docker-compose.yml file] ***********************************************************************************
skipping: [damir-VirtualBox]

TASK [web_app : Remove application directory] *************************************************************************************
skipping: [damir-VirtualBox]

TASK [web_app : (Optional) Remove Docker images] **********************************************************************************
skipping: [damir-VirtualBox]

TASK [web_app : Log wipe completion] **********************************************************************************************
skipping: [damir-VirtualBox]

TASK [web_app : Ensure docker-compose Python library is installed (for older modules)] ********************************************
skipping: [damir-VirtualBox]

TASK [web_app : Create application directory] *************************************************************************************
ok: [damir-VirtualBox]

TASK [web_app : Template docker-compose.yml] **************************************************************************************
ok: [damir-VirtualBox]

TASK [web_app : Log in to Docker Hub (if credentials provided)] *******************************************************************
changed: [damir-VirtualBox]

TASK [web_app : Deploy with Docker Compose] ***************************************************************************************
ok: [damir-VirtualBox]

TASK [web_app : Wait for application to be ready] *********************************************************************************
skipping: [damir-VirtualBox]

TASK [web_app : Check health endpoint] ********************************************************************************************
skipping: [damir-VirtualBox]

TASK [web_app : Show deployment status] *******************************************************************************************
ok: [damir-VirtualBox] => {
    "msg": [
        "Application deployed with Docker Compose",
        "Project directory: /opt/devops-info-service",
        "Container status: run 'docker ps' on target"
    ]
}

TASK [Verify deployment] **********************************************************************************************************
ok: [damir-VirtualBox]

TASK [Show deployment status] *****************************************************************************************************
ok: [damir-VirtualBox] => {
    "msg": [
        "Deployment completed successfully!",
        "Container status: exited",
        "Container started: 2026-03-05T19:37:33.021512188Z"
    ]
}

PLAY RECAP ************************************************************************************************************************
damir-VirtualBox           : ok=23   changed=1    unreachable=0    failed=0    skipped=8    rescued=0    ignored=0   

$ ssh damir@192.168.1.8 "docker ps"
CONTAINER ID   IMAGE                                     COMMAND           CREATED         STATUS         PORTS                                         NAMES
3e63ccceb7b4   damirsadykov/devops-info-service:latest   "python app.py"   9 minutes ago   Up 9 minutes   0.0.0.0:5000->5000/tcp, [::]:5000->5000/tcp   devops-info-service-compose
```

## Output of Scenario 2 showing wipe-only operation
``` sh
$ ansible-playbook playbooks/deploy.yml --ask-vault-pass -e "web_app_wipe=true" --tags web_app_wipe
Vault password: 

PLAY [Deploy application to web servers] ******************************************************************************************

TASK [Gathering Facts] ************************************************************************************************************
[WARNING]: Platform linux on host damir-VirtualBox is using the discovered Python interpreter at /usr/bin/python3.10, but future
installation of another Python interpreter could change the meaning of that path. See https://docs.ansible.com/ansible-
core/2.17/reference_appendices/interpreter_discovery.html for more information.
ok: [damir-VirtualBox]

TASK [web_app : Include wipe tasks] ***********************************************************************************************
included: /home/damir/Desktop/DevOps/DevOps-Core-Course/ansible/roles/web_app/tasks/wipe.yml for damir-VirtualBox

TASK [web_app : Stop and remove containers with Docker Compose] *******************************************************************
changed: [damir-VirtualBox]

TASK [web_app : Remove docker-compose.yml file] ***********************************************************************************
changed: [damir-VirtualBox]

TASK [web_app : Remove application directory] *************************************************************************************
changed: [damir-VirtualBox]

TASK [web_app : (Optional) Remove Docker images] **********************************************************************************
skipping: [damir-VirtualBox]

TASK [web_app : Log wipe completion] **********************************************************************************************
ok: [damir-VirtualBox] => {
    "msg": "Application devops-info-service wiped successfully from /opt/devops-info-service"
}

TASK [web_app : Ensure docker-compose Python library is installed (for older modules)] ********************************************
skipping: [damir-VirtualBox]

PLAY RECAP ************************************************************************************************************************
damir-VirtualBox           : ok=6    changed=3    unreachable=0    failed=0    skipped=2    rescued=0    ignored=0   

damir@damir-VB:~/Desktop/DevOps/DevOps-Core-Course/ansible$ ssh damir@192.168.1.8 "docker ps"
CONTAINER ID   IMAGE     COMMAND   CREATED   STATUS    PORTS     NAMES
damir@damir-VB:~/Desktop/DevOps/DevOps-Core-Course/ansible$ ssh damir@192.168.1.8 "ls /opt"
containerd
```

## Output of Scenario 3 showing clean reinstall (wipe → deploy)
``` sh
$ ansible-playbook playbooks/deploy.yml --ask-vault-pass -e "web_app_wipe=true"
Vault password: 

PLAY [Deploy application to web servers] ******************************************************************************************

TASK [Gathering Facts] ************************************************************************************************************
[WARNING]: Platform linux on host damir-VirtualBox is using the discovered Python interpreter at /usr/bin/python3.10, but future
installation of another Python interpreter could change the meaning of that path. See https://docs.ansible.com/ansible-
core/2.17/reference_appendices/interpreter_discovery.html for more information.
ok: [damir-VirtualBox]

TASK [Display deployment information] *********************************************************************************************
ok: [damir-VirtualBox] => {
    "msg": [
        "Starting deployment of devops-info-service",
        "Image: damirsadykov/devops-info-service:latest",
        "Target host: damir-VirtualBox"
    ]
}

TASK [docker : Include Docker installation tasks] *********************************************************************************
included: /home/damir/Desktop/DevOps/DevOps-Core-Course/ansible/roles/docker/tasks/install.yml for damir-VirtualBox

TASK [docker : Remove old Docker packages (if any)] *******************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install required system packages] **********************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add Docker GPG key] ************************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add Docker repository] *********************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install Docker packages] *******************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Include Docker configuration tasks] ********************************************************************************
included: /home/damir/Desktop/DevOps/DevOps-Core-Course/ansible/roles/docker/tasks/config.yml for damir-VirtualBox

TASK [docker : Ensure Docker service is running and enabled] **********************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add user to docker group] ******************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install Python Docker module for Ansible] **************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Verify Docker installation] ****************************************************************************************
ok: [damir-VirtualBox]

TASK [docker : Display Docker version] ********************************************************************************************
ok: [damir-VirtualBox] => {
    "msg": "Docker installed: Docker version 29.3.0, build 5927d80"
}

TASK [docker : Ensure Docker service is enabled and started] **********************************************************************
ok: [damir-VirtualBox]

TASK [web_app : Include wipe tasks] ***********************************************************************************************
included: /home/damir/Desktop/DevOps/DevOps-Core-Course/ansible/roles/web_app/tasks/wipe.yml for damir-VirtualBox

TASK [web_app : Stop and remove containers with Docker Compose] *******************************************************************
fatal: [damir-VirtualBox]: FAILED! => {"changed": false, "msg": "\"/opt/devops-info-service\" is not a directory"}
...ignoring

TASK [web_app : Remove docker-compose.yml file] ***********************************************************************************
ok: [damir-VirtualBox]

TASK [web_app : Remove application directory] *************************************************************************************
ok: [damir-VirtualBox]

TASK [web_app : (Optional) Remove Docker images] **********************************************************************************
skipping: [damir-VirtualBox]

TASK [web_app : Log wipe completion] **********************************************************************************************
ok: [damir-VirtualBox] => {
    "msg": "Application devops-info-service wiped successfully from /opt/devops-info-service"
}

TASK [web_app : Ensure docker-compose Python library is installed (for older modules)] ********************************************
skipping: [damir-VirtualBox]

TASK [web_app : Create application directory] *************************************************************************************
changed: [damir-VirtualBox]

TASK [web_app : Template docker-compose.yml] **************************************************************************************
changed: [damir-VirtualBox]

TASK [web_app : Log in to Docker Hub (if credentials provided)] *******************************************************************
changed: [damir-VirtualBox]

TASK [web_app : Deploy with Docker Compose] ***************************************************************************************
changed: [damir-VirtualBox]

TASK [web_app : Wait for application to be ready] *********************************************************************************
ok: [damir-VirtualBox]

TASK [web_app : Check health endpoint] ********************************************************************************************
ok: [damir-VirtualBox]

TASK [web_app : Show deployment status] *******************************************************************************************
ok: [damir-VirtualBox] => {
    "msg": [
        "Application deployed with Docker Compose",
        "Project directory: /opt/devops-info-service",
        "Container status: run 'docker ps' on target"
    ]
}

TASK [Verify deployment] **********************************************************************************************************
ok: [damir-VirtualBox]

TASK [Show deployment status] *****************************************************************************************************
ok: [damir-VirtualBox] => {
    "msg": [
        "Deployment completed successfully!",
        "Container status: exited",
        "Container started: 2026-03-05T19:37:33.021512188Z"
    ]
}

PLAY RECAP ************************************************************************************************************************
damir-VirtualBox           : ok=29   changed=4    unreachable=0    failed=0    skipped=2    rescued=0    ignored=1   

damir@damir-VB:~/Desktop/DevOps/DevOps-Core-Course/ansible$ ssh damir@192.168.1.8 "docker ps"
CONTAINER ID   IMAGE                                     COMMAND           CREATED          STATUS          PORTS                                         NAMES
d93503af48c6   damirsadykov/devops-info-service:latest   "python app.py"   35 seconds ago   Up 34 seconds   0.0.0.0:5000->5000/tcp, [::]:5000->5000/tcp   devops-info-service-compose
```

## Output of Scenario 4a showing wipe blocked by when condition
```sh
 ansible-playbook playbooks/deploy.yml --ask-vault-pass --tags web_app_wipe
Vault password: 

PLAY [Deploy application to web servers] ******************************************************************************************

TASK [Gathering Facts] ************************************************************************************************************
[WARNING]: Platform linux on host damir-VirtualBox is using the discovered Python interpreter at /usr/bin/python3.10, but future
installation of another Python interpreter could change the meaning of that path. See https://docs.ansible.com/ansible-
core/2.17/reference_appendices/interpreter_discovery.html for more information.
ok: [damir-VirtualBox]

TASK [web_app : Include wipe tasks] ***********************************************************************************************
included: /home/damir/Desktop/DevOps/DevOps-Core-Course/ansible/roles/web_app/tasks/wipe.yml for damir-VirtualBox

TASK [web_app : Stop and remove containers with Docker Compose] *******************************************************************
skipping: [damir-VirtualBox]

TASK [web_app : Remove docker-compose.yml file] ***********************************************************************************
skipping: [damir-VirtualBox]

TASK [web_app : Remove application directory] *************************************************************************************
skipping: [damir-VirtualBox]

TASK [web_app : (Optional) Remove Docker images] **********************************************************************************
skipping: [damir-VirtualBox]

TASK [web_app : Log wipe completion] **********************************************************************************************
skipping: [damir-VirtualBox]

TASK [web_app : Ensure docker-compose Python library is installed (for older modules)] ********************************************
skipping: [damir-VirtualBox]

PLAY RECAP ************************************************************************************************************************
damir-VirtualBox           : ok=2    changed=0    unreachable=0    failed=0    skipped=6    rescued=0    ignored=0   
```

# 4. CI/CD Integration - skipped

# Research Answers
## Q: What happens if rescue block also fails?
If the rescue block itself fails, Ansible stops execution on that host and marks the task as failed (unless ignore_errors is set). The always block still runs afterwards.

## Q: Can you have nested blocks?
Yes, blocks can be nested. Inner blocks can have their own rescue/always, and errors propagate to outer blocks.

## Q: How do tags inherit to tasks within blocks?
Tags applied to a block are inherited by all tasks inside that block. Tags can also be overridden at task level.

##  Q: What's the difference between restart: always and restart: unless-stopped?
    always – container restarts regardless of exit status, even if manually stopped.

    unless-stopped – restarts unless the container was explicitly stopped. Better for production because manual stops are honoured.

## Q: How do Docker Compose networks differ from Docker bridge networks?

Compose automatically creates a dedicated bridge network for the project, providing service discovery via container names. It isolates containers from other projects.
## Q: Can you reference Ansible Vault variables in the template?

Yes, as long as the playbook has access to the vault password, variables like {{ app_secret_key }} are decrypted and can be used in templates.
## Q: Why use both variable AND tag? (Double safety mechanism)

    Variable ensures wipe cannot happen accidentally even if the tag is present.

    Tag allows limiting execution to wipe tasks only (no deployment).
    Together they give fine‑grained control and prevent mistakes.

## Q: What's the difference between never tag and this approach?

The never tag permanently excludes tasks from running unless explicitly requested. Our approach uses a variable condition, which is more flexible (e.g., can be overridden per run with -e).
## Q: Why must wipe logic come BEFORE deployment in main.yml?

To support clean reinstallation – remove the old application before deploying the new one, ensuring a fresh state.
## Q: When would you want clean reinstallation vs. rolling update?

    Clean reinstall – when application state is broken or you need a guaranteed fresh start (e.g., schema changes).

    Rolling update – for zero‑downtime updates; not implemented here but could be done with Compose.


## Q: How would you extend this to wipe Docker images and volumes too?

Add tasks in wipe.yml:
yaml

- name: Remove Docker images
  docker_image:
    name: "{{ dockerhub_username }}/{{ app_name }}:{{ docker_tag }}"
    state: absent
- name: Remove named volumes
  docker_volume:
    name: "{{ item }}"
    state: absent
  loop: "{{ volumes_to_remove }}"