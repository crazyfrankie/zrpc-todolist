#!/bin/bash

set -e

# 默认配置值 - 根据项目需要修改这些默认值
DEFAULT_APP_NAME="goim"
DEFAULT_APP_CODE="1"
DEFAULT_IMPORT_PATH="github.com/crazyfrankie/zrpc-todolist/pkg/errorx/code"
DEFAULT_OUTPUT_DIR="types/errno"
DEFAULT_SCRIPT_DIR="scripts/error"

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# 代码生成器路径
CODE_GEN_PATH="scripts/error/code_gen.go"

# 显示使用说明
show_usage() {
    echo "用法: $0 [选项]"
    echo ""
    echo "Robot Web 项目错误码生成工具的便捷包装器"
    echo ""
    echo "选项:"
    echo "  --biz <名称>              业务域名称 (必需)"
    echo "  --app-name <名称>         应用名称 (默认: $DEFAULT_APP_NAME)"
    echo "  --app-code <代码>         应用代码 (默认: $DEFAULT_APP_CODE)"
    echo "  --import-path <路径>      错误码包的导入路径 (默认: $DEFAULT_IMPORT_PATH)"
    echo "  --output-dir <目录>       输出目录 (默认: $DEFAULT_OUTPUT_DIR)"
    echo "  --script-dir <目录>       脚本目录 (默认: $DEFAULT_SCRIPT_DIR)"
    echo "  --help, -h                显示此帮助信息"
    echo ""
    echo "使用示例:"
    echo "  # 为 admin 域生成错误码 (使用默认配置)"
    echo "  $0 --biz admin"
    echo ""
    echo "  # 为 auth 域生成错误码，自定义应用代码"
    echo "  $0 --biz auth --app-code 4"
    echo ""
    echo "  # 为所有域生成错误码"
    echo "  $0 --biz \"*\""
    echo ""
    echo "  # 完整参数示例 (等同于原始长命令)"
    echo "  $0 --biz task --app-name todolist --app-code 1 \\"
    echo "     --import-path \"github.com/transairobot/robot-web/pkg/errorx/code\" \\"
    echo "     --output-dir \"./types/errno\" --script-dir \"./scripts/errorx\""
    echo ""
    echo "可用的业务域:"
    echo "  - common, auth, robot, application, upload, admin"
    echo "  - \"*\" (所有域)"
}

# 检查是否有 errorgen 命令可用
check_errorgen_command() {
    if command -v gen >/dev/null 2>&1; then
        return 0  # errorgen 命令存在
    else
        return 1  # errorgen 命令不存在
    fi
}

# 检查代码生成器是否存在
check_code_generator() {
    if [[ ! -f "$CODE_GEN_PATH" ]]; then
        echo "错误: 在 $CODE_GEN_PATH 找不到代码生成器"
        echo "请确保错误码生成工具已正确安装。"
        exit 1
    fi
}

# 解析命令行参数
BIZ=""
APP_NAME="$DEFAULT_APP_NAME"
APP_CODE="$DEFAULT_APP_CODE"
IMPORT_PATH="$DEFAULT_IMPORT_PATH"
OUTPUT_DIR="$DEFAULT_OUTPUT_DIR"
SCRIPT_DIR_ARG="$DEFAULT_SCRIPT_DIR"

while [[ $# -gt 0 ]]; do
    case $1 in
        --biz)
            BIZ="$2"
            shift 2
            ;;
        --app-name)
            APP_NAME="$2"
            shift 2
            ;;
        --app-code)
            APP_CODE="$2"
            shift 2
            ;;
        --import-path)
            IMPORT_PATH="$2"
            shift 2
            ;;
        --output-dir)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        --script-dir)
            SCRIPT_DIR_ARG="$2"
            shift 2
            ;;
        --help|-h)
            show_usage
            exit 0
            ;;
        *)
            echo "错误: 未知选项 $1"
            echo "使用 --help 查看使用说明。"
            exit 1
            ;;
    esac
done

# 检查是否提供了业务域
if [[ -z "$BIZ" ]]; then
    echo "错误: 必须指定业务域 (--biz)"
    echo "使用 --help 查看使用说明。"
    exit 1
fi

# 转换相对路径为绝对路径
if [[ "$OUTPUT_DIR" != /* ]]; then
    OUTPUT_DIR="$PROJECT_ROOT/$OUTPUT_DIR"
fi

if [[ "$SCRIPT_DIR_ARG" != /* ]]; then
    SCRIPT_DIR_ARG="$PROJECT_ROOT/$SCRIPT_DIR_ARG"
fi

# 显示配置信息
echo "=== 错误码生成配置 ==="
echo "业务域: $BIZ"
echo "应用名称: $APP_NAME"
echo "应用代码: $APP_CODE"
echo "导入路径: $IMPORT_PATH"
echo "输出目录: $OUTPUT_DIR"
echo "脚本目录: $SCRIPT_DIR_ARG"
echo "====================="
echo ""

# 切换到项目根目录
cd "$PROJECT_ROOT"

# 构建命令参数
CMD_ARGS=(
    "--biz" "$BIZ"
    "--app-name" "$APP_NAME"
    "--app-code" "$APP_CODE"
    "--import-path" "$IMPORT_PATH"
    "--output-dir" "$OUTPUT_DIR"
    "--script-dir" "$SCRIPT_DIR_ARG"
)

# 检查是否有 gen 命令，如果没有则使用 go run 方式
if check_errorgen_command; then
    # 使用 gen 命令
    echo "执行命令: gen ${CMD_ARGS[*]}"
    echo ""

    if gen "${CMD_ARGS[@]}"; then
        echo ""
        echo "✅ 错误码生成成功!"

        if [[ "$BIZ" == "*" ]]; then
            echo "已为所有业务域生成错误码。"
        else
            echo "已为业务域 '$BIZ' 生成错误码。"
        fi

        echo "输出目录: $OUTPUT_DIR"
    else
        echo ""
        echo "❌ 错误码生成失败!"
        exit 1
    fi
else
    # 使用 go run 方式
    check_code_generator

    # 重新构建参数，第一个参数是业务域名称
    GO_CMD_ARGS=(
        "--biz" "$BIZ"
        "--app-name" "$APP_NAME"
        "--app-code" "$APP_CODE"
        "--import-path" "$IMPORT_PATH"
        "--output-dir" "$OUTPUT_DIR"
        "--script-dir" "$SCRIPT_DIR_ARG"
    )

    echo "执行命令: go run $CODE_GEN_PATH ${GO_CMD_ARGS[*]}"
    echo ""

    if go run "$CODE_GEN_PATH" "${GO_CMD_ARGS[@]}"; then
        echo ""
        echo "✅ 错误码生成成功!"

        if [[ "$BIZ" == "*" ]]; then
            echo "已为所有业务域生成错误码。"
        else
            echo "已为业务域 '$BIZ' 生成错误码。"
        fi

        echo "输出目录: $OUTPUT_DIR"
    else
        echo ""
        echo "❌ 错误码生成失败!"
        exit 1
    fi
fi
