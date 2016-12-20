package sailor

import (
  "log"
  "github.com/PuerkitoBio/goquery"
)

func proxyMap() map[int]string  {
  urls := []string {"http://www.us-proxy.org", "http://free-proxy-list.net/uk-proxy.html", "http://free-proxy-list.net/anonymous-proxy.html", "http://free-proxy-list.net"}
  m := make(map[int]string)

  for num, url  := range urls {
    // doc, err := goquery.NewDocument("http://www.us-proxy.org") 
    doc, err := goquery.NewDocument(url ) 
    if err != nil {
      log.Fatal(err)
    }
 
    proxyAddr := ""
    doc.Find("tbody").Find("tr").Each(func(i int, s *goquery.Selection) {
        s.Find("td").Each(func(j int, se *goquery.Selection){
         if ( 0 == j ) {
           proxyAddr = se.Text()
         }
         if ( 1 == j ) {
           proxyAddr += ":" + se.Text()
         }
        })
       cc := num*10000 + i
       m[cc] = proxyAddr
     })
  }
  return m 
}
