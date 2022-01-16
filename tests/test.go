package main

/**
	@info The example structure
	@property {string} [name] The name of the structure
	@property {int} [money] The money of the structure
*/
type ExampleStructure struct {
	name  string
	money int
}

/**
	@info The example function
	@param {string} [name] The name to return
	@returns {string}
*/
func Example(name string) string {
	return name
}

/**
	@info The example function
	@param {int} [number] The number to return
	@returns {int}
*/
func ExampleTwo(number int) int {
	return number
}

/**
	@info The example function
	@param {[]byte} [bytes] The bytes to return
	@returns {[]byte}
*/
func ExampleThree(bytes []byte) []byte {
	return bytes
}

/**
	@info The example function that returns the example structure
	@param {string} [name] The name of the structure
	@param {int} [money] The money of the structure
	@returns {ExampleStructure}
*/
func ExampleFour(name string, money int) ExampleStructure {
	return ExampleStructure{name: name, money: money}
}
