package commands

import (
  "fmt"
  "net/http"
  "html/template"
  "github.com/c2h5oh/datasize"
)

import _ "embed"

//go:embed web/list_bundles.html
var listBundlesHtml string
var listBundlesTemplate *template.Template

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

  err = listBundlesTemplate.Execute(w, bl)
  if err != nil {
    http.Error(w, fmt.Sprintf("Internal error: %s", err), http.StatusInternalServerError)
    return
  }
}

func WebAppStart(bindAddr string)  {

  listBundlesTemplate = template.Must(template.New("list_bundles").Parse(listBundlesHtml))

  http.HandleFunc("/bundles/", listBundlesHandler)

  fmt.Printf("Starting server @ %s\n", bindAddr)
  http.ListenAndServe(bindAddr, nil)
}
