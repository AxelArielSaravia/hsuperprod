package main

import (
    "log"
    "fmt"
    "flag"
    "io"

    "github.com/axelarielsaravia/hsuperprod/supermami"
    "github.com/axelarielsaravia/hsuperprod/disco"
    "github.com/axelarielsaravia/hsuperprod/roldanonline"
    "github.com/axelarielsaravia/hsuperprod/hiperlibertad"
)

type ScrapeFunc func(r io.Reader) (string, error)


const (
    OP_SUPERMAMI = iota
    OP_DISCO
    OP_ROLDANIONLINE
    OP_HIPERLIBERTAD
    OP_LEN
)


var SHORT_OPS = [OP_LEN]string{
    OP_SUPERMAMI: "m",
    OP_DISCO: "d",
    OP_ROLDANIONLINE: "r",
    OP_HIPERLIBERTAD: "h",
}

var LONG_OPS = [OP_LEN]string{
    OP_SUPERMAMI: "supermami",
    OP_DISCO: "disco",
    OP_ROLDANIONLINE: "roldanonline",
    OP_HIPERLIBERTAD: "hiperlibertad",
}

func main() {
    var short *string = flag.String("s", "", "short define a supermarket")
    var long *string = flag.String("S", "", "long define a supermarket")
    var url *string = flag.String("url", "", "define url")

    var helpShort *bool = flag.Bool("h", false, "help")
    var help *bool = flag.Bool("help", false, "help")
    flag.Parse()

    if *helpShort || *help {
        fmt.Print(helpText)
        return
    }

    if *url == ""{
        fmt.Print(goodFormat)
        return
    }

    var err error

    if *short == SHORT_OPS[OP_SUPERMAMI] || *long == LONG_OPS[OP_SUPERMAMI] {
        fmt.Println(supermami.Header)
        err = supermami.Search(*url)
        if err != nil {
            log.Fatal(err)
        }
    } else if *short == SHORT_OPS[OP_DISCO] || *long == LONG_OPS[OP_DISCO] {
        fmt.Println(disco.Header)
        err = disco.Search(*url)
        if err != nil {
            log.Fatal(err)
        }
    } else if *short == SHORT_OPS[OP_ROLDANIONLINE] ||
    *long == LONG_OPS[OP_ROLDANIONLINE] {
        fmt.Println(roldanonline.Header)
        err = roldanonline.Search(*url)
        if err != nil {
            log.Fatal(err)
        }
    } else if *short == SHORT_OPS[OP_HIPERLIBERTAD] ||
    *long == LONG_OPS[OP_HIPERLIBERTAD] {
        fmt.Println(*url)
        fmt.Println(hiperlibertad.Header)
        err = hiperlibertad.Search(*url)
        if err != nil {
            log.Fatal(err)
        }
    } else {
        log.Fatal(goodFormat)
    }
}
