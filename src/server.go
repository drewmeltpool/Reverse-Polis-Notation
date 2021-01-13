package main

import (
	"fmt"
	"encoding/json"
	"log"
	"os"
	"net/http"
	"strconv"
	"math"
)
//GlobalVariables ...
type GlobalVariables struct {
    Port	string	`json:"port"`
    Host	string	`json:"host"`
}
//Answer ...
type Answer struct {
    Postfix	string	`json:"postfix"`
    Prefix	string	`json:"prefix"`
    Infix	string	`json:"infix"`
    Answer	string	`json:"answer"`
}
//Answers ...
type Answers struct {
	Items []Answer
}
//AddItem ...
func (box *Answers) AddItem(item Answer) {
	box.Items = append(box.Items, item)
}

func main(){

	globalVariablesFile, err := os.Open("../config/default.json")
	if err != nil {log.Fatal(err)}
    defer globalVariablesFile.Close()

	var globalVariablesDecoder *json.Decoder = json.NewDecoder(globalVariablesFile)
    if err != nil {
        log.Fatal(err)
	}
	
	var GlobalVariable GlobalVariables

	err = globalVariablesDecoder.Decode(&GlobalVariable)
    if err != nil {
        log.Fatal(err)
    }
	
	serverPort := GlobalVariable.Port
	serverHost := GlobalVariable.Host

	http.HandleFunc("/",ServeFiles)
	fmt.Println("Serving @ : ",serverHost + serverPort)
	log.Fatal(http.ListenAndServe(serverPort,nil))
}

//ServeFiles ...
func ServeFiles(w http.ResponseWriter, r *http.Request){

	switch r.Method{

	case "GET":

		path := r.URL.Path

		fmt.Println(path)

		if path == "/"{
			path = "./static/index.html"
		}else{
			path = "." + path
		}

		http.ServeFile(w,r,path)

	case "POST":

		r.ParseMultipartForm(0)

		message := r.FormValue("message")
		res := splite(message,";")


		fmt.Println("----------------------------------")
		fmt.Println("Message from Client: ", res)
		box := Answers{}
		for i := 0; i < len(res); i++ {
			answer := Answer{Postfix: Postfix(res[i]),Prefix:PostfixToPrefix(res[i]),Infix:PostfixToInfix(res[i]),Answer:CalcPostfix(res[i])}
			box.AddItem(answer)
		}
		fmt.Println("Message from Client: ", box.Items)
		JSON, _ := json.Marshal(box)
		fmt.Println(string(JSON))
		
		fmt.Fprintf(w, "%s", string(JSON))
	
		default:
		
			fmt.Fprintf(w,"Request type other than GET or POSt not supported")

	}

}

//Postfix ...
func Postfix(elem string) string {
	if include(elem) == 0 {
		return "Incorrect statement..."
	}
	if len(splite(elem, " ")) < 3 {
		return "dont enought args..."
	}
	return elem
}
//CalcPostfix ...
func CalcPostfix(elem string) string {
	if include(elem) == 0 {
		return "Incorrect statement..."
	}
	if len(splite(elem, " ")) < 3 {
		return "dont enought args..."
	}
	return fmt.Sprintf("%f",Calc(elem))
}
//PostfixToPrefix ...
func PostfixToPrefix(elem string) string {
	if include(elem) == 0 {
		return "Incorrect statement..."
	}
	if len(splite(elem, " ")) < 3 {
		return "dont enought args..."
	}
	return prefix(infix(elem))
}
//PostfixToInfix ...
func PostfixToInfix(elem string) string {
	if include(elem) == 0 {
		return "Incorrect statement..."
	}
	if len(splite(elem, " ")) < 3 {
		return "dont enought args..."
	}
	return infix(elem)
}


func include(state string) int {
	if len(state) == 0 {
		return 0
	}
	for i := 0; i < len(state); i++ {
		if haveItem(string(state[i])) == 0 {
			return 0
		}
	}
	return 1
}
func haveItem(char string) int {
	input := " .0123456789+-*/^()"
	for j := 0; j < len(input); j++ {
		if char == string(input[j]) {
			return 1
		}
	}
	return 0
}

func splite(elem string, space string) (result []string) {

	input := string([]byte(elem))
	arrayItem := ""
	for i := 0; i < len(input); i++ {
		if space == string("") {
			arrayItem += string(input[i])
			result = append(result, arrayItem)
			arrayItem = ""
		} else {
			if string(input[i]) != space {
				arrayItem += string(input[i])
			} else {
				result = append(result, arrayItem)
				arrayItem = ""
			}
			if i == len(input)-1 {
				result = append(result, arrayItem)
			}
		}
	}
	return
}

func toString(elem []string) (res string) {
	for _, v := range elem {
		res += v
	}
	return
}

func reverse(str string) (result string) {
	for _, v := range str {
		if string(v) == string("(") {
			result = ")" + result
		} else if string(v) == string(")") {
			result = "(" + result
		} else {
			result = string(v) + result
		}

	}
	return
}

func getPriority(elem string) int {

	if elem == string("^") {
		return 4
	} else if elem == string("*") || elem == string("/") {
		return 3
	} else if elem == string("+") || elem == string("-") {
		return 2
	} else if elem == string("(") {
		return 1
	} else if elem == string(")") {
		return -1
	} else {
		return 0
	}
}

func postfix(elem string) string {
	current := ""
	stack := ""
	for i := 0; i < len(elem); i++ {
		if getPriority(string(elem[i])) == 0 {
			current += string(elem[i])
		}
		if getPriority(string(elem[i])) == 1 {
			stack += string(elem[i])
		}
		if getPriority(string(elem[i])) == -1 {
			for j := len(stack) - 1; j >= 0; j-- {
				if getPriority(string(stack[j])) != 1 {
					current += " "
					current += string(stack[j])
					stack = stack[:len(stack)-1]
				} else {
					stack = stack[:len(stack)-1]
					break
				}
			}
		}
		if getPriority(string(elem[i])) > 1 {
			current += " "
			for j := len(stack) - 1; j >= 0; j-- {
				if getPriority(string(stack[j])) >= getPriority(string(elem[i])) {
					current += string(stack[j])
					current += " "
					stack = stack[:len(stack)-1]
				} else {
					break
				}
			}
			stack += string(elem[i])
		}
	}
	for j := len(stack) - 1; j >= 0; j-- {
		current += " "
		current += string(stack[j])
	}
	return current
}

func prefix(elem string) string {
	return reverse(postfix(reverse(elem)))
}

func infix(elem string) string {
	array := splite(elem, " ")
	stack := make([]string, 0)
	str := ""
	for i := 0; i < len(array); i++ {
		if getPriority(array[i]) == 0 {
			stack = append(stack, array[i])
		} else {
			length := len(stack)
			str += "("
			str += stack[len(stack)-2]
			str += array[i]
			str += stack[len(stack)-1]
			str += ")"
			stack = stack[0 : length-2]
			stack = append(stack, str)
			str = ""
		}
	}

	return toString(stack)
}

//Calc ...
func Calc(elem string) float64 {
	array := splite(elem, " ")
	stack := make([]float64, 0)
	for i := 0; i < len(array); i++ {
		
		if getPriority(array[i]) == 0 {
			numb, _ := strconv.ParseFloat(array[i], 64)
			stack = append(stack, numb)
		}else {
			length := len(stack)
			a := stack[len(stack)-1]
			b := stack[len(stack)-2]

			stack = stack[0 : length-2]

			if(array[i] == "+") {stack = append(stack,b+a)}
			if(array[i] == "-") {stack = append(stack,b-a)}
			if(array[i] == "*") {stack = append(stack,b*a)}
			if(array[i] == "/") {stack = append(stack,b/a)}
			if(array[i] == "^") {stack = append(stack,math.Pow(a,b))}

		}
	}
	if(len(stack) < 1){
		length := len(stack)
			a := stack[len(stack)-1]
			b := stack[len(stack)-2]

			stack = stack[0 : length-2]
			// fmt.Println("a:",a)

			// fmt.Println("b:",b)
			if(array[len(array)-1] == "+") {stack = append(stack,b+a)}
			if(array[len(array)-1] == "-") {stack = append(stack,b-a)}
			if(array[len(array)-1] == "*") {stack = append(stack,b*a)}
			if(array[len(array)-1] == "/") {stack = append(stack,b/a)}
			if(array[len(array)-1] == "^") {stack = append(stack,math.Pow(a,b))}	
	}
	
	return stack[0]
}