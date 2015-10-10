package herd

import (
	"fmt"
	"os/exec"
	"strings"
)

// Run rabbitmqctl with the specified args. Prints command without execution
// if `dryRun` is set to True
func RunRabbitmqctl(args []string, dryRun bool) (err error) {
	err = nil
	fmt.Printf(
		"Running command: rabbitmqctl %s\n",
		strings.Trim(fmt.Sprint(args), "[]"),
	)
	if dryRun == false {
		err = exec.Command("rabbitmqctl", args...).Run()
	}
	return err
}

// Attempts to join the cluster. Errors are printed and the app is restarted
func JoinCluster(ips []string, dryRun bool) {
	if len(ips) == 0 {
		fmt.Println("No other hosts to join!")
		return
	}
	err := RunRabbitmqctl([]string{"stop_app"}, dryRun)
	if err != nil {
		fmt.Printf("Error stopping the app: %s\n", err)
		return
	}
	for _, ip := range ips {
		hostname := fmt.Sprintf("rabbit@%s", ip)
		err = RunRabbitmqctl([]string{"join_cluster", hostname}, dryRun)
		if err != nil {
			fmt.Printf("Error joining %s\n", hostname)
		} else {
			fmt.Printf("Successfully joined %s\n", hostname)
			break
		}
	}
	err = RunRabbitmqctl([]string{"start_app"}, dryRun)
	if err != nil {
		fmt.Printf("Error restarting the app: %s\n", err)
		return
	}
}
