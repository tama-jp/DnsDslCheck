package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net"
	"net/url"
	"strings"
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

func (a *App) GetDNSCheck(host string) string {
	output := getDNSCheck(host)

	return output
}

/////////////////////////////////////////////////////////

func getDNSCheck(input string) string {
	// 分割して処理
	lines := strings.Split(input, "\n")
	var results []map[string]interface{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// URLを解析
		parsed, err := url.Parse(line)
		if err != nil || parsed.Host == "" {
			// ホスト名がない場合、直接使用
			if strings.Contains(line, ":") {
				line = line
			} else {
				line = line + ":443"
			}
		} else {
			// ポート番号が明示されていない場合を補完
			if !strings.Contains(parsed.Host, ":") {
				if parsed.Scheme == "http" {
					line = parsed.Host + ":80"
				} else {
					line = parsed.Host + ":443"
				}
			} else {
				line = parsed.Host
			}
		}

		// SSL証明書の詳細を取得
		details, err := getSSLCertificateDetails(line)
		result := map[string]interface{}{
			"host":   line,
			"status": "success",
		}
		if err != nil {
			result["status"] = "error"
			result["error"] = err.Error()
		} else {
			result["details"] = details
		}
		results = append(results, result)
	}

	// JSON形式に変換
	compressedJSON, err := json.Marshal(results)
	if err != nil {
		return fmt.Sprintf("JSON encode error: %v", err)
	}

	return string(compressedJSON)
}

func getSSLCertificateDetails(host string) (map[string]interface{}, error) {
	// タイムアウト設定付きの Dialer を作成
	dialer := &net.Dialer{
		Timeout: 10 * time.Second, // 10秒のタイムアウトを設定
	}

	// TCPコネクションを確立し、TLSハンドシェイクを実行
	conn, err := tls.DialWithDialer(dialer, "tcp", host, &tls.Config{
		InsecureSkipVerify: true, // 証明書の検証をスキップ（自己署名証明書の対応用）
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer func(conn *tls.Conn) {
		err := conn.Close()
		if err != nil {
			// エラーハンドリング（必要に応じて追加）
		}
	}(conn)

	// ピア証明書（サーバー証明書）を取得
	cert := conn.ConnectionState().PeerCertificates[0]

	// JSTタイムゾーンの設定
	jst := time.FixedZone("JST", 9*60*60)

	// 証明書の詳細情報をマップに格納
	certDetails := map[string]interface{}{
		"Subject":        cert.Subject,               // サブジェクト（発行先情報）
		"Issuer":         cert.Issuer,                // 発行者（CA情報）
		"NotBefore":      cert.NotBefore,             // 有効開始日時（UTC）
		"NotBeforeJST":   cert.NotBefore.In(jst),     // 有効開始日時（JST）
		"NotAfter":       cert.NotAfter,              // 有効期限（UTC）
		"NotAfterJST":    cert.NotAfter.In(jst),      // 有効期限（JST）
		"DNSNames":       cert.DNSNames,              // サブジェクト代替名（DNS名）
		"EmailAddresses": cert.EmailAddresses,        // サブジェクト代替名（Emailアドレス）
		"IPAddresses":    cert.IPAddresses,           // サブジェクト代替名（IPアドレス）
		"SerialNumber":   cert.SerialNumber.String(), // シリアル番号
		"PublicKeyAlgo":  cert.PublicKeyAlgorithm,    // 公開鍵アルゴリズム
		"SignatureAlgo":  cert.SignatureAlgorithm,    // 署名アルゴリズム
	}

	return certDetails, nil
}

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
