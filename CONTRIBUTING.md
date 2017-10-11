# Contributing

We accept pull requests from everyone. 
To contribute to this project you should follow the following steps:

1. Fork the repo.

2. Clone the project from your repository to your computer:

```
git clone git@github.com:your-username/ergo.git
```
3. Prepare your machine for development:

_Please be aware that this project is written in golang, so you need golang if you want to contribute with code_

* on MacOS, Linux and other unixes:
Open your terminal, cd to your project directory and then run:
```
make tools
make dep
```

*on Windows:
Open PowerShell, cd to your project directory and then run:
```
.\.make.ps1 -tools
.\.make.ps1 -dep
```

4. Make sure that you are all set up by running the tests:

* On OsX, Linux and other unixes:
```
make test-integration
```

* On Windows
```
./.make.ps1 -test_integration
```

5. Make your change. Add tests for your change. Make the tests pass:

Push to your fork and submit a pull request.

6. Wait for us to review your pull request. We will try to review it as soon as possible. 

Some things that will increase the chance that your pull request is accepted:

* Write tests.
* Write a good commit message.

Thank you
