#!/bin/bash

# Record E2E Test Milestone Script
# This script handles milestone recording logic for E2E test progress tracking

set -e

# Parse command line arguments
MILESTONE_TYPE="$1"
SUCCESS_RATE="$2"
CURRENT_DATE="$3"
TOTAL_TESTS="$4"
PASSED_TESTS="$5"
DATETIME="$6"

# Validate required parameters
if [ -z "$MILESTONE_TYPE" ] || [ -z "$SUCCESS_RATE" ] || [ -z "$CURRENT_DATE" ] || [ -z "$TOTAL_TESTS" ] || [ -z "$PASSED_TESTS" ] || [ -z "$DATETIME" ]; then
    echo "Error: Missing required parameters"
    echo "Usage: $0 <milestone_type> <success_rate> <current_date> <total_tests> <passed_tests> <datetime>"
    exit 1
fi

# Generate milestone content based on type
case "$MILESTONE_TYPE" in
    "completion")
        EMOJI="🎉"
        TITLE="E2Eテスト100%達成"
        DESCRIPTION="### 🎯 完全な成功達成!" \
            $'\n\n' \
            "すべてのE2Eテストが成功し、" \
            "アプリケーションの品質が最高水準に達しました！" \
            $'\n\n' \
            "**達成指標:**" \
            $'\n' \
            "- ✅ 総テスト数: $TOTAL_TESTS" \
            $'\n' \
            "- ✅ 成功テスト: $PASSED_TESTS" \
            $'\n' \
            "- ✅ 失敗テスト: 0" \
            $'\n' \
            "- 🎯 **成功率: 100%**" \
            $'\n\n' \
            "**技術的成果:**" \
            $'\n' \
            "- フロントエンド品質の完全保証" \
            $'\n' \
            "- ユーザーワークフローの完全な動作確認" \
            $'\n' \
            "- リリース準備完了状態の確立"
        ;;
    "near_completion")
        EMOJI="🌟"
        TITLE="E2Eテスト95%以上達成"
        DESCRIPTION="### 🌟 優秀な品質レベル達成!" \
            $'\n\n' \
            "**達成指標:**" \
            $'\n' \
            "- ✅ 総テスト数: $TOTAL_TESTS" \
            $'\n' \
            "- ✅ 成功テスト: $PASSED_TESTS" \
            $'\n' \
            "- ⚠️ 残り修正: $((TOTAL_TESTS - PASSED_TESTS))" \
            $'\n' \
            "- 🎯 **成功率: $SUCCESS_RATE%**" \
            $'\n\n' \
            "**次のステップ:**" \
            $'\n' \
            "- 残りの少数テストの修正完了" \
            $'\n' \
            "- 100%達成への最終段階"
        ;;
    "excellent")
        EMOJI="✨"
        TITLE="E2Eテスト90%以上達成"
        DESCRIPTION="### ✨ 高品質レベル達成!" \
            $'\n\n' \
            "**達成指標:**" \
            $'\n' \
            "- ✅ 総テスト数: $TOTAL_TESTS" \
            $'\n' \
            "- ✅ 成功テスト: $PASSED_TESTS" \
            $'\n' \
            "- 🔧 修正対象: $((TOTAL_TESTS - PASSED_TESTS))" \
            $'\n' \
            "- 🎯 **成功率: $SUCCESS_RATE%**"
        ;;
    *)
        echo "Error: Unknown milestone type: $MILESTONE_TYPE"
        exit 1
        ;;
esac

# Create milestone entry
printf -v MILESTONE_ENTRY \
    "## %s %s - %s\n\n%s\n\n" \
    "### 自動記録情報\n- **測定日時**: %s\n" \
    "- **測定方法**: GitHub Actions自動実行\n" \
    "- **Playwright Version**: Latest\n\n---\n\n" \
    "$EMOJI" "$CURRENT_DATE" "$TITLE" "$DESCRIPTION" \
    "$DATETIME"

# Insert milestone into implementation log
if [ -f "docs/IMPLEMENTATION_LOG.md" ]; then
    if grep -n "## 📅 主要実装マイルストーン" docs/IMPLEMENTATION_LOG.md > /dev/null; then
        LINE_NUM=$(grep -n "## 📅 主要実装マイルストーン" docs/IMPLEMENTATION_LOG.md | cut -d: -f1)
        head -n $LINE_NUM docs/IMPLEMENTATION_LOG.md > temp_log.md
        echo "" >> temp_log.md
        echo "$MILESTONE_ENTRY" >> temp_log.md
        tail -n +$((LINE_NUM + 1)) docs/IMPLEMENTATION_LOG.md >> temp_log.md
        mv temp_log.md docs/IMPLEMENTATION_LOG.md
    else
        echo "" >> docs/IMPLEMENTATION_LOG.md
        echo "## 📅 主要実装マイルストーン" >> docs/IMPLEMENTATION_LOG.md
        echo "" >> docs/IMPLEMENTATION_LOG.md
        echo "$MILESTONE_ENTRY" >> docs/IMPLEMENTATION_LOG.md
    fi
    
    echo "🎉 E2E milestone recorded!"
else
    echo "Warning: docs/IMPLEMENTATION_LOG.md not found"
    exit 1
fi