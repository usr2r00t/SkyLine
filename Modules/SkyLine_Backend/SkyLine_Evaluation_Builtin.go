package SkyLine

import (
	"bufio"
	"fmt"
	SkyLine_BuiltIn_System "main/Modules/SkyLine_Builtin/SystemFunctions"
	"os"
	"strings"
)

var builtins = map[string]*Builtin{
	"OS_": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			switch args[0].Inspect() {
			case "os_name":
				name, _ := SkyLine_BuiltIn_System.GrabOperatingSystemDataBasedOnKey["os_name"]()
				return &String{Value: name}
			case "os_arch":
				arch, _ := SkyLine_BuiltIn_System.GrabOperatingSystemDataBasedOnKey["os_arch"]()
				return &String{Value: arch}
			default:
				return &String{Value: "unknown value | run SkyLine__('OS') for more information"}
			}
		},
	},
	"USER_": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			switch args[0].Inspect() {
			case "name":
				name, x := SkyLine_BuiltIn_System.GrabUserInformationFromOS["name"]()
				if x != nil {
					return newError("SkyLine backend (ERR_OS_INFO) => got error when working with OS information %s", x)
				} else {
					return &String{Value: name}
				}
			case "gid":
				gid, x := SkyLine_BuiltIn_System.GrabUserInformationFromOS["gid"]()
				if x != nil {
					return newError("SkyLine backend (ERR_OS_INFO) => got error when working with OS information %s", x)
				} else {
					return &String{Value: gid}
				}
			case "uid":
				uid, x := SkyLine_BuiltIn_System.GrabUserInformationFromOS["uid"]()
				if x != nil {
					return newError("SkyLine backend (ERR_OS_INFO) => got error when working with OS information %s", x)
				} else {
					return &String{Value: uid}
				}
			case "username":
				username, x := SkyLine_BuiltIn_System.GrabUserInformationFromOS["username"]()
				if x != nil {
					return newError("SkyLine backend (ERR_OS_INFO) => got error when working with OS information %s", x)
				} else {
					return &String{Value: username}
				}
			case "hdir":
				hdir, x := SkyLine_BuiltIn_System.GrabUserInformationFromOS["hdir"]()
				if x != nil {
					return newError("SkyLine backend (ERR_OS_INFO) => got error when working with OS information %s", x)
				} else {
					return &String{Value: hdir}
				}
			default:
				return &String{Value: "unknown value | run SkyLine__('USER') for more information"}
			}
		},
	},
	"SkyLine__": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			var msg string
			switch args[0].Inspect() {
			case "OS":
				msg += `
				OS or Operating System is a standard SkyLine 
				function to grab or view information about 
				the current operating system in which the 
				SkyLine interpreter is running on. This 
				function has the following values

				OS_("os_name")    | Grabs the current operating system
				OS_("os_arch")    | Grabs the current operating system architecture
				`
			case "USER":
				msg += `
                USER or Username is a standard SkyLine 
                function to grab or view information about 
                the current user in which the SkyLine 
                interpreter is running on. This function has 
                the following values
				
				USER_("username")    | Grabs the current username
				USER_("uid")         | Grabs the current uid
				USER_("gid")         | Grabs the current gid
				USER_("name")        | Grabs the name
				USER_("hdir")        | Grabs the home directory of the user
				`
			default:
				msg += `METHOD DOES NOT EXIST -> `
				msg += args[0].Inspect()
			}
			return &String{Value: msg}
		},
	},

	"length": {
		Fn: func(args ...Object) Object {
			if l := len(args); l != 1 {
				return newError("wrong number of arguments. want=1, got=%d", l)
			}

			switch arg := args[0].(type) {
			case *String:
				return &Integer{Value: int64(len(arg.Value))}
			case *Array:
				return &Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", arg.Type_Object())
			}
		},
	},

	"first": {
		Fn: func(args ...Object) Object {
			if l := len(args); l != 1 {
				return newError("wrong number of arguments. want=1, got=%d", l)
			}

			if typ := args[0].Type_Object(); typ != ArrayType {
				return newError("argument to `first` must be Array, got %s", typ)
			}

			arr := args[0].(*Array)
			if len(arr.Elements) == 0 {
				return NilValue
			}
			return arr.Elements[0]
		},
	},

	"last": {
		Fn: func(args ...Object) Object {
			if l := len(args); l != 1 {
				return newError("wrong number of arguments. want=1, got=%d", l)
			}

			if typ := args[0].Type_Object(); typ != ArrayType {
				return newError("argument to `last` must be Array, got %s", typ)
			}

			arr := args[0].(*Array)
			l := len(arr.Elements)
			if l == 0 {
				return NilValue
			}
			return arr.Elements[l-1]
		},
	},

	"rest": {
		Fn: func(args ...Object) Object {
			if l := len(args); l != 1 {
				return newError("wrong number of arguments. want=1, got=%d", l)
			}

			if typ := args[0].Type_Object(); typ != ArrayType {
				return newError("argument to `last` must be Array, got %s", typ)
			}

			arr := args[0].(*Array)
			l := len(arr.Elements)
			if l == 0 {
				return NilValue
			}

			newElems := make([]Object, l-1)
			copy(newElems, arr.Elements[1:l])
			return &Array{Elements: newElems}
		},
	},

	"push": {
		Fn: func(args ...Object) Object {
			if l := len(args); l != 2 {
				return newError("wrong number of arguments. want=%d, got=%d", 2, l)
			}

			if typ := args[0].Type_Object(); typ != ArrayType {
				return newError("first argument to `push` must be Array, got %s", typ)
			}

			arr := args[0].(*Array)
			l := len(arr.Elements)

			newElems := make([]Object, l+1)
			copy(newElems, arr.Elements)
			newElems[l] = args[1]
			return &Array{Elements: newElems}
		},
	},

	"print": {
		Fn: func(args ...Object) Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect())
			}
			return &String{Value: ""}
		},
	},
	"println": {
		Fn: func(args ...Object) Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return &String{Value: ""}
		},
	},
	"sprint": {
		Fn: func(args ...Object) Object {
			if l := len(args); l == 0 {
				return newError("SkyLine Builtin (sprint): SPRINT function requires 1 (argument), the argument should be a variable you want to convert to a string")
			}
			return &String{Value: fmt.Sprint(args[0].Inspect())}
		},
	},
	"input": {
		Fn: func(args ...Object) Object {
			if l := len(args); l != 2 {
				return newError("wrong number of arguments. want=1, got=%d | SkyLine's builtin functions such as INPUT require you to enter a character and the name of the input such as 'input' and 'n' where n is the second argument to tell the parser when to use that input. Current supported characters are (n) -> newline ", l)
			}
			input := bufio.NewReader(os.Stdin)
			var Payload string
			fmt.Print(args[0].Inspect())
			for {
				switch args[1].Inspect() {
				case "newline":
					Payload, _ = input.ReadString('\n')              // read input until new line
					Payload = strings.Replace(Payload, "\n", "", -1) // read and replace input state
				case "n":
					Payload, _ = input.ReadString('\n')              // read input until new line
					Payload = strings.Replace(Payload, "\n", "", -1) // read and replace input state
				}
				if Payload != "" {
					break
				}
			}
			return &String{Value: Payload}
		},
	},
}
