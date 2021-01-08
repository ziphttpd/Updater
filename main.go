package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	fpath "path/filepath"

	"github.com/xorvercom/util/pkg/fileutil"
	"github.com/xorvercom/util/pkg/json"
	"github.com/xorvercom/util/pkg/zip"
	"github.com/ziphttpd/zhsig/pkg/zhsig"
)

const (
	ziphttpdexe = "ziphttpd.exe"
	ziphttpdold = "ziphttpd.old"
	zhgetexe    = "zhget.exe"
	zhsignexe   = "zhsign.exe"
	updaterexe  = "updater.exe"
	launcherexe = "launcher.exe"
)

func main() {
	// 実行ファイルのディレクトリ
	exe, _ := os.Executable()
	dir := fpath.Dir(exe)

	// zhget でダウンロード
	err := exec.Command(fpath.Join(dir, zhgetexe), "-host ziphttpd.com -group windows").Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 解凍用の一時フォルダ
	err = zhsig.TempSpace(func(tempdir string) error {
		// ziphttpd.com のファイルへのアクセスの準備
		filename, err := downloadedFile(dir)
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
			zhsignexe,
			updaterexe,
			launcherexe,
		}

		// クリーニング
		for _, name := range files {
			err = fileutil.FileIfDelete(fpath.Join(dir, name+".copy"))
			if err != nil {
				return err
			}
			err = fileutil.FileIfDelete(fpath.Join(dir, name+".old"))
			if err != nil {
				return err
			}
		}
		// zipから取り出し
		for _, name := range files {
			err = fileExtract(dic, name, fpath.Join(tempdir, name))
			if err != nil {
				return err
			}
		}
		// 適用
		for _, name := range files {
			// コピー
			err = fileutil.FileCopy(fpath.Join(tempdir, name), fpath.Join(dir, name+".copy"))
			if err != nil {
				return err
			}
			// 以前の実行ファイルをoldに
			err = fileutil.FileIfMove(fpath.Join(dir, name), fpath.Join(dir, name+".old"))
			if err != nil {
				// old を復旧
				for _, name := range files {
					fileutil.FileIfMove(fpath.Join(dir, name+".old"), fpath.Join(dir, name))
				}
				return err
			}
			// コピーした実行ファイルを新に
			err = fileutil.FileIfMove(fpath.Join(dir, name+".copy"), fpath.Join(dir, name))
			if err != nil {
				// old を復旧
				for _, name := range files {
					fileutil.FileIfMove(fpath.Join(dir, name+".old"), fpath.Join(dir, name))
				}
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
func downloadedFile(dir string) (string, error) {
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

// zipエントリをファイルに出力
func fileExtract(dic zip.Dictionary, entry, exetemp string) error {
	// 出力先を消去
	err := fileutil.FileIfDelete(exetemp)
	if err != nil {
		return err
	}
	// zipエントリを取得
	rc, err := dic.GetReader(entry)
	if err != nil {
		return err
	}
	defer rc.Close()
	// zipエントリを読み込み
	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return err
	}
	// ファイルに出力
	return ioutil.WriteFile(exetemp, b, 0644)
}
