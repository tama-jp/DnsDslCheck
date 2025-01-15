package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GetSSLCertificateExpiry(host string) string {
	expiry, err := getSSLCertificateExpiry(host)
	if err != nil {
		return "error" + err.Error()
	}

	// 日本標準時（JST）のLocationを取得
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return "error" + err.Error()
	}
	// JSTに変換
	jstTime := expiry.In(jst)

	return "JST:" + jstTime.Format("2006-01-02 15:04:05")
}

func (a *App) GetParseCertificateDatesFromPEM(certFile string) string {
	from, to, err := getParseCertificateDatesFromPEM(certFile)
	if err != nil {
		return "error" + err.Error()
	}

	// 日本標準時（JST）のLocationを取得
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return "error" + err.Error()
	}
	// 期待される日付
	fromJstTime := from.In(jst)
	toJstTime := to.In(jst)

	return "from(JST):" + fromJstTime.Format("2006-01-02 15:04:05") + "\n" +
		"to(JST):" + toJstTime.Format("2006-01-02 15:04:05")

}

/////////////////////////////////////////////////////////

// GetSSLCertificateExpiry retrieves the SSL certificate expiry date of a given host.
func getSSLCertificateExpiry(host string) (time.Time, error) {
	conn, err := tls.Dial("tcp", host, nil)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]
	return cert.NotAfter, nil
}

func getParseCertificateDatesFromPEM(certPEM string) (time.Time, time.Time, error) {
	// PEM ブロックのデコード
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil || block.Type != "CERTIFICATE" {
		return time.Time{}, time.Time{}, fmt.Errorf("failed to decode PEM blck")
	}

	// 証明書のパース
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("failed to parse certificate: %v", err)
	}

	// `from` と `to` を返す
	return cert.NotBefore, cert.NotAfter, nil
}
