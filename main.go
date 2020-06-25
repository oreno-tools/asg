package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	_ "github.com/aws/aws-sdk-go/aws/credentials"
	_ "github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	_ "github.com/kyokomi/emoji"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
)

const (
	AppVersion = "0.0.5"
)

var (
	argVersion    = flag.Bool("version", false, "Print version number.")
	argGroup      = flag.String("group", "", "Set a AutoScaling Group Name.")
	argDesired    = flag.String("desired", "", "Set a Desired capacity number.")
	argAppend     = flag.String("append", "", "Set a Append capacity number.")
	argMax        = flag.String("max", "", "Set a Max capacity number.")
	argPercentage = flag.String("per", "", "Set a OnDemand percentage number (%).")
	argDryrun     = flag.Bool("dryrun", false, "Show a update execution.")
	// argMin     = flag.Int64("min", 0, "Set a Min capacity number.")

	// mega  = emoji.Sprint(":mega:")
	// sushi = emoji.Sprint(":sushi:")
	// warn  = emoji.Sprint(":beer:")

	asg = autoscaling.New(session.New())
)

func main() {
	flag.Parse()

	if *argVersion {
		fmt.Println(AppVersion)
		os.Exit(0)
	}

	var asgGroups *autoscaling.DescribeAutoScalingGroupsOutput
	asgGroups = getGroups(*argGroup)

	if len(asgGroups.AutoScalingGroups) == 0 {
		fmt.Println(fmt.Sprintf("AutoScaling Group `%s` does not exists.", *argGroup))
		os.Exit(1)
	}

	if *argGroup != "" && *argPercentage != "" {
		_per, _ := strconv.ParseInt(*argPercentage, 10, 64)
		if *argDryrun {
			fmt.Printf("Will be updated as follows...\n")
			fmt.Printf("  OnDemand Percentage : %d\n", _per)
			os.Exit(0)
		}
		fmt.Printf("Change the ondemand percentage of AutoScaling Group: \x1b[31m%s\x1b[0m.\n", *argGroup)
		fmt.Printf("Do you want to continue processing? (y/n): ")
		var stdin string
		fmt.Scan(&stdin)
		switch stdin {
		case "y", "Y":
			result := setOndemandPercent(*argGroup, _per)
			if !result {
				fmt.Println("Update Ondemand percentage Failure!!")
				os.Exit(1)
			}
		case "n", "N":
			fmt.Println("Interrupted.")
			os.Exit(0)
		default:
			fmt.Println("Interrupted.")
			os.Exit(0)
		}
		asgGroups = getGroups(*argGroup)
	}

	if *argGroup != "" && *argDesired == "" && *argAppend != "" {
		_, _, _, instances := getDetectedSize(asgGroups)
		_append, _ := strconv.Atoi(*argAppend)
		_desired := instances + _append
		*argDesired = strconv.Itoa(_desired)
	}

	if *argGroup != "" && *argDesired == "" && *argMax != "" {
		_, _, desired, _ := getDetectedSize(asgGroups)
		*argDesired = strconv.FormatInt(desired, 10)
	}

	if *argGroup != "" && *argDesired != "" {
		min, max, _, _ := getDetectedSize(asgGroups)
		// fmt.Printf("min: %d\n", min)
		// fmt.Printf("max: %d\n", max)
		_desired, _ := strconv.ParseInt(*argDesired, 10, 64)
		_max, _ := strconv.ParseInt(*argMax, 10, 64)

		// Minimun が 0 又は, Desired Capacity が 0 である場合には,  Desired Capacity を代入する
		// Minumun が Desired Capacity よりも大きい場合には, Desired Capacity を代入
		if min == 0 || _desired == 0 || min > _desired {
			min = _desired
		}

		// argDesired が 0 以上で且つ, 現在の max よりも大きい場合には, argDesired を max に代入
		if _desired > 0 && max < _desired {
			max = _desired
		} else if _max > 0 && _desired >= 0 {
			// argDesired が 0 以上で且つ, argMax が 0 以上の場合には, argMax を max に代入
			max = _max
		}

		if *argDryrun {
			fmt.Printf("Will be updated as follows...\n")
			fmt.Printf("  Min              : %d\n", min)
			fmt.Printf("  Max              : %d\n", max)
			fmt.Printf("  Desired Capacity : %d\n", _desired)
			os.Exit(0)
		}

		// fmt.Printf("%sChange the capacity of AutoScaling Group: \x1b[31m%s\x1b[0m.\n", mega, *argGroup)
		fmt.Printf("Change the capacity of AutoScaling Group: \x1b[31m%s\x1b[0m.\n", *argGroup)
		fmt.Printf("Do you want to continue processing? (y/n): ")
		var stdin string
		fmt.Scan(&stdin)
		switch stdin {
		case "y", "Y":
			if max >= _desired {
				result := setCapacity(*argGroup, min, max, _desired)
				if !result {
					fmt.Println("Update Capacity Failure!!")
					os.Exit(1)
				}
				fmt.Println("\n")
			} else {
				fmt.Println("Max size must be greater than desired capacity!!")
				os.Exit(1)
			}
		case "n", "N":
			fmt.Println("Interrupted.")
			os.Exit(0)
		default:
			fmt.Println("Interrupted.")
			os.Exit(0)
		}
		asgGroups = getGroups(*argGroup)
	}

	outputGroups(asgGroups)
}

func getDetectedSize(asgGroup *autoscaling.DescribeAutoScalingGroupsOutput) (min int64, max int64, desired int64, ins int) {
	min = *asgGroup.AutoScalingGroups[0].MinSize
	max = *asgGroup.AutoScalingGroups[0].MaxSize
	desired = *asgGroup.AutoScalingGroups[0].DesiredCapacity
	ins = len(asgGroup.AutoScalingGroups[0].Instances)

	return min, max, desired, ins
}

func getGroups(groupName string) *autoscaling.DescribeAutoScalingGroupsOutput {
	params := &autoscaling.DescribeAutoScalingGroupsInput{}
	if groupName != "" {
		params.SetAutoScalingGroupNames(
			[]*string{
				aws.String(groupName),
			},
		)
	}

	asgGroups, err := asg.DescribeAutoScalingGroups(params)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		} else {
			fmt.Println(err.Error())
		}
	}

	return asgGroups
}

func setCapacity(groupName string, min int64, max int64, desiredCap int64) bool {
	params := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(groupName),
		MaxSize:              aws.Int64(max),
		MinSize:              aws.Int64(min),
		DesiredCapacity:      aws.Int64(desiredCap),
	}

	_, err := asg.UpdateAutoScalingGroup(params)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
			return false
		} else {
			fmt.Println(err.Error())
			return false
		}
		return false
	}
	return true
}

func setOndemandPercent(groupName string, percentage int64) bool {
	instanceDist := &autoscaling.InstancesDistribution{
		OnDemandPercentageAboveBaseCapacity: aws.Int64(percentage),
	}
	mixedPolicy := &autoscaling.MixedInstancesPolicy{
		InstancesDistribution: instanceDist,
	}
	params := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(groupName),
		MixedInstancesPolicy: mixedPolicy,
	}

	_, err := asg.UpdateAutoScalingGroup(params)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
			return false
		} else {
			fmt.Println(err.Error())
			return false
		}
		return false
	}
	return true
}

func outputGroups(asgGroups *autoscaling.DescribeAutoScalingGroupsOutput) {
	allASG := [][]string{}
	for _, g := range asgGroups.AutoScalingGroups {
		var percentage string
		if g.MixedInstancesPolicy == nil {
			percentage = "N/A"
		} else {
			percentage = strconv.FormatInt(*g.MixedInstancesPolicy.InstancesDistribution.OnDemandPercentageAboveBaseCapacity, 10)
		}
		ASGroup := []string{
			*g.AutoScalingGroupName,
			strconv.Itoa(len(g.Instances)),
			strconv.FormatInt(*g.DesiredCapacity, 10),
			strconv.FormatInt(*g.MinSize, 10),
			strconv.FormatInt(*g.MaxSize, 10),
			percentage,
		}
		allASG = append(allASG, ASGroup)
	}

	printTable(allASG)
}

func printTable(data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"AutoScaling Group Name", "Running Instances", "Desired Capacity", "Min Size", "Max Size", "Ondemand Percentage"})
	table.AppendBulk(data)
	table.Render()
}
