package supermami

import (
    "fmt"
    "strings"
    "io"
    "net/http"
)


const BaseUrl = "https://www.supermami.com.ar/"

const Header = "category,name,brand,price,sku"
const (
    propId = iota
    propName
    propBrand
    propCategory
    propPrice
    propLen
)
var propsPrefix = [propLen]string{
    propId: "\"id\":\"",
    propName: "\"name\":\"",
    propBrand: "\"brand\":\"",
    propCategory: "\"category\":\"",
    propPrice: "\"price\":",
}

func scrapeItem(items string) {
    var props [propLen]string

    for items != "]" {
        var fst int = strings.Index(items, "{")
        if fst == -1 {
            return
        }
        items = items[fst+1:]
        var lst int = strings.Index(items, "}")
        if lst == -1 {
            return
        }
        var item string = items[:lst+1]

        for i, s := range propsPrefix {
            var ifst int = strings.Index(item, s)
            if ifst == -1 {
                return
            }
            item = item[ifst+len(s):]
            var ilst int = strings.Index(item, "\"")
            if ilst == -1 {
                ilst = strings.Index(item, "}")
                if ilst == -1 {
                    return
                }
            }
            if i == propName {
                props[i] = strings.Replace(item[:ilst], ",", "",-1)
            } else {
                props[i] = item[:ilst]
            }
            item = item[ilst+1:]
        }

        fmt.Printf(
            "%s,%s,%s,%s,%s\n",
            props[propCategory],
            props[propName],
            props[propBrand],
            props[propPrice],
            props[propId],
        )

        items = items[lst+1:]
    }
}

func searchItems(b string) {
    var lst int = strings.Index(b,"categoryProduct")
    if lst == -1 {
        return
    }
    b = b[:lst]

    var fst int = strings.LastIndex(b, "<script>")
    if fst == -1 {
        return
    }
    b = b[fst:]

    fst = strings.Index(b, "[")
    if fst == -1 {
        return
    }
    b = b[fst:]

    lst = strings.LastIndex(b, "]")
    if lst == -1 {
        return
    }
    b = b[:lst]

    if b != "" {
        scrapeItem(b)
    }
}

func scrapeNextPath(b string) string {
    var fst int = strings.Index(b, "pagination")
    if fst == -1 {
        return ""
    }

    b = b[fst:]
    var lst int = strings.Index(b, "container")
    if lst == -1 {
        return ""
    }
    b = b[:lst]

    fst = strings.Index(b, "active")
    if fst == -1 {
        //the active li is not found
        return ""
    }
    b = b[fst:]

    fst = strings.Index(b, "<li>");
    if fst == -1 {
        //no more pages
        return ""
    }
    b = b[fst:]
    fst = strings.Index(b, "href=")
    if fst == -1 {
        return ""
    }
    b = b[fst+6:]
    lst = strings.Index(b, "\"")
    if lst == -1 {
        return ""
    }
    b = b[:lst]

    return b
}

func scrape(r io.Reader) (string, error) {
    var body []byte
    var err error
    body, err = io.ReadAll(r)
    if err != nil {
        return "", err
    }
    var sbody string = string(body)
    var nextPath string = scrapeNextPath(sbody)

    searchItems(sbody)
    return nextPath, nil
}

func Search(url string) error {
    var res *http.Response
    var err error

    //WARNING: we don't check if there are queries already in the url
    //         and if Nrpp exist
    url += "?Nrpp=100"

    res, err = http.Get(url)
    if err != nil {
        return err
    }

    if res.StatusCode < 200 || 299 < res.StatusCode {
        return fmt.Errorf("Response failed with status code: %d\n", res.StatusCode)
    }

    var nextPath string
    nextPath, err = scrape(res.Body)
    if err != nil {
        return err
    }

    err = res.Body.Close()
    if err != nil {
        return err
    }

    for nextPath != "" {
        var nextUrl string = BaseUrl + nextPath
        res, err = http.Get(nextUrl)
        if err != nil {
            return err
        }
        if res.StatusCode < 200 || 299 < res.StatusCode {
            return fmt.Errorf("Response failed with status code: %d\n", res.StatusCode)
        }

        nextPath, err = scrape(res.Body)
        if err != nil {
            return err
        }

        err = res.Body.Close()
        if err != nil {
            return err
        }
    }

    return nil
}
