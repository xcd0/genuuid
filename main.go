package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/google/uuid"
	"github.com/oklog/ulid"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

var (
	version  string = "debug build" // makefileからビルドされると上書きされる。
	revision string
	Revision = func() string { // {{{
		revision := ""
		modified := false
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					//return setting.Value
					revision = setting.Value
					if len(setting.Value) > 7 {
						revision = setting.Value[:7] // 最初の7文字にする
					}
				}
				if setting.Key == "vcs.modified" {
					modified = setting.Value == "true"
				}
			}
		}
		if modified {
			revision = "develop+" + revision
		}
		return revision
	}() // }}}
	parser *arg.Parser // ShowHelp() で使う

	toJsonFromXml = true // jsonからxmlに変換する。
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile) // ログの出力書式を設定する

	args := ArgParse()

	num := 1
	if len(args.Num) > 0 {
		num = args.Num[0]
	}
	if args.Debug {
		log.Printf("num: %v", num)
	}
	str := ""
	for i := 0; i < num; i++ {
		id := ""
		if args.Ulid {
			id = GenULID()
		} else if args.Xid {
			id = GenXid()
		} else if args.GenVersion == 1 {
			id = GenUUIDv1()
		} else if args.GenVersion == 2 {
			if args.V2Domain < 0 || 2 < args.V2Domain { // ドメインは[0,2]
				panic(errors.Errorf("Domains are specified using numbers from 0 to 2. (0:Person, 1:Group, 2: Org). Specify the domain using the --v2-domain DOMAIN argument. See also the --v2-id argument."))
			}
			domain := uuid.Domain(args.V2Domain)
			id = GenUUIDv2(domain, args.V2Id)
		} else if args.GenVersion == 3 {
			if len(args.V35Namespace) != 36 {
				panic(errors.Errorf("Specify the domain using the --namespace UUID argument."))
			}
			if len(args.V35Namespace) != 36 {
				panic(errors.Errorf("Specify the domain using the --namespace UUID argument."))
			}
			namespace := uuid.MustParse(args.V35Namespace)
			data := []byte(args.V35Data)
			id = GenUUIDv3(namespace, data)
		} else if args.GenVersion == 4 {
			id = GenUUIDv4()
		} else if args.GenVersion == 5 {
			if len(args.V35Namespace) != 36 {
				panic(errors.Errorf("Specify the domain using the --namespace UUID argument."))
			}
			namespace := uuid.MustParse(args.V35Namespace)
			data := []byte(args.V35Data)
			id = GenUUIDv5(namespace, data)
		} else if args.GenVersion == 6 {
			id = GenUUIDv6()
		} else if args.GenVersion == 7 {
			id = GenUUIDv7()
		} else {
			panic(errors.Errorf("version %d is not implemented", args.GenVersion))
		}
		if args.Win {
			str += fmt.Sprintf("{%s}\n", id)
		} else {
			str += fmt.Sprintf("%s\n", id)
		}
	}

	if args.ToUpper {
		str = strings.ToUpper(str)
	} else if args.ToLower {
		str = strings.ToLower(str)
	}

	//fmt.Printf("%s\n", str)
	fmt.Printf("%s", str)
	if len(args.Output) > 0 {
		WriteText(args.Output, str)
	}
}

// optionalな引数はV3,V5の時必要。
func GenUUID(
	version int, // UUIDのバージョン番号。現状1-7を実装している。
	optional_v2_domain *uuid.Domain, // バージョン2で指定する必要があるドメイン。ドメインは、Person、Group、Org のいずれかである必要があります。
	optional_v2_id *uint32, // バージョン2で指定する必要があるID。POSIX システムでは、Person ドメインの場合はユーザーの UID、Group ドメインの場合はユーザーの GID である必要があります。Org ドメインまたは非 POSIX システムにおける id の意味は、サイトによって異なります。
	optional_v35_namespace *uuid.UUID, // バージョン3と5で指定する必要がある名前空間。
	optional_v35_data *[]byte, // バージョン3と5で指定する必要がある名前空間。
) uuid.UUID {
	switch version {
	case 1:
		{
			uidv1, err := uuid.NewUUID() // NewUUID returns a Version 1 UUID based on the current NodeID and clock sequence, and the current time.
			if err != nil {
				panic(err)
			}
			return uidv1
		}
	case 2:
		{
			if optional_v2_domain == nil {
				panic(errors.Errorf("v3 requires namespaces and data."))
			}
			if optional_v2_id == nil {
				panic(errors.Errorf("v3 requires namespaces and data."))
			}
			uidv2, err := uuid.NewDCESecurity(*optional_v2_domain, *optional_v2_id) // NewDCESecurity returns a DCE Security (Version 2) UUID.
			if err != nil {
				panic(err)
			}
			return uidv2
		}
	case 3:
		{
			if optional_v35_namespace == nil {
				panic(errors.Errorf("v3 requires namespaces and data."))
			}
			if optional_v35_data == nil {
				panic(errors.Errorf("v3 requires namespaces and data."))
			}
			uidv3 := uuid.NewMD5(*optional_v35_namespace, *optional_v35_data)
			return uidv3
		}
	case 4:
		{
			uidv4, err := uuid.NewRandom() // NewRandom returns a Random (Version 4) UUID.
			if err != nil {
				panic(err)
			}
			return uidv4
		}
	case 5:
		{
			if optional_v35_namespace == nil {
				panic(errors.Errorf("v5 requires namespaces and data."))
			}
			if optional_v35_data == nil {
				panic(errors.Errorf("v5 requires namespaces and data."))
			}
			uidv5 := uuid.NewSHA1(*optional_v35_namespace, *optional_v35_data) // NewSHA1 returns a new SHA1 (Version 5) UUID based on the supplied name space and data.
			return uidv5
		}
	case 6:
		{
			uidv6, err := uuid.NewV6()
			if err != nil {
				panic(err)
			}
			return uidv6
		}
	case 7:
		{
			uidv7, err := uuid.NewV7()
			if err != nil {
				panic(err)
			}
			return uidv7
		}
	}
	panic(errors.Errorf("version %d is not implemented", version))
}

func GenUUIDv1() string {
	return fmt.Sprintf("%s", GenUUID(1, nil, nil, nil, nil))
}
func GenUUIDv2(domain uuid.Domain, id uint32) string {
	return fmt.Sprintf("%s", GenUUID(2, &domain, &id, nil, nil))
}
func GenUUIDv3(namespace uuid.UUID, data []byte) string {
	return fmt.Sprintf("%s", GenUUID(3, nil, nil, &namespace, &data))
}
func GenUUIDv4() string {
	return fmt.Sprintf("%s", GenUUID(4, nil, nil, nil, nil))
}
func GenUUIDv5(namespace uuid.UUID, data []byte) string {
	return fmt.Sprintf("%s", GenUUID(5, nil, nil, &namespace, &data))
}
func GenUUIDv6() string {
	return fmt.Sprintf("%s", GenUUID(6, nil, nil, nil, nil))
}
func GenUUIDv7() string {
	return fmt.Sprintf("%s", GenUUID(7, nil, nil, nil, nil))
}

func GenULID() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return fmt.Sprintf("%s", id)
}

func GenXid() string {
	guid := xid.New()
	return guid.String()
}

func WriteText(file, str string) {
	f, err := os.Create(file)
	defer f.Close()
	if err != nil {
		panic(errors.Errorf("%v", err))
	} else {
		if _, err := f.Write([]byte(str)); err != nil {
			panic(errors.Errorf("%v", err))
		}
	}
}

type Args struct {
	Num        []int  `arg:"positional"   help:"生成個数。" default: [1]`
	GenVersion int    `arg:"--gen-version,-v" help:"生成するUUIDのバージョンを指定する。[1-7]" default:"4" placeholder:"NUM"`
	Ulid       bool   `arg:"--ulid"       help:"ULIDを生成する。生成順でソートできるUUIDのようなもの。"`
	Xid        bool   `arg:"--xid"        help:"xidを生成する。生成順でソートでき、かつURLエンコードが不要なUUIDのようなもの。"`
	Win        bool   `arg:"-w,--win"     help:"Windowsのレジストリで使用される形式で出力する。{}で囲まれる。" default:false`
	ToLower    bool   `arg:"--lower"      help:"小文字で出力する。" default:false`
	ToUpper    bool   `arg:"--upper"      help:"大文字で出力する。" default:false`
	Output     string `arg:"-o,--output"  help:"ファイル出力する。" default:"" placeholder:"OUT"`

	// バージョン2用。
	V2Domain int    `arg:"--v2-domain"  help:"(v2) 0,1,2のいずれかを指定する。(0:Person、1:Group、2:Org)" placeholder:"DOMAIN"`
	V2Id     uint32 `arg:"--v2-id"      help:"(v2) ドメイン内でのID。PersonはUID、GroupはGIDである必要があります。Orgまたは非POSIXではidの意味は、サイトによって異なります。" placeholder:"ID"`

	// バージョン3/5用。
	V35Namespace string `arg:"--namespace" help:"(v3/v5) 名前空間。UUID文字列を指定する。" default:"" placeholder:"UUID"`
	V35Data      string `arg:"--data"      help:"(v3/v5) データ。" default:""`

	Version bool `arg:"--version" help:"このプログラムのバージョン情報を出力する。"`
	Debug   bool `arg:"-d,--debug"      help:"デバッグ用。ログが詳細になる。"`
}

type ArgsVersion struct {
}

func (args *Args) Print() {
	log.Printf(`
	Num          : %v
	GenVersion   : %v
	Ulid         : %v
	Xid          : %v
	ToLower      : %v
	ToUpper      : %v
	Output       : %v
	Win          : %v
	Version      : %v
	V2Domain     : %v
	V2Id         : %v
	V35Namespace : %v
	V35Data      : %v
	`,
		args.Num,
		args.GenVersion,
		args.Ulid,
		args.Xid,
		args.ToLower,
		args.ToUpper,
		args.Output,
		args.Win,
		args.Version,
		args.V2Domain,
		args.V2Id,
		args.V35Namespace,
		args.V35Data,
	)
}

func ShowHelp(post string) {
	buf := new(bytes.Buffer)
	parser.WriteHelp(buf)
	fmt.Printf("%v\n", strings.ReplaceAll(buf.String(), "display this help and exit", "ヘルプを出力する。"))
	if len(post) != 0 {
		fmt.Println(post)
	}
	os.Exit(1)
}
func GetFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}
func ShowVersion() {
	if len(Revision) == 0 {
		// go installでビルドされた場合、gitの情報がなくなる。その場合v0.0.0.のように末尾に.がついてしまうのを避ける。
		fmt.Printf("%v version %v\n", GetFileNameWithoutExt(os.Args[0]), version)
	} else {
		fmt.Printf("%v version %v.%v\n", GetFileNameWithoutExt(os.Args[0]), version, revision)
	}
	os.Exit(0)
}

func ArgParse() *Args {
	args := &Args{
		Num:        []int{1}, //[]int  `arg:"positional"   help:"生成個数。" default: [1]`
		GenVersion: 4,        //int    `arg:"-v"           help:"生成するUUIDのバージョンを指定する。[1-7]" default:4`
		Ulid:       false,    //bool   `arg:"--ulid"       help:"ULIDを生成する。生成順でソートできるUUIDのようなもの。"`
		Xid:        false,    //bool   `arg:"--xid"        help:"xidを生成する。生成順でソートでき、かつURLエンコードが不要なUUIDのようなもの。"`
		Win:        false,    //bool   `arg:"-w,--win"     help:"Windowsのレジストリで使用される形式で出力する。{}で囲まれる。" default:false`
		ToLower:    false,    //bool   `arg:"--lower"      help:"小文字で出力する。" default:false`
		ToUpper:    false,    //bool   `arg:"--upper"      help:"大文字で出力する。" default:false`
		Output:     "",       //string `arg:"-o,--output"  help:"ファイル出力する。"`
		Version:    false,    //bool   `arg:"--version" help:"バージョン情報を出力する。"`
		Debug:      false,    //bool   `arg:"-d,--debug"      help:"デバッグ用。ログが詳細になる。"`

		// バージョン2用。
		V2Domain: -1, // int    `arg:"--v2-domain"  help:"(v2専用) 0,1,2のいずれかを指定する。(0:Person、1:Group、2:Org)"`
		// V2Id:     0, // uint32 `arg:"--v2-id"      help:"(v2専用) ドメイン内でのID。PersonはUID、GroupはGIDである必要があります。Orgまたは非POSIXではidの意味は、サイトによって異なります。"` // バージョン2で指定する必要があるID。POSIX システムでは、Person ドメインの場合はユーザーの UID、Group ドメインの場合はユーザーの GID である必要があります。Org ドメインまたは非 POSIX システムにおける id の意味は、サイトによって異なります。

		// // バージョン3/5用。
		// V35Namespace: "", // string // バージョン3と5で指定する必要がある名前空間。
		// V35Data:      "", // string // バージョン3と5で指定する必要があるデータ。
	}

	var err error
	parser, err = arg.NewParser(arg.Config{Program: GetFileNameWithoutExt(os.Args[0]), IgnoreEnv: false}, args)
	if err != nil {
		ShowHelp(fmt.Sprintf("%v", errors.Errorf("%v", err)))
	}

	if err := parser.Parse(os.Args[1:]); err != nil {
		if err.Error() == "help requested by user" {
			ShowHelp("")
		} else if err.Error() == "version requested by user" {
			ShowVersion()
		} else {
			panic(errors.Errorf("%v", err))
		}
	}

	if args.Version {
		ShowVersion()
		os.Exit(0)
	}

	if args.Debug {
		args.Print()
	}
	return args
}
