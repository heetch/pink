package pink

import "fmt"

// Help will list the plugins installed below a given root
func Help(root string) {
	plugins, err := findPlugins(root)
	if err != nil {
		panic(err.Error())
	}
	if len(plugins) == 0 {
		fmt.Println("No pink plug-ins are currently installed")
		return
	}
	fmt.Println("The following pink plugins are installed:")
	for _, p := range plugins {
		fmt.Printf("\t%s\n", p)
	}
}
