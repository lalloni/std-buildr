package sql

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	isatty "github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/cert"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/context"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/credentials"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/httpe"
)

const defaultNexusURL = "https://nexus.cloudint.afip.gob.ar/nexus/repository"

func star(_ rune) rune {
	return '*'
}

func Publish(cfg *config.Config, ctx *context.Context) error {

	istty := isatty.IsTerminal(os.Stdout.Fd())

	if istty {
		log.SetHandler(cli.New(os.Stdout))
	}

	trust := viper.GetString("buildr.trust")

	// prepare http client
	pool, err := cert.DefaultTrustPool()
	if trust != "" {
		bytes, err := ioutil.ReadFile(trust)
		if err != nil {
			return errors.Wrapf(err, "reading trusted certificate chain from '%s'", trust)
		}
		if !pool.AppendCertsFromPEM(bytes) {
			return errors.Errorf("no certificate was found in '%s': aborting", trust)
		}
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: pool,
			},
		},
		Timeout: 15 * time.Minute,
	}

	// get credentials
	shouldask := credentials.NeverAsk
	if runtime.GOOS == "windows" || isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		shouldask = credentials.CanAsk
	}
	creds, err := credentials.GetUsernamePassword(shouldask, viper.GetString("nexus.username"), viper.GetString("nexus.password"))
	if err != nil {
		return errors.Wrapf(err, "getting credentials")
	}

	if creds.Username == "" {
		return errors.New("nexus username can not be empty")
	}
	if creds.Password == "" {
		return errors.New("nexus password can not be empty")
	}

	log.Debugf("using username %q and password %q", creds.Username, strings.Map(star, creds.Password))

	for _, file := range ctx.Artifacts {

		log.Infof("publishing %+v", file.Path)

		// open input file
		bodyFile, err := os.Open(file.Path)
		if err != nil {
			return errors.Wrapf(err, "opening '%s' for reading", file.Path)
		}

		var body io.Reader = bodyFile

		// wrap md5 & sha1 digesters
		md5sum := md5.New()
		body = io.TeeReader(body, md5sum)
		sha1sum := sha1.New()
		body = io.TeeReader(body, sha1sum)

		// put file
		log.Debug("putting file...")

		base := cfg.Nexus.URL
		if base == "" {
			base = defaultNexusURL
		}

		u := path.Join(base, cfg.SystemID+"-raw", cfg.SystemID, cfg.ApplicationID, ctx.Build.String(), file.File)

		log.Infof("uploading in %s", u)

		_, err = httpe.Put(client, creds, u, body, http.StatusCreated)
		if err != nil {
			return errors.Wrapf(err, "uploading file")
		}

		// put md5 digest
		log.Debug("putting md5 digest file...")
		body = strings.NewReader(hex.EncodeToString(md5sum.Sum(nil)))
		_, err = httpe.Put(client, creds, u+".md5", body, http.StatusCreated)
		if err != nil {
			return errors.Wrapf(err, "uploading md5 digest file")
		}

		// put sha1 digest
		log.Debug("putting sha1 digest file...")
		body = strings.NewReader(hex.EncodeToString(sha1sum.Sum(nil)))
		_, err = httpe.Put(client, creds, u+".sha1", body, http.StatusCreated)
		if err != nil {
			return errors.Wrapf(err, "uploading md5 digest file")
		}
		if istty {
			log.Info(fmt.Sprintf("Artifact %s succesfully uploaded", file.File))
		}

	}

	if istty {
		log.Info("all files succesfully uploaded")
	}

	return nil
}
