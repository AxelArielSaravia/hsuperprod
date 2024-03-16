package hiperlibertad

import (
    "fmt"
    "io"
    "net/http"
    "strconv"
    "strings"
)

const Header = "category,name,brand,price,sku"

const (
    propName = iota
    propBrand
    propId
    propPrice
    propLen
)

var propsPrefix = [propLen]string{
    propName:   "\"productName\":\"",
    propBrand:  "\"brand\":\"",
    propId:     "\"ean\":\"",
    propPrice:  "\"Price\":",
}

var propsPostfix = [propLen]string{
    propName:   "\"",
    propBrand:  "\"",
    propId:     "\"",
    propPrice:  ",",

}

var category string

func scrapeItem(sitems string) int {
    var items int = 0
    var props [propLen]string
    var fst int = strings.Index(sitems, "\"Product:")
    for fst != -1 {
        sitems = sitems[fst:]
        for i := 0; i < propLen; i += 1 {
            var ifst int = strings.Index(sitems, propsPrefix[i])
            if ifst == -1 {
                return 0
            }
            sitems = sitems[ifst+len(propsPrefix[i]):]
            var ilst int = strings.Index(sitems, propsPostfix[i])
            if ilst == -1 {
                return 0
            }
            if i == propName {
                props[i] = strings.Replace(sitems[:ilst], ",", "",-1)
            } else {
                props[i] = sitems[:ilst]
            }

            sitems = sitems[ilst+1:]
        }

        fmt.Printf(
            "%s,%s,%s,%s,%s\n",
            category,
            props[propName],
            props[propBrand],
            props[propPrice],
            props[propId],
        )
        items += 1
        fst = strings.Index(sitems, "\"Product:")
    }
    return items
}

func searchItems(b string) int {
    var fst int = strings.Index(b,"__STATE__")
    if fst == -1 {
        return 0
    }
    b = b[fst:]

    var lst = strings.Index(b, "\"$ROOT_QUERY")
    if lst == -1 {
        return 0
    }
    b = b[:lst]

    if b != "" {
        return 0
    }
    return scrapeItem(b)
}

func foundItemsLen(b string) int {
    var fst int = strings.Index(b, "\"recordsFiltered\":")
    if fst == -1 {
        return 0
    }
    b = b[fst+18:]

    var lst int = strings.Index(b, ",")
    if lst == -1 {
        return 0
    }
    b = b[:lst]

    if (b == "") {
        return 0
    }
    var err error
    var n int
    n, err = strconv.Atoi(b)
    if err != nil {
        return 0
    }
    return n
}

func getBody(r io.Reader) (string, error) {
    var body []byte
    var err error
    body, err = io.ReadAll(r)
    if err != nil {
        return "", err
    }
    return string(body), nil
}

func Search(url string) error {
    var pathI int = strings.LastIndex(url, "/")
    if pathI == -1 {
        return fmt.Errorf("Category not found")
    }
    category = url[pathI+1:]
    var ends = strings.Index(category, "?")
    if ends != -1 {
        category = category[:ends]
    }
    category = strings.Replace(category, "-", " ", -1)

    var res *http.Response
    var err error
    res, err = http.Get(url)
    if err != nil {
        return err
    }

    if res.StatusCode < 200 || 299 < res.StatusCode {
        return fmt.Errorf("Response failed with status code: %d\n", res.StatusCode)
    }

    var body string
    body, err = getBody(res.Body)
    if err != nil {
        return err
    }
    if body == "" {
        return fmt.Errorf("No string body")
    }
    err = res.Body.Close()
    if err != nil {
        return err
    }

    var page int = 1
    var itemsLen int = foundItemsLen(body);
    if itemsLen == 0 {
        //no found
        return nil
    }

    var items int = searchItems(body)

    url += "?page="
    for items < itemsLen {
        page += 1
        var pageNum string = strconv.Itoa(page)
        var nextUrl string = url + pageNum

        res, err = http.Get(nextUrl)
        if err != nil {
            return err
        }
        if res.StatusCode < 200 || 299 < res.StatusCode {
            return fmt.Errorf("Response failed with status code: %d\n", res.StatusCode)
        }

        body, err = getBody(res.Body)
        if err != nil {
            return err
        }
        if body == "" {
            return fmt.Errorf("No string body")
        }
        err = res.Body.Close()
        if err != nil {
            return err

        }
        items += searchItems(body)
    }

    return nil
}
