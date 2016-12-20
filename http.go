package sailor

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"errors"
	"path"
	//     "strconv"
	"path/filepath"
)

// func HttpHead(url string) (string, error) {
func HttpHead(url string) (http.Header, error) {
	//  map[string][]string
	resp, err := http.Head(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Http Head method err: %v, try use get mothod range 0-1\n", err)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Range", "bytes=0-1")
		var client http.Client
		resp, err = client.Do(req)
		if resp.StatusCode != http.StatusOK {
		   err = errors.New("http status not ok")
		}
	}
	// log.Println( "=====", resp, resp.StatusCode)
	// if resp.StatusCode == http.StatusOK || resp.StatusCode != http.StatusPartialContent{
	// 	err = errors.New("http status code is not OK! [http 200|206]\n")
	// 	return resp.Header, err
	// }

	return resp.Header, err
}

func HttpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("HttpGet err:", err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("HttpGet read body", err)
		return "", err
	}
	return string(body), err
}

func HttpNewRequest(method, url string, postbody io.Reader) (string, error) {
	client := &http.Client{}
	// req, err := http.NewRequest(method, url, strings.NewReader(v.URLEncode()))
	req, err := http.NewRequest(method, url, postbody)

	// must set "Content-Type" as "application/x-www-form-urlencoded"
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		log.Println("err is :", err)
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println( "client do err: ", err )
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("HttpGet read body", err)
		return "", err
	}
	return string(body), err

}

func ProxyHttpRequest(method, url_addr, proxy_addr string) (response *http.Response, err error) {
	request, err := http.NewRequest(method, url_addr, nil)
	proxy, err := url.Parse(proxy_addr)
        client := &http.Client{
                Timeout: time.Duration(time.Second * 30),
                Transport: &http.Transport{
                        Proxy: http.ProxyURL(proxy),
                },
        }
        response, err = client.Do(request)
        if err != nil {
                log.Println("===", err)
        }
	return
}

// func ProxyHttpGet(proxy_addr, url_addr string) (*http.Response, error) {
func ProxyHttpGet(proxy_addr, url_addr string) (string, error) {
	returnStr := ""
	request, _ := http.NewRequest("GET", url_addr, nil)
	// request, _ := http.NewRequest("PURGE", url_addr, nil)
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")

	proxy, err := url.Parse(proxy_addr)
	if err != nil {
		log.Printf("%s proxy http get failed\n", url_addr)
	}

	client := &http.Client{
		Timeout: time.Duration(time.Second * 30),
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}

	resp, err := client.Do(request)

	if err != nil {
		log.Println(err)
		return "", err
	}
	if resp.StatusCode == http.StatusOK {
		log.Printf("%s resp status code: %v, ok\n", url_addr, resp.StatusCode)
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		returnStr = string(body)
		// }  else if resp.StatusCode == http.StatusConflict {
		//    log.Printf( "%s resp status code: %v, retry my self\n", url_addr, resp.StatusCode)
		//    ProxyHttpGet(proxy_addr, url_addr)
	} else {
		log.Printf("%s resp status code: %v, failed\n", url_addr, resp.StatusCode)
		returnStr = ""
	}

	log.Println("====== return string ", returnStr)
	return returnStr, err
}

func Download(url, filePath string) {
	tokens := strings.Split(filePath, "/")
	fileName := tokens[len(tokens)-1]
	log.Println("Downloading", url, "to", fileName)

	// TODO: check file existence first with io.IsExist
	// absFileName := strings.TrimRight(path, "/") + "/" + fileName
	// output, err := os.Create(path + fileName)
	pdir := filepath.Dir( filePath )
	if !IsFileExists( pdir ) {
		err := os.MkdirAll( pdir,  0755)
		if err != nil { 
			 log.Printf("%s", err)
		} else{
			log.Println("Create Directory OK!")
		}
	}
	tmpFilePath := filePath +"."+ RandString()
	//tmpFilePath := filePath
	output, err := os.Create(tmpFilePath)
	if err != nil {
		log.Println("Error while creating", filePath, "-", err)
		return
	}
	// output.l.Lock()
	defer output.Close()
	// defer output.l.Unlock()
	

	response, err := http.Get(url)
	if err != nil {
		log.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	size, err := io.Copy(output, response.Body)
	if err != nil {
		log.Println("Error while downloading", url, "-", err)
		return
	}

	log.Printf("%s with %v bytes downloaded", tmpFilePath , size)
	FileRename(tmpFilePath, filePath)
	log.Printf("%s with %v bytes downloaded",  filePath, size)
}

func HttpFsname(url_addr string) string {
	u, err := url.Parse(url_addr)
	if err != nil {
		panic(err)
	}
	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		panic(err)
	}
	return m["fsname"][0]
}

func HTTPDownload(uri string) (d []byte, err error) {
    log.Printf("HTTPDownload From: %s.\n", uri)
    timeout := time.Duration(30 * time.Second)
    client := http.Client{
        Timeout: timeout,
    }
    res, err := client.Get(uri)
    // res, err := http.Get(uri)
    if err != nil {
	return
        log.Fatal(err)
    }
    defer res.Body.Close()
    d, err = ioutil.ReadAll(res.Body)
    if err != nil {
	return
        log.Fatal(err)
    }
    log.Printf("ReadFile: Size of download: %d\n", len(d))
    return
}

func DownloadToFile(uri string, dst string) (err error) {
    if !path.IsAbs( dst ) {
       log.Printf("dst filename is not abs, pls check.") 
    }
    log.Printf("DownloadToFile From: %s.\n", uri)
    if d, err := HTTPDownload(uri); err == nil {
        log.Printf("downloaded %s.\n", uri)
        if WriteByteToFile(dst, d) == nil {
            log.Printf("saved %s as %s\n", uri, dst)
        }
    }
   return
}







/*
func Download() {
    f, err := os.OpenFile("./file.exe", os.O_CREATE|os.O_RDWR, 0666)
    if err != nil { panic(err) }
    stat, err := f.Stat()
    if err != nil { panic(err) }
    f.Seek(stat.Size, 0)
    url := "http://dl.google.com/chrome/install/696.57/chrome_installer.exe"
    var req http.Request
    req.Method = "GET"
    // req.UserAgent = UA
    req.Close = true
    req.URL, err = http.ParseURL(url)
    if err != nil { panic(err) }
    header := http.Header{}
    header.Set("Range", "bytes=" + strconv.Itoa64(stat.Size) + "-")
    req.Header = header
    resp, err := http.DefaultClient.Do(&req)
    if err != nil { panic(err) }
    written, err := io.Copy(f, resp.Body)
    if err != nil { panic(err) }
    println("written: ", written)
}
*/

/*
func httpPost() {
	resp, err := http.Post("http://www.01happy.com/demo/accept.php",
		"application/x-www-form-urlencoded",
		strings.NewReader("name=cjb"))
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	log.Println(string(body))
}

func httpPostForm() {
	resp, err := http.PostForm("http://www.01happy.com/demo/accept.php",
		url.Values{"key": {"Value"}, "id": {"123"}})

	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	log.Println(string(body))

}
*/

func HttpPost(urla string, parm string ) (string,  error) {
  client := &http.Client{
    Timeout: time.Duration(time.Second * 5 ),
  }
  // req, err := http.NewRequest("POST", url, parm)
  log.Println( "=======", parm )
  req, err := http.NewRequest("POST", urla, strings.NewReader(parm) )
  if err != nil {
    log.Println(err)
  }
  req.Header.Set("Content-Type","application/json;charset=UTF-8")
  req.Header.Set("User-Agent", "MIG-Patrol")
  req.Header.Set("Referer", "Deanlzhang" )


  defer func() {
    if err:=recover();err!=nil{
      log.Println(  err )
    }
  }()

  resp, err := client.Do(req)
  if err != nil {
    log.Println( err )
  }

  returnStr := ""
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Println( err )
  }
  returnStr = strings.TrimRight( string(body), "\n")

  return  returnStr, err
}
