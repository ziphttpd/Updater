package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	fpath "path/filepath"

	"github.com/xorvercom/util/pkg/json"
	"github.com/xorvercom/util/pkg/zip"
	"github.com/ziphttpd/zhsig/pkg/zhsig"
)

const (
	ziphttpdexe = "ziphttpd.exe"
	ziphttpdold = "ziphttpd.old"
	zhgetexe    = "zhget.exe"
	zhsigexe    = "zhsig.exe"
)

func main() {
	// 実行ファイルのディレクトリ
	exe, _ := os.Executable()
	dir := fpath.Dir(exe)
	err := zhsig.TempSpace(func(tempdir string) error {
		// ziphttpd.com のファイルへのアクセスの準備
		filename, err := download(dir)
		if err != nil {
			return err
		}
		dic, err := zip.OpenDictionary(filename, true)
		if err != nil {
			return err
		}
		defer dic.Close()

		files := []string{
			ziphttpdexe,
			zhgetexe,
			zhsigexe,
		}

		// 取り出し
		for _, name := range files {
			exetemp := fpath.Join(tempdir, name)
			rc, err := dic.GetReader(ziphttpdexe)
			if err != nil {
				return err
			}
			b, err := ioutil.ReadAll(rc)
			rc.Close()
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(exetemp, b, 0644)
			if err != nil {
				return err
			}
		}
		// コピー
		for _, name := range files {
			exetemp := fpath.Join(tempdir, name)
			copyname := fpath.Join(dir, name+"~")
			if err != nil {
				return err
			}
			b, err := ioutil.ReadFile(exetemp)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(copyname, b, 0644)
			if err != nil {
				return err
			}
		}
		// リネーム
		for _, name := range files {
			copyname := fpath.Join(dir, name+"~")
			exename := fpath.Join(dir, name)
			err = os.Rename(copyname, exename)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

// プログラムをダウンロード
func download(dir string) (string, error) {
	host := zhsig.NewHost(dir, "ziphttpd.com")
	clientsig := host.SigFile("client")
	sigelem, err := json.LoadFromJSONFile(clientsig)
	if err != nil {
		return "", err
	}
	if cliurl, ok := json.QueryElemString(sigelem, "url"); ok {
		if url, err := url.Parse(cliurl.Text()); err == nil {
			clientname := fpath.Base(url.Path)
			return host.File("windows", clientname), nil
		}
	}
	return "", fmt.Errorf("no file")
}
