# asg

## これは

* AutoScaling Group の Desired Capacity, Max Size をコマンドラインで操作するツールです

## Install


## Help

```sh
$ asg --help
Usage of asg:
  -desired string
        Set a Desired capacity number.
  -dryrun
        Show a update execution.
  -group string
        Set a AutoScaling Group Name.
  -max string
        Set a Max capacity number.
  -version
        Print version number.
```

## Usage

### list up AutoScaling Group

```sh
$ asg
+--------------------------+-------------------+------------------+----------+----------+
|  AUTOSCALING GROUP NAME  | RUNNING INSTANCES | DESIRED CAPACITY | MIN SIZE | MAX SIZE |
+--------------------------+-------------------+------------------+----------+----------+
| oreno-autoscaling-demo   |                 0 |                0 |        0 |        2 |
| oreno-autoscaling-demo1  |                 0 |                0 |        0 |        1 |
+--------------------------+-------------------+------------------+----------+----------+
```

### update desired capacity

```sh
$ asg --group=oreno-autoscaling-demo --desired=2 --dryrun
Will be updated as follows...
  Min              : 2
  Max              : 2
  Desired Capacity : 2

$ asg --group=oreno-autoscaling-demo --desired=2
Change the capacity of AutoScaling Group: oreno-autoscaling-demo
Do you want to continue processing? (y/n): y
+--------------------------+-------------------+------------------+----------+----------+
|  AUTOSCALING GROUP NAME  | RUNNING INSTANCES | DESIRED CAPACITY | MIN SIZE | MAX SIZE |
+--------------------------+-------------------+------------------+----------+----------+
| oreno-autoscaling-demo   |                 0 |                2 |        2 |        2 |
+--------------------------+-------------------+------------------+----------+----------+
```

### update max size

```sh
$ asg --group=oreno-terraform-dev-demo --max=10 --dryrun
Will be updated as follows...
  Min              : 2
  Max              : 10
  Desired Capacity : 2

$ asg --group=oreno-terraform-dev-demo --max=10
Change the capacity of AutoScaling Group: oreno-terraform-dev-demo.
Do you want to continue processing? (y/n): y
+--------------------------+-------------------+------------------+----------+----------+
|  AUTOSCALING GROUP NAME  | RUNNING INSTANCES | DESIRED CAPACITY | MIN SIZE | MAX SIZE |
+--------------------------+-------------------+------------------+----------+----------+
| oreno-terraform-dev-demo |                 2 |                2 |        2 |       10 |
+--------------------------+-------------------+------------------+----------+----------+
```

