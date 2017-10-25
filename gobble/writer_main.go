// The writer takes in the MessagesData struct produced by the
// translator and reads through the arrays of structs that it contains
// to generate a working Go API and implementation of the Scribble
// protocol. If the user has passed Gobble one Scribble protocol to
// process, Gobble will output an inter-system implementation with code
// allowing it to communicate with other roles over a TCP/IP connection.
// If the user has passed Gobble multiple Scribble protocols Gobble will
// combine them to form an intra-system implementation in which all roles
// will be launched by a single process and will communicate directly
// with one another across channels.

// input: []MessagesData containing arrays ([]RecData, []ParallelData,
// etc) and maps (map[string]*RecData, map[string]*ParallelData, etc.)

// output: .go files in an ``output'' directory localed within the
// current directory.

package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Translation struct {
	Package                 string
	Imports                 string
	Channels                string
	Structs                 string
	Methods                 string
	Functions               string
	Main                    string
	Network                 string
	HasNetConn              bool
	NetConnType             string
	NetConnAddress          string
	NetConnPort             string
	StructSlice             []string
	ChannelDefSlice         []string
	ChannelConstructorSlice []string
}

func GetSuffixWithRecUnderscoresAddedAsString(slice []string) string {
	s := ""
	for i := 0; i < len(slice)-2; i++ {
		if slice[i] == "_rec" && !strings.Contains(slice[i+2], "_par") {
			s += slice[i] + slice[i+1] + "_" + slice[i+2]
			i += 2
		} else {
			s += slice[i]
		}
	}
	return s
}

func AddUnderscoreSeparatorBetweenIntStrings(slice []string) []string {
	newSlice := make([]string, 0)
	for i := 0; i < len(slice); i++ {
		if i < len(slice)-1 {
			cur := slice[i]
			next := slice[i+1]
			newSlice = append(newSlice, cur)
			if _, err := strconv.Atoi(cur); err == nil {
				if _, err := strconv.Atoi(next); err == nil {
					newSlice = append(newSlice, "_")
				}
			}
		} else {
			newSlice = append(newSlice, slice[len(slice)-1])
		}
	}
	return newSlice
}

func WriteMessageStruct(t *Translation, mess *MessageData, choiceName string) {
	strct := "type "
	if choiceName == "" {
		strct += mess.Protagonist + GetStringSliceAsString(mess.Suffix)
	} else {
		strct += choiceName
	}
	strct += " struct {\n"
	strct += "\tChannels *Channels\n"
	strct += "\tUsed bool\n}\n\n"
	t.StructSlice = append(t.StructSlice, strct)
	t.Structs += strct
}

func SliceContainsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func IsInitialChoice(mess *MessageData) bool {
	if len(mess.Suffix) > 3 {
		if mess.Suffix[len(mess.Suffix)-4] == "_choice" && mess.Suffix[len(mess.Suffix)-1] == "1" {
			return true
		}
	}
	return false
}

func GetInitialChoiceName(mess *MessageData) string {
	name := mess.Protagonist
	for i := 0; i < len(mess.Suffix)-2; i++ {
		name += mess.Suffix[i]
	}
	return name
}

func WriteMessageStructs(t *Translation, m *MessagesData) {
	writtenChoices := make([]string, 0)
	for _, mess := range m.Messages {
		duplicateIntialChoice := false
		choiceName := ""
		if IsInitialChoice(mess) {
			choiceName = GetInitialChoiceName(mess)
			if SliceContainsString(writtenChoices, choiceName) {
				duplicateIntialChoice = true
			} else {
				writtenChoices = append(writtenChoices, choiceName)
			}
		}
		if !duplicateIntialChoice {
			WriteMessageStruct(t, mess, choiceName)
		}
	}
}

func WriteParStartStruct(t *Translation, par *ParallelData) {
	strct := "type " + par.Protagonist
	strct += par.ChanNameBase
	strct += "_start struct {\n"
	strct += "\tChannels *Channels\n"
	strct += "\tUsed bool\n"
	strct += "}\n\n"
	t.StructSlice = append(t.StructSlice, strct)
	t.Structs += strct
}

func WriteParEndStruct(t *Translation, par *ParallelData) {
	strct := "type " + par.Protagonist
	strct += par.ChanNameBase
	strct += "_end struct {\n"
	strct += "\tChannels *Channels\n"
	strct += "\tUsed bool\n"
	strct += "}\n\n"
	t.StructSlice = append(t.StructSlice, strct)
	t.Structs += strct
}

func WriteParallelStructs(t *Translation, m *MessagesData) {
	for _, par := range m.Parallels {
		WriteParStartStruct(t, par)
		WriteParEndStruct(t, par)
	}
}

func GetMultiValueTypeStruct(t *Translation, mess *MessageData) string {
	strct := "type "
	for i, param := range mess.Parameters {
		if i > 0 {
			strct += "_"
		}
		strct += param
	}
	strct += " struct {\n"
	for i, param := range mess.Parameters {
		strct += "\tparam" + strconv.Itoa(i+1) + " " + param + "\n"
	}
	strct += "}\n\n"
	return strct
}

func WriteSequentialStructs(t *Translation, m *MessagesData) {
	WriteMessageStructs(t, m)
	WriteParallelStructs(t, m)
}

func WriteDataStruct(mess *MessageData) string {
	strct := "type " + GetSendValName(mess) + " struct {\n"
	if len(mess.Parameters) == 0 {
		strct += "\tParam1  struct{}\n"
	} else {
		for i, param := range mess.Parameters {
			strct += "\tParam" + strconv.Itoa(i+1) + " "
			strct += param + "\n"
		}
	}
	strct += "}\n\n"
	return strct
}

func WriteDataStructs(t *Translation, m *MessagesData) string {
	dataStructs := ""
	for _, mess := range m.Messages {
		newDataStruct := WriteDataStruct(mess)
		dataStructs += newDataStruct
		t.StructSlice = append(t.StructSlice, newDataStruct)
	}
	return dataStructs
}

func WriteStructCommentTitle(t *Translation, m *MessagesData) {
	title := "//Structs\n\n"
	t.Structs += title
}

func WriteStructs(t *Translation, m *MessagesData, wg *sync.WaitGroup) {
	defer wg.Done()
	t.StructSlice = make([]string, 0)
	WriteStructCommentTitle(t, m)
	dataStructs := WriteDataStructs(t, m)
	t.Structs += dataStructs
	WriteSequentialStructs(t, m)
}

func WriteChannels(t *Translation, m *MessagesData, wg *sync.WaitGroup) {
	defer wg.Done()
	WriteChannelStruct(t, m)
	WriteChannelConstructor(t, m)
}

func WriteChannelStruct(t *Translation, m *MessagesData) {
	t.ChannelDefSlice = make([]string, 0)
	t.Channels = "//Channels\n\n"
	t.Channels += "type Channels struct {\n"
	WriteParameterChannels(t, m)
	WriteParallelChannels(t, m)
	WriteChoiceChannels(t, m)
	WriteCloseNetworkConnectionChannels(t, m)
	t.Channels += "}\n\n"
}

func GetParallelChannelName(par *ParallelData, index int) string {
	name := "done"
	name += par.ChanNameBase + GetNextLetter(index)
	return name
}

func GetParallelChannelType() string {
	return "bool"
}

func WriteParallelChannel(t *Translation, par *ParallelData) {
	for i, _ := range par.OptionSuffixes {
		line := "\t"
		line += GetParallelChannelName(par, i)
		line += " chan " + GetParallelChannelType() + "\n"
		t.ChannelDefSlice = append(t.ChannelDefSlice, line)
		t.Channels += line
	}
}

func WriteParallelChannels(t *Translation, m *MessagesData) {
	for _, par := range m.Parallels {
		WriteParallelChannel(t, par)
	}
}

func GetChannelName(mess *MessageData) string {
	chanName := mess.MethodNameBase
	chanName += "From" + mess.FromBase
	chanName += "To" + mess.ToBase
	if len(mess.Parameters) == 0 {
		chanName += "_Empty"
	} else {
		for _, param := range mess.Parameters {
			chanName += "_" + param
		}
	}
	return chanName
}

func GetChoiceChannelName(mess *MessageData) string {
	chanName := mess.MethodNameBase
	chanName += "From" + mess.ToBase
	chanName += "To" + mess.ToBase
	if len(mess.Parameters) == 0 {
		chanName += "_Empty"
	} else {
		for _, param := range mess.Parameters {
			chanName += "_" + param
		}
	}
	return chanName
}

func GetChannelType(mess *MessageData) string {
	t := ""
	if len(mess.Parameters) == 0 {
		t += "Empty"
	} else {
		for i, param := range mess.Parameters {
			if i > 0 {
				t += "_"
			}
			t += param
		}
	}
	return t
}

func WriteChoiceChannels(t *Translation, m *MessagesData) {
	for _, mess := range m.Messages {
		if mess.IsChoiceOption && mess.Protagonist == mess.ToBase {
			line := "\t"
			line += GetChoiceChannelName(mess)
			line += " chan "
			line += GetSendValName(mess) + "\n"
			t.ChannelDefSlice = append(t.ChannelDefSlice, line)
			t.Channels += line
		}
	}
}

func WriteParameterChannels(t *Translation, m *MessagesData) {
	for _, mess := range m.Messages {
		line := "\t"
		line += GetChannelName(mess)
		line += " chan "
		line += GetSendValName(mess) + "\n"
		t.ChannelDefSlice = append(t.ChannelDefSlice, line)
		t.Channels += line
	}
}

func WriteConstructorParameterLines(t *Translation, m *MessagesData) {
	for _, mess := range m.Messages {
		line := "\tc."
		line += GetChannelName(mess)
		line += " = make(chan "
		line += GetSendValName(mess)
		line += ")\n"
		t.ChannelConstructorSlice = append(t.ChannelConstructorSlice, line)
		t.Channels += line
	}
}

func GetNonProtagonistRoles(m *MessagesData) []string {
	otherRoles := make([]string, 0)
	rolesDict := make(map[string]bool)
	for _, mess := range m.Messages {
		other := ""
		if mess.ToBase != mess.Protagonist {
			other = mess.ToBase
		} else if mess.FromBase != mess.Protagonist {
			other = mess.FromBase
		}
		if _, ok := rolesDict[other]; ok {
			// Key already added, do nothing
		} else {
			otherRoles = append(otherRoles, other)
			rolesDict[other] = true
		}
	}
	return otherRoles
}

func WriteCloseNetworkConnectionChannels(t *Translation, m *MessagesData) {
	otherRoles := GetNonProtagonistRoles(m)
	for _, role := range otherRoles {
		line := "\tdoneCommunicatingWith" + role + " chan bool\n"
		t.ChannelDefSlice = append(t.ChannelDefSlice, line)
		t.Channels += line
	}
}

func WriteConstructorParallelLine(t *Translation, par *ParallelData) {
	for i, _ := range par.OptionSuffixes {
		line := "\tc."
		line += GetParallelChannelName(par, i)
		line += " = make(chan "
		line += GetParallelChannelType()
		line += ", 1)\n"
		t.ChannelConstructorSlice = append(t.ChannelConstructorSlice, line)
		t.Channels += line
	}
}

func WriteConstructorParallelLines(t *Translation, m *MessagesData) {
	for _, par := range m.Parallels {
		WriteConstructorParallelLine(t, par)
	}
}

func WriteConstructorChoiceLines(t *Translation, m *MessagesData) {
	for _, mess := range m.Messages {
		if mess.IsChoiceOption && mess.Protagonist == mess.ToBase {
			line := "\tc."
			line += GetChoiceChannelName(mess)
			line += " = make(chan "
			line += GetSendValName(mess)
			line += ", 1"
			line += ")\n"
			t.ChannelConstructorSlice = append(t.ChannelConstructorSlice, line)
			t.Channels += line
		}
	}
}

func GetImportsForString(s string) string {
	noImports := true
	importString := "import (\n"
	potentialImports := []string{"errors", "sync", "log", "fmt", "net", "gob", "time"}
	for _, potimp := range potentialImports {
		if strings.Contains(s, potimp+".") {
			importString += "\t\""
			if potimp == "gob" {
				importString += "encoding/"
			}
			importString += potimp + "\"\n"
			noImports = false
		}
	}
	importString += ")\n\n"
	if noImports {
		return ""
	}
	return importString
}

func SetupNetworkConnection(t *Translation) {
	t.HasNetConn = true
	t.NetConnType = "tcp"
	t.NetConnAddress = "localhost"
	t.NetConnPort = ":8080"
}

func PrintTranslation(t *Translation) {
	fmt.Println(t.Network)
	fmt.Println(t.Package)
	fmt.Println(t.Imports)
	fmt.Println(t.Channels)
	fmt.Println(t.Structs)
	fmt.Println(t.Methods)
	fmt.Println(t.Functions)
}

func WriteToFile(t *Translation, fileName string) {
	file, err := os.Create("./output/" + fileName + ".go")
	if err != nil {
		log.Fatal("Cannot create file: ", err)
	}
	defer file.Close()
	fmt.Fprint(file, t.Package)
	fmt.Fprint(file, t.Imports)
	fmt.Fprint(file, t.Network)
	fmt.Fprint(file, t.Channels)
	fmt.Fprint(file, t.Structs)
	fmt.Fprint(file, t.Methods)
	fmt.Fprint(file, t.Functions)
}

func WriteAPIToFile(t *Translation, moduleName string, wg *sync.WaitGroup) {
	defer wg.Done()
	sep := string(os.PathSeparator)
	path := "." + sep + "output" + sep + moduleName + "_Gobble" + sep
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatal("Error creating directory: ", err)
	}
	file, err := os.Create(path + "api.go")
	if err != nil {
		log.Fatal("Cannot create file: ", err)
	}
	defer file.Close()
	output := t.Channels + t.Structs + t.Methods
	imports := GetImportsForString(output)
	output = t.Package + imports + output
	fmt.Fprint(file, output)
}

func WriteNetworkToFile(t *Translation, moduleName string, wg *sync.WaitGroup) {
	defer wg.Done()
	sep := string(os.PathSeparator)
	path := "." + sep + "output" + sep + moduleName + "_Gobble" + sep
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatal("Error creating directory: ", err)
	}
	fileName := path + "network.go"
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal("Cannot create file: ", err)
	}
	defer file.Close()
	output := t.Network
	imports := GetImportsForString(output)
	output = t.Package + imports + output
	fmt.Fprint(file, output)
}

func WriteFunctionsToFile(t *Translation, moduleName string, wg *sync.WaitGroup) {
	defer wg.Done()
	sep := string(os.PathSeparator)
	path := "." + sep + "output" + sep + moduleName + "_Gobble" + sep
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatal("Error creating directory: ", err)
	}
	fileName := path + "main.go"
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal("Cannot create file: ", err)
	}
	defer file.Close()
	output := t.Functions + t.Main
	imports := GetImportsForString(output)
	output = t.Package + imports + output
	fmt.Fprint(file, output)
}

func WriteCombinedTranslationToFileInstance(pkg string, content string, path string, name string) {
	fileName := path + name
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal("Cannot create file: ", err)
	}
	defer file.Close()
	output := content
	imports := GetImportsForString(output)
	output = pkg + imports + output
	fmt.Fprint(file, output)
}

func WriteCombinedTranslationToFile(t *Translation, moduleName string) {
	sep := string(os.PathSeparator)
	path := "." + sep + "output" + sep + moduleName + "_Gobble" + sep + "combined" + sep
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatal("Error creating directory: ", err)
	}
	pkg := "package main\n\n"
	WriteCombinedTranslationToFileInstance(pkg, t.Functions+t.Main, path, "main.go")
	WriteCombinedTranslationToFileInstance(pkg, t.Structs, path, "structs.go")
	WriteCombinedTranslationToFileInstance(pkg, t.Channels, path, "channels.go")
	WriteCombinedTranslationToFileInstance(pkg, t.Methods, path, "methods.go")
}

func GetCombinedChannels(t *Translation) string {
	channels := "// Channels\n\n"
	channels += "type Channels struct {\n"
	channelDefMap := make(map[string]bool)
	for _, line := range t.ChannelDefSlice {
		if _, ok := channelDefMap[line]; ok {
			// Do nothing
		} else {
			channels += line
			channelDefMap[line] = true
		}
	}
	channels += "}\n\n"
	channels += "func NewChannels() *Channels {\n"
	channels += "\tc := new(Channels)\n"
	channelConstMap := make(map[string]bool)
	for _, line := range t.ChannelConstructorSlice {
		if _, ok := channelConstMap[line]; ok {
			// Do nothing
		} else {
			channels += line
			channelConstMap[line] = true
		}
	}
	channels += "\treturn c\n"
	channels += "}\n\n"
	return channels
}

func GetCombinedStructs(t *Translation) string {
	structs := "// Structs\n\n"
	structMap := make(map[string]bool)
	for _, strct := range t.StructSlice {
		if _, ok := structMap[strct]; ok {
			// Do nothing
		} else {
			structs += strct
			structMap[strct] = true
		}
	}
	return structs
}

func GenerateCombinedTranslation(t *Translation) {
	channels := GetCombinedChannels(t)
	t.Channels = channels
	structs := GetCombinedStructs(t)
	t.Structs = structs
}

func GetCombinedTranslationMain(m []*MessagesData, t *Translation) string {
	roles := make([]string, 0)
	for _, md := range m {
		newRole := md.Messages[0].Protagonist
		roles = append(roles, newRole)
	}
	main := "func main() {\n"
	main += "\tchans := NewChannels()\n"
	for _, role := range roles {
		main += "\tstartStruct" + role + " := Start" + role + "(chans)\n"
	}
	main += "\tvar newWg sync.WaitGroup\n"
	for _, role := range roles {
		main += "\tnewWg.Add(1)\n"
		main += "\tgo run" + role + "(&newWg, startStruct" + role + ")\n"
	}
	main += "\tnewWg.Wait()\n"
	main += "}\n"
	return main
}

func WriteCombinedTranslation(mds []*MessagesData, fileName string) {
	t := &Translation{}
	main := GetCombinedTranslationMain(mds, t)
	channelDefs := make([]string, 0)
	channelConstructors := make([]string, 0)
	structs := make([]string, 0)
	for _, m := range mds {
		var wg sync.WaitGroup
		SetupNetworkConnection(t)
		wg.Add(1)
		WriteChannels(t, m, &wg)
		channelDefs = append(channelDefs, t.ChannelDefSlice...)
		channelConstructors = append(channelConstructors, t.ChannelConstructorSlice...)
		wg.Add(1)
		WriteStructs(t, m, &wg)
		structs = append(structs, t.StructSlice...)
		wg.Add(1)
		WriteMethods(t, m, &wg)
		vals := AssignConversationsFunctionValues(t, m)
		t.Functions += GetStringConversationsFuncVals(vals, m)
		wg.Add(1)
		WriteFunctions(t, m, &wg)
		wg.Wait()
	}
	t.ChannelDefSlice = channelDefs
	sort.Strings(t.ChannelDefSlice)
	t.ChannelConstructorSlice = channelConstructors
	sort.Strings(t.ChannelConstructorSlice)
	t.StructSlice = structs
	sort.Strings(t.StructSlice)
	t.Main = main
	GenerateCombinedTranslation(t)
	WriteCombinedTranslationToFile(t, fileName)
}

func WriteTranslation(m *MessagesData, fileName string) {
	t := &Translation{}
	SetupNetworkConnection(t)
	WritePackage(t, m)
	// Generate output strings concurrently
	var wg sync.WaitGroup
	wg.Add(1)
	WriteChannels(t, m, &wg)
	wg.Add(1)
	WriteStructs(t, m, &wg)
	wg.Add(1)
	WriteMethods(t, m, &wg)
	WriteMainFunction(t, m)
	wg.Add(1)
	network := WriteNetwork(t, m, &wg)
	t.Network += network
	vals := AssignConversationsFunctionValues(t, m)
	t.Functions += GetStringConversationsFuncVals(vals, m)
	wg.Add(1)
	WriteFunctions(t, m, &wg)
	wg.Wait()
	WriteImports(t, m)
	// Write to file concurrently
	var wg2 sync.WaitGroup
	wg2.Add(1)
	WriteAPIToFile(t, fileName, &wg2)
	wg2.Add(1)
	WriteNetworkToFile(t, fileName, &wg2)
	wg2.Add(1)
	WriteFunctionsToFile(t, fileName, &wg2)
	wg2.Wait()
}
