package main

import (
	"encoding/xml" // XML parsing support
	"flag"         // Command-line flag parsing
	"fmt"          // I/O formatting
	"os"           // Operating system functionalities
	"os/exec"      // External command execution
	"strconv"
	"strings" // String manipulation functions

	"github.com/fatih/color" // Colorized output in terminal
)

// StringResource defines the structure for parsing strings.xml files.
type StringResource struct {
	XMLName xml.Name `xml:"resources"` // The root XML element
	Strings []String `xml:"string"`    // Slice of String elements within the file
}

// String holds data for a single string element within strings.xml.
type String struct {
	Name string `xml:"name,attr"` // String's name attribute
	Text string `xml:",chardata"` // Text content of the string element
}

type Manifest struct {
	XMLName    xml.Name `xml:"manifest"`
	Activities []App    `xml:"application>activity"`
	Aliases    []App    `xml:"application>activity-alias"`
	Services   []App    `xml:"application>service"`
	Receivers  []App    `xml:"application>receiver"`
}

// App encapsulates an application component like an activity or service, including its intent filters.
type App struct {
	Name     string         `xml:"name,attr"`     // Component name
	Exported string         `xml:"exported,attr"` // Exported status
	Filters  []IntentFilter `xml:"intent-filter"` // Intent filters
}

// IntentFilter contains actions and data elements for filtering intents.
type IntentFilter struct {
	Actions []Action `xml:"action"` // Actions within the filter
	Data    []Data   `xml:"data"`   // Data elements specifying URI patterns
}

// Action defines an action element within an intent-filter.
type Action struct {
	Name string `xml:"name,attr"` // Action name
}

// Data represents a data element within an intent-filter, detailing URI handling.
type Data struct {
	Scheme      string `xml:"scheme,attr"`      // URI scheme
	Host        string `xml:"host,attr"`        // Hostname
	Port        string `xml:"port,attr"`        // Port number
	Path        string `xml:"path,attr"`        // Exact path
	PathPrefix  string `xml:"pathPrefix,attr"`  // Path prefix
	PathPattern string `xml:"pathPattern,attr"` // Path pattern
}

// IsSchemeData checks if the Data struct represents a URI scheme.
func (d Data) IsSchemeData() bool {
	return d.Scheme != "" || d.Host != "" || d.Port != "" || d.Path != "" || d.PathPrefix != "" || d.PathPattern != ""
}

// Uses apktool to decompile an APK file to a specified output directory.
func decompileAPK(apkPath string) (string, error) {
	outputDir := strings.TrimSuffix(apkPath, ".apk") + "_decompiled"    // Naming the output directory
	cmd := exec.Command("apktool", "d", apkPath, "-o", outputDir, "-f") // Constructing the apktool command
	err := cmd.Run()                                                    // Executing the command
	if err != nil {
		return "", err // Error handling for command execution failure
	}
	return outputDir, nil // Successful decompilation returns the output directory
}

// displayHelp
func displayHelp() {
	color.Yellow("Usage: deeeeper [OPTIONS]\n")
	color.Yellow("Options:\n")
	color.Yellow("  -apk <path>       Path to the APK file to be decompiled\n")
	color.Yellow("  -folder <path>    Folder to search in if APK is already decompiled\n")
	color.Yellow("  -h, --help        Display this help and exit\n")
}

// displayBanner
func displayBanner() {
	banner := `
    dMMMMb  dMMMMMP dMMMMMP dMMMMMP dMMMMMP dMMMMb  dMMMMMP dMMMMb 
   dMP VMP dMP     dMP     dMP     dMP     dMP.dMP dMP     dMP.dMP 
  dMP dMP dMMMP   dMMMP   dMMMP   dMMMP   dMMMMP" dMMMP   dMMMMK"  
 dMP.aMP dMP     dMP     dMP     dMP     dMP     dMP     dMP"AMF   
dMMMMP" dMMMMMP dMMMMMP dMMMMMP dMMMMMP dMP     dMMMMMP dMP dMP    

 	Deeeeper - Decompile, find activities and deeplinks
 	Version: 1.0.1
	`
	color.Magenta(banner)
}

// processComponents processes each application component and prints detailed info with colors
func processComponents(components []App) {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	for _, component := range components {
		// Convert the exported attribute to a boolean for easier handling
		exported, err := strconv.ParseBool(component.Exported)
		if err != nil {
			// If the exported attribute is missing or invalid, treat the component as not exported
			exported = false
		}

		// Only process and display components that are exported
		if exported {
			fmt.Printf("%s (exported=%t)\n", cyan(component.Name), exported)

			// Process each intent filter within the component
			for _, filter := range component.Filters {
				for _, action := range filter.Actions {
					fmt.Printf("  %s\n", green(action.Name))
				}
				for _, data := range filter.Data {
					uri := constructURI(data)
					if uri != "" {
						fmt.Printf("  %s\n", green(uri))
					}
				}
			}
		}
	}
}

// constructURI builds a URI string from Data struct
func constructURI(data Data) string {
	if !data.IsSchemeData() {
		return ""
	}
	// Construct the path correctly, considering all attributes (path, pathPrefix, pathPattern)
	var path string
	if data.Path != "" {
		path = data.Path
	} else if data.PathPrefix != "" {
		path = data.PathPrefix
	} else if data.PathPattern != "" {
		path = data.PathPattern
	}

	// Ensure the path starts with a "/"
	if path != "" && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	uri := fmt.Sprintf("%s://%s%s", data.Scheme, data.Host, path)
	return uri
}

func main() {
	displayBanner()

	// Command-line flags definition
	apkPath := flag.String("apk", "", "Path to the APK file to be decompiled")
	folderPath := flag.String("folder", "", "Folder to search in if APK is already decompiled")
	help := flag.Bool("help", false, "Display help")
	flag.BoolVar(help, "h", false, "Display help (shorthand)")

	flag.Parse() // Parsing the command-line flags

	if *help { // If help flag is invoked, display help menu
		displayHelp()
		return // Exit after displaying help
	}

	// Variables for manifest and strings file paths
	var manifestPath, stringsPath string

	if *apkPath != "" { // Proceed if APK path is provided
		color.Green("Decompiling APK...")
		outputDir, err := decompileAPK(*apkPath)
		if err != nil { // Handling errors from APK decompilation
			color.Red("Error decompiling APK: %s\n", err)
			os.Exit(1) // Exiting with error code
		}
		// Setting paths for manifest and strings within the decompiled directory
		manifestPath = fmt.Sprintf("%s/AndroidManifest.xml", outputDir)
		stringsPath = fmt.Sprintf("%s/res/values/strings.xml", outputDir)
	} else if *folderPath != "" { // If only the folder path is provided
		color.Green("Using provided folder for search...")
		// Directly set paths assuming the standard structure within the folder
		manifestPath = fmt.Sprintf("%s/AndroidManifest.xml", *folderPath)
		stringsPath = fmt.Sprintf("%s/res/values/strings.xml", *folderPath)
	} else {
		color.Red("Please provide either an APK file or a folder to proceed.")
		os.Exit(1) // Exit if neither flag is provided
	}

	// Reading and parsing strings.xml
	stringsFile, err := os.ReadFile(stringsPath)
	if err != nil { // Error handling for file reading failure
		color.Red("Error reading strings file: %s\n", err)
		os.Exit(1) // Exiting with error code
	}

	var stringResources StringResource
	xml.Unmarshal(stringsFile, &stringResources) // Unmarshalling XML into struct
	stringMap := make(map[string]string)         // Map for string name-value pairs
	for _, s := range stringResources.Strings {
		stringMap[s.Name] = s.Text // Populating the map
	}

	// Reading and preprocessing AndroidManifest.xml
	manifestFile, err := os.ReadFile(manifestPath)
	if err != nil { // Error handling for file reading failure
		color.Red("Error reading manifest file: %s\n", err)
		os.Exit(1) // Exiting with error code
	}

	rawManifest := string(manifestFile) // Converting file content to string
	for key, value := range stringMap { // Replacing placeholders with actual string values
		placeholder := fmt.Sprintf("@string/%s", key)
		rawManifest = strings.ReplaceAll(rawManifest, placeholder, value)
	}

	var manifest Manifest
	err = xml.Unmarshal([]byte(rawManifest), &manifest) // Unmarshalling manifest XML
	if err != nil {                                     // Error handling for XML unmarshalling failure
		color.Red("Error parsing manifest: %s\n", err)
		os.Exit(1) // Exiting with error code
	}

	// Process components
	color.Yellow("\nProcessing Activities:")
	processComponents(manifest.Activities)

	color.Yellow("\nProcessing Aliases:")
	processComponents(manifest.Aliases)

	color.Yellow("\nProcessing Services:")
	processComponents(manifest.Services)

	color.Yellow("\nProcessing Receivers:")
	processComponents(manifest.Receivers)

	color.Green("Done.")
}
