package commands

import (
  "fmt"
  "io"
  "net/http"
  "html/template"
  "regexp"
  "github.com/c2h5oh/datasize"
)

import _ "embed"

//go:embed web/list_bundles.html
var listBundlesHtml string

//go:embed web/list_bundle.html
var listBundleHtml string


var templates struct {
  listBundles *template.Template
  listBundle  *template.Template
}

func handleListBundles(w http.ResponseWriter, req *http.Request) {

  bundles, err := loadBundles()
  if err != nil {
    logWarn("Error loading bundles: %s", err)
    http.Error(w, fmt.Sprintf("Internal error: %s", err), http.StatusInternalServerError)
  }

  type bundleListItem struct {
    Name      string
    NumFiles  int
    Size      string
    Href      string
  }

  type bundlesList struct {
    Bundles     []*bundleListItem
  }

  bl := &bundlesList{
    Bundles: []*bundleListItem{},
  }

  for _, bundle := range bundles {
    b := &bundleListItem{
      Name: bundle.name,
      NumFiles: len(bundle.files),
      Size: datasize.ByteSize(bundle.objStoreStatus.Size()).HumanReadable(),
      Href: fmt.Sprintf("./%s", bundle.name),
    }
    bl.Bundles = append(bl.Bundles, b)
  }

  err = templates.listBundles.Execute(w, bl)
  if err != nil {
    http.Error(w, fmt.Sprintf("Internal error: %s", err), http.StatusInternalServerError)
    return
  }
}

func handleListBundle(w http.ResponseWriter, req *http.Request, bundleName string) {

  bundle, err := loadBundle(bundleName)
  if err == kErrBundleNotFound {
    http.NotFound(w, req)
    return
  } else if err != nil {
    logWarn("Error loading bundle: %s", err)
    http.Error(w, fmt.Sprintf("Internal error: %s", err), http.StatusInternalServerError)
    return
  }

  type _bundleFile struct {
    Name      string
    Size      string
    Href      string
  }

  type _bundle struct {
    Name      string
    Files     []*_bundleFile
  }

  b := _bundle{
    Name: bundle.name,
    Files: []*_bundleFile{},
  }

  for _, file := range bundle.files {
    f := &_bundleFile{
      Name: file.fileName,
      Size: datasize.ByteSize(file.objInfo.Size).HumanReadable(),
      Href: fmt.Sprintf("./%s/%s", bundle.name, file.fileName),
    }
    b.Files = append(b.Files, f)
  }

  err = templates.listBundle.Execute(w, b)
  if err != nil {
    http.Error(w, fmt.Sprintf("Internal error: %s", err), http.StatusInternalServerError)
    return
  }
}

func handleBundleFileDownload(w http.ResponseWriter, req *http.Request, bundleName, fileName string) {
  bFile, err := loadBundleFile(bundleName, fileName)
  if err == kErrBundleNotFound || err == kErrBundleFileNotFound {
    logWarn("Error loading bundle file: %s", err)
    http.NotFound(w, req)
    return
  }

  reader, err := getBundleFileReader(bFile)
  if err != nil {
    logWarn("Error getting bundle file reader: %s", err)
    http.Error(w, fmt.Sprintf("Internal error: %s", err), http.StatusInternalServerError)
    return
  }

  w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", bFile.fileName))
  w.Header().Add("Content-Length", fmt.Sprintf("%d", bFile.objInfo.Size))

  io.Copy(w, reader)
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
      logDebug("List bundles")
      handleListBundles(w, req)

    } else if (listBundleRe.MatchString(requestPath)) && (isGet(req)) {
      bundleName := listBundleRe.FindStringSubmatch(requestPath)[1]
      logDebug("List bundle: %s", bundleName)
      handleListBundle(w, req, bundleName)

    } else if (downloadFileRe.MatchString(requestPath)) && (isGet(req)) {
      matches := downloadFileRe.FindStringSubmatch(requestPath)
      bundleName, fileName := matches[1], matches[2]
      logDebug("Download %s/%s", bundleName, fileName)
      handleBundleFileDownload(w, req, bundleName, fileName)

    } else if (uploadFileRe.MatchString(requestPath)) && (isGet(req)) {
      bundleName := uploadFileRe.FindStringSubmatch(requestPath)[1]
      logInfo("Upload form for bundle: %s", bundleName)

    } else if (uploadFileRe.MatchString(requestPath)) && (isPost(req)) {
      bundleName := uploadFileRe.FindStringSubmatch(requestPath)[1]
      logInfo("Upload file to bundle: %s", bundleName)

    } else {
      logWarn("Unhandled request: %s", requestPath)
      http.Error(w, "Invalid request", http.StatusBadRequest)
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
  templates.listBundle  = template.Must(template.New("list_bundle").Parse(listBundleHtml))

  http.HandleFunc(prefix, makeHandler(prefix))

  fmt.Printf("Starting server @ %s\n", bindAddr)
  http.ListenAndServe(bindAddr, nil)
}
