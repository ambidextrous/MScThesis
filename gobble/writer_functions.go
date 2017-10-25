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
	"errors"
	"log"
	"strconv"
	"strings"
	"sync"
)

func WriteMainFunctionHeader(t *Translation, m *MessagesData) {
	header := "func main() {\n"
	t.Main += header
}

func WriteMainFunctionChannelDef(t *Translation) {
	def := "\tchans := NewChannels()\n"
	t.Main += def
}

func WriteMainWaitGroupDef(t *Translation) {
	def := "\tvar newWg sync.WaitGroup\n"
	t.Main += def
}

func WriteMainFunctionLaunchParty(t *Translation, m *MessagesData) {
	launch := "\tnewWg.Add(1)\n"
	launch += "\tgo run"
	launch += m.Conversations[0].Protagonist
	launch += "(&newWg, startStruct)\n"
	t.Main += launch
}

func WriteMainWaitLine(t *Translation) {
	line := "\tnewWg.Wait()\n"
	t.Main += line
}

func WriteMainFunctionGetStartStruct(t *Translation, m *MessagesData) {
	firstMess := m.Messages[0]
	protagonist := firstMess.Protagonist
	line := "\tstartStruct := "
	line += "Start" + protagonist + "(chans)\n"
	t.Main += line
}

func WriteMainFunctionSetNetworkVariables(t *Translation) {
	variables := ""
	variables += "\tconnType := \"" + t.NetConnType + "\"\n"
	variables += "\taddress := \"" + t.NetConnAddress + "\"\n"
	variables += "\tport := \"" + t.NetConnPort + "\"\n"
	t.Main += variables
}

func WriteMainFunctionSetupNetworkConnection(t *Translation) {
	line := "\tSetupNetworkConnections(chans, connType, address, port)\n"
	t.Main += line
}

func WriteMainFunctionCloseNetworkConnections(t *Translation) {
	line := "\tCloseNetworkConnections(chans)\n"
	t.Main += line
}

func WriteMainFunction(t *Translation, m *MessagesData) {
	WriteMainFunctionHeader(t, m)
	WriteMainFunctionChannelDef(t)
	if t.HasNetConn {
		WriteMainFunctionSetNetworkVariables(t)
		WriteMainFunctionSetupNetworkConnection(t)
	}
	WriteMainFunctionGetStartStruct(t, m)
	WriteMainWaitGroupDef(t)
	WriteMainFunctionLaunchParty(t, m)
	WriteMainWaitLine(t)
	if t.HasNetConn {
		WriteMainFunctionCloseNetworkConnections(t)
	}
	t.Main += "}\n\n"
}

func WriteStartFunctionHeader(t *Translation, m *MessagesData) {
	header := "func Start"
	protagonist := GetProtagonist(m)
	header += protagonist
	header += "(chans *Channels) (*"
	name, typ := GetFirstStep(m)
	if typ == "parallel" {
		name += "_start"
	}
	header += name
	header += ") {"
	header += "\n"
	t.Functions += header
}

func GetFirstStep(m *MessagesData) (string, string) {
	firstStep := ""
	firstId := ""
	var messageData *MessageData
	var parallelData *ParallelData
	var choiceData *ChoiceData
	var recData *RecData
	var err error
	for _, conv := range m.Conversations {
		if IsFirstConv(conv) {
			firstId = conv.ConversationElementsData[0].UniqueId
			messageData, err = GetMessageDataWithId(m, firstId)
			if err == nil {
				firstStep = messageData.Protagonist + GetStringSliceAsString(AddUnderscoreSeparatorBetweenIntStrings(messageData.Suffix))
				return firstStep, "message"
			}
			parallelData, err = GetParallelDataWithId(m, firstId)
			if err == nil {
				firstStep = parallelData.Protagonist + GetStringSliceAsString(AddUnderscoreSeparatorBetweenIntStrings(parallelData.Suffix))
				return firstStep, "parallel"
			}
			choiceData, err = GetChoiceDataWithId(m, firstId)
			if err == nil {
				firstStep = choiceData.Protagonist + GetStringSliceAsString(AddUnderscoreSeparatorBetweenIntStrings(choiceData.Suffix))
				return firstStep, "choice"
			}
			recData, err = GetRecDataWithId(m, firstId)
			if err == nil {
				firstStep = GetDataStructNameWithIdAndType(recData.FirstStructId, recData.FirstStructType, m)
				return firstStep, "rec"
			}
		}
	}
	newErr := errors.New("GetFirstStep() Could not find MessageData, ParallelData, ChoiceData or RecData with id " + firstId)
	log.Fatal(newErr)
	return firstStep, ""
}

func GetProtagonist(m *MessagesData) string {
	protagonist := m.Messages[0].Protagonist
	return protagonist
}

func WriteStartFunctionDef(t *Translation, m *MessagesData) {
	def := "\t"
	name, typ := GetFirstStep(m)
	if typ == "parallel" {
		name += "_start"
	}
	def += "start := &" + name
	def += "{Channels: chans}\n"
	t.Functions += def
}

func WriteStartFunctionReturn(t *Translation) {
	line := "\treturn start\n"
	t.Functions += line
}

func WriteStartFunction(t *Translation, m *MessagesData) {
	WriteStartFunctionHeader(t, m)
	WriteStartFunctionDef(t, m)
	WriteStartFunctionReturn(t)
	t.Functions += "}\n\n"
}

func GetWhereToIfBranchEndsStructName(w WhereToIfBranchEnds, m *MessagesData) string {
	name := ""
	if w.Type == "message" {
		mess, err := GetMessageDataWithId(m, w.Id)
		if err != nil {
			log.Fatal(err)
		}
		name = mess.Protagonist + GetStringSliceAsString(mess.Suffix)
		return name
	} else if w.Type == "par" {
		par, err := GetParallelDataWithId(m, w.Id)
		if err != nil {
			log.Fatal(err)
		}
		name = par.Protagonist + GetStringSliceAsString(par.Suffix) + "_start"
		return name
	} else if w.Type == "choice" {
		choi, err := GetChoiceDataWithId(m, w.Id)
		if err != nil {
			log.Fatal(err)
		}
		name = choi.Protagonist + CutStringAfterLetter(choi.OptionSuffixes[0], "_")
		return name
	} else if w.Type == "rec" {
		rec, err := GetRecDataWithId(m, w.Id)
		if err != nil {
			log.Fatal(err)
		}
		name = GetDataStructNameWithIdAndType(rec.FirstStructId, rec.FirstStructType, m)
	} else {
		log.Fatal("Unknown type encountered in GetWhereToIfBranchEndsStructName(): ", w.Type)
	}
	return name
}

func WriteFunctionErrCheck(t *Translation) {
	t.Functions += GetFunctionErrCheck()
}

func GetFunctionErrCheck() string {
	check := "\tif err != nil {\n"
	check += "\t\tlog.Fatal(err)\n"
	check += "\t}\n"
	return check
}

func IsFirstConv(conv *ConversationData) bool {
	if len(conv.Suffix) == 0 {
		return true
	}
	return false
}

func WriteParallelFunctionHeader(t *Translation, parallelData *ParallelData, m *MessagesData) {
	header := "func run"
	header += parallelData.Protagonist + CutStringAfterLetter(parallelData.OptionSuffixes[0], "_")
	header += "("
	self := parallelData.Protagonist + CutStringAfterLetter(parallelData.OptionSuffixes[0], "_") + "_start"
	header += self + " *" + self
	header += ") ("
	header += "interface{}, "
	header += "error"
	header += ") {\n"
	t.Functions += header
}

func WriteStartParallelFunction(t *Translation, parallelData *ParallelData, m *MessagesData) {
	start := "\t"
	nameBase := parallelData.Protagonist + CutStringAfterLetter(parallelData.OptionSuffixes[0], "_")
	start += nameBase + "_end"
	for _, par := range parallelData.OptionSuffixes {
		start += ", " + parallelData.Protagonist + par
	}
	start += ", err := " + nameBase + "_start.StartPar()\n"
	t.Functions += start
}

func WriteDeclareNewWaitGroup(t *Translation) {
	line := "\tvar newWg sync.WaitGroup\n"
	t.Functions += line
}

func WriteParGoRoutine(t *Translation, nameBase string, index int, paramVal string) {
	routine := "\tnewWg.Add(1)\n"
	routine += "\tgo run" + nameBase + GetNextLetter(index)
	routine += "(" + paramVal
	routine += ", &newWg)\n"
	t.Functions += routine
}

func WriteParGoRoutines(t *Translation, parallelData *ParallelData, m *MessagesData) {
	nameBase := ""
	for index, _ := range parallelData.OptionSuffixes {
		nameBase = parallelData.Protagonist + CutStringAfterLetter(parallelData.OptionSuffixes[0], "_")
		paramVal := parallelData.Protagonist + parallelData.OptionSuffixes[index]
		WriteParGoRoutine(t, nameBase, index, paramVal)
	}
}

func WriteEndParallelFunction(t *Translation, parallelData *ParallelData, m *MessagesData) {
	end := "\t"
	nameBase := parallelData.Protagonist + CutStringAfterLetter(parallelData.OptionSuffixes[0], "_")
	if len(parallelData.SubsequentSuffix) > 0 {
		end += "\t" + parallelData.Protagonist + GetStringSliceAsString(parallelData.SubsequentSuffix) + ", "
	}
	end += "err_end := "
	end += nameBase + "_end.EndPar()\n"
	t.Functions += end
}

func WriteParallelFunctionReturn(t *Translation, parallelData *ParallelData, m *MessagesData) {
	ret := "\treturn "
	originalReturn := ret
	if len(parallelData.SubsequentSuffix) > 0 {
		ret += parallelData.Protagonist + GetStringSliceAsString(parallelData.SubsequentSuffix) + ", "
	}
	if ret == originalReturn {
		ret += "\"\", "
	}
	ret += "nil"
	ret += "\n"
	t.Functions += ret
}

func WriteParallelFunctionEndErrCheck(t *Translation, parallelData *ParallelData) {
	check := "\tif err_end != nil {\n"
	check += "\t\tlog.Fatal(err_end)\n"
	check += "\t}\n"
	t.Functions += check
}

func WriteFunctionWaitLine(t *Translation) {
	t.Functions += "\tnewWg.Wait()\n"
}

func WriteParallelFunction(t *Translation, parallelData *ParallelData, m *MessagesData) {
	WriteParallelFunctionHeader(t, parallelData, m)
	WriteStartParallelFunction(t, parallelData, m)
	WriteFunctionErrCheck(t)
	WriteDeclareNewWaitGroup(t)
	WriteParGoRoutines(t, parallelData, m)
	WriteFunctionWaitLine(t)
	WriteEndParallelFunction(t, parallelData, m)
	WriteParallelFunctionEndErrCheck(t, parallelData)
	WriteParallelFunctionReturn(t, parallelData, m)
	t.Functions += "}\n\n"
}

func WriteParallelFunctions(t *Translation, m *MessagesData) {
	for _, parallelData := range m.Parallels {
		WriteParallelFunction(t, parallelData, m)
	}
}

func WriteChoiceFunctionHeader(t *Translation, choiceData *ChoiceData, m *MessagesData) {
	header := "func Make_" + choiceData.Protagonist + CutStringAfterLetter(choiceData.OptionSuffixes[0], "_") + "_Choices("
	header += choiceData.Protagonist + CutStringAfterLetter(choiceData.OptionSuffixes[0], "_") + " *" + choiceData.Protagonist + CutStringAfterLetter(choiceData.OptionSuffixes[0], "_")
	header += ") ("
	header += "interface{}, "
	header += "error) {\n"
	t.Functions += header
}

func WriteChooserFunctionBranch(t *Translation, choiceData *ChoiceData, m *MessagesData, branch string, index int) {
	line := "\t"
	if index != 0 {
		line += "//"
	}
	self := choiceData.Protagonist + CutStringAfterLetter(choiceData.OptionSuffixes[0], "_")
	line += "retVal, err := "
	line += "run" + choiceData.Protagonist
	line += RemoveLastRuneFromString(branch) + "("
	line += self
	line += ")"
	if index != 0 {
		line += " // Uncomment to use"
	}
	line += "\n"
	t.Functions += line
}

func RemoveLastRuneFromString(s string) string {
	length := len(s)
	if length > 0 {
		s = s[:length-1]
	}
	return s
}

func WriteChooserFunctionBranches(t *Translation, choiceData *ChoiceData, m *MessagesData) {
	for index, branch := range choiceData.OptionSuffixes {
		WriteChooserFunctionBranch(t, choiceData, m, branch, index)
	}
}

func RemoveLastStringFromSlice(oldSlice []string) []string {
	newSlice := make([]string, 0)
	for i := 0; i < len(oldSlice)-1; i++ {
		newSlice = append(newSlice, oldSlice[i])
	}
	return newSlice
}

func WriteChooserFunctionReturnLine(t *Translation) {
	line := "\treturn retVal, err\n"
	t.Functions += line
}

func WriteChooserFunction(t *Translation, choiceData *ChoiceData, m *MessagesData) {
	WriteChoiceFunctionHeader(t, choiceData, m)
	WriteChooserFunctionBranches(t, choiceData, m)
	WriteChooserFunctionReturnLine(t)
}

func WriteChosenFunctionBranch(t *Translation, choiceData *ChoiceData, m *MessagesData, branch string, index int) {
	messageData, err := GetMessageDataWithId(m, choiceData.OptionIds[index])
	if err != nil {
		log.Fatal(err)
	}
	chanName := messageData.ChanName
	self := choiceData.Protagonist + CutStringAfterLetter(choiceData.OptionSuffixes[0], "_")
	bran := "\tcase received" + GetNextLetter(index) + " := <- "
	bran += self + ".Channels." + chanName + ":\n"
	if messageData.ChoiceChanName != "" {
		chanName = messageData.ChoiceChanName
	}
	bran += "\t\t" + self + ".Channels." + chanName + " <- received" + GetNextLetter(index) + "\n"
	bran += "\t\t"
	bran += "retVal, "
	bran += "err = "
	bran += "run" + self + GetNextLetter(index) + "("
	bran += self
	bran += ")\n"
	bran += "\t\tif err != nil {\n"
	bran += "\t\t\tlog.Fatal(err)\n"
	bran += "\t\t}\n"
	t.Functions += bran
}

func WriteChoiceFunctionReturn(t *Translation, choiceData *ChoiceData, m *MessagesData) {
	t.Functions += "\treturn "
	t.Functions += "retVal, "
	t.Functions += "nil\n"
}

func WriteChoiceFunctionRetValDef(t *Translation, choiceData *ChoiceData, m *MessagesData) {
	defs := ""
	defs += "\tvar retVal interface{}\n"
	errDeclaration := "\tvar err error\n"
	t.Functions += defs + errDeclaration
}

func WriteChosenFunctionBranches(t *Translation, choiceData *ChoiceData, m *MessagesData) {
	WriteChoiceFunctionRetValDef(t, choiceData, m)
	t.Functions += "\tselect {\n"
	for index, branch := range choiceData.OptionSuffixes {
		WriteChosenFunctionBranch(t, choiceData, m, branch, index)
	}
	t.Functions += "\t}\n"
	WriteChoiceFunctionReturn(t, choiceData, m)
}

func WriteChosenFunction(t *Translation, choiceData *ChoiceData, m *MessagesData) {
	WriteChoiceFunctionHeader(t, choiceData, m)
	WriteChosenFunctionBranches(t, choiceData, m)
}

func WriteChoiceFunction(t *Translation, choiceData *ChoiceData, m *MessagesData) {
	if choiceData.Protagonist == choiceData.Chooser {
		WriteChooserFunction(t, choiceData, m)
	} else {
		WriteChosenFunction(t, choiceData, m)
	}
	t.Functions += "}\n\n"
}

func WriteChoiceFunctions(t *Translation, m *MessagesData) {
	for _, choiceData := range m.Choices {
		WriteChoiceFunction(t, choiceData, m)
	}
}

func WriteRecFunctionHeader(t *Translation, r *RecData, m *MessagesData, firstStructType string) {
	header := "func "
	funcName := "Loop" + r.Protagonist + GetStringSliceAsString(r.Suffix)
	header += funcName
	header += "("
	header += firstStructType + "_old *" + firstStructType
	header += ") ("
	header += "interface{}, "
	header += "error) {\n"
	t.Functions += header
}

func WriteRecFunctionLoopOpening(t *Translation, r *RecData, m *MessagesData) {
	opening := "\tlooping := true\n"
	opening += "\tfor looping {\n"
	t.Functions += opening
}

func WriteRecFunctionRetValDec(t *Translation, r *RecData, m *MessagesData) {
	dec := "\tvar retVal interface{}\n"
	dec += "\tvar err error\n"
	t.Functions += dec
}

func WriteRecFunctionRunCall(t *Translation, r *RecData, m *MessagesData, firstStructType string) {
	call := "\t\tretVal, err "
	call += "= run" + CutStringAfterLetter(firstStructType, "_") + "("
	call += firstStructType + "_old)\n"
	t.Functions += call
}

func WriteRecFunctionErrCheck(t *Translation) {
	check := "\t\tif err != nil {\n"
	check += "\t\t\tlog.Fatal(err)\n"
	check += "\t\t}\n"
	t.Functions += check
}

func WriteRecFunctionSwitch(t *Translation, r *RecData, m *MessagesData, firstStructType string) {
	sw := "\t\tswitch t := retVal.(type) {\n"
	sw += "\t\tcase *" + firstStructType + ":\n"
	sw += "\t\t\t" + firstStructType + "_old = t\n"
	sw += "\t\tdefault:\n"
	sw += "\t\t\tlooping = false\n"
	sw += "\t\t}\n"
	sw += "\t}\n"
	t.Functions += sw
}

func WriteRecFunctionReturn(t *Translation, r *RecData, m *MessagesData) {
	line := "\treturn "
	line += "retVal, "
	line += "nil\n"
	t.Functions += line
}

func WriteRecFunction(t *Translation, r *RecData, m *MessagesData) {
	firstStructType := GetDataStructNameWithIdAndType(r.FirstStructId, r.FirstStructType, m)
	WriteRecFunctionHeader(t, r, m, firstStructType)
	WriteRecFunctionRetValDec(t, r, m)
	WriteRecFunctionLoopOpening(t, r, m)
	WriteRecFunctionRunCall(t, r, m, firstStructType)
	WriteRecFunctionErrCheck(t)
	WriteRecFunctionSwitch(t, r, m, firstStructType)
	WriteRecFunctionReturn(t, r, m)
	t.Functions += "}\n\n"
}

func WriteRecFunctions(t *Translation, m *MessagesData) {
	for _, recData := range m.Recs {
		WriteRecFunction(t, recData, m)
	}
}

func WriteFunctions(t *Translation, m *MessagesData, wg *sync.WaitGroup) {
	defer wg.Done()
	t.Functions += "// Functions\n\n"
	WriteParallelFunctions(t, m)
	WriteChoiceFunctions(t, m)
	WriteStartFunction(t, m)
	WriteRecFunctions(t, m)
}

func GetConvFuncHeaderRetVals(conv *ConversationData) []string {
	headerRetVals := make([]string, 0)
	headerRetVals = append(headerRetVals, "interface{}")
	headerRetVals = append(headerRetVals, "error")
	return headerRetVals
}

func GetConvFuncBottomRetVals(conv *ConversationData, m *MessagesData) []string {
	bottomRetVals := make([]string, 0)
	potentialConv := conv.ConversationElementsData[len(conv.ConversationElementsData)-1]
	if potentialConv.Type == "continue" {
		cont, err := GetContinueDataWithId(m, potentialConv.UniqueId)
		if err != nil {
			log.Fatal(err)
		}
		contStructName := cont.FirstStructName
		contStructName += "_new"
		bottomRetVals = append(bottomRetVals, contStructName)
	} else {
		bottomRetVals = append(bottomRetVals, "\"\"")
	}
	bottomRetVals = append(bottomRetVals, "nil")
	return bottomRetVals
}

func GetConvFuncFinalStruct(conv *ConversationData, m *MessagesData) string {
	structName := ""
	finalConvElem := conv.ConversationElementsData[len(conv.ConversationElementsData)-1]
	if finalConvElem.Type == "continue" {
		cont, err := GetContinueDataWithId(m, finalConvElem.UniqueId)
		if err != nil {
			log.Fatal(err)
		}
		structName += cont.FirstStructName
	} else if finalConvElem.Type == "message" {
		mess, err := GetMessageDataWithId(m, finalConvElem.UniqueId)
		if err != nil {
			log.Fatal(err)
		}
		whereTo := mess.WhereToIfBranchEnds
		if whereTo.Id != "" && whereTo.Type != "" {
			structName = GetDataStructNameWithIdAndType(whereTo.Id, whereTo.Type, m)
		}
	} else if finalConvElem.Type == "choice" {
		choi, err := GetChoiceDataWithId(m, finalConvElem.UniqueId)
		if err != nil {
			log.Fatal(err)
		}
		whereTo := choi.WhereToIfBranchEnds
		if whereTo.Id != "" && whereTo.Type != "" {
			structName = GetDataStructNameWithIdAndType(whereTo.Id, whereTo.Type, m)
		}
	} else if finalConvElem.Type == "par" || finalConvElem.Type == "parallel" {
		par, err := GetParallelDataWithId(m, finalConvElem.UniqueId)
		if err != nil {
			log.Fatal(err)
		}
		whereTo := par.WhereToIfBranchEnds
		if whereTo.Id != "" && whereTo.Type != "" {
			structName = GetDataStructNameWithIdAndType(whereTo.Id, whereTo.Type, m)
		}
	} else if finalConvElem.Type == "rec" {
		rec, err := GetRecDataWithId(m, finalConvElem.UniqueId)
		if err != nil {
			log.Fatal(err)
		}
		whereTo := rec.WhereToIfBranchEnds
		if whereTo.Id != "" && whereTo.Type != "" {
			structName = GetDataStructNameWithIdAndType(whereTo.Id, whereTo.Type, m)
		}
	} else {
		log.Fatal("Unknown type " + finalConvElem.Type + " encountered in GetConvFuncFinalStruct().")
	}
	return structName
}

func AssignConversationFunctionValues(conv *ConversationData, t *Translation, m *MessagesData) *ConvFuncVals {
	id := conv.UniqueId
	paramName := ""
	paramName = GetDataStructNameWithIdAndType(conv.ConversationElementsData[0].UniqueId, conv.ConversationElementsData[0].Type, m)
	headerRetVals := GetConvFuncHeaderRetVals(conv)
	finalStruct := GetConvFuncFinalStruct(conv, m)
	bottomRetVals := GetConvFuncBottomRetVals(conv, m)
	IsFirstConv := IsFirstConv(conv)
	convFuncVals := &ConvFuncVals{ConvId: id, ParamName: paramName, HeaderRetVals: headerRetVals, BottomRetVals: bottomRetVals, IsFirstConv: IsFirstConv, FinalStruct: finalStruct, Protagonist: conv.Protagonist, InParBlock: conv.InParBlock, ParentIsRecBlock: conv.ParentIsRecBlock}
	convFuncStepVals := make([]*ConvFuncStepVals, 0)
	for index, elem := range conv.ConversationElementsData {
		stepVals := AssignConversationStepValues(elem, conv, m, index)
		convFuncStepVals = append(convFuncStepVals, stepVals)
	}
	convFuncVals.Steps = convFuncStepVals
	return convFuncVals
}

func GetMethodNameForDataStructureWithIdAndType(id string, typ string, m *MessagesData) string {
	methodName := ""
	if typ == "message" {
		mess, err := GetMessageDataWithId(m, id)
		if err != nil {
			log.Fatal(err)
		}
		if mess.Protagonist == mess.ToBase {
			methodName += "Receive_"
		} else {
			methodName += "Send_"
		}
		methodName += mess.MethodNameBase
		for _, param := range mess.Parameters {
			methodName += "_" + param
		}
	} else if typ == "choice" {
		choi, err := GetChoiceDataWithId(m, id)
		if err != nil {
			log.Fatal(err)
		}
		methodName += "Make_"
		methodName += choi.Protagonist
		methodName += CutStringAfterLetter(choi.OptionSuffixes[0], "_")
		methodName += "_Choices"
	} else if typ == "par" || typ == "parallel" {
		par, err := GetParallelDataWithId(m, id)
		if err != nil {
			log.Fatal(err)
		}
		methodName += "run" + par.Protagonist
		methodName += CutStringAfterLetter(par.OptionSuffixes[0], "_")
	} else if typ == "rec" {
		rec, err := GetRecDataWithId(m, id)
		if err != nil {
			log.Fatal(err)
		}
		methodName += "Loop" + rec.Protagonist + GetStringSliceAsString(rec.Suffix)
	} else if typ == "continue" {
		// Leave methodName blank
	} else {
		log.Fatal("GetMethodNameForDataStructureWithIdAndType() unable to find data struct with type " + typ + " and id " + id)
	}
	return methodName
}

func GetRetValsForDataStructureWithIdAndType(id string, typ string, m *MessagesData, index int) []string {
	retVals := make([]string, 0)
	if typ == "message" {
		mess, err := GetMessageDataWithId(m, id)
		if err != nil {
			log.Fatal(err)
		}
		if mess.Protagonist == mess.ToBase {
			for i, param := range mess.Parameters {
				val := "received_" + strconv.Itoa(index+1) + "_" + strconv.Itoa(i+1) + "_" + param
				retVals = append(retVals, val)
			}
		}
	} else if typ == "choice" {
		// Add nothing to retVals
	} else if typ == "par" || typ == "parallel" {
		// Add nothing to retVals
	} else if typ == "rec" {
		// Add nothing to retVals
	} else if typ == "continue" {
		// Add nothing to retVals
	} else {
		log.Fatal("GetRetValsForDataStructureWithIdAndType() unable to find data struct with type " + typ + " and id " + id)
	}
	retVals = append(retVals, "err"+strconv.Itoa(index+1))
	return retVals
}

func GetParamsForDataStructureWithIdAndType(id string, typ string, m *MessagesData, index int) []string {
	params := make([]string, 0)
	if typ == "message" {
		mess, err := GetMessageDataWithId(m, id)
		if err != nil {
			log.Fatal(err)
		}
		if mess.Protagonist == mess.FromBase {
			for i, p := range mess.Parameters {
				val := "sending_" + strconv.Itoa(index+1) + "_" + strconv.Itoa(i+1) + "_" + p
				params = append(params, val)
			}
		}
	} else if typ == "choice" {
		// Add nothing to params
	} else if typ == "par" || typ == "parallel" {
		// Add nothing to params
	} else if typ == "rec" {
		// Add nothing to params
	} else if typ == "continue" {
		// Add nothing to params
	} else {
		log.Fatal("GetParamssForDataStructureWithIdAndType() unable to find data struct with type " + typ + " and id " + id)
	}
	return params
}

func GetConversationStepErrCheck(index int) string {
	check := "\tif err" + strconv.Itoa(index+1) + " != nil {\n"
	check += "\t\tlog.Fatal(err" + strconv.Itoa(index+1) + ")\n"
	check += "\t}\n"
	return check
}

func GetConversationStepVarDecs(id string, typ string, m *MessagesData, index int) []string {
	decs := make([]string, 0)
	if typ == "message" {
		mess, err := GetMessageDataWithId(m, id)
		if err != nil {
			log.Fatal(err)
		}
		if mess.Protagonist == mess.FromBase {
			for i, p := range mess.Parameters {
				dec := "\tvar sending_" + strconv.Itoa(index+1) + "_" + strconv.Itoa(i+1) + "_" + p + " " + p + "\n"
				decs = append(decs, dec)
			}
		}
	} else if typ == "choice" {
		// Add nothing to decs
	} else if typ == "par" || typ == "parallel" {
		// Add nothing to decs
	} else if typ == "rec" {
		// Add nothing to decs
	} else if typ == "continue" {
		// Add nothing to decs
	} else {
		log.Fatal("GetConversationStepVarDecs() unable to find data struct with type " + typ + " and id " + id)
	}
	return decs
}

func GetConversationStepVarPrints(id string, typ string, m *MessagesData, index int) []string {
	decs := make([]string, 0)
	if typ == "message" {
		mess, err := GetMessageDataWithId(m, id)
		if err != nil {
			log.Fatal(err)
		}
		if mess.Protagonist == mess.ToBase {
			for i, p := range mess.Parameters {
				dec := "\tfmt.Println(\"" + mess.Protagonist + GetStringSliceAsString(mess.Suffix) + " received value type: " + p + ": \"" + ", received_" + strconv.Itoa(index+1) + "_" + strconv.Itoa(i+1) + "_" + p + ")\n"
				decs = append(decs, dec)
			}
		}
	} else if typ == "choice" {
		// Add nothing to prints
	} else if typ == "par" || typ == "parallel" {
		// Add nothing to prints
	} else if typ == "rec" {
		// Add nothing to prints
	} else if typ == "continue" {
		// Add nothing to prints
	} else {
		log.Fatal("GetConversationStepVarPrints() unable to find data struct with type " + typ + " and id " + id)
	}
	return decs
}

func GetIsChoiceOption(messId string, m *MessagesData) bool {
	mess, err := GetMessageDataWithId(m, messId)
	if err != nil {
		log.Fatal(err)
	}
	return mess.IsChoiceOption
}

func AssignConversationStepValues(elem *ConversationElementData, conv *ConversationData, m *MessagesData, index int) *ConvFuncStepVals {
	stepId := elem.UniqueId
	stepType := elem.Type
	structName := ""
	structName += GetDataStructNameWithIdAndType(elem.UniqueId, elem.Type, m)
	methodName := GetMethodNameForDataStructureWithIdAndType(stepId, stepType, m)
	retVals := GetRetValsForDataStructureWithIdAndType(stepId, stepType, m, index)
	params := GetParamsForDataStructureWithIdAndType(stepId, stepType, m, index)
	errCheck := GetConversationStepErrCheck(index)
	varDecs := GetConversationStepVarDecs(stepId, stepType, m, index)
	varPrints := GetConversationStepVarPrints(stepId, stepType, m, index)
	step := &ConvFuncStepVals{StepId: stepId, StepType: stepType, Index: index, StructName: structName, MethodName: methodName, RetVals: retVals, Params: params, ErrCheck: errCheck, VarDecs: varDecs, VarPrints: varPrints}
	return step
}

func AssignConversationsFunctionValues(t *Translation, m *MessagesData) []*ConvFuncVals {
	convFuncVals := make([]*ConvFuncVals, 0)
	for _, conv := range m.Conversations {
		convFuncVal := AssignConversationFunctionValues(conv, t, m)
		convFuncVals = append(convFuncVals, convFuncVal)
	}
	return convFuncVals
}

func GetConvFuncHeader(convVals *ConvFuncVals, m *MessagesData) string {
	header := "func "
	firstStep := convVals.Steps[0]
	firstStruct := firstStep.StructName
	funcName := firstStruct
	cutFuncNameAfterLastUnderscore := false
	if convVals.IsFirstConv {
		funcName = convVals.Protagonist
	} else if firstStep.StepType == "message" {
		if GetIsChoiceOption(firstStep.StepId, m) {
			funcName = RemoveLastRuneFromString(funcName)
			firstStruct = CutStringAfterLetter(firstStruct, "_")
		}
	} else if firstStep.StepType == "choice" {
		cutFuncNameAfterLastUnderscore = true
	}
	if convVals.ParentIsRecBlock {
		cutFuncNameAfterLastUnderscore = true
	}
	if convVals.InParBlock {
		funcName = RemoveLastRuneFromString(funcName)
	}
	if cutFuncNameAfterLastUnderscore {
		funcName = CutStringAfterLetter(funcName, "_")
	}
	header += "run" + funcName + "("
	if convVals.IsFirstConv {
		header += "wg *sync.WaitGroup, "
	}
	header += firstStruct + " *" + firstStruct
	if !convVals.IsFirstConv && convVals.InParBlock {
		header += ", wg *sync.WaitGroup"
	}
	header += ") ("
	header += "interface{}, "
	header += "error) {\n"
	if convVals.IsFirstConv || convVals.InParBlock {
		header += "\tdefer wg.Done()\n"
	}
	return header
}

func GetConvFuncStepFinalStruct(step *ConvFuncStepVals, convVals *ConvFuncVals) string {
	structName := ""
	if convVals.FinalStruct != "" {
		structName += convVals.FinalStruct + ", "
	}
	return structName
}

func GetCommaSepList(slice []string) string {
	s := ""
	for i, item := range slice {
		s += item
		if i < len(slice)-1 {
			s += ", "
		}
	}
	return s
}

func GetConvFuncStep(convVals *ConvFuncVals, step *ConvFuncStepVals, m *MessagesData) string {
	s := ""
	if step.StepType != "continue" {
		isMethCall := true
		if step.StepType == "choice" || step.StepType == "parallel" || step.StepType == "par" || step.StepType == "rec" {
			isMethCall = false
		}
		for _, varDec := range step.VarDecs {
			s += varDec
		}
		s += "\t"
		isFinalStep := step.Index == len(convVals.Steps)-1
		retValAdded := false
		retValText := "retVal, "
		if !isMethCall {
			s += retValText
			retValAdded = true
		}
		if isFinalStep {
			s += GetConvFuncStepFinalStruct(step, convVals)
		} else if convVals.Steps[step.Index+1].StepType == "continue" {
			cont, err := GetContinueDataWithId(m, convVals.Steps[step.Index+1].StepId)
			if err != nil {
				log.Fatal(err)
			}
			s += cont.ContinueToStruct + "_new, "
		} else {
			nextStep := convVals.Steps[step.Index+1]
			if nextStep.StructName != "" {
				if retValAdded {
					s = s[:len(s)-len(retValText)]
				}
				if GetNeedTypeCheck(convVals, step) {
					s += nextStep.StructName + "_candidate, "
				} else {
					s += nextStep.StructName + ", "
				}
			}
		}
		if strings.Count(s, ",") > 1 {
			s = strings.Replace(s, retValText, "", 1)
		}
		s += GetCommaSepList(step.RetVals)
		s += " := "
		structName := step.StructName
		if step.StepType == "message" {
			if GetIsChoiceOption(step.StepId, m) {
				structName = CutStringAfterLetter(structName, "_")
			}
		}
		if step.Index != 0 {
			if GetNeedTypeCheck(convVals, convVals.Steps[step.Index-1]) {
				structName += "_step" + strconv.Itoa(step.Index)
			}
		}
		if isMethCall {
			s += structName + "." + step.MethodName
			s += "("
			s += GetCommaSepList(step.Params)
		} else {
			s += step.MethodName + "("
			s += structName
		}
		s += ")\n"
		s += step.ErrCheck
		for _, p := range step.VarPrints {
			s += p
		}
	}
	return s
}

func GetNeedTypeCheck(convVals *ConvFuncVals, step *ConvFuncStepVals) bool {
	if step.Index < len(convVals.Steps)-1 {
		if step.StepType == "choice" || step.StepType == "parallel" || step.StepType == "par" || step.StepType == "rec" {
			return true
		}
	}
	return false
}

func GetConvFuncRetValTypeCheck(convVals *ConvFuncVals, step *ConvFuncStepVals, m *MessagesData) string {
	check := ""
	if step.Index < len(convVals.Steps)-1 {
		nextStep := convVals.Steps[step.Index+1]
		if GetNeedTypeCheck(convVals, step) {
			check += "\tvar " + nextStep.StructName + "_step" + strconv.Itoa(step.Index+1) + " *" + nextStep.StructName + "\n"
			check += "\tswitch v := "
			check += nextStep.StructName + "_candidate"
			check += ".(type) {\n"
			check += "\tcase *" + nextStep.StructName
			check += ":\n"
			check += "\t\t" + nextStep.StructName + "_step" + strconv.Itoa(step.Index+1) + " = v\n"
			check += "\tdefault:\n"
			check += "\t\tlog.Fatalf(\"Expected type " + nextStep.StructName + ", received type %T\", v)\n"
			check += "\t}\n"
		}
	}
	return check
}

func GetConvFuncReturnLine(convVals *ConvFuncVals, m *MessagesData) string {
	line := "\treturn "
	if convVals.FinalStruct != "" {
		line += convVals.FinalStruct
	} else if convVals.Steps[len(convVals.Steps)-1].StepType == "continue" {
		cont, err := GetContinueDataWithId(m, convVals.Steps[len(convVals.Steps)-1].StepId)
		if err != nil {
			log.Fatal(err)
		}
		line += cont.ContinueToStruct + "_new"
	} else if convVals.Steps[len(convVals.Steps)-1].StepType == "rec" || convVals.Steps[len(convVals.Steps)-1].StepType == "choice" || convVals.Steps[len(convVals.Steps)-1].StepType == "par" || convVals.Steps[len(convVals.Steps)-1].StepType == "parallel" {
		line += "retVal"
	} else {
		line += "\"\""
	}
	line += ", nil\n"
	return line
}

func GetStringConversationFuncVals(convVals *ConvFuncVals, m *MessagesData) string {
	fun := ""
	fun += GetConvFuncHeader(convVals, m)
	for i := 0; i < len(convVals.Steps); i++ {
		step := convVals.Steps[i]
		fun += GetConvFuncStep(convVals, step, m)
		fun += GetConvFuncRetValTypeCheck(convVals, step, m)
	}
	fun += GetConvFuncReturnLine(convVals, m)
	fun += "}\n\n"
	return fun
}

func GetStringConversationsFuncVals(vals []*ConvFuncVals, m *MessagesData) string {
	funcs := ""
	for _, convVals := range vals {
		funcs += GetStringConversationFuncVals(convVals, m)
	}
	return funcs
}

type ConvFuncStepVals struct {
	StepId     string
	StepType   string
	Index      int
	StructName string
	MethodName string
	RetVals    []string
	Params     []string
	ErrCheck   string
	VarDecs    []string
	VarPrints  []string
}

type ConvFuncVals struct {
	ConvId           string
	Protagonist      string
	ParamName        string
	Steps            []*ConvFuncStepVals
	HeaderRetVals    []string
	BottomRetVals    []string
	IsFirstConv      bool
	FinalStruct      string
	InParBlock       bool
	ParentIsRecBlock bool
}
