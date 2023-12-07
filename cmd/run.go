package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Jiang-Gianni/ssd/parse"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RunRoot(cmd *cobra.Command, args []string) {
	if err := Run(os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func RunHTML(cmd *cobra.Command, args []string) {
	out := viper.GetString("out")

	outfile, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = outfile.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := Run(outfile); err != nil {
		log.Fatal(err)
	}
}

func RunServer(cmd *cobra.Command, args []string) {
	port := viper.GetString("port")
	timeout := 3 * time.Second
	handle := func(w http.ResponseWriter, r *http.Request) {
		if err := Run(w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}
	srv := http.Server{
		Addr:         port,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		IdleTimeout:  timeout,
		Handler:      http.HandlerFunc(handle),
	}

	fmt.Printf("http://localhost%s/\n", port)
	log.Fatal(srv.ListenAndServe())
}

func Run(w io.Writer) (err error) {
	sql := viper.GetString("sql")
	filenames := viper.GetStringSlice("md")
	entry := viper.GetString("entry")

	file, err := os.OpenFile(sql, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}

	defer func() {
		closeErr := file.Close()
		if err != nil {
			err = closeErr
		}
	}()

	schema := parse.Read(file)

	b := bytes.NewBuffer(nil)

	schema.Parse(filenames, entry, b)

	return parse.Markdown(b.Bytes(), w)

	// return schema.Parse(filenames, entry, w)
	// return nil
}
