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
	"log"
	"strconv"
	"strings"
	"sync"
)

func WriteSendMethodHeader(t *Translation, mess *MessageData, m *MessagesData) {
	header := "func (self *"
	if IsInitialChoice(mess) {
		header += GetInitialChoiceName(mess)
	} else {
		header += mess.Protagonist + GetStringSliceAsString(mess.Suffix)
	}
	header += ") Send"
	header += "_" + mess.MethodNameBase
	for _, param := range mess.Parameters {
		header += "_" + param
	}
	header += "("
	for i, param := range mess.Parameters {
		if i > 0 {
			header += ", "
		}
		header += "param" + strconv.Itoa(i+1) + " " + param
	}
	header += ") ("
	if len(mess.SubsequentSuffix) > 0 {
		if mess.SubsequentSuffix[len(mess.SubsequentSuffix)-1] != "_end" {
			header += "*" + mess.Protagonist + GetStringSliceAsString(AddUnderscoreSeparatorBetweenIntStrings(mess.SubsequentSuffix)) + ", "
		}
	}
	whereTo := mess.WhereToIfBranchEnds
	if whereTo != (WhereToIfBranchEnds{}) {
		retVal := GetWhereToIfBranchEndsStructName(whereTo, m)
		header += "*" + retVal + ", "
	}
	if mess.ContinueToStruct != "" {
		name := mess.ContinueToStruct
		header += "*" + name + ", "
	}
	header += "error) {\n"
	t.Methods += header
}

func WriteReceiveMethodHeader(t *Translation, mess *MessageData, m *MessagesData) {
	header := "func (self *"
	if IsInitialChoice(mess) {
		header += GetInitialChoiceName(mess)
	} else {
		header += mess.Protagonist + GetStringSliceAsString(mess.Suffix)
	}
	header += ") Receive"
	header += "_" + mess.MethodNameBase
	for _, param := range mess.Parameters {
		header += "_" + param
	}
	header += "() ("
	if len(mess.SubsequentSuffix) > 0 {
		if mess.SubsequentSuffix[len(mess.SubsequentSuffix)-1] != "_end" {
			header += "*" + mess.Protagonist + GetStringSliceAsString(mess.SubsequentSuffix) + ", "
		}
	}
	whereTo := mess.WhereToIfBranchEnds
	if whereTo != (WhereToIfBranchEnds{}) {
		retVal := GetWhereToIfBranchEndsStructName(whereTo, m)
		header += "*" + retVal + ", "
	}
	if mess.ContinueToStruct != "" {
		name := mess.ContinueToStruct
		header += "*" + name + ", "
	}
	for _, param := range mess.Parameters {
		header += param + ", "
	}
	header += "error) {\n"
	t.Methods += header
}

func GetDataStructNameWithIdAndType(id string, typ string, m *MessagesData) string {
	name := ""
	if typ == "message" {
		mess, err := GetMessageDataWithId(m, id)
		if err != nil {
			log.Fatal(err)
		}
		name = mess.Protagonist + GetStringSliceAsString(mess.Suffix)
	} else if typ == "choice" {
		choi, err := GetChoiceDataWithId(m, id)
		if err != nil {
			log.Fatal(err)
		}
		name = choi.Protagonist + CutStringAfterLetter(choi.OptionSuffixes[0], "_")
	} else if typ == "continue" {
		cont, err := GetContinueDataWithId(m, id)
		if err != nil {
			log.Fatal(err)
		}
		name = cont.FirstStructName
	} else if typ == "rec" {
		rec, err := GetRecDataWithId(m, id)
		if err != nil {
			log.Fatal(err)
		}
		name = GetDataStructNameWithIdAndType(rec.FirstStructId, rec.FirstStructType, m)
	} else if typ == "parallel" || typ == "par" {
		par, err := GetParallelDataWithId(m, id)
		if err != nil {
			log.Fatal(err)
		}
		name = par.Protagonist + GetStringSliceAsString(par.Suffix) + "_start"
	} else {
		panic("GetDataStructNameWithIdAndType() unable to find data struct with type " + typ + " and id " + id)
	}
	return name
}

func WriteSetUsedFunction(t *Translation) {
	setUsedFunction := "\tdefer func() { self.Used = true }()\n"
	t.Methods += setUsedFunction
}

func GetSendValName(mess *MessageData) string {
	name := ""
	name += strings.Title(mess.MethodNameBase)
	name += "_from_" + mess.FromBase
	name += "_to_" + mess.ToBase + "_"
	if len(mess.Parameters) == 0 {
		name += "Empty"
	}
	for i, param := range mess.Parameters {
		if i > 0 {
			name += "_"
		}
		name += param
	}
	return name
}

func GetSendValContents(mess *MessageData) string {
	contents := ""
	for i, _ := range mess.Parameters {
		if i > 0 {
			contents += ", "
		}
		contents += "Param" + strconv.Itoa(i+1) + ": "
		contents += "param" + strconv.Itoa(i+1)
	}
	return contents
}

func WriteSendMultiParamsMethodLine(t *Translation, mess *MessageData) {
	line := "\tsendVal := " + GetSendValName(mess) + "{"
	line += GetSendValContents(mess) + "}\n"
	t.Methods += line
}

func WriteSendMethodSendToChannelLine(t *Translation, mess *MessageData) {
	line := "\tself.Channels."
	line += GetChannelName(mess) + " <- "
	line += "sendVal"
	line += "\n"
	t.Methods += line
}

func WriteReceiveMethodAssignInputsToVariables(t *Translation, mess *MessageData) {
	assignments := ""
	for i, param := range mess.Parameters {
		assignments += "\tin_" + strconv.Itoa(i+1) + "_" + param
		assignments += " = in.Param" + strconv.Itoa(i+1) + "\n"
	}
	t.Methods += assignments
}

func WriteReceiveMethodInputChannel(t *Translation, mess *MessageData) {
	line := "\t"
	if len(mess.Parameters) > 0 {
		line += "in := "
	}
	if mess.ChoiceChanName == "" {
		line += "<- self.Channels." + GetChannelName(mess) + "\n"
	} else {
		line += "<- self.Channels." + mess.ChoiceChanName + "\n"
	}
	t.Methods += line
}

func WriteReceiveMethodInputChannels(t *Translation, mess *MessageData) {
	WriteReceiveMethodInputChannel(t, mess)
}

func WriteMethodRetValDefinitionLine(t *Translation, mess *MessageData, m *MessagesData) {
	if len(mess.SubsequentSuffix) > 0 {
		if mess.SubsequentSuffix[len(mess.SubsequentSuffix)-1] != "_end" {
			line := "\tretVal := "
			line += "&" + mess.Protagonist + GetStringSliceAsString(AddUnderscoreSeparatorBetweenIntStrings(mess.SubsequentSuffix))
			line += "{Channels: self.Channels}\n"
			t.Methods += line
		}
	}
	whereTo := mess.WhereToIfBranchEnds
	if whereTo != (WhereToIfBranchEnds{}) {
		retVal := GetWhereToIfBranchEndsStructName(whereTo, m)
		line := "\tretVal := &" + retVal + "{Channels: self.Channels}\n"
		t.Methods += line
	}
	if mess.ContinueToStruct != "" {
		name := mess.ContinueToStruct
		line := "\tretVal := "
		line += "&" + name + "{Channels: self.Channels}\n"
		t.Methods += line
	}
}

func WriteMethodStructAlreadyUsedCheck(t *Translation, mess *MessageData, m *MessagesData) {
	check := "\tif self.Used {\n"
	check += "\t\treturn "
	if len(mess.SubsequentSuffix) > 0 {
		if mess.SubsequentSuffix[len(mess.SubsequentSuffix)-1] != "_end" {
			check += "retVal, "
		}
	}
	whereTo := mess.WhereToIfBranchEnds
	if whereTo != (WhereToIfBranchEnds{}) {
		check += "retVal, "
	}
	if mess.ContinueToStruct != "" {
		check += "retVal, "
	}
	if mess.Protagonist == mess.ToBase {
		for i, param := range mess.Parameters {
			check += "in_" + strconv.Itoa(i+1) + "_" + param + ", "
		}
	}
	check += "errors.New(\"Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: "
	if IsInitialChoice(mess) {
		check += GetInitialChoiceName(mess)
	} else {
		check += mess.Protagonist + GetStringSliceAsString(mess.Suffix)
	}
	check += "; method: "
	if mess.Protagonist == mess.FromBase {
		check += "Send"
	} else {
		check += "Receive"
	}
	for _, param := range mess.Parameters {
		check += "_" + param
	}
	check += "()\")\n"
	check += "\t}\n"
	t.Methods += check
}

func WriteSendMethodReturnLine(t *Translation, mess *MessageData) {
	line := "\treturn "
	if len(mess.SubsequentSuffix) > 0 {
		if mess.SubsequentSuffix[len(mess.SubsequentSuffix)-1] != "_end" {
			line += "retVal, "
		}
	}
	if mess.ContinueToStruct != "" {
		line += "retVal, "
	}

	whereTo := mess.WhereToIfBranchEnds
	if whereTo != (WhereToIfBranchEnds{}) {
		line += "retVal, "
	}

	line += "nil\n"
	t.Methods += line
}

func WriteReceiveMethodReturnLine(t *Translation, mess *MessageData, m *MessagesData) {
	line := "\treturn "
	if len(mess.SubsequentSuffix) > 0 {
		if mess.SubsequentSuffix[len(mess.SubsequentSuffix)-1] != "_end" {
			line += "retVal, "
		}
	}
	if mess.ContinueToStruct != "" {
		line += "retVal, "
	}
	whereTo := mess.WhereToIfBranchEnds
	if whereTo != (WhereToIfBranchEnds{}) {
		line += "retVal, "
	}
	for i, param := range mess.Parameters {
		line += "in_" + strconv.Itoa(i+1) + "_" + param + ", "
	}
	line += "nil\n"
	t.Methods += line
}

func WriteSendParDoneSignalLine(t *Translation, mess *MessageData) {
	if len(mess.SubsequentSuffix) > 0 {
		if mess.SubsequentSuffix[len(mess.SubsequentSuffix)-1] == "_end" {
			line := "\tself.Channels.done"
			for i := 0; i < len(mess.Suffix)-1; i++ {
				line += mess.Suffix[i]
			}
			line += " <- true\n"
			t.Methods += line
		}
	}
}

func WriteSendMethodSendValDefinition(t *Translation, mess *MessageData) {
	WriteSendMultiParamsMethodLine(t, mess)
}

func WriteSendMethod(t *Translation, mess *MessageData, m *MessagesData) {
	WriteSendMethodHeader(t, mess, m)
	WriteSetUsedFunction(t)
	WriteSendMethodSendValDefinition(t, mess)
	WriteMethodRetValDefinitionLine(t, mess, m)
	WriteMethodStructAlreadyUsedCheck(t, mess, m)
	WriteSendParDoneSignalLine(t, mess)
	WriteSendMethodSendToChannelLine(t, mess)
	WriteSendMethodReturnLine(t, mess)
}

func WriteReceiveMethodInputVariableDeclarations(t *Translation, mess *MessageData) {
	line := ""
	for i, param := range mess.Parameters {
		line += "\tvar in_" + strconv.Itoa(i+1) + "_" + param + " " + param + "\n"
	}
	t.Methods += line
}

func WriteReceiveMethod(t *Translation, mess *MessageData, m *MessagesData) {
	WriteReceiveMethodHeader(t, mess, m)
	WriteSetUsedFunction(t)
	WriteReceiveMethodInputVariableDeclarations(t, mess)
	WriteMethodRetValDefinitionLine(t, mess, m)
	WriteMethodStructAlreadyUsedCheck(t, mess, m)
	WriteReceiveMethodInputChannels(t, mess)
	WriteReceiveMethodAssignInputsToVariables(t, mess)
	WriteSendParDoneSignalLine(t, mess)
	WriteReceiveMethodReturnLine(t, mess, m)
}

func WriteParStartMethodHeader(t *Translation, par *ParallelData) {
	header := "func (self *"
	header += par.Protagonist + CutStringAfterLetter(par.OptionSuffixes[0], "_") + "_start) StartPar() ("
	header += "*" + par.Protagonist + CutStringAfterLetter(par.OptionSuffixes[0], "_") + "_end"
	for _, param := range par.OptionSuffixes {
		if param[len(param)-5:len(param)] == "_par1" {
			param += "_start"
		}
		header += ", *" + par.Protagonist + param
	}
	header += ", error) {\n"
	t.Methods += header
}

func CutStringAfterLetter(myString string, letter string) string {
	newString := ""
	newLength := 0
	for i := len(myString) - 1; i > -1; i-- {
		if string(myString[i]) == letter {
			newLength = i
			for j := 0; j < newLength; j++ {
				newString += string(myString[j])
			}
			return newString
		}
	}
	return myString
}

func WriteParEndMethodHeader(t *Translation, par *ParallelData) {
	header := "func (self *"
	header += par.Protagonist + CutStringAfterLetter(par.OptionSuffixes[0], "_") + "_end) EndPar() ("
	if len(par.SubsequentSuffix) > 0 {
		header += "*" + par.Protagonist + GetStringSliceAsString(AddUnderscoreSeparatorBetweenIntStrings(par.SubsequentSuffix)) + ", "
	}
	header += "error) {\n"
	t.Methods += header
}

func WriteParStartMethodEndDefinition(t *Translation, par *ParallelData) {
	def := "\t" + par.Protagonist + CutStringAfterLetter(par.OptionSuffixes[0], "_")
	def += "_end := &" + par.Protagonist + CutStringAfterLetter(par.OptionSuffixes[0], "_")
	def += "_end{Channels: self.Channels}\n"
	t.Methods += def
}

func WriteParStartMethodParDefinitions(t *Translation, par *ParallelData) {
	for _, suffix := range par.OptionSuffixes {
		if suffix[len(suffix)-5:len(suffix)] == "_par1" {
			suffix += "_start"
		}
		def := "\t" + par.Protagonist + suffix + " := "
		def += "&" + par.Protagonist + suffix
		def += "{Channels: self.Channels}\n"
		t.Methods += def
	}
}

func WriteParStartMethodAlreadyUsedCheck(t *Translation, par *ParallelData) {
	check := "\tif self.Used {\n"
	check += "\t\treturn "
	check += par.Protagonist + CutStringAfterLetter(par.OptionSuffixes[0], "_") + "_end"
	for _, param := range par.OptionSuffixes {
		if param[len(param)-5:len(param)] == "_par1" {
			param += "_start"
		}
		check += ", " + par.Protagonist + param
	}
	check += ", errors.New(\"Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: "
	check += par.Protagonist + GetStringSliceAsString(par.Suffix) + "_start"
	check += "; method: StartPar()\")\n\t}\n"
	t.Methods += check
}

func WriteParEndMethodAlreadyUsedCheck(t *Translation, par *ParallelData) {
	check := "\tif self.Used {\n"
	check += "\t\treturn "
	if len(par.SubsequentSuffix) > 0 {
		check += "retVal, "
	}
	check += "errors.New(\"Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: "
	check += par.Protagonist + GetStringSliceAsString(par.Suffix) + "_start"
	check += "; method: EndPar()\")\n\t}\n"
	t.Methods += check
}

func WriteParStartMethodReturnLine(t *Translation, par *ParallelData) {
	line := "\treturn "
	line += par.Protagonist + CutStringAfterLetter(par.OptionSuffixes[0], "_") + "_end"
	for _, param := range par.OptionSuffixes {
		if param[len(param)-5:len(param)] == "_par1" {
			param += "_start"
		}
		line += ", " + par.Protagonist + param
	}
	line += ", nil\n"
	t.Methods += line
}

func WriteParEndMethodRetValDefinition(t *Translation, par *ParallelData) {
	if len(par.SubsequentSuffix) > 0 {
		def := "\tretVal := "
		def += "&" + par.Protagonist + GetStringSliceAsString(AddUnderscoreSeparatorBetweenIntStrings(par.SubsequentSuffix))
		def += "{Channels: self.Channels}\n"
		t.Methods += def
	}
}

func WriteParEndDoneSignalReceivers(t *Translation, par *ParallelData) {
	for i := 0; i < len(par.OptionSuffixes); i++ {
		receiver := "\tpar" + GetNextLetter(i) + "_done := false\n"
		receiver += "\tselect{\n"
		receiver += "\tcase <- self.Channels.done" + CutStringAfterLetter(par.OptionSuffixes[0], "_") + GetNextLetter(i) + ":\n"
		receiver += "\t\tpar" + GetNextLetter(i) + "_done = true\n"
		receiver += "\tdefault:\n"
		receiver += "\t}\n"
		t.Methods += receiver
	}
}

func WriteParEndDoneCheckers(t *Translation, par *ParallelData) {
	for i := 0; i < len(par.OptionSuffixes); i++ {
		checker := "\tif !par" + GetNextLetter(i) + "_done {\n"
		checker += "\t\treturn "
		if len(par.SubsequentSuffix) > 0 {
			checker += "retVal, "
		}
		checker += "errors.New(\"Dynamic session type checking error: attempted call to Seller" + GetStringSliceAsString(par.Suffix) + "_end.EndPar() prior to completion of parallel process Seller" + GetStringSliceAsString(par.Suffix) + GetNextLetter(i) + "\")\n"
		checker += "\t}\n"
		t.Methods += checker
	}
}

func WriteParEndReturnLine(t *Translation, par *ParallelData) {
	line := "\treturn "
	if len(par.SubsequentSuffix) > 0 {
		line += "retVal, "
	}
	line += "nil\n"
	t.Methods += line
}

func WriteParStartMethod(t *Translation, par *ParallelData) {
	WriteParStartMethodHeader(t, par)
	WriteSetUsedFunction(t)
	WriteParStartMethodEndDefinition(t, par)
	WriteParStartMethodParDefinitions(t, par)
	WriteParStartMethodAlreadyUsedCheck(t, par)
	WriteParStartMethodReturnLine(t, par)
}

func WriteParEndMethod(t *Translation, par *ParallelData) {
	WriteParEndMethodHeader(t, par)
	WriteSetUsedFunction(t)
	WriteParEndMethodRetValDefinition(t, par)
	WriteParEndMethodAlreadyUsedCheck(t, par)
	WriteParEndDoneSignalReceivers(t, par)
	WriteParEndDoneCheckers(t, par)
	WriteParEndReturnLine(t, par)
}

func WriteMethods(t *Translation, m *MessagesData, wg *sync.WaitGroup) {
	defer wg.Done()
	t.Methods += "//Methods\n\n"
	for _, mess := range m.Messages {
		if mess.Protagonist == mess.ToBase {
			WriteReceiveMethod(t, mess, m)
		}
		if mess.Protagonist == mess.FromBase {
			WriteSendMethod(t, mess, m)
		}
		t.Methods += "}\n\n"
	}
	for _, par := range m.Parallels {
		WriteParStartMethod(t, par)
		t.Methods += "}\n\n"
		WriteParEndMethod(t, par)
		t.Methods += "}\n\n"
	}
}
