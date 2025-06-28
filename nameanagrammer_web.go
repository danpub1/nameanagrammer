//go:build web
package main

import (
    "crypto/tls"
    "embed"
    "flag"
    "fmt"
    "io"
    "log"
    "net/http"
    "os/exec"
    "runtime"
    "strings"
)

type headerResponseWriter struct {
    http.ResponseWriter
    req *http.Request
}

func (w *headerResponseWriter) WriteHeader(status int) {
    if strings.HasSuffix(fmt.Sprintf("%v", w.req.URL), ".js") {
        w.ResponseWriter.Header().Set("Content-Type", "application/javascript")
    }
    //w.ResponseWriter.Header().Set("Content-Security-Policy", "frame-ancestors https: 'self'; font-src https: data: 'self'; default-src https: 'unsafe-inline' 'self';");
    //w.ResponseWriter.Header().Set("X-Content-Security-Policy", "font-src https: data: 'self'; default-src https: 'unsafe-inline' 'self';");
    //w.ResponseWriter.Header().Set("X-Frame-Options", "ALLOW-FROM http://localhost");
    w.ResponseWriter.WriteHeader(status)
}

func headerHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
        next.ServeHTTP(&headerResponseWriter{ResponseWriter: rw, req: req}, req)
    })
}

func allowedIPsHandler(next http.Handler, allowedIPs string) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
        remoteAddr := req.RemoteAddr[:strings.LastIndex(req.RemoteAddr, ":")]
        if strings.Compare(allowedIPs, remoteAddr) == 0 || strings.Contains(allowedIPs, "," + remoteAddr + ",") ||
           strings.HasPrefix(allowedIPs, remoteAddr + ",") || strings.HasSuffix(allowedIPs, "," + remoteAddr) {
            next.ServeHTTP(rw, req)
        } else {
            log.Printf("Blocked %v for %v", req.URL, remoteAddr)
            http.NotFoundHandler().ServeHTTP(rw, req)
        }
    })
}

func logHandler(next http.Handler, verbose bool) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
        header := ""
        for key := range req.Header {
            for val := range req.Header[key] {
                header = header + "\n" + key + ": " + string(req.Header[key][val])
            }
        }

        bodystr := ""
        if req.ContentLength > 0 {
            lengthToRead := int64(1024*1024)
            if req.ContentLength < lengthToRead {
                lengthToRead = req.ContentLength
            }
            body := make([]byte, lengthToRead, lengthToRead)
            n, err := req.Body.Read(body)
            if (err != nil && err != io.EOF) || int64(n) != lengthToRead {
                bodystr = fmt.Sprintf("\nERROR: read %v bytes, err: %v\n", n, err)
            } else {
                bodystr = "\n" + string(body)
            }
        } 

        if verbose {
            log.Printf("%v %v %v%v%v\n\n", req.RemoteAddr, req.Method, req.URL, header, bodystr)
        } else {
            log.Printf("%v %v %v\n", req.RemoteAddr, req.Method, req.URL)
        }
        next.ServeHTTP(rw, req)
    })
}

type selfFileSystem struct {
    selfFS             http.FileSystem
    redirectToFilename string
}

func (fs selfFileSystem) Open(name string) (http.File, error) {
    file, err := fs.selfFS.Open(name)
    if err != nil && len(fs.redirectToFilename) != 0 {
        file, err = fs.selfFS.Open(fs.redirectToFilename)
    }
    return file, err
}

// open opens the specified URL in the default browser of the user.
func openBrowser(url string) error {
    var cmd string
    var args []string

    switch runtime.GOOS {
    case "windows":
        cmd = "cmd"
        args = []string{"/c", "start"}
    case "darwin":
        cmd = "open"
    default: // "linux", "freebsd", "openbsd", "netbsd"
        cmd = "xdg-open"
    }
    args = append(args, url)
    return exec.Command(cmd, args...).Start()
}

//go:embed nameanagrammer.wasm.gz
//go:embed index.html
//go:embed wasm_exec.js
var self embed.FS
var redirect string = ""
var root string = ""

func main() {
    var port int
	var host string
    var certFile string
    var keyFile string
    var allowedIPs string
	var noclient bool
    var verbose bool

    var redirect string
    var err error
    var fs selfFileSystem

    // Default is to serve embeddedd FS as / on port 8080 over http

    flag.IntVar(&port, "port", 8080, "Port Number")
    flag.StringVar(&certFile, "certFile", "", "Certificate File (PEM)")
    flag.StringVar(&keyFile, "keyFile", "", "Key File (PEM)")
    flag.StringVar(&allowedIPs, "allowedIPs", "127.0.0.1,[::1]", "Allowed IP Addresses")
    flag.BoolVar(&verbose, "verbose", false, "Log in verbose mode")
	flag.StringVar(&host, "host", "localhost", "Host")
    flag.BoolVar(&noclient, "noclient", false, "Don't launch the browser")
    flag.Parse()

    fs = selfFileSystem{}
    fs.selfFS = http.FS(self)
    fs.redirectToFilename = redirect

    handler := headerHandler(allowedIPsHandler(logHandler(http.FileServer(fs.selfFS), verbose), allowedIPs))

	mux := http.NewServeMux();
	mux.Handle("GET /", handler);
	// handlers access named portions of path via Request.PathValue(name)
	// can now easily add handlers for POST messages

    if certFile != "" && keyFile != "" {
        _, err2 := tls.LoadX509KeyPair(certFile, keyFile)
        if err2 != nil {
            log.Printf("error '%v' reading key pair from %v and %v", err2, certFile, keyFile)
        } else {
            log.Printf("SERVING HTTPS ON PORT %v to %v", port, allowedIPs)
			if (!noclient) {
                go openBrowser(fmt.Sprintf("https://%v:%v/%v", host, port, root))
			}
            err = http.ListenAndServeTLS(fmt.Sprintf(":%v", port), certFile, keyFile, mux)
            // If we ever get here, err is not nil
            log.Printf("error '%v' serving files on port %v", err, port)
        }
    } else {
        log.Printf("SERVING HTTP ON PORT %v to %v", port, allowedIPs)
		if (!noclient) {
			go openBrowser(fmt.Sprintf("http://%v:%v/%v", host, port, root));
		}
        err = http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
        // If we ever get here, err is not nil
        log.Printf("error '%v' serving files on port %v", err, port)
    }
}
