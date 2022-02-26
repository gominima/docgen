package main

import "strings"

/**
 * Get the type from a documentation command line
 * @param {string} [line] The comment line to parse
 * @returns {string,string}
 */
func GetType(line string) (string, string) {
	array := Split(line, " ")
	Type := ""
	for _, word := range array {
		if StartsWith(word, "{") && EndsWith(word, "}") {
			Type = Remove(Remove(word, "{"), "}")
		}
	}
	line = Replace(line, "{"+Type+"}", "")
	return Type, line
}

/**
 * Get the name from a documentation command line
 * @param {string} [line] The comment line to parse
 * @returns {string,string}
 */
func GetName(line string) (string, string) {
	array := Split(line, " ")
	Name := ""
	for _, word := range array {
		if StartsWith(word, "[") && EndsWith(word, "]") {
			Name = Remove(Remove(word, "["), "]")
		}
	}
	line = Replace(line, "["+Name+"]", "")
	return Name, line
}

/**
 * Split a line by a seperator
 * @param {string} [line] The comment line to split
 * @param {string} [seperator] The seperator to split line with
 * @returns {[]string}
 */
func Split(line string, seperator string) []string {
	return strings.Split(line, seperator)
}

/**
 * Check if a line by starts with a prefix
 * @param {string} [line] The comment line to check
 * @param {string} [prefix] The prefix to check
 * @returns {bool}
 */
func StartsWith(line string, prefix string) bool {
	return strings.HasPrefix(line, prefix)
}

/**
 * Check if a line by ends with a suffix
 * @param {string} [line] The comment line to check
 * @param {string} [suffix] The suffix to check
 * @returns {bool}
 */
func EndsWith(line string, suffix string) bool {
	return strings.HasSuffix(line, suffix)
}

/**
 * Replace a word by another one
 * @param {string} [line] The comment line to replace the word in
 * @param {string} [replace] The word to replace
 * @param {string} [with] The word to replace with
 * @returns {string}
 */
func Replace(line string, replace string, with string) string {
	return strings.ReplaceAll(line, replace, with)
}

/**
 * Replace a word by another one
 * @param {string} [line] The comment line to replace the word in
 * @param {string} [replace] The word to remove
 * @returns {string}
 */
func Remove(line string, remove string) string {
	return strings.ReplaceAll(line, remove, "")
}

/**
 * Trim spaces from a line
 * @param {string} [line] The comment line to trim
 * @returns {string}
 */
func Trim(line string) string {
	return strings.TrimSpace(line)
}

/**
 * Check if given line starts with func
 * @param {string} [data] The comment line to check
 * @returns {bool}
 */
func IsFunctionLine(data string) bool {
	return strings.HasPrefix(data, "func")
}

/**
 * Check if given line starts with type
 * @param {string} [data] The comment line to check
 * @returns {bool}
 */
func IsStructureLine(data string) bool {
	return strings.HasPrefix(data, "type")
}

/**
 * Check if given comment documents a function of structure
 * @param {string} [data] The comment to check
 * @returns {string}
 */
func IsFunctionOfStructureLine(data string) bool {
	array := Split(Trim(Remove(data, "func")), " ")
	if len(array) == 0 {
		return false
	}
	return StartsWith(array[0], "(") && EndsWith(array[1], ")") 
}

/**
 * Check if given comment documents a function
 * @param {string} [data] The comment to check
 * @returns {bool}
 */
func IsFunction(data string) bool {
	array := strings.Split(data, "\n")
	return IsFunctionLine(array[len(array)-1])
}

/**
 * Check if given comment documents a structure
 * @param {string} [data] The comment to check
 * @returns {bool}
 */
func IsStructure(data string) bool {
	array := strings.Split(data, "\n")
	return IsStructureLine(array[len(array)-1])
}