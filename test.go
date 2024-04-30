package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "check-configs",
		Short: "Check and modify tailwind.config.js and postcss.config.js files if necessary",
		Run:   checkAndModifyConfigs,
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func checkAndModifyConfigs(cmd *cobra.Command, args []string) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	tailwindConfigPath := filepath.Join(dir, "tailwind.config.js")
	if _, err := os.Stat(tailwindConfigPath); os.IsNotExist(err) {
		fmt.Println("tailwind.config.js not found in the current directory")
		fmt.Println("Installing required npm packages...")
		installPackages()
		initTailwindConfig()
		modifyTailwindConfig(tailwindConfigPath)
		modifyIndexFile(dir)
	} else if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("tailwind.config.js found in the current directory")
		modifyTailwindConfig(tailwindConfigPath)
		modifyIndexFile(dir)
	}
}

func installPackages() {
	cmd := exec.Command("npm", "install", "-D", "tailwindcss", "postcss", "autoprefixer")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error installing npm packages:", err)
		os.Exit(1)
	}
}

func initTailwindConfig() {
	cmd := exec.Command("npx", "tailwindcss", "init", "-p")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error initializing tailwind.config.js:", err)
		os.Exit(1)
	}
}
func modifyTailwindConfig(filePath string) {
	if isViteReactApp() {
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading tailwind.config.js file:", err)
			return
		}

		modifiedContent := strings.Replace(string(content), "content: [],", "content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],", 1)

		err = os.WriteFile(filePath, []byte(modifiedContent), 0644)
		if err != nil {
			fmt.Println("Error writing to tailwind.config.js file:", err)
			return
		}

		fmt.Println("Modified tailwind.config.js for Vite React app")
	}
	if isNextJSApp() {
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading tailwind.config.js file:", err)
			return
		}

		modifiedContent := strings.Replace(string(content), "content: [],", `content: ['./app/**/*.{js,ts,jsx,tsx,mdx}','./pages/**/*.{js,ts,jsx,tsx,mdx}','./components/**/*.{js,ts,jsx,tsx,mdx}'],`, 1)

		err = os.WriteFile(filePath, []byte(modifiedContent), 0644)
		if err != nil {
			fmt.Println("Error writing to tailwind.config.js file:", err)
			return
		}

		fmt.Println("Modified tailwind.config.js for Vite React app")
	}
}

func isViteReactApp() bool {
	content, err := os.ReadFile("package.json")
	if err != nil {
		fmt.Println("Error reading package.json file:", err)
		return false
	}

	return strings.Contains(string(content), `"vite"`) && strings.Contains(string(content), `"react"`)
}

func isNextJSApp() bool {
	content, err := os.ReadFile("package.json")
	if err != nil {
		fmt.Println("Error reading package.json file:", err)
		return false
	}

	return strings.Contains(string(content), `"next"`) && strings.Contains(string(content), `"react"`)
}

func modifyIndexFile(dir string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %q: %v\n", path, err)
			return err
		}
		if !info.IsDir() && info.Name() == "index.css" {
			content, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("Error reading index.css file at %q: %v\n", path, err)
				return err
			}
			modifiedContent := fmt.Sprintf("@tailwind base;\n@tailwind components;\n@tailwind utilities;\n\n%s", string(content))

			err = os.WriteFile(path, []byte(modifiedContent), 0644)
			if err != nil {
				fmt.Printf("Error writing to index.css file at %q: %v\n", path, err)
				return err
			}
			fmt.Printf("Modified index.css at %q to include Tailwind CSS directives\n", path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking directory %q: %v\n", dir, err)
	}
}
