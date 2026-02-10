# Lab 02

## Testing
### Choosing a Testing Framework
I choosed to use `pytest`, because I already a little bit familiar with it. In this semester, but in different course we studied it too, so I will use it here to practice more.

### Test structure explanation
1. Make a test client
2. Request all pathes (health, default and error)
3. Fully check presence of each response section on all depths. Also check types where possible

### How to run tests locally
```sh
$cd DevOps-Core-Course/app_python
$pytest .
```

### Terminal output
```sh
damir@damir-VB:~/Desktop/DevOps/DevOps-Core-Course/app_python$ pytest .
============================================= test session starts ==============================================
platform linux -- Python 3.10.12, pytest-9.0.2, pluggy-1.6.0
rootdir: /home/damir/Desktop/DevOps/DevOps-Core-Course/app_python
collected 3 items                                                                                              

tests/test_app.py ...                                                                                    [100%]

============================================== 3 passed in 0.30s ===============================================
```

## GitHub Actions CI Workflow
### Workflow trigger strategy
The workflow will trigger on every push into lab** branches. Since in this course we need to push only there, it is more than okay practice. After each last (or even single) push to lab**, we make to PRs, so it is not so good to start workflow again 2 times.

### Actions used
| Action | Reason |
-------------------
| actions/checkout@v4 | Clone repo |
| actions/setup-python@v5 | Setup Python |
| docker/login-action@v3 | Login to DockerHub |
| docker/build-push-action@v5 | Build and Push Docker Image |

### Docker tagging strategy
I choose "Option B: Calendar Versioning (CalVer)" with format: "YYYY.MM.DD", so I can easily distinguish each versions: newer date - newer assignment. Also, when I will do lab assignment and make some really small changes in one day - I will save only last version with all fixes and bugs found. Also, I ofcourse use `latest` tag

### Successful workflow run
https://github.com/SamuelAnton/DevOps-Core-Course/actions/runs/21873318288

### Green checkmark - see in screenshots


## CI Best Practices & Security

### Caching implementation
```yaml
- uses: actions/setup-python@v5
        with:
          python-version: '3.12'
          cache: 'pip'
          cache-dependency-path: 'app_python/requirements.txt'
```

### Speed improvement with caching
- Before caching: test stage take 10-16s
- After caching: test stage take 