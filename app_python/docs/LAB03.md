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

## Versioning Strategy
I choose "Option B: Calendar Versioning (CalVer)" with format: "YYYY.MM.DD", so I can easily distinguish each versions: newer date - newer assignment. Also, when I will do lab assignment and make some really small changes in one day - I will save only last version with all fixes and bugs found.

https://github.com/SamuelAnton/DevOps-Core-Course/pull/3