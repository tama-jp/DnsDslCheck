<script lang="ts">
    import { GetSSLCertificateExpiry, GetParseCertificateDatesFromPEM } from '$lib/wailsjs/go/main/App';
    import { onMount } from "svelte";

    let activeTab: string ; // 現在のアクティブなタブ

    let inputDNS: string;
    let inputPort: string;
    let outputText: string;

    let inputText: string;
    let outputText2: string;

    onMount(() => {
        activeTab = "dns";
        inputDNS = '';
        inputPort = '443';
        outputText = '';
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
        <button class:active={activeTab === "dns"} on:click={() => activeTab = "dns"}>DNS情報取得</button>
        <button class:active={activeTab === "pem"} on:click={() => activeTab = "pem"}>証明書解析</button>
    </div>

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
