package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

/**
@info The generic data structure extended by other documentation structures
@property {string} [Name] The Name of the structure
@property {string} [Type] The type of the structure
@property {string} [Description] The Description of the structure
*/
type Data struct {
	Name        string `json:"Name,omitempty"`
	Type        string `json:"Type,omitempty"`
	Description string `json:"Description,omitempty"`
}

/**
@info The function data structure used for functions
@property {string} [Name] The Name of the function
@property {string} [Description] The Description of the function
@property {[]Data} [Parameters] The Parameters of the function
@property {Data} [Returns] The Return value of the function
*/
type FunctionData struct {
	Name        string `json:"Name,omitempty"`
	Line        string `json:"Line,omitempty"`
	Description string `json:"Description,omitempty"`
	Example     string `json:"Example,omitempty"`
	Parameters  []Data `json:"Parameters,omitempty"`
	Returns     Data   `json:"Returns,omitempty"`
}

/**
@info The structure data structure used for structures
@property {string} [Name] The Name of the structure
@property {string} [Description] The Description of the structure
@property {[]Data} [Properties] The properties of the structure
*/
type StructureData struct {
	Name        string         `json:"Name,omitempty"`
	Line        string         `json:"Line,omitempty"`
	Description string         `json:"Description,omitempty"`
	Functions   []FunctionData `json:"Functions,omitempty"`
	Properties  []Data         `json:"Properties,omitempty"`
}

/**
@info The general meta information of the documentation
@property {string} [Generator] The Name of the structure
@property {string} [Format] The Description of the structure
@property {string} [Date] The properties of the structure
*/
type Meta struct {
	Generator string `json:"Generator,omitempty"`
	Format    string `json:"Format,omitempty"`
	Date      string `json:"Date,omitempty"`
}

/**
@info The DocgenData used to make the docs JSON
@property {Meta} [Meta] The general meta information of the documentation
@property {[]FunctionData} [Functions] The Functions of the project
@property {[]StructureData} [Date] The Structures of the project
*/
type DocgenData struct {
	Meta       Meta
	Functions  []FunctionData  `json:"Functions,omitempty"`
	Structures []StructureData `json:"Structures,omitempty"`
}

var (
	DocsMatcher              = regexp.MustCompile(`\/\*[\s\S]*?\*\/[\r\n]+([^\r\n]+)`)
	TypeMatcher              = regexp.MustCompile(`{.*?}`)
	NameMatcher              = regexp.MustCompile(`\[[a-zA-Z.]{1,}\]`)
	FuncMatcher              = regexp.MustCompile(`func.*{`)
	FuncNameMatcher          = regexp.MustCompile(`func.([a-zA-z]*)+.*`)
	StructureFuncNameMatcher = regexp.MustCompile(`func\s\([a-zA-Z]*\s\*?[a-zA-Z]*\)\s([a-zA-z]*).*{`)
	FuncStructureNameMatcher = regexp.MustCompile(`func\s\([a-zA-Z]*\s\*?([a-zA-Z]*)\)\s[a-zA-z]*.*{`)
	StructureMatcher         = regexp.MustCompile(`type.*{`)
	StructureNameMatcher     = regexp.MustCompile(`type.([a-zA-z]*)+.*`)
	DescriptionMatcher       = regexp.MustCompile(`@info.*`)
	ExampleMatcher           = regexp.MustCompile(`@example\s\x60{3,}\n([\s\S]*)\x60{3,}`)
)

/**
@info Get all files ending in .go from a directory, recursively
@param {string} [root] The root directory
*/
func GetFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), "go") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

/**
@info Parse the description of a comment
@param {string} [data] The comment
*/
func ParseDescription(data string) string {
	description := DescriptionMatcher.FindAllString(data, -1)[len(DescriptionMatcher.FindAllString(data, -1))-1]
	description = strings.TrimSpace(strings.Replace(description, "@info", "", 1))
	return description
}

/**
@info Parse the parameter of a function comment
@param {string} [data] The comment
*/
func ParseParam(data string) Data {
	data = strings.Replace(data, "@param", "", 1)
	Type := TypeMatcher.FindAllString(data, -1)[len(TypeMatcher.FindAllString(data, -1))-1]
	data = strings.ReplaceAll(data, Type, "")
	Type = Type[1 : len(Type)-1]
	Name := NameMatcher.FindAllString(data, -1)[len(NameMatcher.FindAllString(data, -1))-1]
	data = strings.ReplaceAll(data, Name, "")
	Name = Name[1 : len(Name)-1]
	Description := strings.TrimSpace(data)
	return Data{Name: Name, Type: Type, Description: Description}
}

/**
@info Parse the example of a function comment
@param {string} [data] The comment
@example ```
	package main

	import (
		"fmt"
	)

	func main() {
		fmt.Println("Hello world!")
	}
```
*/
func ParseExample(data string) string {
	if ExampleMatcher.MatchString(data) {
		return ExampleMatcher.FindAllStringSubmatch(data, -1)[len(ExampleMatcher.FindAllStringSubmatch(data, -1))-1][1]
	}
	return ""
}

/**
@info Parse the property of a structure comment
@param {string} [data] The comment
*/
func ParseProperty(data string) Data {
	data = strings.Replace(data, "@property", "", 1)
	Type := TypeMatcher.FindAllString(data, -1)[len(TypeMatcher.FindAllString(data, -1))-1]
	data = strings.ReplaceAll(data, Type, "")
	Type = Type[1 : len(Type)-1]
	Name := NameMatcher.FindAllString(data, -1)[len(NameMatcher.FindAllString(data, -1))-1]
	data = strings.ReplaceAll(data, Name, "")
	Name = Name[1 : len(Name)-1]
	Description := strings.TrimSpace(data)
	return Data{Name: Name, Type: Type, Description: Description}
}

/**
@info Parse the return value of a return comment
@param {string} [data] The comment
*/
func ParseReturn(data string) Data {
	data = strings.Replace(data, "@returns", "", 1)
	Type := TypeMatcher.FindAllString(data, -1)[len(TypeMatcher.FindAllString(data, -1))-1]
	data = strings.ReplaceAll(data, Type, "")
	Type = Type[1 : len(Type)-1]
	Description := strings.TrimSpace(data)
	return Data{Type: Type, Description: Description}
}

/**
@info Parse a single line of a structure comment
@param {string} [line] The line of comment
@param {StructureData} [StructureDocs] The Structure Docs
*/
func ParseStructure(line string, StructureDocs StructureData) StructureData {
	line = strings.TrimSpace(line)
	parsed := Data{}
	if strings.HasPrefix(line, "@property") {
		parsed = ParseProperty(line)
		StructureDocs.Properties = append(StructureDocs.Properties, parsed)
	}
	return StructureDocs
}

/**
@info Parse a single line of a function comment
@param {string} [line] The line of comment
@param {FunctionData} [FunctionDocs] The Function Docs
*/
func ParseFunction(line string, FunctionDocs FunctionData) FunctionData {
	line = strings.TrimSpace(line)
	parsed := Data{}
	if strings.HasPrefix(line, "@param") {
		parsed = ParseParam(line)
		FunctionDocs.Parameters = append(FunctionDocs.Parameters, parsed)
	} else if strings.HasPrefix(line, "@returns") {
		parsed = ParseReturn(line)
		FunctionDocs.Returns = parsed
	}
	return FunctionDocs
}

func main() {
	DocJson := DocgenData{Meta: Meta{Generator: "1",
		Format: "1",
		Date:   time.Now().String()}}

	args := "."
	outFile := "output.json"

	if len(os.Args) > 1 {
		args = os.Args[1]
	}

	if len(os.Args) > 2 {
		outFile = os.Args[2]
	}

	files, err := GetFiles(args)

	if err != nil {
		log.Fatal(err)
	}

	if len(args) > 0 {
		for _, file := range files {
			data, err := os.Open(file)

			if err != nil {
				log.Fatal(err)
			}

			defer data.Close()

			output, err := io.ReadAll(data)

			if err != nil {
				log.Fatal(err)
			}

			content := string(output)

			matches := DocsMatcher.FindAllString(content, -1)

			for _, data := range matches {
				lines := strings.Split(data, "\n")
				if FuncMatcher.MatchString(data) {
					line := strings.TrimSpace(strings.ReplaceAll(FuncMatcher.FindAllString(data, -1)[len(FuncMatcher.FindAllString(data, -1))-1], "{", ""))
					name := FuncNameMatcher.FindAllStringSubmatch(line, -1)[len(FuncNameMatcher.FindAllStringSubmatch(line, -1))-1][1]
					description := ParseDescription(data)
					example := ParseExample(data)
					if StructureFuncNameMatcher.MatchString(data) {
						name := FuncStructureNameMatcher.FindAllStringSubmatch(data, -1)[len(FuncStructureNameMatcher.FindAllStringSubmatch(data, -1))-1][1]
						funcname := StructureFuncNameMatcher.FindAllStringSubmatch(data, -1)[len(StructureFuncNameMatcher.FindAllStringSubmatch(data, -1))-1][1]
						parsed := FunctionData{Name: funcname, Line: line, Description: description, Example: example}
						for _, line := range lines {
							parsed = ParseFunction(line, parsed)
						}
						for i := range DocJson.Structures {
							if DocJson.Structures[i].Name == name {
								DocJson.Structures[i].Functions = append(DocJson.Structures[i].Functions, parsed)
							}
						}
						continue
					}
					parsed := FunctionData{Name: name, Line: line, Description: description, Example: example}
					for _, line := range lines {
						parsed = ParseFunction(line, parsed)
					}
					DocJson.Functions = append(DocJson.Functions, parsed)
				} else if StructureMatcher.MatchString(data) {
					line := strings.TrimSpace(strings.ReplaceAll(StructureMatcher.FindAllString(data, -1)[len(StructureMatcher.FindAllString(data, -1))-1], "{", ""))
					name := StructureNameMatcher.FindAllStringSubmatch(line, -1)[len(StructureNameMatcher.FindAllStringSubmatch(line, -1))-1][1]
					description := ParseDescription(data)
					parsed := StructureData{Name: name, Line: line, Description: description}
					for _, line := range lines {
						parsed = ParseStructure(line, parsed)
					}
					DocJson.Structures = append(DocJson.Structures, parsed)
				}
			}
			file, _ := json.MarshalIndent(DocJson, "", "\t")
			_ = ioutil.WriteFile(outFile, file, 0644)
		}
	}
}
