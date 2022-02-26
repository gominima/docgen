package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

/**
 * @info The generic data structure extended by other documentation structures
 * @property {string} [Name] The Name of the structure
 * @property {string} [Type] The type of the structure
 * @property {string} [Description] The Description of the structure
 */
type Data struct {
	Name        string `json:"Name,omitempty"`
	Type        string `json:"Type,omitempty"`
	Description string `json:"Description,omitempty"`
}

/**
 * @info The function data structure used for functions
 * @property {string} [Name] The Name of the function
 * @property {string} [Description] The Description of the function
 * @property {[]Data} [Parameters] The Parameters of the function
 * @property {Data} [Returns] The Return value of the function
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
 * @info The structure data structure used for structures
 * @property {string} [Name] The Name of the structure
 * @property {string} [Description] The Description of the structure
 * @property {[]Data} [Properties] The properties of the structure
 */
type StructureData struct {
	Name        string         `json:"Name,omitempty"`
	Line        string         `json:"Line,omitempty"`
	Description string         `json:"Description,omitempty"`
	Functions   []FunctionData `json:"Functions,omitempty"`
	Properties  []Data         `json:"Properties,omitempty"`
}

/**
 * @info The general meta information of the documentation
 * @property {string} [Generator] The Name of the structure
 * @property {string} [Format] The Description of the structure
 * @property {string} [Date] The properties of the structure
 */
type Meta struct {
	Generator string `json:"Generator,omitempty"`
	Format    string `json:"Format,omitempty"`
	Date      string `json:"Date,omitempty"`
}

/**
 * @info The DocgenData used to make the docs JSON
 * @property {Meta} [Meta] The general meta information of the documentation
 * @property {[]FunctionData} [Functions] The Functions of the project
 * @property {[]StructureData} [Date] The Structures of the project
 */
type DocgenData struct {
	Meta       Meta
	Functions  []FunctionData  `json:"Functions,omitempty"`
	Structures []StructureData `json:"Structures,omitempty"`
}

var DocsMatcher = regexp.MustCompile(`\/\*[\s\S]*?\*\/[\r\n]+([^\r\n]+)`)

/**
 * @info Get all files ending in .go from a directory, recursively
 * @param {string} [root] The root directory
 */
func GetFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && EndsWith(info.Name(), "go") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

/**
 * @info Parse the description of a comment
 * @param {string} [data] The comment to parse
 */
func ParseDescription(data string) string {
	description := Trim(Remove(data, "* @info"))
	return description
}

/**
 * @info Parse the parameter of a function comment
 * @param {string} [data] The comment to parse
 */
func ParseParam(data string) Data {
	data = Trim(Remove(data, "* @param"))
	Type, data := GetType(data)
	Name, data := GetName(data)
	Description := Trim(data)
	return Data{Name: Name, Type: Type, Description: Description}
}

/**
 * @info Parse the property of a structure comment
 * @param {string} [data] The comment to parse
 */
func ParseProperty(data string) Data {
	data = Remove(data, "* @property")
	Type, data := GetType(data)
	Name, data := GetName(data)
	Description := Trim(data)
	return Data{Name: Name, Type: Type, Description: Description}
}

/**
 * @info Parse the return value of a return comment
 * @param {string} [data] The comment to parse
 */
func ParseReturn(data string) Data {
	data = Remove(data, "* @returns")
	Type, data := GetType(data)
	Description := Trim(data)
	return Data{Type: Type, Description: Description}
}

/**
 * @info Parse a single line of a structure comment
 * @param {string} [line] The comment to parse
 * @param {StructureData} [StructureDocs] The Structure Docs for adding data
 */
func ParseStructure(line string, StructureDocs StructureData) StructureData {
	line = Trim(line)
	if IsStructureLine(line) {
		StructureDocs.Line = Trim(Remove(line, "{"))
		array := Split(StructureDocs.Line, " ")
		for _, word := range array {
			if IsStructureLine(word) {
				continue
			}
			StructureDocs.Name = word
			break
		}
	}
	if StartsWith(line, "* @info") {
		StructureDocs.Description = ParseDescription(line)
	}
	if StartsWith(line, "* @property") {
		parsed := ParseProperty(line)
		StructureDocs.Properties = append(StructureDocs.Properties, parsed)
	}
	return StructureDocs
}

/**
 * @info Parse a single line of a function comment
 * @param {string} [line] The line of comment
 * @param {FunctionData} [FunctionDocs] The Function Docs for adding data
 */
func ParseFunction(line string, FunctionDocs FunctionData) (FunctionData, string) {
	line, name := Trim(line), ""
	FunctionDocs.Line = Trim(Remove(line, "{"))
	array := Split(FunctionDocs.Line, " ")
	if IsFunctionOfStructureLine(line) {
		name = Remove(Remove(array[2], "*"), ")")
	}
        if IsFunctionLine(line) {
		for _, word := range array {
			if IsFunctionLine(word) {
				continue
			}
			if StartsWith(word, "(") || EndsWith(word, ")") {
				continue
			}
			FunctionDocs.Name = Split(word, "(")[0]
			break
		}
	}
	if StartsWith(line, "* @info") {
		FunctionDocs.Description = ParseDescription(line)
	}
	if StartsWith(line, "* @param") {
		parsed := ParseParam(line)
		FunctionDocs.Parameters = append(FunctionDocs.Parameters, parsed)
	}
	if StartsWith(line, "* @returns") {
		parsed := ParseReturn(line)
		FunctionDocs.Returns = parsed
	}
	return FunctionDocs, name
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
				lines := Split(data, "\n")
				if IsFunction(data) {
					parsed, name := FunctionData{}, ""
					for _, line := range lines {
						parsed, name = ParseFunction(line, parsed)
					}
					if name != "" {
						for i := range DocJson.Structures {
							if DocJson.Structures[i].Name == name {
								DocJson.Structures[i].Functions = append(DocJson.Structures[i].Functions, parsed)
							}
						}
						continue
					}
					DocJson.Functions = append(DocJson.Functions, parsed)
				} else if IsStructure(data) {
					parsed := StructureData{}
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
