package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

type Data struct {
	Type        string
	Name        string
	Description string
}

type FunctionData struct {
	Function   string
	Parameters []Data
	Returns    Data
}

type StructureData struct {
	Structure  string
	Properties []Data
}

type Meta struct {
	Generator string
	Format    string
	Date      string
}

type DocgenData struct {
	Meta       Meta
	Functions  []FunctionData
	Structures []StructureData
}

var TypeMatcher = regexp.MustCompile(`{.*?}`)
var NameMatcher = regexp.MustCompile(`\[[a-z]{1,}\]`)
var FuncMatcher = regexp.MustCompile(`func.*`)
var StructureMatcher = regexp.MustCompile(`type.*`)

func ParseParam(data string) Data {
	data = strings.Replace(data, "@param", "", 1)
	Type := TypeMatcher.FindAllString(data, -1)[0]
	data = strings.ReplaceAll(data, Type, "")
	Name := NameMatcher.FindAllString(data, -1)[0]
	data = strings.ReplaceAll(data,
		Name, "")
	Description := strings.TrimSpace(data)
	return Data{Type: Type,
		Name: Name, Description: Description}
}

func ParseProperty(data string) Data {
	data = strings.Replace(data, "@property", "", 1)
	Type := TypeMatcher.FindAllString(data, -1)[0]
	data = strings.ReplaceAll(data, Type, "")
	Name := NameMatcher.FindAllString(data, -1)[0]
	data = strings.ReplaceAll(data,
		Name, "")
	Description := strings.TrimSpace(data)
	return Data{Type: Type,
		Name: Name, Description: Description}
}

func ParseReturn(data string) Data {
	data = strings.Replace(data, "@returns", "", 1)
	Type := TypeMatcher.FindAllString(data, -1)[0]
	data = strings.ReplaceAll(data,
		Type, "")
	Description := strings.TrimSpace(data)
	return Data{Type: Type, Description: Description}
}

func ParseStructure(line string, doc StructureData) StructureData {
	line = strings.TrimSpace(line)
	parsed := Data{}
	if strings.HasPrefix(line, "@property") {
		parsed = ParseProperty(line)
		doc.Properties = append(doc.Properties, parsed)
	}
	return doc
}

func ParseFunction(line string, doc FunctionData) FunctionData {
	line = strings.TrimSpace(line)
	parsed := Data{}
	if strings.HasPrefix(line, "@param") {
		parsed = ParseParam(line)
		doc.Parameters = append(doc.Parameters, parsed)
	} else if strings.HasPrefix(line, "@returns") {
		parsed = ParseReturn(line)
		doc.Returns = parsed
	}
	return doc
}

func main() {
	DocJson := DocgenData{Meta: Meta{Generator: "1",
		Format: "1",
		Date:   time.Now().String()}}
	args := os.Args[1:]
	matcher := regexp.MustCompile(`\/\*[\s\S]*?\*\/[\r\n]+([^\r\n]+)`)

	if len(args) > 0 {
		for _, file := range args {
			data, err := os.Open(file)

			if err != nil {
				log.Panic(err)
			}

			defer data.Close()

			output, err := io.ReadAll(data)

			if err != nil {
				log.Panic(err)
			}

			content := string(output)

			matches := matcher.FindAllString(content, -1)

			for _, data := range matches {
				lines := strings.Split(data, "\n")
				object := ""
				if FuncMatcher.MatchString(data) {
					object = strings.ReplaceAll(FuncMatcher.FindAllString(data, -1)[0], "{", "")
					parsed := FunctionData{Function: strings.TrimSpace(object)}
					for _, line := range lines {
						parsed = ParseFunction(line, parsed)
					}
					DocJson.Functions = append(DocJson.Functions, parsed)
				} else if StructureMatcher.MatchString(data) {
					object = strings.ReplaceAll(StructureMatcher.FindAllString(data, -1)[0], "{", "")
					parsed := StructureData{Structure: strings.TrimSpace(object)}
					for _, line := range lines {
						parsed = ParseStructure(line, parsed)
					}
					DocJson.Structures = append(DocJson.Structures, parsed)
				}
			}
			file, _ := json.MarshalIndent(DocJson, "", "\t")
			_ = ioutil.WriteFile("output.json", file, 0644)
		}
	}
}
