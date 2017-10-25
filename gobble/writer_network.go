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
	"sort"
	"strings"
	"sync"
)

func WriteCloseNetworkConnectionLines(t *Translation, m *MessagesData) {
	otherRoles := GetNonProtagonistRoles(m)
	for _, other := range otherRoles {
		line := "\tc.doneCommunicatingWith" + other
		line += " = make(chan bool)\n"
		t.ChannelConstructorSlice = append(t.ChannelConstructorSlice, line)
		t.Channels += line
	}
}

func WriteChannelConstructor(t *Translation, m *MessagesData) {
	t.ChannelConstructorSlice = make([]string, 0)
	t.Channels += "func NewChannels() *Channels {\n"
	t.Channels += "\tc := new(Channels)\n"
	WriteConstructorParameterLines(t, m)
	WriteCloseNetworkConnectionLines(t, m)
	WriteConstructorParallelLines(t, m)
	WriteConstructorChoiceLines(t, m)
	t.Channels += "\treturn c\n}\n\n"
}

func WriteImports(t *Translation, m *MessagesData) {
	imports := []string{"errors", "sync", "log"}
	if strings.Contains(t.Functions, "fmt.") {
		imports = append(imports, "fmt")
	}
	if strings.Contains(t.Network, "net.") {
		imports = append(imports, "net")
	}
	if strings.Contains(t.Network, "gob.") {
		imports = append(imports, "encoding/gob")
	}
	if strings.Contains(t.Network, "time.") {
		imports = append(imports, "time")
	}
	t.Imports += "import (\n"
	for _, imp := range imports {
		t.Imports += "\t\"" + imp + "\"\n"
	}
	t.Imports += ")\n\n"
}

func WritePackage(t *Translation, m *MessagesData) {
	packageName := "main"
	t.Package = "package " + packageName + "\n\n"
}

func WriteTransmitterStruct() string {
	trans := "type Transmitter struct {\n"
	trans += "\tEncoder *gob.Encoder\n"
	trans += "\tDecoder *gob.Decoder\n"
	trans += "\tConnection net.Conn\n"
	trans += "}\n\n"
	return trans
}

func WriteClientTransmitterConstructor() string {
	constructor := "func NewTransmitter(conType string, serverAddress string) *Transmitter {\n"
	constructor += "\tt := new(Transmitter)\n"
	constructor += "\tvar connToServer net.Conn\n"
	constructor += "\tvar err error\n"
	constructor += "\testablishingConnection := true\n"
	constructor += "\tfor establishingConnection {\n"
	constructor += "\t\tconnToServer, err = net.Dial(conType, serverAddress)\n"
	constructor += "\t\tif err != nil {\n"
	constructor += "\t\t\tfmt.Println(\"Connection error: \" + err.Error())\n"
	constructor += "\t\t\tfmt.Println(\"Retrying...\")\n"
	constructor += "\t\t\ttime.Sleep(time.Second * 3)\n"
	constructor += "\t\t} else {\n"
	constructor += "\t\t\testablishingConnection = false\n"
	constructor += "\t\t\tfmt.Println(\"Successfully established \" + conType + \" connection with \" + serverAddress)\n"
	constructor += "\t\t}\n"
	constructor += "\t}\n"
	constructor += "\tEncoder := gob.NewEncoder(connToServer)\n"
	constructor += "\tDecoder := gob.NewDecoder(connToServer)\n"
	constructor += "\tt.Connection = connToServer\n"
	constructor += "\tt.Encoder = Encoder\n"
	constructor += "\tt.Decoder = Decoder\n"
	constructor += "\treturn t\n"
	constructor += "}\n\n"
	return constructor
}

type DialogueRecord struct {
	Partner      string
	Protagonist  string
	Server       string
	Client       string
	ToMessages   bool
	FromMessages bool
	Messages     []*MessageData
}

func GetDialoguePartners(m *MessagesData) []string {
	partners := make([]string, 0)
	partnerNames := make(map[string]bool)
	for _, mess := range m.Messages {
		name := ""
		if mess.ToBase != mess.Protagonist {
			name = mess.ToBase
		} else if mess.FromBase != mess.Protagonist {
			name = mess.FromBase
		}
		if name == "" {
			// Do nothing, name empty string
		} else if _, ok := partnerNames[name]; ok {
			// Do nothing name, already in map
		} else {
			partnerNames[name] = true
			partners = append(partners, name)
		}
	}
	return partners
}

func GetDialoguesWithPartner(m *MessagesData, partnerName string) *DialogueRecord {
	msgs := make([]*MessageData, 0)
	toMessages := false
	fromMessages := false
	for _, mess := range m.Messages {
		if mess.ToBase == partnerName {
			toMessages = true
		} else if mess.FromBase == partnerName {
			fromMessages = true
		}
		if mess.ToBase == partnerName || mess.FromBase == partnerName {
			msgs = append(msgs, mess)
		}
	}
	names := []string{m.Messages[0].Protagonist, partnerName}
	// First alphabetically designated as client, second alphabetically
	// designated as server
	sort.Strings(names)
	client := names[0]
	server := names[1]
	dialogue := DialogueRecord{Partner: partnerName, Messages: msgs, Client: client, Server: server, ToMessages: toMessages, FromMessages: fromMessages, Protagonist: msgs[0].Protagonist}
	return &dialogue
}

func GetDialogues(m *MessagesData) []*DialogueRecord {
	dialogues := make([]*DialogueRecord, 0)
	partners := GetDialoguePartners(m)
	for _, p := range partners {
		newDialogue := GetDialoguesWithPartner(m, p)
		dialogues = append(dialogues, newDialogue)
	}
	return dialogues
}

func WriteNetworkSentToFunctionHeader(m *MessagesData, dialogue *DialogueRecord) string {
	header := "func SendTo" + dialogue.Partner
	header += "(chans *Channels, trans *Transmitter) {\n"
	return header
}

func WriteNetworkSentToFunctionDefer(m *MessagesData, dialogue *DialogueRecord) string {
	def := "\tdefer func() {\n"
	def += "\t\ttrans.Connection.Close()\n"
	def += "\t}()\n"
	return def
}

func WriteNetworkSentToFunctionIdentifySelf(m *MessagesData, dialogue *DialogueRecord) string {
	slice := []*DialogueRecord{dialogue}
	isClient := GetIsClient(slice)
	id := "\t"
	if !isClient {
		id += "// "
	}
	id += "identifier := Identifier{Id: \""
	id += m.Messages[0].Protagonist + "\"}"
	id += " // Comment out line if connecting to " + dialogue.Partner + " as server; uncomment line if connecting to " + dialogue.Partner + " as client."
	id += "\n"
	id += "\t"
	if !isClient {
		id += "// "
	}
	id += "trans.Encoder.Encode(identifier)"
	id += " // Comment out line if connecting to " + dialogue.Partner + " as server; uncomment line if connecting to " + dialogue.Partner + " as client."
	id += "\n"
	return id
}

func WriteNetworkSendToFunctionCase(mess *MessageData, m *MessagesData, dialogue *DialogueRecord) string {
	self := dialogue.Protagonist
	other := dialogue.Partner
	roleNames := []string{self, other}
	// Order alphabetically
	sort.Strings(roleNames)
	gobName := strings.Title(roleNames[0]) + roleNames[1] + "Gob"
	c := "\t\tcase out := <-chans."
	c += mess.ChanName + ":\n"
	c += "\t\t\tg := " + gobName + "{Name: \""
	c += GetSendValName(mess) + "\", "
	c += GetSendValName(mess) + ": out}\n"
	c += "\t\t\terr := trans.Encoder.Encode(g)\n"
	c += "\t\t\tif err != nil {\n"
	c += "\t\t\t\tlog.Fatal(err)\n"
	c += "\t\t\t}\n"
	return c
}

func WriteNetworkSendToFunctionDoneCase(dialogue *DialogueRecord) string {
	c := "\t\tcase <- chans.doneCommunicatingWith" + dialogue.Partner + ":\n"
	c += "\t\t\treturn\n"
	return c
}

func WriteNetworkSentToFunctionLoop(m *MessagesData, dialogue *DialogueRecord) string {
	loop := "\tfor {\n"
	loop += "\t\tselect {\n"
	for _, mess := range dialogue.Messages {
		if mess.ToBase != mess.Protagonist {
			loop += WriteNetworkSendToFunctionCase(mess, m, dialogue)
		}
	}
	loop += WriteNetworkSendToFunctionDoneCase(dialogue)
	loop += "\t\tdefault:\n"
	loop += "\t\t\t// Keep looping\n"
	loop += "\t\t}\n"
	loop += "\t}\n"
	return loop
}

func WriteNetworkSendToFunction(m *MessagesData, dialogue *DialogueRecord) string {
	fun := WriteNetworkSentToFunctionHeader(m, dialogue)
	fun += WriteNetworkSentToFunctionDefer(m, dialogue)
	fun += WriteNetworkSentToFunctionIdentifySelf(m, dialogue)
	fun += WriteNetworkSentToFunctionLoop(m, dialogue)
	fun += "}\n\n"
	return fun
}

func WriteNetworkSendToFunctions(m *MessagesData, dialogues []*DialogueRecord) string {
	funcs := ""
	for _, dialogue := range dialogues {
		if dialogue.ToMessages {
			funcs += WriteNetworkSendToFunction(m, dialogue)
		}
	}
	return funcs
}

func WriteNetworkReceiveFromFunctionHeader(dialogue *DialogueRecord) string {
	header := "func ReceiveFrom" + dialogue.Partner
	header += "(chans *Channels, trans *Transmitter) {\n"
	return header
}

func WriteNetworkReceiveFromFunctionGobDec(dialogue *DialogueRecord) string {
	self := dialogue.Protagonist
	other := dialogue.Partner
	roleNames := []string{self, other}
	// Order alphabetically
	sort.Strings(roleNames)
	varType := strings.Title(roleNames[0]) + roleNames[1] + "Gob"
	dec := "\tvar in " + varType + "\n"
	return dec
}

func WriteNetworkReceiveFromFunctionLoopDecode(dialogue *DialogueRecord) string {
	dec := "\t\terr := trans.Decoder.Decode(&in)\n"
	dec += "\t\tif err != nil {\n"
	dec += "\t\t\tlog.Fatal(err)\n"
	dec += "\t\t}\n"
	return dec
}

func WriteNetworkReceiveFromFunctionLoopMessage(firstInIfThenElseChain bool, mess *MessageData) string {
	cond := ""
	if firstInIfThenElseChain {
		cond = "\t\tif"
	} else {
		cond = "\t\t} else if"
	}
	msg := cond + " in.Name == \"" + GetSendValName(mess) + "\" {\n"
	msg += "\t\t\tchans." + mess.ChanName
	msg += " <- in." + GetSendValName(mess) + "\n"
	return msg
}

func WriteNetworkReceiveFromFunctionLoopUnknownGob() string {
	unknown := "\t\t} else {\n"
	unknown += "\t\t\tlog.Fatal(\"ReceiveFromServer() received unknown gob: \", in)\n"
	unknown += "\t\t}\n"
	return unknown
}

func WriteNetworkReceiveFromFunctionLoop(dialogue *DialogueRecord) string {
	loop := ""
	loop += "\tfor {\n"
	loop += WriteNetworkReceiveFromFunctionLoopDecode(dialogue)
	firstInIfThenElseChain := true
	for _, mess := range dialogue.Messages {
		if mess.FromBase != mess.Protagonist {
			loop += WriteNetworkReceiveFromFunctionLoopMessage(firstInIfThenElseChain, mess)
			firstInIfThenElseChain = false
		}
	}
	loop += WriteNetworkReceiveFromFunctionLoopUnknownGob()
	loop += "\t}\n"
	return loop
}

func WriteNetworkReceiveFromFunction(m *MessagesData, dialogue *DialogueRecord) string {
	fun := ""
	fun += WriteNetworkReceiveFromFunctionHeader(dialogue)
	fun += WriteNetworkReceiveFromFunctionGobDec(dialogue)
	fun += WriteNetworkReceiveFromFunctionLoop(dialogue)
	fun += "}\n\n"
	return fun
}

func WriteNetworkReceiveFromFunctions(m *MessagesData, dialogues []*DialogueRecord) string {
	funcs := ""
	for _, dialogue := range dialogues {
		if dialogue.FromMessages {
			funcs += WriteNetworkReceiveFromFunction(m, dialogue)
		}
	}
	return funcs
}

func WriteNetworkIdentifierStruct() string {
	strct := "type Identifier struct {\n"
	strct += "\tId string\n"
	strct += "}\n\n"
	return strct
}

func WriteNetworkGobStruct(m *MessagesData, dialogue *DialogueRecord) string {
	self := dialogue.Protagonist
	other := dialogue.Partner
	roleNames := []string{self, other}
	// Order alphabetically
	sort.Strings(roleNames)
	strct := "type " + strings.Title(roleNames[0]) + roleNames[1] + "Gob struct {\n"
	strct += "\tName string\n"
	typeNames := make([]string, 0)
	for _, mess := range dialogue.Messages {
		nextType := GetSendValName(mess)
		typeNames = append(typeNames, nextType)
	}
	// Order alphabetically
	sort.Strings(typeNames)
	for _, t := range typeNames {
		strct += "\t" + t + " " + t + "\n"
	}
	strct += "}\n\n"
	return strct
}

func WriteNetworkGobStructs(m *MessagesData, dialogues []*DialogueRecord) string {
	strcts := ""
	for _, dialogue := range dialogues {
		strcts += WriteNetworkGobStruct(m, dialogue)
	}
	return strcts
}

func WriteNetworkStartConnectionFunctionHeader(m *MessagesData, dialogue *DialogueRecord) string {
	header := "func ConnectTo" + dialogue.Partner + "AsClient"
	header += "(chans *Channels, conType string, serverAddress string) {\n"
	return header
}

func WriteNetworkStartConnectionFunctionTransmitterDef() string {
	def := "\ttrans := NewTransmitter(conType, serverAddress)\n"
	return def
}

func WriteNetworkStartConnectionFunctionLaunchGoRoutines(m *MessagesData, dialogue *DialogueRecord) string {
	launchers := ""
	if dialogue.ToMessages {
		launchers += "\tgo SendTo" + dialogue.Partner + "(chans, trans)\n"
	}
	if dialogue.FromMessages {
		launchers += "\tgo ReceiveFrom" + dialogue.Partner + "(chans, trans)\n"
	}
	return launchers
}

func WriteNetworkStartConnectionFunction(m *MessagesData, dialogue *DialogueRecord) string {
	fun := ""
	fun += WriteNetworkStartConnectionFunctionHeader(m, dialogue)
	fun += WriteNetworkStartConnectionFunctionTransmitterDef()
	fun += WriteNetworkStartConnectionFunctionLaunchGoRoutines(m, dialogue)
	fun += "}\n\n"
	return fun
}

func WriteNetworkStartConnectionFunctions(m *MessagesData, dialogues []*DialogueRecord) string {
	funcs := ""
	for _, dialogue := range dialogues {
		funcs += WriteNetworkStartConnectionFunction(m, dialogue)
	}
	return funcs
}

func WriteNetworkCloseConnectionFunctionHeader(m *MessagesData, dialogue *DialogueRecord) string {
	header := "func CloseConnectAsClientWith" + dialogue.Partner
	header += "(chans *Channels) {\n"
	return header
}

func WriteNetworkCloseConnectionFunctionSend(m *MessagesData, dialogue *DialogueRecord) string {
	line := "\tchans.doneCommunicatingWith" + dialogue.Partner
	line += "<- true\n"
	return line
}

func WriteNetworkCloseConnectionFunction(m *MessagesData, dialogue *DialogueRecord) string {
	fun := WriteNetworkCloseConnectionFunctionHeader(m, dialogue)
	fun += WriteNetworkCloseConnectionFunctionSend(m, dialogue)
	fun += "}\n\n"
	return fun
}

func WriteNetworkCloseConnectionFunctions(m *MessagesData, dialogues []*DialogueRecord) string {
	funcs := ""
	for _, dialogue := range dialogues {
		funcs += WriteNetworkCloseConnectionFunction(m, dialogue)
	}
	return funcs
}

func WriteNetworkAcceptConnectionsFunction() string {
	fun := "func AcceptConnections(conType string, port string, chans *Channels) {\n"
	fun += "\tln, err := net.Listen(conType, port)\n"
	fun += "\tif err != nil {\n"
	fun += "\t\tlog.Fatal(err)\n"
	fun += "\t}\n"
	fun += "\tfor {\n"
	fun += "\t\tconn, err := ln.Accept()\n"
	fun += "\t\tif err != nil {\n"
	fun += "\t\t\tlog.Fatal(err)\n"
	fun += "\t\t}\n"
	fun += "\t\tHandleConnection(conn, chans)\n"
	fun += "\t}\n"
	fun += "}\n\n"
	return fun
}

func WriteNetworkHandleConnectionsFunction(m *MessagesData, dialogues []*DialogueRecord) string {
	fun := "func HandleConnection(conn net.Conn, chans *Channels) {\n"
	fun += "\tdecoder := gob.NewDecoder(conn)\n"
	fun += "\tencoder := gob.NewEncoder(conn)\n"
	fun += "\ttrans := &Transmitter{Decoder: decoder, Encoder: encoder, Connection: conn}\n"
	fun += "\tvar identifier Identifier\n"
	fun += "\terr := trans.Decoder.Decode(&identifier)\n"
	fun += "\tif err != nil {\n"
	fun += "\t\tlog.Fatal(err)\n"
	fun += "\t}\n"
	for i, dialogue := range dialogues {
		if i == 0 {
			fun += "\tif "
		} else {
			fun += "\t} else if "
		}
		fun += "identifier.Id == \"" + dialogue.Partner
		fun += "\" {\n"
		fun += "\t\tHandle" + dialogue.Partner
		fun += "ConnectionAsServer(trans, chans)\n"
	}
	fun += "\t} else {\n"
	fun += "\t\tlog.Fatal(\"HandleConnection received "
	fun += "unknown identifier: \", identifier)\n"
	fun += "\t}\n"
	fun += "}\n\n"
	return fun
}

func WriteNetworkHandleConnectionsAsServerFunction(m *MessagesData, dialogue *DialogueRecord) string {
	fun := "func Handle" + dialogue.Partner
	fun += "ConnectionAsServer(trans *Transmitter, chans *Channels) {\n"
	if dialogue.ToMessages {
		fun += "\tgo SendTo" + dialogue.Partner + "(chans, trans)\n"
	}
	if dialogue.FromMessages {
		fun += "\tgo ReceiveFrom" + dialogue.Partner + "(chans, trans)\n"
	}
	fun += "}\n\n"
	return fun
}

func WriteNetworkHandleConnectionsAsServerFunctions(m *MessagesData, dialogues []*DialogueRecord) string {
	funcs := ""
	for _, dialogue := range dialogues {
		funcs += WriteNetworkHandleConnectionsAsServerFunction(m, dialogue)
	}
	return funcs
}

func WriteNetworkSetupFunctionHeader(m *MessagesData, dialogues []*DialogueRecord) string {
	header := "func SetupNetworkConnections("
	header += "chans *Channels, "
	header += "connType string, "
	header += "address string,"
	header += "port string"
	header += ") {\n"
	return header
}

func GetIsServer(dialogues []*DialogueRecord) bool {
	isServer := false
	for _, dialogue := range dialogues {
		pair := []string{dialogue.Protagonist, dialogue.Partner}
		// Order alphabetically
		sort.Strings(pair)
		if pair[1] == dialogue.Protagonist {
			isServer = true
		}
	}
	return isServer
}

func GetIsClient(dialogues []*DialogueRecord) bool {
	isClient := false
	for _, dialogue := range dialogues {
		pair := []string{dialogue.Protagonist, dialogue.Partner}
		// Order alphabetically
		sort.Strings(pair)
		if pair[0] == dialogue.Protagonist {
			isClient = true
		}
	}
	return isClient
}

func WriteNetworkLaunchAcceptConnections(m *MessagesData, dialogues []*DialogueRecord, isServer bool) string {
	line := "\t"
	if !isServer {
		line += "//"
	}
	line += "go AcceptConnections(connType, port, chans)"
	if !isServer {
		line += " // Uncomment to accept connections as Server"
	} else {
		line += " // Comment out to stop accepting connections as Server"
	}
	line += "\n"
	return line
}

func WriteNetworkLaunchMakeConnectionAsClient(m *MessagesData, dialogue *DialogueRecord) string {
	slice := []*DialogueRecord{dialogue}
	isClient := GetIsClient(slice)
	line := "\t"
	if !isClient {
		line += "// "
	}
	line += "ConnectTo" + dialogue.Partner + "AsClient("
	line += "chans, "
	line += "connType, "
	line += "address + port"
	line += ")"
	if !isClient {
		line += " // Uncomment to connect as client"
	} else {
		line += " // Comment out to stop connecting as client"
	}
	line += "\n"
	return line
}

func WriteNetworkSetupFunction(m *MessagesData, dialogues []*DialogueRecord) string {
	isServer := GetIsServer(dialogues)
	_ = GetIsClient(dialogues)
	fun := WriteNetworkSetupFunctionHeader(m, dialogues)
	fun += WriteNetworkLaunchAcceptConnections(m, dialogues, isServer)
	for _, dialogue := range dialogues {
		fun += WriteNetworkLaunchMakeConnectionAsClient(m, dialogue)
	}
	fun += "}\n\n"
	return fun
}

func WriteNetworkCloseConnectionsFunction(m *MessagesData, dialogues []*DialogueRecord) string {
	fun := "func CloseNetworkConnections(chans *Channels) {\n"
	for _, dialogue := range dialogues {
		fun += "\tchans.doneCommunicatingWith"
		fun += dialogue.Partner + " <- true\n"
		fun += "\t<-chans.doneCommunicatingWith" + dialogue.Partner + "\n"
	}
	fun += "}\n\n"
	return fun
}

func WriteNetwork(t *Translation, m *MessagesData, wg *sync.WaitGroup) string {
	defer wg.Done()
	network := "// Network \n\n"
	network += WriteTransmitterStruct()
	network += WriteClientTransmitterConstructor()
	network += WriteNetworkIdentifierStruct()
	dialogues := GetDialogues(m)
	network += WriteNetworkGobStructs(m, dialogues)
	network += WriteNetworkSendToFunctions(m, dialogues)
	network += WriteNetworkReceiveFromFunctions(m, dialogues)
	network += WriteNetworkStartConnectionFunctions(m, dialogues)
	network += WriteNetworkCloseConnectionFunctions(m, dialogues)
	network += WriteNetworkAcceptConnectionsFunction()
	network += WriteNetworkHandleConnectionsFunction(m, dialogues)
	network += WriteNetworkHandleConnectionsAsServerFunctions(m, dialogues)
	network += WriteNetworkSetupFunction(m, dialogues)
	network += WriteNetworkCloseConnectionsFunction(m, dialogues)
	return network
}
