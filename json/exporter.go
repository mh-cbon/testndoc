package json

import (
	"encoding/json"
	"os"

	"github.com/mh-cbon/testndoc"
)

// Export the documentation to MD format.
func Export(recorded testndoc.APIDoc, dest string) error {
	j := &JExport{EndPoints: map[string]map[string]testndoc.APIRequestReponse{}}

	for _, ep := range recorded.EndPoints {
		for _, r := range ep.ReqRes {
			_, Type, Func := r.TestFn.Split()
			y := Func
			if Type != "" {
				y = Type + "." + Func
			}
			j.Add(ep.ParameterizedPath, y, r)
		}
	}

	os.Remove(dest)
	f, err2 := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0755)
	if err2 != nil {
		return err2
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")
	return enc.Encode(j)
}

type JExport struct {
	EndPoints map[string]map[string]testndoc.APIRequestReponse
}

func (j JExport) Add(path, testHandler string, r testndoc.APIRequestReponse) {
	if _, ok := j.EndPoints[path]; ok == false {
		j.EndPoints[path] = map[string]testndoc.APIRequestReponse{}
	}
	j.EndPoints[path][testHandler] = r
}
