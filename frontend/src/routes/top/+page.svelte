<script lang="ts">
    import {GetSSLCertificateExpiry, GetParseCertificateDatesFromPEM, GetDNSCheck} from '$lib/wailsjs/go/main/App';
    import {onMount} from "svelte";

    let activeTab: string; // 現在のアクティブなタブ

    let inputDNS: string;
    let inputPort: string;
    let outputText: string;

    let inputText: string;
    let inputText3: string;
    let inputText4: string;
    let outputText2: string;
    let outputText3: string;
    let outputText4: string;

    onMount(() => {
        activeTab = "dns_list2";
        inputText = '';
        inputText3 = '';
        inputText4 = '';
        inputDNS = '';
        inputPort = '443';
        outputText = '';
        outputText2 = '';
        outputText3 = '';
        outputText4 = '';
    });

    function getInfo(): void {
        console.log(inputDNS);
        console.log(inputPort);

        let inputText = inputDNS + ":" + inputPort;
        GetSSLCertificateExpiry(inputText).then(result => {
            console.log(result);
            outputText = result;
        });
    }

    function getInfo2(): void {
        console.log(inputText);

        GetParseCertificateDatesFromPEM(inputText).then(result => {
            console.log(result);
            outputText2 = result;
        });
    }

    function getInfo4(): void {
        console.log(inputText4);

        GetDNSCheck(inputText4).then(result => {
            console.log(result);

            // JSONを整形して出力用テキストに設定
            try {
                // JSON文字列をオブジェクトにパース
                const parsedResult = JSON.parse(result);

                outputText4 = ""

                // 各エントリをループして処理
                parsedResult.forEach((entry: any) => {

                    outputText4 += `----------- [${entry.host}] -------------- \n`;
                    if (entry.status === "success" && entry.details) {
                        const details = entry.details;
                        outputText4 += `  Status: ${entry.status}\n`
                        outputText4 += `  Valid From (JST): ${details.NotBeforeJST}\n`
                        outputText4 += `  Valid Until (JST): ${details.NotAfterJST}\n`
                    } else if (entry.status === "error") {
                        outputText3 += `  Status: ${entry.status}\n`
                        outputText3 += `  Error: ${entry.error}\n`
                    }
                });
            } catch (error) {
                console.error("JSONの整形に失敗しました:", error);
                outputText3 = "エラー: JSONを整形できませんでした。";
            }
        }).catch(error => {
            console.error("エラーが発生しました:", error);
            outputText3 = `エラーが発生しました: ${error.message || error}`;
        });
    }

    function getInfo3(): void {
        console.log(inputText3);

        GetDNSCheck(inputText3).then(result => {
            console.log(result);

            // JSONを整形して出力用テキストに設定
            try {
                // JSON文字列をオブジェクトにパース
                const parsedResult = JSON.parse(result);

                outputText3 = ""

                // 各エントリをループして処理
                parsedResult.forEach((entry: any) => {

                    outputText3 += `----------- [${entry.host}] -------------- \n`;
                    if (entry.status === "success" && entry.details) {
                        const details = entry.details;
                        outputText3 += `  Status: ${entry.status}\n`
                        outputText3 += `  DNS Names: ${details.DNSNames?.join(", ")}\n`
                        outputText3 += `  Issuer: ${details.Issuer.CommonName}\n`
                        outputText3 += `  Valid From: ${details.NotBefore}\n`
                        outputText3 += `  Valid From (JST): ${details.NotBeforeJST}\n`
                        outputText3 += `  Valid Until: ${details.NotAfter}\n`
                        outputText3 += `  Valid Until (JST): ${details.NotAfterJST}\n`
                        outputText3 += `  Subject: ${details.Subject.CommonName}\n`
                        outputText3 += `  Serial Number: ${details.SerialNumber}\n`
                    } else if (entry.status === "error") {
                        outputText3 += `  Status: ${entry.status}\n`
                        outputText3 += `  Error: ${entry.error}\n`
                    }
                });
            } catch (error) {
                console.error("JSONの整形に失敗しました:", error);
                outputText3 = "エラー: JSONを整形できませんでした。";
            }
        }).catch(error => {
            console.error("エラーが発生しました:", error);
            outputText3 = `エラーが発生しました: ${error.message || error}`;
        });
    }


</script>

<style>
    .tab-buttons {
        display: flex;
        gap: 1rem;
        margin-bottom: 1rem;
    }

    .tab-buttons button {
        padding: 0.5rem 1rem;
        border: none;
        background-color: #e0e0e0;
        cursor: pointer;
        border-radius: 4px;
    }

    .tab-buttons button.active {
        background-color: #007bff;
        color: white;
    }

    .form-inline {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }

    .dns-input {
        flex: 2; /* DNSフィールドを大きく */
    }

    .port-input {
        flex: 1; /* Portフィールドを小さく */
        max-width: 80px; /* Portフィールドの最大幅 */
    }

    .button-container {
        margin-top: 1rem;
    }
</style>

<div>
    <!-- タブボタン -->
    <div class="tab-buttons">
        <button class:active={activeTab === "dns_list2"} on:click={() => activeTab = "dns_list2"}>DNS情報取得<br/>(list 期限のみ)</button>
        <button class:active={activeTab === "dns_list"} on:click={() => activeTab = "dns_list"}>DNS情報取得(list)</button>
        <button class:active={activeTab === "dns"} on:click={() => activeTab = "dns"}>DNS情報取得</button>
        <button class:active={activeTab === "pem"} on:click={() => activeTab = "pem"}>証明書解析</button>
    </div>

    {#if activeTab === "dns_list2"}
        <form>
            <div class="mb-3">
                <label class="form-label">DNS一覧 (文字列)</label>
                <textarea class="form-control" rows="6" bind:value={inputText4}></textarea>
            </div>

            <div class="button-container">
                <button type="button" class="btn btn-success" on:click={getInfo4}>取得</button>
            </div>

            <div class="mb-3">
                <label class="form-label">出力</label>
                <textarea class="form-control" rows="6" bind:value={outputText4}></textarea>
            </div>
        </form>
    {/if}

    {#if activeTab === "dns_list"}
        <form>
            <div class="mb-3">
                <label class="form-label">DNS一覧 (文字列)</label>
                <textarea class="form-control" rows="6" bind:value={inputText3}></textarea>
            </div>

            <div class="button-container">
                <button type="button" class="btn btn-success" on:click={getInfo3}>取得</button>
            </div>

            <div class="mb-3">
                <label class="form-label">出力</label>
                <textarea class="form-control" rows="6" bind:value={outputText3}></textarea>
            </div>
        </form>
    {/if}

    <!-- DNS情報取得フォーム -->
    {#if activeTab === "dns"}
        <form>
            <div class="form-inline">
                <input
                        class="form-control dns-input"
                        type="text"
                        placeholder="example google.co.jp"
                        aria-label="DNS"
                        bind:value={inputDNS}>
                <span>:</span>
                <input
                        class="form-control port-input"
                        type="text"
                        placeholder="port"
                        aria-label="port"
                        bind:value={inputPort}>
            </div>

            <div class="button-container">
                <button type="button" class="btn btn-success" on:click={getInfo}>取得</button>
            </div>

            <div class="mb-3">
                <label class="form-label">有効期限</label>
                <input
                        class="form-control"
                        type="text"
                        aria-label="DNS"
                        bind:value={outputText}>
            </div>
        </form>
    {/if}

    <!-- PEM形式証明書データフォーム -->
    {#if activeTab === "pem"}
        <form>
            <div class="mb-3">
                <label class="form-label">PEM形式の証明書データ (文字列)</label>
                <textarea class="form-control" rows="6" bind:value={inputText}></textarea>
            </div>

            <div class="button-container">
                <button type="button" class="btn btn-success" on:click={getInfo2}>取得</button>
            </div>

            <div class="mb-3">
                <label class="form-label">出力</label>
                <textarea class="form-control" rows="6" bind:value={outputText2}></textarea>
            </div>
        </form>
    {/if}
</div>
