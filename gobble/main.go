// Gobble is a source-to-source compiler which converts well-formed Scribble protocols 
// into Go code in which protocol adherence is ensured by the use of session types: 
// formal, structured descriptions of a protocol which specify for each message sent a 
// data type, a sequential position and a direction.

// input: one or more validated Scribble .scr local protocol files, passed as arguments
// from the command line

// output: .go files in an ``output'' directory localed within the current directory.

package main

import (
	"log"
	//"os"
	"strconv"
)

func ProcessFiles(fileNames []string, moduleName string) {
	if len(fileNames) == 0 {
		log.Fatal("Please specicfy one or more .scr files to process.")
	} else if len(fileNames) == 1 {
		fileName := fileNames[0]
		tokens := GetTokens(fileName)
		tree := ProcessTokens(tokens)
		translation := TranslateTree(tree)
		WriteTranslation(translation, moduleName)
	} else {
		translations := make([]*MessagesData, 0)
		for _, name := range fileNames {
			tokens := GetTokens(name)
			tree := ProcessTokens(tokens)
			translation := TranslateTree(tree)
			translations = append(translations, translation)
		}
		WriteCombinedTranslation(translations, moduleName)
	}
}

func runCaseStudy() {
	fileNames := []string{"aggregatorLocal.scr", "clientLocal.scr", "brutishAirwaysLocal.scr", "queasyJetLocal.scr"}
	ProcessFiles(fileNames, "Aggregator")
	for _, fileName := range fileNames {
		slice := []string{fileName}
		ProcessFiles(slice, "Aggregator"+"_"+fileName[:len(fileName)-4])
	}
}

func RunTest(testRoles []string, index int) {
	testNameBase := "test" + strconv.Itoa(index+1)
	ProcessFiles(testRoles, testNameBase)
	for _, role := range testRoles {
		slice := []string{role}
		ProcessFiles(slice, role[:len(role)-4])
	}
}

func RunTests() {
	one := []string{"test1_client.scr", "test1_server.scr"}
	two := []string{"test2_client.scr", "test2_server.scr"}
	three := []string{"test3_A.scr", "test3_B.scr", "test3_C.scr"}
	four := []string{"test4_A.scr", "test4_B.scr"}
	five := []string{"test5_A.scr", "test5_B.scr"}
	six := []string{"test6_A.scr", "test6_B.scr"}
	seven := []string{"test7_C.scr", "test7_P.scr"}
	eight := []string{"test8_A.scr", "test8_B.scr", "test8_C.scr"}
	tests := [][]string{one, two, three, four, five, six, seven, eight}
	for i, test := range tests {
		RunTest(test, i)
	}
}

func GenerateSpeedTest() {
	slice := []string{"speedTestServer.scr", "speedTestClient.scr"}
	ProcessFiles(slice, "SpeedTest")
}

func main() {
	allArgs := os.Args
	args := allArgs[1:]
	if len(args) == 0 {
		log.Fatal("Enter the name of one or more .scr files after \"sumProj\"")
	}
	ProcessFiles(args, "Protocol")
}
