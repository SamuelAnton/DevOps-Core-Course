# LAB 05

## 1. Architecture Overview
### Ansible version used
ansible 2.10.8

### Target VM OS and version
Ubuntu 22.04.5 LTS

### Role structure diagram
```
ansible/
├── inventory/
│   └── hosts.ini                    # Static inventory with VM details
├── group_vars/
│   └── all.yml                       # Encrypted variables (Vault)
├── roles/
│   ├── common/                        # Base system configuration
│   │   ├── tasks/
│   │   │   └── main.yml              # System updates and packages
│   │   └── defaults/
│   │       └── main.yml               # Default package list
│   ├── docker/                         # Docker installation
│   │   ├── tasks/
│   │   │   └── main.yml               # Docker setup tasks
│   │   ├── handlers/
│   │   │   └── main.yml               # Docker restart handler
│   │   └── defaults/
│   │       └── main.yml               # Docker version defaults
│   └── app_deploy/                     # Application deployment
│       ├── tasks/
│       │   └── main.yml               # Container deployment
│       ├── handlers/
│       │   └── main.yml               # Container restart
│       └── defaults/
│           └── main.yml                # App configuration
├── playbooks/
│   ├── provision.yml                   # System provisioning
│   └── deploy.yml                       # Application deployment
└── ansible.cfg                           # Ansible configuration
```

### Why roles instead of monolithic playbooks?
Roles provide modularity, reusability, and maintainability. Instead of having one large playbook with hundreds of tasks, roles allow you to:
- Separate concerns: Each role handles a specific functionality (system setup, Docker, application)
- Reuse code: Same role can be used across different projects
- Share with community: Roles can be easily shared via Ansible Galaxy
- Parallel development: Multiple team members can work on different roles simultaneously
- Easier testing: Each role can be tested independently

## 2. Roles Documentation
### Common Role
**Purpose**: Configures the base system with essential packages and timezone settings. This role ensures all servers have consistent basic configuration.

**Variables**:

| Variable | Default | Description
|----------|--------|-------------|
|common_packages | [python3-pip, curl, git, vim, htop, ...] | List of essential packages
|common_timezone | "UTC" | 	System timezone

**Handlers**: None defined in this role.

**Dependencies**: No dependencies on other roles.

### Docker Role
**Purpose**: Installs and configures Docker CE on Ubuntu systems, including all necessary dependencies and user permissions.

**Variables**:
| Variable | Default | Description
|----------|--------|-------------|
| docker_version | ["latest"] | Docker version to install
| docker_user | "{{ ansible_user }}" | User to add to docker group
| dockdocker_packageser_user | [docker-ce, docker-ce-cli, ...] | Docker packages to install

**Handlers**:
- `restart docker`: Restarts Docker service when configuration changes

**Dependencies**: None, but requires common role for base packages.

### App_Deploy Role
**Purpose**: Deploys the application container from Docker Hub, manages container lifecycle, and performs health checks.

**Variables**:
| Variable | Default | Description
|----------|--------|-------------|
| app_name | "devops-info-service" | Application name
| app_port | "5000" | Container port
| restart_policy | "unless-stopped" | Docker restart policy
| app_environment | {...} | Environment variables

**Handlers**:
- `restart app container`: Restarts the application container

**Dependencies**: Requires Docker role to be executed first.


## 3. Idempotency Demonstration

### Terminal output from FIRST provision.yml run
```
$ ansible-playbook playbooks/provision.yml

PLAY [Provision web servers with common tools and Docker] ******************************************************

TASK [Gathering Facts] *****************************************************************************************
ok: [damir-VirtualBox]

TASK [common : Update apt cache] *******************************************************************************
changed: [damir-VirtualBox]

TASK [common : Install essential packages] *********************************************************************
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/core.py) as it seems to be
invalid: cannot import name 'environmentfilter' from 'jinja2.filters' (/home/damir/.local/lib/python3.10/site-
packages/jinja2/filters.py)
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/mathstuff.py) as it seems to
be invalid: cannot import name 'environmentfilter' from 'jinja2.filters'
(/home/damir/.local/lib/python3.10/site-packages/jinja2/filters.py)
changed: [damir-VirtualBox]

TASK [common : Set timezone] ***********************************************************************************
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/core.py) as it seems to be
invalid: cannot import name 'environmentfilter' from 'jinja2.filters' (/home/damir/.local/lib/python3.10/site-
packages/jinja2/filters.py)
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/mathstuff.py) as it seems to
be invalid: cannot import name 'environmentfilter' from 'jinja2.filters'
(/home/damir/.local/lib/python3.10/site-packages/jinja2/filters.py)
ok: [damir-VirtualBox]

TASK [docker : Remove old Docker packages (if any)] ************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install required system packages] ***************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add Docker GPG key] *****************************************************************************
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/core.py) as it seems to be
invalid: cannot import name 'environmentfilter' from 'jinja2.filters' (/home/damir/.local/lib/python3.10/site-
packages/jinja2/filters.py)
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/mathstuff.py) as it seems to
be invalid: cannot import name 'environmentfilter' from 'jinja2.filters'
(/home/damir/.local/lib/python3.10/site-packages/jinja2/filters.py)
changed: [damir-VirtualBox]

TASK [docker : Add Docker repository] **************************************************************************
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/core.py) as it seems to be
invalid: cannot import name 'environmentfilter' from 'jinja2.filters' (/home/damir/.local/lib/python3.10/site-
packages/jinja2/filters.py)
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/mathstuff.py) as it seems to
be invalid: cannot import name 'environmentfilter' from 'jinja2.filters'
(/home/damir/.local/lib/python3.10/site-packages/jinja2/filters.py)
changed: [damir-VirtualBox]

TASK [docker : Install Docker packages] ************************************************************************
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/core.py) as it seems to be
invalid: cannot import name 'environmentfilter' from 'jinja2.filters' (/home/damir/.local/lib/python3.10/site-
packages/jinja2/filters.py)
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/mathstuff.py) as it seems to
be invalid: cannot import name 'environmentfilter' from 'jinja2.filters'
(/home/damir/.local/lib/python3.10/site-packages/jinja2/filters.py)
changed: [damir-VirtualBox]

TASK [docker : Ensure Docker service is running and enabled] ***************************************************
ok: [damir-VirtualBox]

TASK [docker : Add user to docker group] ***********************************************************************
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/core.py) as it seems to be
invalid: cannot import name 'environmentfilter' from 'jinja2.filters' (/home/damir/.local/lib/python3.10/site-
packages/jinja2/filters.py)
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/mathstuff.py) as it seems to
be invalid: cannot import name 'environmentfilter' from 'jinja2.filters'
(/home/damir/.local/lib/python3.10/site-packages/jinja2/filters.py)
changed: [damir-VirtualBox]

TASK [docker : Install Python Docker module for Ansible] *******************************************************
changed: [damir-VirtualBox]

TASK [docker : Verify Docker installation] *********************************************************************
ok: [damir-VirtualBox]

TASK [docker : Display Docker version] *************************************************************************
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/core.py) as it seems to be
invalid: cannot import name 'environmentfilter' from 'jinja2.filters' (/home/damir/.local/lib/python3.10/site-
packages/jinja2/filters.py)
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/mathstuff.py) as it seems to
be invalid: cannot import name 'environmentfilter' from 'jinja2.filters'
(/home/damir/.local/lib/python3.10/site-packages/jinja2/filters.py)
ok: [damir-VirtualBox] => {
    "msg": "Docker installed: Docker version 29.2.1, build a5c7197"
}

RUNNING HANDLER [docker : restart docker] **********************************************************************
changed: [damir-VirtualBox]

PLAY RECAP *****************************************************************************************************
damir-VirtualBox           : ok=15   changed=8    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0   
```

### Terminal output from SECOND provision.yml run
```
$ ansible-playbook playbooks/provision.yml

PLAY [Provision web servers with common tools and Docker] ******************************************************

TASK [Gathering Facts] *****************************************************************************************
ok: [damir-VirtualBox]

TASK [common : Update apt cache] *******************************************************************************
ok: [damir-VirtualBox]

TASK [common : Install essential packages] *********************************************************************
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/core.py) as it seems to be
invalid: cannot import name 'environmentfilter' from 'jinja2.filters' (/home/damir/.local/lib/python3.10/site-
packages/jinja2/filters.py)
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/mathstuff.py) as it seems to
be invalid: cannot import name 'environmentfilter' from 'jinja2.filters'
(/home/damir/.local/lib/python3.10/site-packages/jinja2/filters.py)
ok: [damir-VirtualBox]

TASK [common : Set timezone] ***********************************************************************************
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/core.py) as it seems to be
invalid: cannot import name 'environmentfilter' from 'jinja2.filters' (/home/damir/.local/lib/python3.10/site-
packages/jinja2/filters.py)
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/mathstuff.py) as it seems to
be invalid: cannot import name 'environmentfilter' from 'jinja2.filters'
(/home/damir/.local/lib/python3.10/site-packages/jinja2/filters.py)
ok: [damir-VirtualBox]

TASK [docker : Remove old Docker packages (if any)] ************************************************************
ok: [damir-VirtualBox]

TASK [docker : Install required system packages] ***************************************************************
ok: [damir-VirtualBox]

TASK [docker : Add Docker GPG key] *****************************************************************************
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/core.py) as it seems to be
invalid: cannot import name 'environmentfilter' from 'jinja2.filters' (/home/damir/.local/lib/python3.10/site-
packages/jinja2/filters.py)
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/mathstuff.py) as it seems to
be invalid: cannot import name 'environmentfilter' from 'jinja2.filters'
(/home/damir/.local/lib/python3.10/site-packages/jinja2/filters.py)
ok: [damir-VirtualBox]

TASK [docker : Add Docker repository] **************************************************************************
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/core.py) as it seems to be
invalid: cannot import name 'environmentfilter' from 'jinja2.filters' (/home/damir/.local/lib/python3.10/site-
packages/jinja2/filters.py)
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/mathstuff.py) as it seems to
be invalid: cannot import name 'environmentfilter' from 'jinja2.filters'
(/home/damir/.local/lib/python3.10/site-packages/jinja2/filters.py)
ok: [damir-VirtualBox]

TASK [docker : Install Docker packages] ************************************************************************
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/core.py) as it seems to be
invalid: cannot import name 'environmentfilter' from 'jinja2.filters' (/home/damir/.local/lib/python3.10/site-
packages/jinja2/filters.py)
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/mathstuff.py) as it seems to
be invalid: cannot import name 'environmentfilter' from 'jinja2.filters'
(/home/damir/.local/lib/python3.10/site-packages/jinja2/filters.py)
ok: [damir-VirtualBox]

TASK [docker : Ensure Docker service is running and enabled] ***************************************************
ok: [damir-VirtualBox]

TASK [docker : Add user to docker group] ***********************************************************************
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/core.py) as it seems to be
invalid: cannot import name 'environmentfilter' from 'jinja2.filters' (/home/damir/.local/lib/python3.10/site-
packages/jinja2/filters.py)
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/mathstuff.py) as it seems to
be invalid: cannot import name 'environmentfilter' from 'jinja2.filters'
(/home/damir/.local/lib/python3.10/site-packages/jinja2/filters.py)
ok: [damir-VirtualBox]

TASK [docker : Install Python Docker module for Ansible] *******************************************************
ok: [damir-VirtualBox]

TASK [docker : Verify Docker installation] *********************************************************************
ok: [damir-VirtualBox]

TASK [docker : Display Docker version] *************************************************************************
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/core.py) as it seems to be
invalid: cannot import name 'environmentfilter' from 'jinja2.filters' (/home/damir/.local/lib/python3.10/site-
packages/jinja2/filters.py)
[WARNING]: Skipping plugin (/usr/lib/python3/dist-packages/ansible/plugins/filter/mathstuff.py) as it seems to
be invalid: cannot import name 'environmentfilter' from 'jinja2.filters'
(/home/damir/.local/lib/python3.10/site-packages/jinja2/filters.py)
ok: [damir-VirtualBox] => {
    "msg": "Docker installed: Docker version 29.2.1, build a5c7197"
}

PLAY RECAP *****************************************************************************************************
damir-VirtualBox           : ok=14   changed=0    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0   
```

### First Run Analysis
What changed:
- Common role: Updated apt cache, installed packages
- Docker role: Added GPG key, repository, installed Docker, added user to group, installed Python modules
- Handlers: Restarted Docker service

### Second Run Analysis
What didn't change:
- No package installations (all already present)
- No repository additions (already configured)
- No user modifications (already in docker group)
- No service restarts (configuration stable)

### Why Idempotency Works
Ansible modules are designed to be idempotent:
- apt module: Checks if package is installed before attempting installation
- docker_container: Checks container state before creating/modifying
- user module: Verifies group membership before adding
- systemd: Confirms service state before changing
Each task checks the current state of the system and only makes changes when necessary, ensuring predictable and safe executions.

## 4. Ansible Vault Usage
### Secure Credential Storage
```
$ cat group_vars/all.yml 
$ANSIBLE_VAULT;1.1;AES256
37363538613332346638663330653662623065303465383039323431613066323366376465393965
3338613132343965316433313431346233343935383430300a346130633437633365343664353132
39353134616263626561356135356133313834343936386664646233306639366463303166303630
6435343166373463370a623635656434353730356366373738643564616432663065653237353334
62393931313930363734666133656463666535306435633765373764303131373462643961633733
36653863633037663766656661316362633965326632303037623634383535336261393463666231
35646431386239353037306535303663323161386536303663386230313165303239636466373264
33343638616264666438376665313839646335633231633532623733306132613765386263363439
39363337643531393634666131623163666432663631363734303333373032653663383666623035
32363138396432613661313630383763303231383233656335656335613163336566646364336239
63333034346632336561633263663136383539666562396635633434623966396166313063393038
31666265633161613537353430333133663732363237636164613765386564663433346539636135
37396237616537623565393566636134396138616538333065656164663166306338633865633464
39333538623033643934373163643330356132383361343934323066303735373961653637663865
31613230623435316330393531646538333162353363656363376132306261346163306430353432
38303737316666333731373065383030633731326638343932653938333131336464353530313765
62653065336138656235323932643831643238306130353731666233366530333361313737353233
65343932353839343665663238376163373233623931343036346163613632353637376237396663
38303266646665626231646463313530363466393335333565643934363833633463303238653764
39613761633234363133383236303039623635386261346266373663313064373962306630313564
35383932666139373336333634396365303364653339613938306238636139386337653331326266
31356134333662306538323465663261356433616433343933633439663663623238363363333065
396137666236663538343338633061376438
```

### Vault Password Management Strategy
1. Local Development: Use .vault_pass file with restricted permissions (600)
2. CI/CD Pipeline: Store password in secure CI/CD variables
3. Team Environment: Use password manager with shared access
4. Backup: Encrypted backup of password in company vault

### Why Ansible Vault is Important
- Security: Prevents exposure of sensitive data in version control
- Compliance: Meets security requirements for credential management
- Auditability: Encrypted files show clear intent to protect data
- Separation of concerns: Keeps configuration separate from secrets

## 5. Deployment Verification

### Terminal output from deploy.yml run
```
$ ansible-playbook playbooks/deploy.yml --ask-vault-pass
Vault password: 

PLAY [Deploy application to web servers] ***********************************************************************

TASK [Gathering Facts] *****************************************************************************************
ok: [damir-VirtualBox]

TASK [Display deployment information] **************************************************************************
ok: [damir-VirtualBox] => {
    "msg": [
        "Starting deployment of devops-info-service",
        "Image: damirsadykov/devops-info-service:latest",
        "Target host: damir-VirtualBox"
    ]
}

TASK [app_deploy : Log in to Docker Hub] ***********************************************************************
ok: [damir-VirtualBox]

TASK [app_deploy : Pull Docker image] **************************************************************************
ok: [damir-VirtualBox]

TASK [app_deploy : Check if container is running] **************************************************************
ok: [damir-VirtualBox]

TASK [app_deploy : Stop existing container if running] *********************************************************
skipping: [damir-VirtualBox]

TASK [app_deploy : Remove old container if exists] *************************************************************
skipping: [damir-VirtualBox]

TASK [app_deploy : Run new container] **************************************************************************
changed: [damir-VirtualBox]

TASK [app_deploy : Wait for application to be ready] ***********************************************************
ok: [damir-VirtualBox]

TASK [app_deploy : Check health endpoint] **********************************************************************
ok: [damir-VirtualBox]

TASK [app_deploy : Display container info] *********************************************************************
ok: [damir-VirtualBox] => {
    "msg": [
        "Application deployed successfully!",
        "Container: devops-info-service",
        "Port: 5000",
        "Health endpoint: http://10.76.148.84:5000/health",
        "Main endpoint: http://10.76.148.84:5000/"
    ]
}

TASK [Verify deployment] ***************************************************************************************
ok: [damir-VirtualBox]

TASK [Show deployment status] **********************************************************************************
ok: [damir-VirtualBox] => {
    "msg": [
        "Deployment completed successfully!",
        "Container status: running",
        "Container started: 2026-02-25T16:12:17.417271155Z"
    ]
}

PLAY RECAP *****************************************************************************************************
damir-VirtualBox           : ok=11   changed=1    unreachable=0    failed=0    skipped=2    rescued=0    ignored=0   

```

### Container status: docker ps output
```
 ansible webservers -a "docker ps" --ask-vault-pass
Vault password: 
damir-VirtualBox | CHANGED | rc=0 >>
CONTAINER ID   IMAGE                                     COMMAND           CREATED         STATUS         PORTS                    NAMES
e26c8e3625c7   damirsadykov/devops-info-service:latest   "python app.py"   3 minutes ago   Up 3 minutes   0.0.0.0:5000->5000/tcp   devops-info-service
```

### Health check verification: curl outputs
```
$ curl http://10.76.148.84:5000/health
{"status":"healthy","timestamp":"2026-02-25T16:13:14.969966Z","uptime_seconds":57}
$ curl http://10.76.148.84:5000/
{"endpoints":[{"description":"Service information","method":"GET","path":"/"},{"description":"Health check","method":"GET","path":"/health"}],"request":{"client_ip":"10.76.148.244","method":"GET","path":"/","user_agent":"curl/7.81.0"},"runtime":{"current_time":"2026-02-25T16:13:30.343720Z","timezone":"UTC","uptime_human":"0 hours, 1 minute","uptime_seconds":72},"service":{"description":"DevOps course info service","framework":"Flask","name":"devops-info-service","version":"1.0.0"},"system":{"architecture":"x86_64","cpu_count":1,"hostname":"e26c8e3625c7","platform":"Linux","platform_version":"#100~22.04.1-Ubuntu SMP PREEMPT_DYNAMIC Mon Jan 19 17:10:19 UTC ","python_version":"3.12.12"}}
```

## Key Decisions
### Why use roles instead of plain playbooks?
Roles provide modular organization, making the codebase easier to understand, maintain, and reuse across different projects. They also enable team collaboration by allowing different people to work on different components simultaneously.

### How do roles improve reusability?
Roles encapsulate specific functionality with well-defined interfaces (variables and defaults), allowing them to be shared via Ansible Galaxy and used in multiple playbooks with different configurations.

### What makes a task idempotent?
A task is idempotent when it checks the current state before making changes and only acts if the desired state differs from the current state, ensuring multiple runs produce the same result.

### How do handlers improve efficiency?
Handlers run only when notified by tasks and execute at the end of the play, preventing multiple unnecessary restarts and ensuring efficient resource usage.

### Why is Ansible Vault necessary?
Ansible Vault protects sensitive information like passwords and API keys from exposure in version control systems, maintaining security while allowing infrastructure to be defined as code.