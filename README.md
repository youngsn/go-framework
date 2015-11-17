# veronica
golang back-end framework &amp; code specification

Veronica is a back-end service framework written by golang. You can build your own back-end service and code specification quickly.

Veronica provides lots of useful infrasturcture and code specification, basically covered most elements back-end service should have. And also provides many convenient components which golang style implements.


## Getting started

> Firstly, getting all depends packages.
> Then, Clone veronica from github.com, change the dirname to your own project name.
> Run script/init.sh project\_name, change your own GOPATH to this dir in .bash\_profile, reload bash.
> Now project has been initialized successfully, you can start to development your own back-end programs.

``` shell
$ git clone https://github.com/youngsn/veronica
$ mv veronica $proj_name
$ chmod +x script/*
$ ./script/init.sh $proj_name
```

> Veronica we have a demo module, this means you can compile the project and start to run:

``` shell
$ ./script/build.sh
$ ./bin/$proj -c conf/conf.development.toml
```

## Components

> veronica provides a lot of basic components & module implements & code specification, you can change to your own quickly:

components:
- log engine (seelog)
- MySQL orm library (gorm)
- config parser (toml)
- module manager (provides effective way to control each module)
- pprof (golang performance monitor)
- monitor interface (monitor module status)
- ticker tasks (run task interval)
