#!/bin/bash

set -e

APP_NAME="DnsDslCheck"      # アプリ名（適宜変更）
BUILD_DIR="build/bin"   # Wails のビルド出力フォルダ
DMG_DIR="dist"             # .dmg ファイルの保存フォルダ
DMG_NAME="${APP_NAME}.dmg" # 作成する .dmg の名前
VOLUME_NAME="${APP_NAME}"  # マウント時に表示されるボリューム名

# 1. Wails ビルド
echo "Building Wails application..."
wails build -platform darwin/universal
wails build -platform windows

# 2. 出力ディレクトリの確認と準備
echo "Preparing directories..."
mkdir -p "${DMG_DIR}"

# 3. .dmg 用の一時ディレクトリ作成
TMP_DIR=$(mktemp -d)
echo "Temporary directory created at: ${TMP_DIR}"

# 4. ビルドされたアプリケーションをコピー
#cp -R "${BUILD_DIR}/${APP_NAME}.app" "${TMP_DIR}"
cp -R "${BUILD_DIR}/${APP_NAME}.app" "${TMP_DIR}"

# 5. .dmg ファイルの作成
echo "Creating DMG file..."
hdiutil create \
  -volname "${VOLUME_NAME}" \
  -srcfolder "${TMP_DIR}" \
  -ov \
  -format UDZO \
  "${DMG_DIR}/${DMG_NAME}"

cp -R "${BUILD_DIR}/${APP_NAME}.exe" "${DMG_DIR}"

# 6. 一時ディレクトリの削除
echo "Cleaning up..."
rm -rf "${TMP_DIR}"

echo "DMG file created at: ${DMG_DIR}/${DMG_NAME}"
