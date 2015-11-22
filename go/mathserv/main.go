package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/skia-dev/glog"
	"go.skia.org/infra/go/common"
	"go.skia.org/infra/go/login"
	"go.skia.org/infra/go/metadata"
	"go.skia.org/infra/go/util"
)

const (
	CA_CERT_CHAIN_FILENAME = "caCertChain.pem"
	CA_KEY_FILENAME        = "caKey.pem"
)

var (
	// indexTemplate is the main index.html page we serve.
	indexTemplate *template.Template = nil
)

// flags
var (
	port           = flag.String("port", ":8000", "HTTP(S) service address (e.g., ':8000')")
	httpPort       = flag.String("http_port", ":8001", "HTTP service address (e.g., ':8000'), only used for redirects, and only if certChainFile is set.")
	local          = flag.Bool("local", false, "Running locally if true. As opposed to in production.")
	graphiteServer = flag.String("graphite_server", "localhost:2003", "Where is Graphite metrics ingestion server running.")
	resourcesDir   = flag.String("resources_dir", "", "The directory to find templates, JS, and CSS files. If blank the current directory will be used.")
	workDir        = flag.String("work_dir", "/tmp", "Directory to keep scratch and work files.")

	certChainFile = flag.String("cert_chain_file", "", "The file name of the TLS certificate chain. If not set then the server only serves HTTP.")
	keyFile       = flag.String("key_file", "", "The file name of the TLS certificate key.")
)

func loadTemplates() {
	indexTemplate = template.Must(template.ParseFiles(
		filepath.Join(*resourcesDir, "templates/index.html"),
		filepath.Join(*resourcesDir, "templates/titlebar.html"),
		filepath.Join(*resourcesDir, "templates/header.html"),
	))
}

func Init() {
	if *resourcesDir == "" {
		_, filename, _, _ := runtime.Caller(0)
		*resourcesDir = filepath.Join(filepath.Dir(filename), "../..")
	}
	loadTemplates()

}

// mainHandler handles the GET of the main page.
func mainHandler(w http.ResponseWriter, r *http.Request) {
	if *local {
		loadTemplates()
	}
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html")
		if err := indexTemplate.Execute(w, struct{}{}); err != nil {
			glog.Errorln("Failed to expand template:", err)
		}
	}
}

func makeResourceHandler() func(http.ResponseWriter, *http.Request) {
	fileServer := http.FileServer(http.Dir(*resourcesDir))
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", string(300))
		fileServer.ServeHTTP(w, r)
	}
}

func redirHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://mathinate.com", 302)
}

func AttemptLoadCertFromMetadata() {
	if *certChainFile == "" {
		// Try loading from GCE project level metadata.
		certChainContents, err := metadata.ProjectGet("ca_cert_chain")
		if err != nil {
			glog.Errorf("Failed to load ca_cert_chain from metadata: %s", err)
			return
		}
		keyContents, err := metadata.ProjectGet("ca_key")
		if err != nil {
			glog.Errorf("Failed to load ca_key from metadata: %s", err)
			return
		}
		fullCertChainFilename := filepath.Join(*workDir, CA_CERT_CHAIN_FILENAME)
		fullKeyFilename := filepath.Join(*workDir, CA_KEY_FILENAME)
		if err := ioutil.WriteFile(fullCertChainFilename, []byte(certChainContents), 0600); err != nil {
			glog.Errorf("Failed to write %s: %s", fullCertChainFilename, err)
			return
		}
		if err := ioutil.WriteFile(fullKeyFilename, []byte(keyContents), 0600); err != nil {
			glog.Errorf("Failed to write %s: %s", fullKeyFilename, err)
			return
		}
		*keyFile = fullKeyFilename
		*certChainFile = fullCertChainFilename
		glog.Infof("SUCCESS: Loaded cert from metadata.")
	}
}

func main() {
	defer common.LogPanic()
	common.InitWithMetrics("mathserv", graphiteServer)
	Init()

	// By default use a set of credentials setup for localhost access.
	var cookieSalt = "notverysecret"
	var clientID = "952643138919-5a692pfevie766aiog15io45kjpsh33v.apps.googleusercontent.com"
	var clientSecret = "QQfqRYU1ELkds90ku8xlIGl1"
	var redirectURL = fmt.Sprintf("http://localhost%s/oauth2callback/", *port)
	if !*local {
		cookieSalt = metadata.Must(metadata.ProjectGet(metadata.COOKIESALT))
		clientID = metadata.Must(metadata.ProjectGet(metadata.CLIENT_ID))
		clientSecret = metadata.Must(metadata.ProjectGet(metadata.CLIENT_SECRET))
		redirectURL = "https://mathinate.com/oauth2callback/"
	}
	login.Init(clientID, clientSecret, redirectURL, cookieSalt, login.DEFAULT_SCOPE, "", *local)

	r := mux.NewRouter()
	r.PathPrefix("/res/").HandlerFunc(util.MakeResourceHandler(*resourcesDir))
	r.HandleFunc("/", mainHandler)
	r.HandleFunc("/loginstatus/", login.StatusHandler)
	r.HandleFunc("/logout/", login.LogoutHandler)
	r.HandleFunc("/oauth2callback/", login.OAuth2CallbackHandler)
	http.Handle("/", util.LoggingGzipRequestResponse(r))
	AttemptLoadCertFromMetadata()
	glog.Infoln("Ready to serve.")

	if *certChainFile != "" {
		glog.Infof("Serving TLS")
		go func() {
			redir := mux.NewRouter()
			redir.HandleFunc("/", redirHandler)
			glog.Fatal(http.ListenAndServe(*httpPort, redir))
		}()
		glog.Fatal(http.ListenAndServeTLS(*port, *certChainFile, *keyFile, nil))
	} else {
		glog.Infof("Only serving HTTP")
		glog.Fatal(http.ListenAndServe(*port, nil))
	}
}
