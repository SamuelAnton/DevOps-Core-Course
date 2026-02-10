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
I choose "Option A: Semantic Versioning" as I will release my app not so often: for not so big period (many years or month) and there will be not so many updates. So for this course (or for this task) I choose to use semantic versioning.