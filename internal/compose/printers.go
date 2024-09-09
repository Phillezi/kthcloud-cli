package compose

import (
	"fmt"
)

func PrintServices(services map[string]Service) {
	for name, service := range services {
		fmt.Printf("Service: %s\n", name)
		fmt.Printf("  Image: %s\n", service.Image)

		fmt.Printf("  Environment Variables:\n")
		if len(service.Environment) == 0 {
			fmt.Println("    None")
		} else {
			for envName, value := range service.Environment {
				fmt.Printf("    %s: %s\n", envName, value)
			}
		}

		fmt.Printf("  Ports:\n")
		if len(service.Ports) == 0 {
			fmt.Println("    None")
		} else {
			for _, port := range service.Ports {
				fmt.Printf("    %s\n", port)
			}
		}

		fmt.Printf("  Volumes:\n")
		if len(service.Volumes) == 0 {
			fmt.Println("    None")
		} else {
			for _, volume := range service.Volumes {
				fmt.Printf("    %s\n", volume)
			}
		}

		fmt.Printf("  Command:\n")
		if len(service.Command) == 0 {
			fmt.Println("    None")
		} else {
			fmt.Printf("    %s\n", service.Command)
		}

		fmt.Println("----------------------------")
	}
}
