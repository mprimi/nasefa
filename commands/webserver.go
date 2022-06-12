package commands

import (
  "fmt"
  "net/http"
  "html/template"
  "regexp"
  "github.com/c2h5oh/datasize"
)

import _ "embed"

//go:embed web/list_bundles.html
var listBundlesHtml string

var templates struct {
  listBundles *template.Template
}

type bundleListItem struct {
  Name      string
  NumFiles  int
  Size      string
}

type bundlesList struct {
  Bundles     []*bundleListItem
}

func listBundlesHandler(w http.ResponseWriter, req *http.Request) {
  bundles, err := loadBundles()
  if err != nil {
    logWarn("Error loading bundles: %s", err)
    http.Error(w, fmt.Sprintf("Internal error: %s", err), http.StatusInternalServerError)
  }

  bl := &bundlesList{
    Bundles: []*bundleListItem{},
  }

  for _, bundle := range bundles {
    b := &bundleListItem{
      Name: bundle.name,
      NumFiles: len(bundle.files),
      Size: datasize.ByteSize(bundle.objStoreStatus.Size()).HumanReadable(),
    }
    bl.Bundles = append(bl.Bundles, b)
  }

  err = templates.listBundles.Execute(w, bl)
  if err != nil {
    http.Error(w, fmt.Sprintf("Internal error: %s", err), http.StatusInternalServerError)
    return
  }
}

func makeHandler(prefix string) (func(w http.ResponseWriter, req *http.Request)) {

  kBundleNameRe := "[a-zA-Z0-9_\\.\\-]{4,128}"
  kFileNameRe := "[a-zA-Z0-9_\\.\\-]{4,128}"

  rootRe := regexp.MustCompile("^" + prefix + "$")
  listBundleRe := regexp.MustCompile("^" + prefix + "bundle/(" + kBundleNameRe + ")/?$")
  downloadFileRe := regexp.MustCompile("^" + prefix + "bundle/(" + kBundleNameRe + ")/(" + kFileNameRe + ")/?$")
  uploadFileRe := regexp.MustCompile("^" + prefix + "upload/(" + kBundleNameRe + ")/?$")

  pathMatch := func(re *regexp.Regexp, req *http.Request) (bool) {
    return re.MatchString(req.URL.Path)
  }
  isGet := func(req *http.Request) (bool) {
    return req.Method == http.MethodGet
  }
  isPost := func(req *http.Request) (bool) {
    return req.Method == http.MethodPost
  }


  return func(w http.ResponseWriter, req *http.Request) {

    requestPath := req.URL.Path
    logDebug("Requested: %s", requestPath)

    if pathMatch(rootRe, req) && isGet(req) {
      logInfo("List bundles")

    } else if (listBundleRe.MatchString(requestPath)) && (isGet(req)) {
      bundleName := listBundleRe.FindStringSubmatch(requestPath)[1]
      logInfo("List bundle: %s", bundleName)

    } else if (downloadFileRe.MatchString(requestPath)) && (isGet(req)) {
      matches := downloadFileRe.FindStringSubmatch(requestPath)
      bundleName, fileName := matches[1], matches[2]
      logInfo("Download %s/%s", bundleName, fileName)

    } else if (uploadFileRe.MatchString(requestPath)) && (isGet(req)) {
      bundleName := uploadFileRe.FindStringSubmatch(requestPath)[1]
      logInfo("Upload form for bundle: %s", bundleName)

    } else if (uploadFileRe.MatchString(requestPath)) && (isPost(req)) {
      bundleName := uploadFileRe.FindStringSubmatch(requestPath)[1]
      logInfo("Upload file to bundle: %s", bundleName)

    } else {
      logWarn("Unhandled request: %s", requestPath)
      http.NotFound(w, req)
    }
  }
}

func WebAppStart(bindAddr, prefix string)  {

  // Absolute and non-empty
  if prefix == "" || prefix[0] != '/' {
    prefix = "/" + prefix
  }
  // Trailing slash (required for matching rules)
  if prefix[len(prefix)-1] != '/' {
    prefix = prefix + "/"
  }

  templates.listBundles = template.Must(template.New("list_bundles").Parse(listBundlesHtml))

  http.HandleFunc(prefix, makeHandler(prefix))

  fmt.Printf("Starting server @ %s\n", bindAddr)
  http.ListenAndServe(bindAddr, nil)
}
