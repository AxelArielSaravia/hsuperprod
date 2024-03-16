package main

const helpText string = `hsuperprod is a tool to scrape product information form supermarkets web pages

    hsuperprod (-s [m|d|r|h] | -S [supermarket]) -url [https://...]

    S       long supermarket name.
                -S [supermami | disco | roldanonline | hiperlibertad]
    s       short supermarket name.
                -s [m | d | r | h]
                m = supermami
                d = disco
                r = roldanonline
                h = hiperlibertad
    url     url of the page where items are show.

examples:
    hsuperprod -S disco -url https://www.disco.com.ar/Bebidas/Aguas
    hsuperprod -s d -url https://www.disco.com.ar/Bebidas/Aguas
`

const goodFormat = `Error: Bad Format
    hsuperprod (-s [m|d|r|h] | -S [supermarket name]) -url [https://...]
`
