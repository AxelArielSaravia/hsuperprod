package roldanonline

import (
    "fmt"
    "strings"
    "io"
    "net/http"
    "strconv"
)


const categoriesURL = "https://roldanonline.com.ar/wp-json/wc/store/products/categories"

var category string

const Header = "category,name,brand,price,sku"
const (
    propName = iota
    propId
    propBrand
    propPrice
    propLen
)

var propsPrefix = [propLen]string{
    propName: "\"name\":\"",
    propId: "\"sku\":\"",
    propBrand: "\"short_description\":\"",
    propPrice: "\"price\":\"",
}

func scrapeItem(b string) {
    var props [propLen]string
    var lst int
    var fst int = strings.Index(b, "\"id\":")
    for fst != -1 {
        b = b[fst:]
        for i := 0; i < propLen; i += 1 {
            fst = strings.Index(b, propsPrefix[i])
            if fst == -1 {
                return
            }
            b = b[fst+len(propsPrefix[i]):]
            lst = strings.Index(b, "\"")
            if lst == -1 {
                return
            }
            var item string = b[:lst]
            if i == propPrice {
                var n int
                var err error
                n, err = strconv.Atoi(item)
                if err != nil {
                    return
                }
                var end  int = n - (n/100*100)
                n /= 100
                if end > 0 {
                    n += 1
                }

                item = strconv.Itoa(n)
            } else if i == propBrand && item != "" {
                fst = strings.Index(item, ">")
                if fst == -1 {
                    return
                }
                item = item[fst+1:]

                lst = strings.Index(item, "<")
                if lst == -1 {
                    return
                }
                item = item[:lst]
            }

            props[i] = string(item)
        }

        fst = strings.Index(b, "\"extensions\":")
        if fst == -1 {
            return
        }
        fmt.Printf(
            "%s,%s,%s,%s,%s\n",
            category,
            props[propName],
            props[propBrand],
            props[propPrice],
            props[propId],
        )

        b = b[fst:]
        fst = strings.Index(b, "\"id\":")
    }
}

func searchItems(apiURL string) error {
    var page = 1

    var err error
    var res *http.Response
    res, err = http.Get(apiURL + strconv.Itoa(page))
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

    err = res.Body.Close()
    if err != nil {
        return err
    }
    scrapeItem(body)

    for len(body) > 2 {
        page += 1

        res, err = http.Get(apiURL + strconv.Itoa(page))
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

        err = res.Body.Close()
        if err != nil {
            return err
        }
        scrapeItem(body)
    }
    return nil
}

func categoryId(b, slug string) (id string) {
    var lst int = strings.Index(b, "\"slug\":\""+slug)
    if lst == -1 {
        return ""
    }
    b = b[:lst]

    var fst int = strings.LastIndex(b, "\"id\":")
    if fst == -1 {
        return ""
    }
    b = b[fst+5:]

    lst = strings.Index(b, ",")
    if lst == -1 {
        return ""
    }

    id = b[:lst]

    b = b[lst:]

    fst = strings.Index(b, "\"name\":\"")
    if fst == -1 {
        return
    }
    b = b[fst+8:]

    lst = strings.Index(b, "\"")
    if lst == -1 {
        return
    }
    category = b[:lst]
    return
}

func getCategory(slug string) (string, error) {
    var res *http.Response
    var err error
    res, err = http.Get(categoriesURL)
    if err != nil {
        return "", err
    }

    if res.StatusCode < 200 || 299 < res.StatusCode {
        return "", fmt.Errorf("Response failed with status code: %d\n", res.StatusCode)
    }

    var body []byte
    body, err = io.ReadAll(res.Body)
    if err != nil {
        return "", err
    }
    err = res.Body.Close()
    if err != nil {
        return "", err
    }
    var b string = string(body)

    var id string = categoryId(b, slug)
    if id == "" {
        return "", fmt.Errorf("Id is not found")
    }
    return id, nil
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
    var slug string
    var str int = strings.LastIndex(url, "/")
    if str != -1 {
        if str == len(url)-1 {
            url = url[:str]
            str = strings.LastIndex(url, "/")
        }
        slug = url[str+1:]
    } else {
        slug = url
    }

    var id string
    var err error
    id, err = getCategory(slug)
    if err != nil {
        return err
    }

    var apiURL string = "https://roldanonline.com.ar/wp-json/wc/store/products?stock_status=instock&orderby=menu_order&order=asc&per_page=100&category="+id+"&page=";
    err = searchItems(apiURL)
    if err != nil {
        return err
    }

    return nil
}
