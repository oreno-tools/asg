# asg

## これは

* AutoScaling Group の Desired Capacity, Max Size をコマンドラインで操作するツールです

## Install


```sh
# Get latest version
v=$(curl -s 'https://api.github.com/repos/oreno-tools/asg/releases' | jq -r '.[0].tag_name')
# For macOS
$ wget https://github.com/oreno-tools/asg/releases/download/${v}/asg_darwin_amd64 -O ~/bin/asg && chmod +x ~/bin/asg
# For Linux
$ wget https://github.com/oreno-tools/asg/releases/download/${v}/asg_linux_amd64 -O ~/bin/asg && chmod +x ~/bin/asg
```

## Help

```sh
$ asg --help
Usage of asg:
  -append string
        Set a Append capacity number.
  -desired string
        Set a Desired capacity number.
  -dryrun
        Show a update execution.
  -group string
        Set a AutoScaling Group Name.
  -max string
        Set a Max capacity number.
  -per string
        Set a OnDemand percentage number (%).
  -version
        Print version number.
  -wait string
        Set a Wait time (sec).
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
| oreno-autoscaling-demo   |                 2 |                2 |        2 |       10 |
+--------------------------+-------------------+------------------+----------+----------+
```

### update OnDemand percentage

```sh
$ asg --group=oreno-terraform-dev-demo --per=10 --dryrun
Will be updated as follows...
  OnDemand Percentage : 10

$ asg --group=oreno-terraform-dev-demo --per=10
Change the ondemand percentage of AutoScaling Group: oreno-terraform-dev-demo.
Do you want to continue processing? (y/n): y
+--------------------------+-------------------+------------------+----------+----------+---------------------+
|  AUTOSCALING GROUP NAME  | RUNNING INSTANCES | DESIRED CAPACITY | MIN SIZE | MAX SIZE | ONDEMAND PERCENTAGE |
+--------------------------+-------------------+------------------+----------+----------+---------------------+
| oreno-autoscaling-demo   |                 0 |                0 |        0 |       10 |                  10 |
+--------------------------+-------------------+------------------+----------+----------+---------------------+
```
