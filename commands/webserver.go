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
//go:embed web/upload_form.html
var uploadFormHtml string
//go:embed web/upload_completed.html
var uploadCompletedHtml string

var templates struct {
  listBundles         *template.Template
  listBundle          *template.Template
  uploadForm          *template.Template
  uploadCompleted     *template.Template
}

var kUploadFormFileInputs = []string{"file1", "file2", "file3"} // TODO numFileInputs int

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

func handleUploadForm(w http.ResponseWriter, req *http.Request, bundleName string) {

  bundle, err := loadBundle(bundleName)
  if err == kErrBundleNotFound {
    http.NotFound(w, req)
    return
  } else if err != nil {
    logWarn("Error loading bundle: %s", err)
    http.Error(w, fmt.Sprintf("Internal error: %s", err), http.StatusInternalServerError)
    return
  }

  type _bundleUpload struct {
    Name            string
    Action          string
    Method          string
    FileInputs      []string
  }

  bu := _bundleUpload{
    Name: bundle.name,
    Action: fmt.Sprintf("/upload/%s", bundle.name),
    Method: "POST",
    FileInputs: kUploadFormFileInputs,
  }

  err = templates.uploadForm.Execute(w, bu)
  if err != nil {
    http.Error(w, fmt.Sprintf("Internal error: %s", err), http.StatusInternalServerError)
    return
  }
}

func handleFilesUpload(w http.ResponseWriter, req *http.Request, bundleName string) {

  bundle, err := loadBundle(bundleName)
  if err == kErrBundleNotFound {
    http.NotFound(w, req)
    return
  } else if err != nil {
    logWarn("Error loading bundle: %s", err)
    http.Error(w, fmt.Sprintf("Internal error: %s", err), http.StatusInternalServerError)
    return
  }

  uploadedFileNames := []string{}
  totalSize := uint64(0)

  for _, fileInput := range kUploadFormFileInputs {
    fileContent, fileHeaders, err := req.FormFile(fileInput)
    if err == http.ErrMissingFile {
      // File field unused
      continue
    } else if err != nil {
      http.Error(w, fmt.Sprintf("Invalid request: %s", err), http.StatusBadRequest)
      return
    }

    fileName := fileHeaders.Filename

    bundleFile, err := _addFileToBundle(bundle, fileContent, fileName)
    if err != nil {
      http.Error(w, fmt.Sprintf("Error uploading file: %s", err), http.StatusInternalServerError)
      return
    }
    uploadedFileNames = append(uploadedFileNames, fileName)
    totalSize += bundleFile.objInfo.Size
  }

  type _uploadInfo struct {
    Name      string
    NumFiles  int
    TotalSize string
  }

  ui := _uploadInfo{
    Name: bundleName,
    NumFiles: len(uploadedFileNames),
    TotalSize: datasize.ByteSize(totalSize).HumanReadable(),
  }

  err = templates.uploadCompleted.Execute(w, ui)
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
      logDebug("Upload form for bundle: %s", bundleName)
      handleUploadForm(w, req, bundleName)

    } else if (uploadFileRe.MatchString(requestPath)) && (isPost(req)) {
      bundleName := uploadFileRe.FindStringSubmatch(requestPath)[1]
      logDebug("Upload file(s) to bundle: %s", bundleName)
      handleFilesUpload(w, req, bundleName)

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

  templates.listBundles      = template.Must(template.New("list_bundles").Parse(listBundlesHtml))
  templates.listBundle       = template.Must(template.New("list_bundle").Parse(listBundleHtml))
  templates.uploadForm       = template.Must(template.New("upload_form").Parse(uploadFormHtml))
  templates.uploadCompleted  = template.Must(template.New("upload_completed").Parse(uploadCompletedHtml))

  http.HandleFunc(prefix, makeHandler(prefix))

  fmt.Printf("Starting server @ %s\n", bindAddr)
  http.ListenAndServe(bindAddr, nil)
}
