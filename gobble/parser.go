// Gobble classifies a Scribble protocol into two basic types of
// structure: ``messages'' (the axiomatic units of the protocol,
// consisting of a message from one communicating party to another) and
// ``conversations'' consisting of ordered sequences of such messages
// and/or other conversations divided into par blocks (indicating where
// multiple conversations may be executed concurrently), choice blocks
// (where multiple options are may be selected by one of the
// conversation's participants) and rec blocks (for which looping
// behaviour is prescribed).

// The parser reads through this tree recursively, launching a new
// goroutine to read through each sub-conversation as it goes and saving
// the information in a tree data structure composed of conversations
// composed of nodes, where each node represents a message, a
// par block, a choice block or a rec block. Each element of this tree
// is represented by a struct, which consists of collection of fields
// which may include other structs. During this part stage each element
// of the tree is given a unique identifier for use in subsequent
// computation.

// input: []string

// output: tree of structs descending from a Protocol root node

package main

import (
	"strconv"
)

type Module struct {
	Name string
}

type Do struct {
	Id    string
	Name  string
	Roles []string
}

type Continue struct {
	Id   string
	Name string
}

type Protocol struct {
	Mod    *Module
	Locals []*Local
	Types  []*Type
}

type Local struct {
	Id          string
	Name        string
	Roles       []*Role
	Conv        *Conversation
	Protagonist string
}

type Conversation struct {
	Id          string
	Name        string
	Nodes       []*Node
	Protagonist string
}

type Choice struct {
	Id          string
	Chooser     string
	Protagonist string
	Convs       []*Conversation
}

type Node struct {
	Id         string
	Name       string
	Cat        string
	Type       string
	Mess       *Message
	Choice     *Choice
	Par        *Parallel
	Conv       *Conversation
	Rec        *Rec
	Do         *Do
	Cont       *Continue
	Subsequent string
}

type Rec struct {
	Id          string
	Name        string
	Conv        *Conversation
	Protagonist string
}

type Parallel struct {
	Id          string
	Protagonist string
	Convs       []*Conversation
}

type Message struct {
	Id          string
	Name        string
	From        string
	RoleFrom    string
	To          string
	RoleTo      string
	Types       []string
	Rec         string
	Protagonist string
}

type Role struct {
	Name string
}

type Type struct {
	Id       string
	Schema   string
	Source   string
	FileName string
	Alias    string
}

func ProcessTokens(tokens []string) *Protocol {
	typeCounter := 0
	typeChan := make(chan *Type)
	globCounter := 0
	globChan := make(chan *Local)
	p := &Protocol{}
	p.Locals = make([]*Local, 0)
	p.Types = make([]*Type, 0)
	writing := false
	inModule := false
	inType := false
	inLocal := false
	currentSection := make([]string, 0)
	openBracketCounter := 0
	closeBracketCounter := 0
	for _, token := range tokens {
		if !writing {
			if token == "module" {
				writing = true
				inModule = true
			} else if token == "type" {
				writing = true
				inType = true
			} else if token == "local" {
				writing = true
				inLocal = true
				openBracketCounter = 0
				closeBracketCounter = 0
			} else {
				panic("Unexpected token: " + token)
			}
		} else {
			if inModule {
				currentSection = append(currentSection, token)
				if token == ";" {
					writing = false
					inModule = false
					p.Mod = ProcessModule(currentSection)
					currentSection = make([]string, 0)
				}
			} else if inType {
				currentSection = append(currentSection, token)
				if token == ";" {
					writing = false
					inType = false
					typeCounter++
					go ProcessConcurrentType(currentSection, typeChan)
					currentSection = make([]string, 0)
				}
			} else if inLocal {
				currentSection = append(currentSection, token)
				if token == "{" {
					openBracketCounter++
				} else if token == "}" {
					closeBracketCounter++
					if openBracketCounter == closeBracketCounter {
						writing = false
						inLocal = false
						globCounter++
						go ProcessConcurrentLocal(currentSection, globChan)
						currentSection = make([]string, 0)
					}
				}
			} else {
				panic("Should be 'in' a module, type or Local, but not. Current token: " + token)
			}
		}
	}
	for i := 0; i < typeCounter; i++ {
		t := <-typeChan
		p.Types = append(p.Types, t)
	}
	for i := 0; i < globCounter; i++ {
		g := <-globChan
		p.Locals = append(p.Locals, g)
	}
	return p
}

func ProcessConcurrentLocal(tokens []string, ch chan *Local) {
	g := ProcessLocal(tokens)
	ch <- g
}

func ProcessConcurrentType(tokens []string, ch chan *Type) {
	t := ProcessType(tokens)
	ch <- t
}

func ProcessModule(tokens []string) *Module {
	name := ""
	for i := 0; i < len(tokens); i++ {
		if tokens[i] == ";" {
			break
		} else {
			name += tokens[i]
		}
	}
	m := &Module{Name: name}
	return m
}

func ProcessType(tokens []string) *Type {
	schema := ""
	source := ""
	fileName := ""
	alias := ""
	writingSchema := false
	writtenSchema := false
	writingSource := false
	writingFileName := false
	writingAlias := false
	token := ""
	prevToken := ""
	prevPrevToken := ""
	for i := 1; i < len(tokens); i++ {
		token = tokens[i]
		prevToken = tokens[i-1]
		// Setting where to write
		// Schema
		if prevToken == "<" && !writtenSchema {
			writingSchema = true
		}
		if writingSchema && token == ">" {
			writingSchema = false
			writtenSchema = true
		}
		// Source
		if prevToken == "\"" && prevPrevToken == ">" {
			writingSource = true
		}
		if writingSource && token == "\"" {
			writingSource = false
		}
		// FileName
		if prevToken == "\"" && prevPrevToken == "from" {
			writingFileName = true
		}
		if writingFileName && token == "\"" {
			writingFileName = false
		}
		// Alias
		if !writingSchema && !writingSource && !writingFileName && !writingAlias && prevToken == "as" {
			writingAlias = true
		}
		if writingAlias && token == ";" {
			writingAlias = false
		}
		// Writing
		if writingSchema {
			schema += token
		} else if writingSource {
			source += token
		} else if writingFileName {
			fileName += token
		} else if writingAlias {
			alias += token
		}
		// Setting second before last token
		prevPrevToken = prevToken
	}
	t := &Type{Schema: schema, Source: source, FileName: fileName, Alias: alias}
	return t
}

func ProcessRoles(tokens []string) []*Role {
	r := make([]*Role, 0)
	token := ""
	prevToken := ""
	for i := 1; i < len(tokens); i++ {
		token = tokens[i]
		prevToken = tokens[i-1]
		if prevToken == "role" {
			newRole := &Role{Name: token}
			r = append(r, newRole)
		}
	}
	return r
}

func ProcessConversation(tokens []string, uniqueIdGen *LocalIdGenerator, protagonist string) *Conversation {
	Nodes := make([]*Node, 0)
	tokenSubsequence := make([]string, 0)
	openingBracketCounter := 0
	closingBracketCounter := 0
	writingRec := false
	writingChoice := false
	writingPar := false
	writingDo := false
	writingContinue := false
	writingMessage := false
	token := ""
	for i := 1; i < len(tokens)-1; i++ { // Ignore opening and closing brackets
		token = tokens[i]
		// Incrementing bracket counters
		if token == "{" {
			openingBracketCounter++
		}
		if token == "}" {
			closingBracketCounter++
		}
		// Setting flags
		if !writingRec && !writingChoice && !writingPar && !writingMessage && !writingDo && !writingContinue {
			openingBracketCounter = 0
			closingBracketCounter = 0
			if token == "rec" {
				writingRec = true
			} else if token == "choice" {
				writingChoice = true
			} else if token == "par" {
				writingPar = true
			} else if token == "do" {
				writingDo = true
			} else if token == "continue" {
				writingContinue = true
			} else {
				writingMessage = true
			}
		}
		// Rec
		if writingRec && openingBracketCounter > 0 && openingBracketCounter == closingBracketCounter {
			tokenSubsequence = append(tokenSubsequence, token)
			writingRec = false
			n := ProcessRec(tokenSubsequence, uniqueIdGen, protagonist)
			Nodes = append(Nodes, n)
			tokenSubsequence = make([]string, 0)
			writingRec = false
		}
		// Choice
		if writingChoice && openingBracketCounter > 0 && openingBracketCounter == closingBracketCounter && token != "or" && tokens[i+1] != "or" {
			tokenSubsequence = append(tokenSubsequence, token)
			n := ProcessChoice(tokenSubsequence, uniqueIdGen, protagonist)
			Nodes = append(Nodes, n)
			tokenSubsequence = make([]string, 0)
			writingChoice = false
		}
		// Par
		if writingPar && openingBracketCounter > 0 && openingBracketCounter == closingBracketCounter && token != "and" && tokens[i+1] != "and" {
			tokenSubsequence = append(tokenSubsequence, token)
			n := ProcessPar(tokenSubsequence, uniqueIdGen, protagonist)
			Nodes = append(Nodes, n)
			tokenSubsequence = make([]string, 0)
			writingPar = false
		}
		// Message
		if writingMessage && token == ";" {
			tokenSubsequence = append(tokenSubsequence, token)
			n := ProcessMessage(tokenSubsequence, uniqueIdGen, protagonist)
			Nodes = append(Nodes, n)
			tokenSubsequence = make([]string, 0)
			writingMessage = false
		}
		// Do
		if writingDo && token == ";" {
			tokenSubsequence = append(tokenSubsequence, token)
			n := ProcessDo(tokenSubsequence, uniqueIdGen, protagonist)
			Nodes = append(Nodes, n)
			tokenSubsequence = make([]string, 0)
			writingDo = false
		}
		// Continue
		if writingContinue && token == ";" {
			tokenSubsequence = append(tokenSubsequence, token)
			n := ProcessContinue(tokenSubsequence, uniqueIdGen, protagonist)
			Nodes = append(Nodes, n)
			tokenSubsequence = make([]string, 0)
			writingContinue = false
		}
		// Writing
		if writingChoice || writingRec || writingMessage || writingDo || writingContinue || writingPar {
			tokenSubsequence = append(tokenSubsequence, token)
		}
	}
	c := &Conversation{Protagonist: protagonist}
	c.Nodes = Nodes
	return c
}

func ProcessMessage(tokens []string, uniqueIdGen *LocalIdGenerator, protagonist string) *Node {
	name := ""
	types := make([]string, 0)
	from := ""
	to := ""
	rec := ""
	writingName := true
	writingTypes := false
	token := ""
	prevToken := ""
	for i := 0; i < len(tokens); i++ {
		token = tokens[i]
		// Name
		if writingName && token == "(" {
			writingName = false
		}
		if writingName {
			name += token
		}
		// Types
		if token == ")" {
			writingTypes = false
		}
		if writingTypes && token != "," {
			types = append(types, token)
		}
		if token == "(" {
			writingTypes = true
		}
		// From
		if prevToken == "from" {
			to = protagonist
			from = token
		}
		// To
		if prevToken == "to" {
			from = protagonist
			to = token
		}
		prevToken = token
	}
	id := uniqueIdGen.GenerateUniqueId(name)
	m := &Message{Id: id, Name: name, Types: types, From: from, RoleFrom: from, To: to, RoleTo: to, Rec: rec, Protagonist: protagonist}
	n := &Node{Cat: "message", Mess: m}
	return n
}

func ProcessDo(tokens []string, uniqueIdGen *LocalIdGenerator, protagonist string) *Node {
	roles := make([]string, 0)
	writingName := false
	writtenName := false
	writingRoles := false
	name := ""
	token := ""
	for i := 0; i < len(tokens); i++ {
		token = tokens[i]
		// Name
		if writingName && token == "(" {
			writingName = false
			writtenName = true
		}
		if writingName {
			name += token
		}
		if token == "do" && !writtenName {
			writingName = true
		}
		// Roles
		if token == ")" {
			writingRoles = false
		}
		if writingRoles && token != "," {
			roles = append(roles, token)
		}
		if token == "(" {
			writingRoles = true
		}
	}
	d := &Do{Name: name, Roles: roles}
	n := &Node{Cat: "do", Do: d}
	return n
}

func ProcessContinue(tokens []string, uniqueIdGen *LocalIdGenerator, protagonist string) *Node {
	name := ""
	token := ""
	for i := 1; i < len(tokens); i++ {
		token = tokens[i]
		if token != ";" {
			name += token
		}
	}
	id := uniqueIdGen.GenerateUniqueId(name)
	c := &Continue{Name: name, Id: id}
	s := &Node{Id: id, Cat: "continue", Cont: c}
	return s
}

func ProcessRec(tokens []string, uniqueIdGen *LocalIdGenerator, protagonist string) *Node {
	s := &Node{}
	name := ""
	currentTokenSubsection := make([]string, 0)
	named := false
	openingBracketCounter := 0
	closingBracketCounter := 0
	token := ""
	prevToken := ""
	for i := 1; i < len(tokens); i++ {
		token = tokens[i]
		prevToken = tokens[i-1]
		if prevToken == "rec" && named == false {
			name = token
			named = true
		}
		if token == "{" {
			openingBracketCounter++
		}
		if token == "}" {
			closingBracketCounter++
		}
		if openingBracketCounter > 0 {
			currentTokenSubsection = append(currentTokenSubsection, token)
		}
		if openingBracketCounter > 0 && openingBracketCounter == closingBracketCounter {
			c := ProcessConversation(currentTokenSubsection, uniqueIdGen, protagonist)
			id := uniqueIdGen.GenerateUniqueId(name)
			r := &Rec{Id: id, Name: name, Conv: c, Protagonist: protagonist}
			s = &Node{Cat: "rec", Rec: r}
			return s
		}
	}
	return s
}

func ProcessChoice(tokens []string, uniqueIdGen *LocalIdGenerator, protagonist string) *Node {
	id := uniqueIdGen.GenerateUniqueId("choice")
	choiceCounter := 0
	convChan := make(chan *Conversation)
	conversations := make([]*Conversation, 0)
	writingConversation := false
	writingChooser := false
	chooser := ""
	chooserNamed := false
	openingBracketCounter := 0
	closingBracketCounter := 0
	currentTokenSubsection := make([]string, 0)
	token := ""
	prevToken := ""
	for i := 1; i < len(tokens); i++ {
		token = tokens[i]
		prevToken = tokens[i-1]
		// Bracket counters
		if token == "{" {
			openingBracketCounter++
		}
		if token == "}" {
			closingBracketCounter++
		}
		// Chooser
		if !chooserNamed && prevToken == "at" {
			writingChooser = true
		}
		if writingChooser && token == "{" {
			writingChooser = false
			chooserNamed = true
		}
		if writingChooser {
			chooser += token
		}
		// Conversation
		if !writingConversation && token == "{" {
			writingConversation = true
		}
		if writingConversation {
			currentTokenSubsection = append(currentTokenSubsection, token)
		}
		if writingConversation && openingBracketCounter == closingBracketCounter {
			choiceCounter++
			go ProcessConcurrentConversation(currentTokenSubsection, convChan, uniqueIdGen, protagonist)
			openingBracketCounter = 0
			closingBracketCounter = 0
			currentTokenSubsection = make([]string, 0)
			writingConversation = false
		}
	}
	for i := 0; i < choiceCounter; i++ {
		c := <-convChan
		conversations = append(conversations, c)
	}
	c := &Choice{Id: id, Chooser: chooser, Convs: conversations, Protagonist: protagonist}
	n := &Node{Cat: "choice", Choice: c}
	return n
}

func ProcessConcurrentConversation(tokens []string, ch chan *Conversation, uniqueIdGen *LocalIdGenerator, protagonist string) {
	c := ProcessConversation(tokens, uniqueIdGen, protagonist)
	ch <- c
}

func ProcessPar(tokens []string, uniqueIdGen *LocalIdGenerator, protagonist string) *Node {
	id := uniqueIdGen.GenerateUniqueId("par")
	parCounter := 0
	convChan := make(chan *Conversation)
	conversations := make([]*Conversation, 0)
	openingBracketCounter := 0
	closingBracketCounter := 0
	currentTokenSubsection := make([]string, 0)
	token := ""
	for i := 1; i < len(tokens); i++ {
		token = tokens[i]
		if token == "{" {
			openingBracketCounter++
		}
		if token == "}" {
			closingBracketCounter++
		}
		if token == "and" && openingBracketCounter == 0 {
			continue
		} else {
			currentTokenSubsection = append(currentTokenSubsection, token)
		}
		if openingBracketCounter > 0 && openingBracketCounter == closingBracketCounter {
			parCounter++
			go ProcessConcurrentConversation(currentTokenSubsection, convChan, uniqueIdGen, protagonist)
			openingBracketCounter = 0
			closingBracketCounter = 0
			currentTokenSubsection = make([]string, 0)
		}
	}
	for i := 0; i < parCounter; i++ {
		c := <-convChan
		conversations = append(conversations, c)
	}
	p := &Parallel{Id: id, Convs: conversations, Protagonist: protagonist}
	s := &Node{Cat: "par", Par: p}
	return s
}

type LocalIdGenerator struct {
	LocalName      string
	RequestIntChan chan bool
	ReceiveIntChan chan int
	StopCounter    chan bool
}

func (idg *LocalIdGenerator) GenerateUniqueId(MessageName string) string {
	idg.RequestIntChan <- true
	counter := <-idg.ReceiveIntChan
	uniqueId := idg.LocalName + MessageName + strconv.Itoa(counter)
	return uniqueId
}

func GenerateCounterValues(idg *LocalIdGenerator) {
	counter := 0
	for {
		select {
		case <-idg.RequestIntChan:
			idg.ReceiveIntChan <- counter
			counter++
		case <-idg.StopCounter:
			return
		default:
			// Keep looping
		}
	}
}

func ProcessLocal(tokens []string) *Local {
	g := &Local{}
	uniqueIdGen := &LocalIdGenerator{}
	uniqueIdGen.ReceiveIntChan = make(chan int)
	uniqueIdGen.RequestIntChan = make(chan bool)
	uniqueIdGen.StopCounter = make(chan bool)
	go GenerateCounterValues(uniqueIdGen)
	name := ""
	protagonist := ""
	rolesTokens := make([]string, 0)
	conversationTokens := make([]string, 0)
	writingName := false
	writingRoles := false
	writingConversation := false
	writingProtagonist := false
	nameWritten := false
	protagonistWritten := false
	rolesWritten := false
	token := ""
	prevToken := ""
	for i := 1; i < len(tokens); i++ {
		token = tokens[i]
		prevToken = tokens[i-1]
		// Name
		if prevToken == "protocol" && !nameWritten {
			writingName = true
			writingRoles = true
		}
		if writingName && token == "at" {
			writingName = false
			nameWritten = true
			uniqueIdGen.LocalName = name
		}
		// Protagonist
		if prevToken == "at" && !protagonistWritten {
			writingProtagonist = true
		}
		// Roles
		if token == "(" && !rolesWritten {
			writingRoles = true
		}
		if writingRoles && token == ")" {
			rolesTokens = append(rolesTokens, token)
			writingRoles = false
			rolesWritten = true
		}
		// Conversation
		if nameWritten && rolesWritten && !writingConversation && token == "{" {
			writingConversation = true
		}
		// Writing
		if writingName {
			name += token
		}
		if writingProtagonist {
			protagonist += token
			writingProtagonist = false
			protagonistWritten = true
		}
		if writingRoles {
			rolesTokens = append(rolesTokens, token)
		}
		if writingConversation {
			conversationTokens = append(conversationTokens, token)
		}
	}
	g.Roles = ProcessRoles(rolesTokens)
	g.Conv = ProcessConversation(conversationTokens, uniqueIdGen, protagonist)
	g.Name = name
	g.Protagonist = protagonist
	uniqueIdGen.StopCounter <- true
	return g
}
